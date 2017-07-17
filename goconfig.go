package goconfig

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

const etcFolder = "/etc/"

// SearchForConfigFile - search for config file
func SearchForConfigFile(filename string, folder string, searchLocally bool, searchUser bool, searchSystem bool) (foundFullPath string, err error) {

	if len(folder) > 0 {
		if folder[:1] == "/" {
			folder = folder[1:]
		}
		if len(folder) > 0 {
			if folder[len(folder)-1:] == "/" {
				folder = folder[:len(folder)]
			}
		}
	}

	if searchLocally {
		// Checking in the current folder
		currentPath, err := GetLocalPath()
		if err == nil {
			fullPath := currentPath + filename
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath, nil
			}
			if len(folder) > 0 {
				fullPath = currentPath + folder + "/" + filename
				if _, err := os.Stat(fullPath); err == nil {
					return fullPath, nil
				}
			}
		}

	}

	if searchUser {
		// Checking in user home folder
		usrHome, err := GetUserHomePath()
		if err == nil {
			fullPath := usrHome + filename
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath, nil
			}
			if len(folder) > 0 {
				fullPath = usrHome + folder + "/" + filename
				if _, err := os.Stat(fullPath); err == nil {
					return fullPath, nil
				}
			}
		}
	}

	if searchSystem {
		// Checking in /etc
		fullPath := etcFolder + filename
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
		if len(folder) > 0 {
			fullPath = etcFolder + folder + "/" + filename
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath, nil
			}
		}
	}

	return "", errors.New("Error: config file is not found")
}

// GetLocalPath - get the path of current folder
func GetLocalPath() (string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if getLastChar(currentPath) != "/" {
		currentPath = currentPath + "/"
	}
	return currentPath, nil
}

// GetUserHomePath - get home folder of the current user
func GetUserHomePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	usrHome := usr.HomeDir
	if getLastChar(usrHome) != "/" {
		usrHome = usrHome + "/"
	}
	return usrHome, nil
}

func getLastChar(checkStr string) string {
	if len(checkStr) < 2 {
		return ""
	}
	return checkStr[len(checkStr)-1:]
}

func folderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// LoadConfig - load configuration file
func LoadConfig(configFile string, config interface{}) (err error) {

	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)

	err = decoder.Decode(config)
	if err != nil {
		return err
	}

	return nil
}

// SaveConfig - saving current config in the file
func SaveConfig(configFile string, config interface{}) error {
	folderPath := filepath.Dir(configFile)
	exists, err := folderExists(folderPath)
	if err != nil {
		return err
	}
	if exists == false {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(configFile)
	if err != nil {
		return err
	}

	defer file.Close()
	encoder := json.NewEncoder(file)

	err = encoder.Encode(&config)
	if err != nil {
		return err
	}

	return nil
}
