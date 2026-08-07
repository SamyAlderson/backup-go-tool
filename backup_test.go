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
	t.Run("backup files", testBackupFiles)
	t.Run("backup docker container", testBackupDocker)
}

func testBackupFiles(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "backup-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := afero.RemoveAll(afero.OsFs, tempDir); err != nil {
			log.Printf("error removing temp dir: %v", err)
		}
	}()

	files := []struct {
		name string
		data []byte
	}{
		{"fichier1.txt", []byte("Contenu du fichier 1")},
		{"fichier2.txt", []byte("Contenu du fichier 2")},
	}

	for _, file := range files {
		if err := afero.WriteFile(afero.OsFs, filepath.Join(tempDir, file.name), file.data, 0644); err != nil {
			t.Fatal(err)
		}
	}

	if err := backupFiles(tempDir); err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		_, err := afero.ReadFile(afero.OsFs, filepath.Join(tempDir, "backup", file.name))
		if err != nil {
			t.Errorf("backup file '%s' not found", file.name)
		}
	}
}

func testBackupDocker(t *testing.T) {
	container, err := runContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := container.Stop(); err != nil {
			log.Printf("error stopping container: %v", err)
		}
		if err := container.Remove(); err != nil {
			log.Printf("error removing container: %v", err)
		}
	}()

	if err := backupDocker(container.ID); err != nil {
		t.Fatal(err)
	}

	files := []struct {
		name string
		path string
	}{
		{"fichier1.txt", filepath.Join("/var/lib/docker/backup", container.ID, "fichier1.txt")},
		{"fichier2.txt", filepath.Join("/var/lib/docker/backup", container.ID, "fichier2.txt")},
	}

	for _, file := range files {
		_, err := afero.ReadFile(afero.OsFs, file.path)
		if err != nil {
			t.Errorf("backup file '%s' not found at '%s'", file.name, file.path)
		}
	}
}

func backupFiles(rootDir string) error {
	// implement backup logic here
	return nil
}

func backupDocker(containerID string) error {
	// implement backup logic here
	return nil
}

func runContainer() (*Container, error) {
	// implement container creation logic here
	return &Container{ID: "container-id"}, nil
}

type Container struct {
	ID string
}