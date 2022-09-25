package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	var path string
	flag.StringVar(&path, "path", dir, "root path")
	flag.Parse()

	lastUpdatesFiles := "files-lists.txt"

	currentTime := time.Now()
	timeString := currentTime.String()
	log.Print("was get currentTime: ", timeString)
	var dateLastCreate time.Time

	if _, err := os.Stat(lastUpdatesFiles); os.IsNotExist(err) {
		// path/to/whatever does not exist
	} else {
		// search date last update
		dateLastCreate, err = GetDateCreateFile(lastUpdatesFiles)

		if err != nil {
			panic("Remove error")
		}

		// delete file
		err = os.Remove(lastUpdatesFiles)
		if err != nil {
			panic("Remove error")
		}
		fmt.Println("Remove.")

		// create new file
	}
	// create file and written into data
	createdFile, err := os.Create(lastUpdatesFiles)

	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer createdFile.Close()

	fmt.Println("---get files from filepath---")
	files, err := FilePathWalkDir(path, dateLastCreate)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	for _, f := range files {
		fmt.Println(f)
		createdFile.WriteString(f)
		createdFile.WriteString("\n")
	}

	fmt.Println("Done.")
}

func FilePathWalkDir(root string, tCompare time.Time) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			tCreate, err := GetDateCreateFile(path)
			if err != nil {
				log.Fatal(err)
				panic("Remove error")
			}
			if tCreate.Before(tCompare) {
				files = append(files, path)
			}
		}
		return nil
	})

	return files, err
}

func GetDateCreateFile(root string) (time.Time, error) {
	f, err := os.Open(root)
	if err != nil {
		panic("GetDateCreateFile error open")
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		panic("GetDateCreateFile error stat")
	}
	dateLastCreate := stat.ModTime()
	return dateLastCreate, err
}
