package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Task struct {
	Method string
	URL    string
	Header http.Header
	Body   []byte
}

type Result struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func main() {
	serverURL := "http://<server-ip>:8080"

	for {
		task, err := getTask(serverURL + "/task")
		if err != nil {
			log.Println("No tasks available, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		result, err := processTask(task)
		if err != nil {
			log.Println("Failed to process task:", err)
			continue
		}

		err = sendResult(serverURL+"/result", result)
		if err != nil {
			log.Println("Failed to send result:", err)
		}
	}
}

func getTask(url string) (Task, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Task{}, nil
	}

	var task Task
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func processTask(task Task) (Result, error) {
	req, err := http.NewRequest(task.Method, task.URL, bytes.NewReader(task.Body))
	if err != nil {
		return Result{}, err
	}

	req.Header = task.Header
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	result := Result{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       body,
	}

	return result, nil
}

func sendResult(url string, result Result) error {
	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	_, err = http.Post(url, "application/json", bytes.NewReader(body))
	return err
}
