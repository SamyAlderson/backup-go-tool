"""
Fichier de test pour le module storage
"""
import unittest
from unittest.mock import Mock
from src.storage import Storage, BackupTypeError

class TestStorage(unittest.TestCase):

    def setUp(self):
        """
        Initialisation des objets à tester
        """
        self.storage = Storage()

    def test_storage_init(self):
        """
        Test de l'initialisation de l'objet Storage
        """
        self.assertIsNotNone(self.storage)

    def test_storage_add_backup(self):
        """
        Test de la méthode add_backup pour ajouter un backup
        """
        self.storage.add_backup("test_backup")
        self.assertEqual(len(self.storage.backups), 1)

    def test_storage_add_backup_error(self):
        """
        Test de la méthode add_backup pour erreur de type
        """
        with self.assertRaises(BackupTypeError):
            self.storage.add_backup(123)

    def test_storage_get_backups(self):
        """
        Test de la méthode get_backups pour récupérer les backups
        """
        self.storage.add_backup("test_backup1")
        self.storage.add_backup("test_backup2")
        self.assertEqual(len(self.storage.get_backups()), 2)

    def test_storage_remove_backup(self):
        """
        Test de la méthode remove_backup pour supprimer un backup
        """
        self.storage.add_backup("test_backup")
        self.storage.remove_backup("test_backup")
        self.assertEqual(len(self.storage.get_backups()), 0)

class TestStorageError(unittest.TestCase):

    def test_backup_type_error(self):
        """
        Test de la création d'un BackupTypeError
        """
        with self.assertRaises(BackupTypeError):
            BackupTypeError("Erreur de type")

if __name__ == '__main__':
    unittest.main()