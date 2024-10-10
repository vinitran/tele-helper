package file

import (
	"bufio"
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
