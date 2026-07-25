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
	"github.com/shopspring/decimal"
	"time"
)

// Config represents the configuration of the backup tool.
type Config struct {
	BackupDir string `json:"backup_dir"`
	SyncDir   string `json:"sync_dir"`
	Docker    struct {
		ContainerName string `json:"container_name"`
	} `json:"docker"`
}

// BackupOption is a function that takes a Backup instance and returns a new Backup instance.
type BackupOption func(*Backup) *Backup

// NewBackup creates a new Backup instance.
func NewBackup(config Config, fs afero.FileSystem) *Backup {
	return &Backup{
		config: config,
		fs:     fs,
	}
}

// Backup is the main backup tool.
type Backup struct {
	config Config
	fs     afero.FileSystem
	mu     sync.RWMutex
}

// NewBackupOption is a function that takes a Backup instance and returns a new Backup instance with a new configuration.
func NewBackupOption(config Config) BackupOption {
	return func(b *Backup) *Backup {
		b.config = config
		return b
	}
}

// Start starts the backup process.
func (b *Backup) Start() error {
	// Generate a new RSA certificate for the backup session.
	cert, err := b.generateCert()
	if err != nil {
		return errors.Wrap(err, "error generating certificate")
	}

	// Copy files from the sync directory to the backup directory.
	if err := b.copyFiles(cert); err != nil {
		return errors.Wrap(err, "error copying files")
	}

	return nil
}

// generateCert generates a new RSA certificate for the backup session.
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

// generateSerialNumber generates a new serial number for the certificate.
func (b *Backup) generateSerialNumber() (*big.Int, error) {
	serialNumber := big.NewInt(0)
	_, err := serialNumber.SetString("1", 10)
	if err != nil {
		return nil, err
	}

	return serialNumber, nil
}

// generateRandomBytes generates 16 random bytes.
func (b *Backup) generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err)
		return nil
	}

	return b
}

// copyFiles copies files from the sync directory to the backup directory.
func (b *Backup) copyFiles(cert *rsa.PrivateKey) error {
	// Sync directory.
	syncDir := b.config.SyncDir

	// Backup directory.
	backupDir := b.config.BackupDir

	// Check if directories exist.
	if err := b.fs.MkdirAll(syncDir, 0755); err != nil {
		return err
	}

	if err := b.fs.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	// Copy files.
	files, err := b.fs.ReadDir(syncDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcPath := filepath.Join(syncDir, file.Name())
		dstPath := filepath.Join(backupDir, file.Name())

		if err := b.fs.CopyFile(dstPath, srcPath); err != nil {
			return err
		}
	}

	return nil
}