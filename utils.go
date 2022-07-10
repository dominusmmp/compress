package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"
)

func cwd() string {
	// unexpectedPath := []string{"/home"}

	path, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

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
