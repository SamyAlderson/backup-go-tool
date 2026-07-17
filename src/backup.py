"""
Fichier src/backup.py

Fichier gère la création des archives de sauvegarde.
"""

import os
import shutil
import logging
import tarfile
from datetime import datetime

from storage import Storage

logger = logging.getLogger(__name__)

class Backup:
    """
    Classe gérant la création des archives de sauvegarde.
    """

    def __init__(self, storage: Storage, backup_dir: str):
        """
        Initialise la classe Backup.

        Args:
            storage (Storage): Objet de stockage de sauvegarde.
            backup_dir (str): Chemin du répertoire de sauvegarde.
        """
        self.storage = storage
        self.backup_dir = backup_dir

    def create_backup(self, backup_name: str, files: list[str]) -> None:
        """
        Crée une archive de sauvegarde.

        Args:
            backup_name (str): Nom de l'archive de sauvegarde.
            files (list[str]): Liste de fichiers à inclure dans l'archive.

        Raises:
            FileNotFoundError: Si un fichier n'existe pas.
            PermissionError: Si un fichier n'est pas accessible.
        """
        try:
            # Crée le répertoire de sauvegarde si nécessaire
            os.makedirs(self.backup_dir, exist_ok=True)

            # Crée l'archive de sauvegarde
            with tarfile.open(os.path.join(self.backup_dir, backup_name), 'w:gz') as tar:
                for file in files:
                    # Récupère le chemin absolu du fichier
                    file_path = os.path.abspath(file)

                    # Vérifie si le fichier existe et est accessible
                    if not os.path.exists(file_path):
                        logger.error(f"Fichier '{file_path}' introuvable.")
                        raise FileNotFoundError(f"Fichier '{file_path}' introuvable.")

                    if not os.access(file_path, os.R_OK):
                        logger.error(f"Fichier '{file_path}' non accessible.")
                        raise PermissionError(f"Fichier '{file_path}' non accessible.")

                    # Ajoute le fichier à l'archive
                    tar.add(file_path, arcname=os.path.basename(file_path))

            # Enregistre la date de création de l'archive
            with open(os.path.join(self.backup_dir, 'backup_timestamp.txt'), 'w') as timestamp_file:
                timestamp_file.write(datetime.now().strftime('%Y-%m-%d %H:%M:%S'))

            # Enregistre l'archive sur le stockage de sauvegarde
            self.storage.upload_backup(backup_name)

            logger.info(f"Archive de sauvegarde '{backup_name}' créée avec succès.")

        except Exception as e:
            logger.error(f"Erreur lors de la création de l'archive de sauvegarde : {str(e)}")
            raise

def main():
    """
    Fichier principal de l'outil de sauvegarde.
    """
    storage = Storage()  # Crée un objet de stockage de sauvegarde
    backup_dir = '/path/au/reperoire/de/sauvegarde'  # Chemin du répertoire de sauvegarde
    backup = Backup(storage, backup_dir)

    # Liste des fichiers à inclure dans l'archive
    files_to_backup = ['/chemin/vers/fichier1.txt', '/chemin/vers/fichier2.txt']

    # Crée une archive de sauvegarde
    backup_name = 'backup_2023-03-15.tar.gz'
    backup.create_backup(backup_name, files_to_backup)

if __name__ == '__main__':
    main()