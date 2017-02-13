package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
    "ktimer"
)

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)

	}
	return string(runes[pos:l])

}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))

}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)

	}
	return strings.Replace(dir, "\\", "/", -1)

}

func main() {
    log.Println("ktimer main")
    confile := GetConfFilePath()
    fmt.Println(confile)

}
