package exporter

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// ImportUserProfilesFromZip imports user profiles from zip files into the Chrome profile directory.
func ImportUserProfilesFromZip(backupFolder, chromeProfileDir string) error {
	// Walk through the backup folder to find zip files
	return filepath.Walk(backupFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", path, err)
		}

		// Check if it's a zip file
		if filepath.Ext(path) == ".zip" {
			// Create a temporary directory to unzip the profile
			tempDir := filepath.Join(backupFolder, "temp")
			log.Println(tempDir)
			err := os.MkdirAll(tempDir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create temp directory: %w", err)
			}

			// Unzip the file into the temporary directory
			log.Printf("Unzipping profile from %s to %s\n", path, tempDir)
			err = Unzip(path, tempDir)
			if err != nil {
				return fmt.Errorf("failed to unzip file %s: %w", path, err)
			}

			// Import the unzipped profiles into the Chrome profile directory
			err = ImportProfilesFromDir(tempDir, chromeProfileDir)
			if err != nil {
				return fmt.Errorf("failed to import profiles from %s: %w", tempDir, err)
			}

			// Clean up the temporary directory
			os.RemoveAll(tempDir)
		}

		return nil
	})
}

// ZipFolder zips the entire folder and saves it as a .zip file.
func ZipFolder(sourceDir, destinationZip string) error {
	// Create the destination zip file
	zipFile, err := os.Create(destinationZip)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	// Initialize the zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the source directory and zip everything
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path to maintain directory structure in the zip file
		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// If it's a directory, add it to the zip
		if info.IsDir() {
			_, err = zipWriter.Create(relativePath + "/")
			if err != nil {
				return fmt.Errorf("failed to add folder to zip: %w", err)
			}
			return nil
		}

		// Otherwise, add the file to the zip
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file for zipping: %w", err)
		}
		defer file.Close()

		// Create a zip entry for the file
		zipEntry, err := zipWriter.Create(relativePath)
		if err != nil {
			return fmt.Errorf("failed to add file to zip: %w", err)
		}

		// Copy the file content into the zip entry
		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return fmt.Errorf("failed to copy file to zip: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to zip folder: %w", err)
	}

	return nil
}

// SaveAllUserProfilesToFolder copies all user profile directories into a single backup folder
func SaveAllUserProfilesToFolder(userProfiles []string, backupFolder string) error {
	// Ensure the backup folder exists
	err := os.MkdirAll(backupFolder, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create backup folder: %w", err)
	}

	// Copy each user profile directory to the backup folder
	for _, profileDir := range userProfiles {
		profileName := filepath.Base(profileDir) // Get the profile name (e.g., user1, user2)
		log.Println("name", profileName)
		destinationDir := filepath.Join(backupFolder, profileName)

		// Copy the profile directory to the backup folder
		log.Printf("Copying profile %s to %s\n", profileName, destinationDir)
		err := CopyDir(profileDir, destinationDir)
		if err != nil {
			return fmt.Errorf("failed to copy profile %s: %w", profileName, err)
		}
	}

	return nil
}

// CopyDir copies the contents of the source directory to the destination directory
func CopyDir(sourceDir, destDir string) error {
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path to recreate the directory structure
		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Create the destination path
		destPath := filepath.Join(destDir, relativePath)

		// If it's a directory, create the corresponding directory in the destination
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Otherwise, copy the file
		return copyFile(path, destPath)
	})
}

// copyFile copies a file from source to destination
func copyFile(sourceFile, destFile string) error {
	src, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Set file permissions
	srcInfo, err := os.Stat(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	err = os.Chmod(destFile, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}

// Unzip extracts a zip file to a specified destination directory.
func Unzip(zipFilePath, destDir string) error {
	// Open the zip file
	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipFile.Close()

	// Iterate through the files in the zip archive
	for _, file := range zipFile.File {
		// Create the full path for the destination file
		destPath := filepath.Join(destDir, file.Name)

		// Check if the file is a directory
		if file.FileInfo().IsDir() {
			// Create the directory if it doesn't exist
			err := os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			continue
		}

		// Create the destination file
		destFile, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", destPath, err)
		}
		defer destFile.Close()

		// Copy the file content from the zip file to the destination file
		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s in zip: %w", file.Name, err)
		}
		defer srcFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return fmt.Errorf("failed to copy file %s: %w", file.Name, err)
		}
	}

	return nil
}

// ImportProfilesFromDir imports profiles from an unzipped directory.
func ImportProfilesFromDir(sourceDir, chromeProfileDir string) error {
	// Walk through the unzipped directory to find user profile directories
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", path, err)
		}

		// Check if it is a directory
		if info.IsDir() {
			// Get the profile name from the directory
			profileName := info.Name()
			destinationDir := filepath.Join(chromeProfileDir, profileName)

			// Copy the profile directory to the Chrome profile directory
			log.Printf("Importing profile from %s to %s\n", path, destinationDir)
			err := CopyDir(path, destinationDir)
			if err != nil {
				return fmt.Errorf("failed to import profile %s: %w", profileName, err)
			}
		}

		return nil
	})
}
