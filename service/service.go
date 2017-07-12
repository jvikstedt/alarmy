package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const BASE_URL = "http://localhost:8080"

var client = &http.Client{
	Timeout: time.Second * 10,
}

func GetResource(path string, resource interface{}) error {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", BASE_URL, path), nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, resource)

	return err
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
