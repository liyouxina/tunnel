package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
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

var (
	taskQueue   = make(chan Task, 100)
	resultQueue = make(chan Result, 100)
	mu          sync.Mutex
	waiting     bool
)

func main() {
	http.HandleFunc("/request", requestHandler)
	http.HandleFunc("/result", resultHandler)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	task := Task{
		Method: r.Method,
		URL:    r.URL.String(),
		Header: r.Header,
		Body:   body,
	}

	taskQueue <- task

	result := <-resultQueue
	for key, values := range result.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(result.StatusCode)
	w.Write(result.Body)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	var result Result
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, "Failed to decode result", http.StatusInternalServerError)
		return
	}

	resultQueue <- result
	w.WriteHeader(http.StatusOK)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case task := <-taskQueue:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	case <-time.After(30 * time.Second):
		http.Error(w, "No tasks available", http.StatusNoContent)
	}
}

func main() {
	http.HandleFunc("/task", getTaskHandler)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
