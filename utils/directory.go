package utils

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: PathExists
//@description: 文件目录是否存在
//@param: path string
//@return: bool, error

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FileExists 文件是否存在
func FileExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return false, nil
		}
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteFile(path string) error {
	has, err := FileExists(path)
	if err != nil {
		return err
	}
	if has {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func AppConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	appPath := filepath.Join(usr.HomeDir, "AppData", "Roaming")
	err = os.MkdirAll(appPath, 0666)
	if err != nil {
		return "", err
	}
	return appPath, nil
}
