package pkg

// HandleHeadRequest sends a HEAD request to the specified URL and prints the response
func HandleHeadRequest(url string) {
	req := APIRequest{
		Method:  "HEAD",
		URL:     url,
		Headers: map[string]string{},
		Body:    "",
	}
	response := MakeHTTPRequest("HEAD", url, "", req.Headers)
	FormatAndPrintResponse(req, response)
}

// HandleHeadRequestAdvanced sends a HEAD request with custom headers
func HandleHeadRequestAdvanced(url string, headers map[string]string) APIResponse {
	return MakeHTTPRequest("HEAD", url, "", headers)
}

// MakeHeadRequest sends a HEAD request and returns structured response data
func MakeHeadRequest(url string, headers map[string]string) APIResponse {
	return MakeHTTPRequest("HEAD", url, "", headers)
}