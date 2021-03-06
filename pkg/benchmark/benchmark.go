package benchmark

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

var apiURL string = net.JoinHostPort("http://localhost", getPort())
var timeFormat = time.UnixDate

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

func StartTransfer(isOffer bool) error {
	isOfferString := "false"
	if isOffer {
		isOfferString = "true"
	}
	resp, err := http.PostForm((apiURL + "/start"), url.Values{
		"time":    {time.Now().Format(timeFormat)},
		"isOffer": {isOfferString},
	})
	if err != nil {
		return fmt.Errorf("Error in StartTransfer post request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Received status %d from server", resp.StatusCode)
	}
	return nil
}

func EndTransfer(isOffer bool) error {
	isOfferString := "false"
	if isOffer {
		isOfferString = "true"
	}
	resp, err := http.PostForm((apiURL + "/end"), url.Values{
		"time":    {time.Now().Format(timeFormat)},
		"isOffer": {isOfferString},
	})
	if err != nil {
		return fmt.Errorf("Error in EndTransfer post request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Received status %d from server", resp.StatusCode)
	}
	return nil
}
