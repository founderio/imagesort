package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	basePath, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		handleFile(basePath, file)
	}
}

func handleFile(basePath string, file os.FileInfo) {
	// Skip directories
	if file.IsDir() {
		return
	}

	// Open File
	f, err := os.Open(file.Name())
	if err != nil {
		log.Println("Skip", file.Name(), "due to error opening file:", err.Error())
		return
	}
	defer f.Close()

	// Read creation date from exif data
	header, err := exif.Decode(f)
	if err != nil {
		log.Println("Skip", file.Name(), "due to EXIF extraction error:", err.Error())
		return
	}

	tm, err := header.DateTime()
	if err != nil {
		log.Println("Skip", file.Name(), "due to date extraction error:", err.Error())
		return
	}

	// Move file into year-month subfolder
	sourceFile := filepath.Join(basePath, file.Name())
	destFolderName := tm.Format("2006-01")
	destFolder := filepath.Join(basePath, destFolderName)
	destFile := filepath.Join(basePath, destFolderName, file.Name())

	log.Println("Moving", file.Name(), "to", destFolder)
	err = os.MkdirAll(destFolder, 0700)
	if err != nil {
		log.Println("Error creating", destFolder, "due to error:", err.Error())
		return
	}
	err = os.Rename(sourceFile, destFile)
	if err != nil {
		log.Println("Error renaming", sourceFile, "due to error:", err.Error())
		return
	}
}
