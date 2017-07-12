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

func PostAsJSON(path string, v interface{}) error {
	asJSON, err := json.Marshal(v)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(asJSON)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", BASE_URL, path), b)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &v)

	return err
}
