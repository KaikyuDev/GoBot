package gobot

import (
	"net/http"
	"encoding/json"
	"fmt"
	"time"
	"bytes"
	"mime/multipart"
)


var client = &http.Client{}

var baseUrl = "https://api.telegram.org/"

func statusCheck(result interface{}, resp *http.Response, status int) bool {
	var old = &result
	var apiError = ApiError{}

	if status != ResponseError {
		if status == RequestNotOk {
			json.NewDecoder(resp.Body).Decode(&apiError)
			fmt.Println("Excecution failed")
			fmt.Printf("Telegram has returned error '%s' with status code '%d'.", apiError.Description, apiError.ErrorCode)
			return false
		} else {
			json.NewDecoder(resp.Body).Decode(result)
		}

		// Defer the closing of the body
		defer resp.Body.Close()

		if result == old {
			fmt.Println("Excecution failed, network error?")
			return false
		}

		return true
	} else {
		return false
	}


}

func makeTimeoutRequest(url string, timeout bool) (*http.Response, int) {
	// safePhone := url.QueryEscape(phone)
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error in the URL '%s'", url)
		fmt.Println(err)
		return nil, ResponseError
	}

	if timeout {
		client.Timeout = time.Duration(time.Second * 122)
	} else {
		client.Timeout = time.Duration(time.Second * 5)
	}

	resp, err := client.Do(req)
	if err != nil {
		println("Error while requesting...")
		fmt.Println(err)
		return nil, ResponseError
	}

	if resp.StatusCode != 200 {
		fmt.Println("Status code not 200.")
		// body, _ := ioutil.ReadAll(resp.Body)
		// result := string(body)
		return resp, RequestNotOk
	}

	return resp, RequestOk
}

func makeRequest(url string) (*http.Response, int) {
	return makeTimeoutRequest(url, false)
}

func makePost(url string, fileType string, content []byte) (*http.Response, int) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(fileType, "file.jpg")
	part.Write(content)
	writer.Close()

	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Printf("Error in the URL '%s'", url)
		fmt.Println(err)
		return nil, WrongRequest
	}

	r.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(r)
	if err != nil {
		println("Error while requesting...")
		fmt.Println(err)
		return nil, ResponseError
	}

	if resp.StatusCode != 200 {
		fmt.Println("Status code not 200.")
		return resp, RequestNotOk
	}

	return resp, RequestOk
}