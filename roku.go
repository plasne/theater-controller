package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

var rokuAddress string
var rokuMutex sync.Mutex

func pressKeyOnRoku(key string) error {
	rokuMutex.Lock()
	defer rokuMutex.Unlock()
	u := fmt.Sprintf("http://%s/keypress/%s", rokuAddress, key)
	resp, err := http.Post(u, "text/plain", nil)
	if err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	return nil
}

func launchOnRoku(appId string) error {
	rokuMutex.Lock()
	defer rokuMutex.Unlock()
	u := fmt.Sprintf("http://%s/launch/%s", rokuAddress, appId)
	resp, err := http.Post(u, "text/plain", nil)
	if err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	return nil
}

func sendTextToRoku(text string) error {
	rokuMutex.Lock()
	defer rokuMutex.Unlock()
	for i := 0; i < len(text); i++ {
		enc := url.QueryEscape(string(text[i]))
		u := fmt.Sprintf("http://%s/keypress/Lit_%s", rokuAddress, enc)
		resp, err := http.Post(u, "text/plain", nil)
		if err != nil {
			return err
		}
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	return nil
}

func clearTextOnRoku() error {
	rokuMutex.Lock()
	defer rokuMutex.Unlock()
	for i := 0; i < 20; i++ {
		u := fmt.Sprintf("http://%s/keypress/Backspace", rokuAddress)
		resp, err := http.Post(u, "text/plain", nil)
		if err != nil {
			return err
		}
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	return nil
}
