// Package backup_config contient les constantes et les fonctions pour la configuration de l'outil de backup.
package backup_config

import (
	"fmt"
	"log"

	"github.com/spf13/afero"
	"github.com/pkg/errors"
)

// Configuration est une structure pour stocker les configurations de l'outil de backup.
type Configuration struct {
	BackupDir       string   // Répertoire de stockage des backups
	BackupFrequency string   // Fréquence des backups (e.g. quotidien, hebdomadaire, mensuel)
	BackupRetention int      // Durée de conservation des backups
	FilesToBackup   []string // Liste de fichiers à sauvegarder
}

// NewConfiguration retourne une nouvelle instance de Configuration.
func NewConfiguration() *Configuration {
	return &Configuration{
		BackupDir:       "/backup",
		BackupFrequency: "quotidien",
		BackupRetention: 30,
		FilesToBackup:   []string{"fichier1.txt", "fichier2.txt"},
	}
}

// LoadConfigurationFromEnv charge les configurations à partir des variables d'environnement.
func LoadConfigurationFromEnv() (*Configuration, error) {
	backupDir := afero.FromEnv("BACKUP_DIR")
	if backupDir == "" {
		return nil, errors.New("BACKUP_DIR est obligatoire")
	}

	backupFrequency := afero.FromEnv("BACKUP_FREQUENCY")
	if backupFrequency == "" {
		backupFrequency = "quotidien"
	}

	backupRetention := afero.AtoiFromEnv("BACKUP_RETENTION")
	if backupRetention <= 0 {
		return nil, errors.New("BACKUP_RETENTION doit être supérieur à 0")
	}

	filesToBackup := afero.FromEnv("FILES_TO_BACKUP")
	if filesToBackup == "" {
		filesToBackup = "fichier1.txt,fichier2.txt"
	}

	return &Configuration{
		BackupDir:       backupDir,
		BackupFrequency: backupFrequency,
		BackupRetention: backupRetention,
		FilesToBackup:   strings.Split(filesToBackup, ","),
	}, nil
}

// LoadConfigurationFromFile charge les configurations à partir d'un fichier JSON.
func LoadConfigurationFromFile(fs afero.Fs) (*Configuration, error) {
	data, err := afero.ReadFile(fs, "/path/to/config.json")
	if err != nil {
		return nil, err
	}

	var config Configuration
	err = afero.UnmarshalJSON(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func init() {
	// Chargement des configurations par défaut
	config := NewConfiguration()
	log.Printf("Configuration par défaut : %v", config)
}