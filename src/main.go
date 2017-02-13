package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
    "ktimer"
)


func main() {
    log.Println("ktimer main")
    confile := GetConfFilePath()
    fmt.Println(confile)

}
