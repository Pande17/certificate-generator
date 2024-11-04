package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// function to get public IP addres from subject
func GetPublicIP() (string, error) {
	// make a request to ipify API
	response, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// read the response body
	ip, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// return the IP address as a string
	return string(ip), nil
}

// function for getting location from IP (implement as per your needs)
func GetLocationFromIP(ip string) string {
	apiToken := os.Getenv("IPINFO")
	if apiToken == "" {
		return "API Token not found"
	}

	url := fmt.Sprintf("https://ipinfo.io/%s/geo?token=%s", ip, apiToken)
	response, err := http.Get(url)
	if err != nil {
		return "Unknown Location"
	}
	defer response.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "Unknown Location"
	}

	city, _ := result["city"].(string)
	country, _ := result["country"].(string)

	if city != "" && country != "" {
		return fmt.Sprintf("%s, %s", city, country)
	}
	return "Unknown Location"
}
