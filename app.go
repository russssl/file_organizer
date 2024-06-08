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

func loadConfig(configPath string) (Config, error) {
	var config Config
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

func organizeFiles(directory string, config Config) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(info.Name())
		if destDir, ok := config.Rules[ext]; ok {
			destPath := filepath.Join(directory, destDir, info.Name())
			os.MkdirAll(filepath.Join(directory, destDir), os.ModePerm)
			err := os.Rename(path, destPath)
			if err != nil {
				return err
			}
			fmt.Printf("Moved %s to %s\n", path, destPath)
		}

		return nil
	})
}

func main() {
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	directory := flag.String("dir", ".", "Directory to organize")
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	err = organizeFiles(*directory, config)
	if err != nil {
		log.Fatalf("Error organizing files: %v\n", err)
	}
}
