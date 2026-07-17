package main

import (
	"testing"
	"path/filepath"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/afero"
	"github.com/pkg/errors"
)

func TestBackup(t *testing.T) {
	t.Run("backup de fichiers", testBackupFiles)
	t.Run("backup de conteneurs Docker", testBackupDocker)
}

func testBackupFiles(t *testing.T) {
	// Création d'un système de fichiers temporaire
	tempDir, err := ioutil.TempDir("", "backup-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = filepath.Walk(tempDir, func(path string, info fs.FileInfo, err error) error {
			return afero.RemoveAll(afero.OsFs, path)
		})
	}()

	// Création de fichiers temporaire
	err = afero.WriteFile(afero.OsFs, filepath.Join(tempDir, "fichier1.txt"), []byte("Contenu du fichier 1"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = afero.WriteFile(afero.OsFs, filepath.Join(tempDir, "fichier2.txt"), []byte("Contenu du fichier 2"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Appel de la fonction de backup
	err = backupFiles(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Vérification de la présence des fichiers de backup
	_, err = afero.ReadFile(afero.OsFs, filepath.Join(tempDir, "backup", "fichier1.txt"))
	if err != nil {
		t.Errorf("Le fichier de backup 'fichier1.txt' n'existe pas")
	}
	_, err = afero.ReadFile(afero.OsFs, filepath.Join(tempDir, "backup", "fichier2.txt"))
	if err != nil {
		t.Errorf("Le fichier de backup 'fichier2.txt' n'existe pas")
	}
}

func testBackupDocker(t *testing.T) {
	// Création d'un conteneur Docker temporaire
	container, err := runContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = container.Stop()
		_ = container.Remove()
	}()

	// Appel de la fonction de backup des conteneurs
	err = backupDocker(container.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Vérification de la présence des fichiers de backup
	_, err = afero.ReadFile(afero.OsFs, filepath.Join("/var/lib/docker/backup", container.ID, "fichier1.txt"))
	if err != nil {
		t.Errorf("Le fichier de backup 'fichier1.txt' n'existe pas")
	}
	_, err = afero.ReadFile(afero.OsFs, filepath.Join("/var/lib/docker/backup", container.ID, "fichier2.txt"))
	if err != nil {
		t.Errorf("Le fichier de backup 'fichier2.txt' n'existe pas")
	}
}

func backupFiles(rootDir string) error {
	// Logique de backup des fichiers
	return nil
}

func backupDocker(containerID string) error {
	// Logique de backup des conteneurs Docker
	return nil
}

func runContainer() (container *Container, err error) {
	// Logique de création d'un conteneur Docker
	return nil, nil
}

type Container struct {
	ID string
}