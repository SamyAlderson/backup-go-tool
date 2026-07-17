// Package backup fournit une API pour l'outil de backup.
package backup

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Config représente la configuration de l'outil de backup.
type Config struct {
	BackupDir string `json:"backup_dir"`
	SyncDir   string `json:"sync_dir"`
	Docker    struct {
		ContainerName string `json:"container_name"`
	} `json:"docker"`
}

// BackupOption est une fonction qui prend une instance de Backup et renvoie une nouvelle instance de Backup.
type BackupOption func(*Backup) *Backup

// NewBackup crée une nouvelle instance de Backup.
func NewBackup(config Config, fs afero.FileSystem) *Backup {
	return &Backup{
		config: config,
		fs:     fs,
	}
}

// Backup est l'outil de backup principal.
type Backup struct {
	config Config
	fs     afero.FileSystem
	mu     sync.RWMutex
}

// NewBackupOption est une fonction qui prend une instance de Backup et renvoie une nouvelle instance de Backup avec une nouvelle configuration.
func NewBackupOption(config Config) BackupOption {
	return func(b *Backup) *Backup {
		b.config = config
		return b
	}
}

// StartStart le processus de backup.
func (b *Backup) Start() error {
	// Générer un nouveau certificat RSA pour la session de backup.
	cert, err := b.generateCert()
	if err != nil {
		return errors.Wrap(err, "erreur lors de la génération du certificat")
	}

	// Copier les fichiers du répertoire de synchronisation vers le répertoire de backup.
	if err := b.copyFiles(cert); err != nil {
		return errors.Wrap(err, "erreur lors de la copie des fichiers")
	}

	return nil
}

// generateCert génère un nouveau certificat RSA pour la session de backup.
func (b *Backup) generateCert() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	serialNumber, err := b.generateSerialNumber()
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		IsCA:         true,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// generateSerialNumber génère un nouveau numéro de série pour le certificat.
func (b *Backup) generateSerialNumber() (*big.Int, error) {
	serialNumber := big.NewInt(0)
	_, err := serialNumber.SetBytes(b.generateRandomBytes(16))
	if err != nil {
		return nil, err
	}

	return serialNumber, nil
}

// generateRandomBytes génère 16 octets aléatoires.
func (b *Backup) generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

// copyFiles copie les fichiers du répertoire de synchronisation vers le répertoire de backup.
func (b *Backup) copyFiles(cert *rsa.PrivateKey) error {
	// Répertoire de synchronisation.
	syncDir := b.config_SYNC_DIR

	// Répertoire de backup.
	backupDir := b.config.BACKUP_DIR

	// Vérifier l'existence des répertoires.
	if err := b.fs.MkdirAll(syncDir, 0755); err != nil {
		return err
	}

	if err := b.fs.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	// Copier les fichiers.
	files, err := b.fs.ReadDir(syncDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcPath := filepath.Join(syncDir, file.Name())
		dstPath := filepath.Join(backupDir, file.Name())

		if err := b.fs.Copy(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}