package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var serverURL string = "https://portkey-server.herokuapp.com/"

func signalExchange(localSignal, remoteSignal *Signal) {
	resp, err := http.PostForm(serverURL, url.Values{
		"connParams": {Encode(localSignal)},
	})
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	key := string(body)
	fmt.Printf("Your Portkey: %s\n", key)

	resp, err = http.PostForm((serverURL + "wait"), url.Values{
		"key": {key},
	})
	if err != nil {
		panic(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	Decode(string(body), remoteSignal)
}

func signalExchangeWithKey(localSignal, remoteSignal *Signal, key *string) {
	resp, err := http.PostForm((serverURL + "key"), url.Values{
		"key":        {*key},
		"connParams": {Encode(localSignal)},
	})
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	Decode(string(body), remoteSignal)
}

func Encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}

func Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}
