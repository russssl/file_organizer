package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Rules map[string]string `json:"rules"`
}

var defaultConfig = Config{
	Rules: map[string]string{
		".png":  "images",
		".jpg":  "images",
		".jpeg": "images",
		".gif":  "images",
		".txt":  "documents",
		".pdf":  "documents",
		".doc":  "documents",
		".docx": "documents",
		".xlsx": "spreadsheets",
		".xls":  "spreadsheets",
		".csv":  "spreadsheets",
		".mp3":  "audio",
		".wav":  "audio",
		".mp4":  "videos",
		".mkv":  "videos",
		".avi":  "videos",
	},
}

// flags
var (
	configPath   = flag.String("config", "", "Path to the configuration file")
	directory    = flag.String("dir", ".", "Directory to organize")
	help         = flag.Bool("help", false, "Show help message")
	printDefault = flag.Bool("print-default", false, "Print default config")
	recursive    = flag.Bool("r", false, "Organize files recursively, including subdirectories")
	v            = flag.Bool("v", false, "Show version")
	dryRun       = flag.Bool("dry-run", false, "Dry run (do not move files)")
)

func loadConfig(configPath string) (Config, error) {
	var config Config
	if configPath == "" {
		return defaultConfig, nil
	}
	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func organizeFiles(directory string, config Config, recursive bool, dryRun bool) error {
	filesMoved := 0
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || (!recursive && filepath.Dir(path) != directory) {
			return nil
		}
		ext := filepath.Ext(info.Name())
		if destDir, ok := config.Rules[ext]; ok {
			if !dryRun {
				destPath := filepath.Join(directory, destDir, info.Name())
				os.MkdirAll(filepath.Join(directory, destDir), os.ModePerm)
				err := os.Rename(path, destPath)
				if err != nil {
					return err
				}
				fmt.Printf("Moved %s to %s\n", path, destPath)
				filesMoved++
			} else {
				fmt.Printf("Would move %s to %s\n", path, filepath.Join(directory, destDir, info.Name()))
			}
		}
		return nil
	})
	if filesMoved == 0 {
		if dryRun {
			fmt.Println("No files to move")
		} else {
			fmt.Println("No files moved")
		}
	}
	return err
}

func main() {

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}
	if *v {
		fmt.Println("v0.0.1")
		return
	}

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	if *printDefault {
		defaultConfigJson, _ := json.MarshalIndent(defaultConfig, "", "  ")
		fmt.Println(string(defaultConfigJson))
		return
	}
	err = organizeFiles(*directory, config, *recursive, *dryRun)

	if err != nil {
		log.Fatalf("Error organizing files: %v\n", err)
	}
}
