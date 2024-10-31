package middleware

// import (
// 	"io"
// 	"net/http"
// 	"os"
// )

// // function to get public IP addres from subject
// func GetPublicIP() (string, error) {
// 	// make a request to ipify API
// 	response, err := http.Get("https://api.ipify.org?format=text")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer response.Body.Close()

// 	// read the response body
// 	ip, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	// return the IP address as a string
// 	return string(ip), nil
// }

// // function for getting location from IP (implement as per your needs)
// func GetLocationFromIP(ip string) string {
// 	apiToken := os.Getenv()
// }
