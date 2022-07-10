package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	projectPath     = flag.String("path", cwd(), "the project path")
	destinationPath = flag.String("dest", "/tmp", "the destination path")
	destinationUrl  = flag.String("destUrl", "", "the destination url")
)

func main() {

	flag.Parse()

	pattern := loadPattern(*projectPath, []string{".gitignore"})

	archivePath := createArchive(*projectPath, pattern, *destinationPath)

	println(archivePath)

}

func loadFiles() []string {

	files := []string{}

	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				files = append(files, path)
				println(path)
			}

			return nil
		})

	if err != nil {
		log.Fatal(err)
	}

	return files

}

func loadPattern(projectPath string, ignoredFiles []string) []string {

	data := []byte{}
	pattern := []string{}
	ignoredFileStatus := map[string]bool{}

	for _, f := range ignoredFiles {

		ignoreFilePath := filepath.Join(projectPath, f)

		if _, err := os.Stat(projectPath); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("%s Doest not exist!", f)
		} else if _, err := os.Stat(ignoreFilePath); err == nil {
			ignoredFileStatus[f] = true
		}

	}

	for f, ok := range ignoredFileStatus {

		ignoreFilePath := filepath.Join(projectPath, f)

		if ok {
			data, _ = os.ReadFile(ignoreFilePath)
		}
	}

	pattern = strings.Split(string(data), "\n")

	return pattern
}

func createArchive(projectPath string, pattern []string, destinationPath string) string {
	archiveName := strconv.Itoa(randomNum()) + ".tar.gz"
	archivePath := filepath.Join(destinationPath, archiveName)

	files := loadFiles()

	tarFile, err := os.Create(archivePath)

	if err != nil {
		log.Fatal(err)
	}

	defer tarFile.Close()

	gw := gzip.NewWriter(tarFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {

		file, err := os.Open(file)

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		info, err := file.Stat()

		if err != nil {
			log.Fatal(err)
		}

		header, err := tar.FileInfoHeader(info, info.Name())

		if err != nil {
			log.Fatal(err)
		}

		header.Name = file.Name()

		if err = tw.WriteHeader(header); err != nil {
			log.Fatal(err)
		}

		if _, err = io.Copy(tw, file); err != nil {
			log.Fatal(err)
		}
	}

	return archivePath
}

func uploadArchive(archivePath string, destinationUrl string) {
	return
}
