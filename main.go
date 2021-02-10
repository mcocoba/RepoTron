package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
)

const folder string = "/"

var logger *logrus.Logger

const (
	accessToken = ""
)

// Repo structure
type Repo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func folderExists(folder string, logger *logrus.Logger) bool {
	_, err := os.Stat(folder)
	logger.Info("folderExists: " + folder)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func createFolder(folder string, logger *logrus.Logger) {
	//Create a folder/directory at a full qualified path
	err := os.Mkdir(folder, 0755)
	logger.Info("createFolder: " + folder)
	if err != nil {
		logger.Fatal(err)
	}
}

func main() {
	start := time.Now()

	// create the logger
	logger := logrus.New()

	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetOutput(os.Stdout)

	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		logger.Fatal(err)
	}

	defer file.Close()
	logger.SetOutput(file)

	logger.Info("RepoTron")

	// read file
	data, err := ioutil.ReadFile("./source.json")

	if err != nil {
		logger.Error(err)
	}

	// json data
	var repos []Repo
	_ = json.Unmarshal([]byte(data), &repos)

	if !folderExists("repos", logger) {
		// Create Repos Folder
		createFolder("repos", logger)
	}

	if !folderExists("backup", logger) {
		// Create Backup Folder
		createFolder("backup", logger)
	}

	// Loop the repos
	for i := 0; i < len(repos); i++ {

		if !folderExists("repos/"+repos[i].Name, logger) {
			createFolder("repos/"+repos[i].Name, logger)
		}

		//Clone Repo
		cmd := exec.Command("git", "clone", repos[i].Path, ".")
		cmd.Dir = "repos/" + repos[i].Name

		err := cmd.Run()
		logger.Info("CMD git clone: " + repos[i].Path)
		if err != nil {
			logger.Fatal(err)
		}

		//compress folder
		comp := exec.Command("/bin/sh", "-c", "zip -r backup/"+repos[i].Name+".zip repos/"+repos[i].Name)
		var out bytes.Buffer
		comp.Stdout = &out
		errs := comp.Run()
		logger.Info("CMD zip: " + repos[i].Path)
		if errs != nil {
			logger.Fatal(errs)
		}

		errFolder := os.RemoveAll("repos/" + repos[i].Name)

		if errFolder != nil {
			logger.Fatal(errFolder)
		}

		logger.Info("Finished: " + repos[i].Name)
	}

	elapsed := time.Since(start)
	logger.Info("Backup took %s", elapsed)
}
