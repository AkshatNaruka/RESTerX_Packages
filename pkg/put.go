package pkg

import (
	"fmt"
)

// HandlePutRequest sends a PUT request to the specified URL and prints the response
func HandlePutRequest(url string) {
	req := APIRequest{
		Method:  "PUT",
		URL:     url,
		Headers: map[string]string{},
		Body:    "",
	}
	response := MakeHTTPRequest("PUT", url, "", req.Headers)
	FormatAndPrintResponse(req, response)
}

// HandlePutRequestAdvanced sends a PUT request with custom headers and body
func HandlePutRequestAdvanced(url string, headers map[string]string, body string) APIResponse {
	return MakeHTTPRequest("PUT", url, body, headers)
}

// MakePutRequest sends a PUT request and returns structured response data
func MakePutRequest(url, body string, headers map[string]string) APIResponse {
	return MakeHTTPRequest("PUT", url, body, headers)
}

func printResponse(response APIResponse) {
	if response.Error != "" {
		// keep behavior via FormatAndPrintResponse in new code path
		// fallback
		fmt.Printf("Error: %s\n", response.Error)
		return
	}

	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Response Time: %v\n", response.ResponseTime)
	fmt.Println("Response Body:")
	fmt.Println(response.Body)
}