package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	ModuleFlag    = "module"
	ModuleNameKey = "moduleName"
	ProjectMapDir = "usersecrets"
)

func GetBaseDirpath() string {
	return path.Join(build.Default.GOPATH, ProjectMapDir)
}

func GetModulename() string {
	return viper.GetString(ModuleNameKey)
}

func GetModuleDir() string {
	return path.Join(GetBaseDirpath(), GetModulename())
}

func GetModuleSecretspath() string {
	fmt.Println(path.Join(GetModuleDir(), "secrets.json"))
	return path.Join(GetModuleDir(), "secrets.json")
}

func ReadSecrets() (map[string]interface{}, error) {
	content, err := ioutil.ReadFile(GetModuleSecretspath())
	if err != nil {
		return nil, err
	}

	values := make(map[string]interface{})
	err = json.Unmarshal(content, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func DumpSecrets(secrets map[string]interface{}) (map[string]interface{}, error) {
	content, err := json.Marshal(secrets)
	if err != nil {
		return nil, err
	}

	return secrets, ioutil.WriteFile(GetModuleSecretspath(), content, os.ModePerm)
}

func InitializeTree(baseDir, valueFile string, value interface{}) error {
	if _, err := os.Stat(baseDir); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(baseDir, os.ModePerm)
	} else if err != nil {
		return err
	}

	if _, err := os.Stat(valueFile); errors.Is(err, os.ErrNotExist) {
		file, err := os.OpenFile(valueFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}

		value, err := json.Marshal(value)
		if err != nil {
			return err
		}

		file.Write(value)
	} else if err != nil {
		return err
	}

	return nil
}
