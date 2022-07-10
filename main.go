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

	ignore "github.com/sabhiram/go-gitignore"
)

var (
	projectPath     = flag.String("path", cwd(), "project path in your computer")
	destinationPath = flag.String("dest", os.TempDir(), "compressed file destination")
)

func main() {

	flag.Parse()

	pattern := loadPattern(*projectPath, []string{".gitignore"})

	archivePath := createArchive(*projectPath, pattern, *destinationPath)

	println(archivePath)

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

func loadFiles(projectPath string, pattern []string) []string {

	files := []string{}
	ignoredOne := ignore.CompileIgnoreLines(pattern...)

	err := filepath.Walk(projectPath,
		func(path string, info os.FileInfo, err error) error {
			checkErr(err)

			if !info.IsDir() {
				filePath, _ := filepath.Rel(projectPath, path)
				isIgnore := ignoredOne.MatchesPath(filePath)
				if !isIgnore {
					files = append(files, filePath)
				}
			}

			return nil
		})

	checkErr(err)

	return files
}

func createArchive(projectPath string, pattern []string, destinationPath string) string {
	archiveName := strconv.Itoa(randomNum()) + ".tar.gz"
	archivePath := filepath.Join(destinationPath, archiveName)

	files := loadFiles(projectPath, pattern)

	tarFile, err := os.Create(archivePath)

	checkErr(err)

	defer tarFile.Close()

	gw := gzip.NewWriter(tarFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {

		os.Chdir(projectPath)

		file, err := os.Open(file)

		checkErr(err)

		defer file.Close()

		info, err := file.Stat()

		checkErr(err)

		header, err := tar.FileInfoHeader(info, info.Name())

		checkErr(err)

		header.Name = file.Name()
		err = tw.WriteHeader(header)
		checkErr(err)
		_, err = io.Copy(tw, file)
		checkErr(err)
	}

	return archivePath
}
