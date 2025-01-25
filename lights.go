package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
)

var lightsAddress string
var lightsGroup string
var lightsKey string
var lightsMutex sync.Mutex

func sendHueHubRequest(method string, url string, body []byte) error {
	// create a custom HTTP client that ignores SSL certificate errors
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// create a new POST request with the given body
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("hue-application-key", lightsKey)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func setLightsToFullOn() error {
	lightsMutex.Lock()
	defer lightsMutex.Unlock()
	err := sendHueHubRequest(
		"PUT",
		fmt.Sprintf("https://%v/clip/v2/resource/grouped_light/%v", lightsAddress, lightsGroup),
		[]byte(`{ "on": { "on": true }, "dimming": { "brightness": 100.0 } }`),
	)
	return err
}

func setLightsToFullOff() error {
	lightsMutex.Lock()
	defer lightsMutex.Unlock()
	err := sendHueHubRequest(
		"PUT",
		fmt.Sprintf("https://%v/clip/v2/resource/grouped_light/%v", lightsAddress, lightsGroup),
		[]byte(`{ "on": { "on": false } }`),
	)
	return err
}

func setLightsToDiningMode() error {
	lightsMutex.Lock()
	defer lightsMutex.Unlock()
	err := sendHueHubRequest(
		"PUT",
		fmt.Sprintf("https://%v/clip/v2/resource/grouped_light/%v", lightsAddress, lightsGroup),
		[]byte(`{ "on": { "on": true }, "dimming": { "brightness": 10.0 } }`),
	)
	return err
}
