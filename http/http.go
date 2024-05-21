package http_rest

import (
	"fmt"
	"net/http"
	"strings"
)

func SendGETRequest(uri string) (*http.Response, error){
	// Create http request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// set header
	req.Header.Set("Accept", "vdn.dac.v1")
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}

	return resp, nil
}

func SendPUTRequest(uri string, jsonBody string) (*http.Response, error) {
    // Create http request
    req, err := http.NewRequest("PUT", uri, strings.NewReader(jsonBody))
    if err != nil {
        fmt.Println("Error creating request:", err)
        return nil, err
    }

    req.Header.Set("Accept", "vdn.dac.v1")
    req.Header.Set("Content-Type", "application/json")

    client := http.DefaultClient
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return nil, err
    }

    return resp, nil
}