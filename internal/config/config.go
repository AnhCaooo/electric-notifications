// AnhCao 2024
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/AnhCaooo/electric-notifications/internal/constants"
	"github.com/AnhCaooo/electric-notifications/internal/helpers"
	"github.com/AnhCaooo/electric-notifications/internal/models"
)

var Config models.Config

// load the configuration from the yaml config file
func ReadFile(cfg *models.Config) error {
	currentDir, err := helpers.GetCurrentDir()
	if err != nil {
		return err
	}
	encryptedConfigFilePath := currentDir + constants.EncryptedConfigFile
	decryptedConfigFilePath := currentDir + constants.DecryptedConfigFile

	if err = helpers.DecryptFile(encryptedConfigFilePath, decryptedConfigFilePath); err != nil {
		return err
	}

	f, err := os.Open(decryptedConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to open config.yml: %s", err.Error())
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return fmt.Errorf("failed to decode config.yml: %s", err.Error())
	}
	return nil
}

// Decrypt the encrypted Firebase configuration and return decrypted config path
func DecryptFirebaseKeyFile() (string, error) {
	currentDir, err := helpers.GetCurrentDir()
	if err != nil {
		return "", err
	}
	encryptedFirebaseFilePath := currentDir + constants.FirebaseKeyEncryptedFile
	decryptedFirebaseFilePath := currentDir + constants.FirebaseKeyDecryptedFile

	if err = helpers.DecryptFile(encryptedFirebaseFilePath, decryptedFirebaseFilePath); err != nil {
		return "", err
	}

	return decryptedFirebaseFilePath, nil
}
