// Package backup_docker fournit des fonctionnalités pour intégrer l'outil de backup avec Docker.
package backup_docker

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/spf13/afero"
	"github.com/pkg/errors"
)

// DockerConfig représente la configuration de Docker pour l'outil de backup.
type DockerConfig struct {
	DockerHost string `json:"docker_host"`
	DockerCert string `json:"docker_cert"`
	DockerKey  string `json:"docker_key"`
}

// NewDockerConfig crée une nouvelle instance de DockerConfig à partir du fichier de configuration.
func NewDockerConfig(configPath string) (*DockerConfig, error) {
	config := &DockerConfig{}
	err := afero.ReadFile(afero.OsFs(), configPath, config)
	if err != nil {
		return nil, errors.Wrap(err, "Impossible de lire le fichier de configuration")
	}
	return config, nil
}

// ListContainers liste tous les conteneurs Docker actifs.
func ListContainers(ctx context.Context, client *client.Client) ([]container.Container, error) {
	filters := filters.NewArgs()
	filters.Add("status", "running")
	filters.Add("status", "paused")
	filters.Add("status", "exited")

	containers, err := client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.Value,
	})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// ListImages liste toutes les images Docker disponibles.
func ListImages(ctx context.Context, client *client.Client) ([]image.Image, error) {
	images, err := client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	return images, nil
}

// BackupConteneurs effectue un backup des conteneurs Docker actifs.
func BackupConteneurs(ctx context.Context, client *client.Client, config *DockerConfig) error {
	// List des conteneurs actifs
	containers, err := ListContainers(ctx, client)
	if err != nil {
		return err
	}

	// Création du répertoire de backup
	backupDir := "backup-conteneurs"
	err = os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Backup des conteneurs
	for _, container := range containers {
		// Récupération des données du conteneur
		containerData, err := client.ContainerInspect(ctx, container.ID)
		if err != nil {
			return err
		}

		// Sauvegarde des données dans un fichier
		containerPath := filepath.Join(backupDir, container.ID)
		err = afero.WriteFile(afero.OsFs(), containerPath, containerData, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

// BackupImages effectue un backup des images Docker disponibles.
func BackupImages(ctx context.Context, client *client.Client, config *DockerConfig) error {
	// List des images disponibles
	images, err := ListImages(ctx, client)
	if err != nil {
		return err
	}

	// Création du répertoire de backup
	backupDir := "backup-images"
	err = os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Backup des images
	for _, image := range images {
		// Récupération des données de l'image
		imageData, err := client.ImageSave(ctx, image.ID)
		if err != nil {
			return err
		}

		// Sauvegarde des données dans un fichier
		imagePath := filepath.Join(backupDir, image.ID)
		err = afero.WriteFile(afero.OsFs(), imagePath, imageData, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

// Main effectue un backup des conteneurs et des images Docker.
func Main(configPath string) error {
	// Lecture de la configuration
	config, err := NewDockerConfig(configPath)
	if err != nil {
		return err
	}

	// Création du client Docker
	client, err := client.NewClientWithOpts(client.FromEnvFile())
	if err != nil {
		return err
	}

	// Backup des conteneurs
	err = BackupConteneurs(context.Background(), client, config)
	if err != nil {
		return err
	}

	// Backup des images
	err = BackupImages(context.Background(), client, config)
	if err != nil {
		return err
	}

	return nil
}
```
```go
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
```
```go
func main() {
	if err := Main("backup_config.json"); err != nil {
		log.Fatal(err)
	}
}