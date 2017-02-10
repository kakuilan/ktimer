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

	var str1, str2 string
	str1 = getCurrentDirectory()

	str2 = getParentDirectory(str1)
	fmt.Println(str1, str2)

    fmt.Println(ktimer.DEFAULT_CONF)
    f := ktimer.GetConfFilePath()
    ck := ktimer.CheckConfFile()
    log.Println("adfadf", ck)
    cj,err := ktimer.CreateConfFile()
    ck = ktimer.CheckConfFile()
    fmt.Println(f, cj,err,ck)
}
