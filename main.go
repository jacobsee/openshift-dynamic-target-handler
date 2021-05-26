package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	cron "github.com/robfig/cron"
)

type TemplateParams struct {
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	c := cron.New()
	c.AddFunc("*/5 * * * *", process)
	c.Start()
	log.Println("Target handler is running")
	for true {
		time.Sleep(time.Hour)
	}
}

func process() {
	token := os.Getenv("AUTH_TOKEN")
	targetPath := os.Getenv("TARGET_PATH")
	if len(token) == 0 {
		log.Println("Missing required environment variable")
		os.Exit(1)
	}
	if len(targetPath) == 0 {
		targetPath = "targets"
	}

	dynamicTargetServer := os.Getenv("SERVER")
	if len(dynamicTargetServer) == 0 {
		dynamicTargetServer = "http://localhost:8081"
	}
	client := http.Client{}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/list", dynamicTargetServer),
		nil,
	)
	q := req.URL.Query()
	q.Add("kind", "openshift")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", token)
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	var results map[string]string
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Println(err)
		return
	}

	// Remove targets that should no longer be here

	dir, err := ioutil.ReadDir(targetPath)
	if err != nil {
		log.Println(err)
		return
	}
	for _, file := range dir {
		filename := file.Name()
		if _, exists := results[filename[:len(filename)-4]]; !exists {
			log.Printf("Removing target: %s\n", filename)
			err = os.Remove(fmt.Sprintf("%s/%s", targetPath, filename))
			if err != nil {
				log.Printf("Could not remove target: %s\n", filename)
			}
		}
	}

	// Write targets that need to be written

	for key, value := range results {
		var params TemplateParams
		err := json.Unmarshal([]byte(value), &params)
		params.Endpoint = key
		tmpl := template.Must(template.ParseFiles("target.tmpl"))
		targetFilename := fmt.Sprintf("%s/%s.yml", targetPath, key)
		f, err := os.Create(targetFilename)
		if err != nil {
			log.Printf("Could not create new target file: %s\n", targetFilename)
		}
		defer f.Close()
		var content bytes.Buffer
		tmpl.Execute(&content, params)
		f.Write(content.Bytes())
	}
}
