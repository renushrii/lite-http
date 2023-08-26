package main

import (
	"encoding/json"
	"fmt"
	gohttp "net/http"
	"poga_gyan/http" // Update this import based on your package structure
)

func main() {
	server := http.NewServer(":8080")

	server.Get("/json", JSONHandler)

	err := server.Start()
	if err != nil {
		panic(err)
	}
}

type Data struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func JSONHandler(req http.Request, resp *http.Response) error {
	var data Data

	body := req.Body()
	defer body.Close()

	// WTF

	// Read and parse JSON from the request body
	if req.Method == "POST" || req.Method == "GET" {
		decoder := json.NewDecoder(body)
		err := decoder.Decode(&data)
		if err != nil {
			// Handle JSON parsing error
			resp.StatusCode(gohttp.StatusBadRequest)
			return err
		}
		fmt.Println("Decoded JSON:", data)
	}

	resp.StatusCode(gohttp.StatusOK)
	resp.AddHeader("Content-Type", "application/json")

	responseData, err := json.Marshal(data)
	if err != nil {
		resp.StatusCode(gohttp.StatusInternalServerError)
		return err
	}
	resp.Write(responseData)

	return nil
}
