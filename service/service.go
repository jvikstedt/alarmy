package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jvikstedt/alarmy/model"
)

const BASE_URL = "http://localhost:8080"

var client = &http.Client{
	Timeout: time.Second * 10,
}

func ProjectAll() ([]model.Project, error) {
	var projects []model.Project

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/projects", BASE_URL), nil)
	resp, err := client.Do(req)
	if err != nil {
		return projects, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return projects, err
	}

	err = json.Unmarshal(data, &projects)

	return projects, err
}

func ProjectNew(project model.Project) (model.Project, error) {
	pJSON, err := json.Marshal(project)
	if err != nil {
		return project, err
	}
	b := bytes.NewBuffer(pJSON)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/projects", BASE_URL), b)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return project, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return project, err
	}

	err = json.Unmarshal(data, &project)

	return project, err
}
