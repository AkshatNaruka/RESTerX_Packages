package pkg

// HandleDeleteRequest sends a DELETE request to the specified URL and prints the response body
func HandleDeleteRequest(url string) {
	req := APIRequest{
		Method:  "DELETE",
		URL:     url,
		Headers: map[string]string{},
		Body:    "",
	}
	response := MakeHTTPRequest("DELETE", url, "", req.Headers)
	FormatAndPrintResponse(req, response)
}

// HandleDeleteRequestAdvanced sends a DELETE request with custom headers
func HandleDeleteRequestAdvanced(url string, headers map[string]string) APIResponse {
	return MakeHTTPRequest("DELETE", url, "", headers)
}

// MakeDeleteRequest sends a DELETE request and returns structured response data
func MakeDeleteRequest(url string, headers map[string]string) APIResponse {
	return MakeHTTPRequest("DELETE", url, "", headers)
}
