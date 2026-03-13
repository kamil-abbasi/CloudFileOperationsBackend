package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	url := "http://localhost:8080/healthcheck"

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Get(url)

	if err != nil {
		fmt.Println("Healthcheck failed:", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode <= 300 {
		fmt.Println("Healthy")
		os.Exit(0)
	} else {
		fmt.Println("Unhealthy, status code:", res.StatusCode)
		os.Exit(1)
	}
}
