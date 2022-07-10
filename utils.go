package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"
)

func verbose(str string) {

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func cwd() string {
	homeDir, _ := os.UserHomeDir()
	unexpectedPath := map[string]bool{homeDir: true}

	path, err := os.Getwd()

	if unexpectedPath[path] {
		log.Fatalf("%s Not Allowed!", path)
	}

	checkErr(err)

	return path
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
