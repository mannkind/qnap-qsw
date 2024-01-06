package qnap

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const apiURL string = "https://%s/api/v3"
const loginURL string = apiURL + "/users/login"
const verificationURL string = apiURL + "/users/verification"
const interfaceURL string = apiURL + "/poe/interface"

type QNAP struct {
	host  string
	token string
}

func New(host string) *QNAP {
	return &QNAP{
		host: host,
	}
}

func NewWithToken(host string, token string) *QNAP {
	q := New(host)
	q.token = token

	return q
}

func (q *QNAP) Login(password string) (string, error) {
	// If we already have a token, verify it and return it
	// Otherwise login again to obtain a new token
	//
	// QNAP is annoying here and obtaining a new token invalidates the existing token
	if q.token != "" {
		if err := q.Verify(); err == nil {
			return q.token, nil
		}
	}

	// Prepare the request
	url := fmt.Sprintf(loginURL, q.host)
	data := fmt.Sprintf(`{ "username": "admin", "password": "%s" }`, base64.StdEncoding.EncodeToString([]byte(password)))

	// Send the request
	resp, err := q.sendHTTPRequest(http.MethodPost, url, data)
	if err != nil {
		return "", fmt.Errorf("login; %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login; status code %d", resp.StatusCode)
	}

	// Read the response and unmarshal the JSON
	login := loginResponse{}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("login; read body of response; %w", err)
	}

	if err := json.Unmarshal([]byte(body), &login); err != nil {
		return "", fmt.Errorf("login; unmarshal response; %w", err)
	}

	// Set the token and return it to allow processes to use it later
	q.token = login.Result.AccessToken
	return login.Result.AccessToken, nil
}

func (q *QNAP) POEInterfaces() (map[string]interfacesValueResponse, error) {
	interfaces := map[string]interfacesValueResponse{}

	// Setup and send the request
	url := fmt.Sprintf(interfaceURL, q.host)
	resp, err := q.sendHTTPRequest(http.MethodGet, url, "")
	if err != nil {
		return interfaces, fmt.Errorf("poe interfaces; %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return interfaces, fmt.Errorf("poe interfaces; status code %d", resp.StatusCode)
	}

	// Read the response and unmarshal the JSON
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return interfaces, fmt.Errorf("poe interfaces; read body of response; %w", err)
	}

	interfacesResult := interfacesResponse{}
	if err := json.Unmarshal([]byte(body), &interfacesResult); err != nil {
		return interfaces, fmt.Errorf("poe interfaces; unmarshal response; %w", err)
	}

	for _, result := range interfacesResult.Result {
		interfaces[result.Key] = result.Val
	}

	return interfaces, nil
}

func (q *QNAP) UpdatePOEInterfaces(wg *sync.WaitGroup, port string, properties interfacesValueResponse) error {
	defer wg.Done()

	// Setup and send the request
	url := fmt.Sprintf(interfaceURL, q.host)
	marshaledProperites, err := json.Marshal(properties)
	if err != nil {
		return fmt.Errorf("update poe interface; marshal data; %w", err)
	}
	data := fmt.Sprintf(`{ "idx": "%s", "data": %s }`, port, marshaledProperites)

	resp, err := q.sendHTTPRequest(http.MethodPut, url, data)
	if err != nil {
		return fmt.Errorf("update poe interface; %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update poe interface; received status code %d", resp.StatusCode)
	}

	// Read the response and unmarshal the JSON
	_, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("update poe interface; read response body; %w", err)
	}

	return nil
}

func (q *QNAP) Verify() error {
	// Setup and send the request
	url := fmt.Sprintf(verificationURL, q.host)
	resp, err := q.sendHTTPRequest(http.MethodGet, url, "")
	if err != nil {
		return fmt.Errorf("verify; %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("verify; status code %d", resp.StatusCode)
	}

	return nil
}

func (q *QNAP) sendHTTPRequest(method string, url string, data string) (*http.Response, error) {
	// Prepare the HTTP client
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Prepare the request itself
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, fmt.Errorf("prepare request; %w", err)
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")
	if q.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", q.token))
	}

	// Send the request and return the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request; %w", err)
	}

	return resp, nil
}
