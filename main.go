package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"flag"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	projectPath     = flag.String("path", cwd(), "the project path")
	destinationPath = flag.String("dest", "/tmp", "the destination path")
	destinationUrl  = flag.String("destUrl", "", "the destination url")
)

func main() {

	flag.Parse()

	// noticeRequired([]string{"path"})

	pattern := loadPattern(*projectPath, []string{".gitignore"})

	archivePath := createArchive(*projectPath, pattern, *destinationPath)

	println(archivePath)
	// uploadArchive(archivePath, *destinationUrl)
}

func verbose(str string) {

}

func randomNum() int {
	rSource := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(rSource)
	return rand.Intn(999999999)
}

func noticeRequired(required []string) {

	seen := map[string]bool{}

	flag.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
	})

	for _, req := range required {
		if !seen[req] {
			log.Fatalf("missing required --%s flag \n", req)
		}
	}
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

	files := []string{filepath.Join(projectPath, "server.js"), filepath.Join(projectPath, "package.json")}

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

		err = tw.WriteHeader(header)

		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(tw, file)
		if err != nil {
			log.Fatal(err)
		}
	}

	return archivePath
}

func uploadArchive(archivePath string, destinationUrl string) {
	return
}

func cwd() string {
	// unexpectedPath := []string{"/home"}

	path, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	return path
}
