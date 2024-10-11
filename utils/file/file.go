package file

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func ReadLines(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func ParseQueryString(query string, target *map[string]interface{}) error {
	decoder := charmap.ISO8859_1.NewDecoder()
	decodedQuery, _, err := transform.String(decoder, query)
	if err != nil {
		return err
	}

	values, err := url.ParseQuery(decodedQuery)
	if err != nil {
		return err
	}

	for key, val := range values {
		if len(val) > 0 {
			(*target)[key] = val[0]
		}
	}
	return nil
}

func CheckExistAndCopy(dst, src string) error {
	if FolderExists(dst) {
		return nil
	}

	return CopyFolder(src, dst)
}

func FolderExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CopyFolder(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return CopyFile(path, dstPath)
	})
}

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return os.Chmod(dst, os.ModePerm)
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

// UnzipAllFilesInFolder finds all .zip files in a given folder and unzips them.
func UnzipAllFilesInFolder(folderPath, destDir string) error {
	// Walk through the directory to find all .zip files
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a .zip file
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zip") {
			fmt.Printf("Found zip file: %s\n", path)

			// Unzip the file
			if err := Unzip(path, destDir); err != nil {
				return fmt.Errorf("failed to unzip file: %w", err)
			}
			fmt.Printf("Successfully unzipped: %s\n", path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error while walking through folder: %w", err)
	}

	return nil
}

// GetFoldersInFolder lists all the folder names inside a given folder
func GetFoldersInFolder(folderPath string) ([]string, error) {
	var folders []string

	// Read directory contents
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Iterate through the directory contents
	for _, entry := range entries {
		// Check if the entry is a directory
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		}
	}

	return folders, nil
}

func DeleteFolder(folderPath string) error {
	// Use os.RemoveAll to delete the folder and all of its contents
	err := os.RemoveAll(folderPath)
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}
	return nil
}
