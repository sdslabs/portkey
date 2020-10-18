package signal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pion/webrtc/v2"
	"github.com/sdslabs/portkey/pkg/utils"
)

type Signal struct {
	ICECandidates  []webrtc.ICECandidate `json:"iceCandidates"`
	ICEParameters  webrtc.ICEParameters  `json:"iceParameters"`
	QuicParameters webrtc.QUICParameters `json:"quicParameters"`
}

var serverURL string = "https://portkey-server.herokuapp.com/"

func SignalExchange(localSignal, remoteSignal *Signal) error {
	connParams, err := utils.Encode(localSignal)
	if err != nil {
		return err
	}
	resp, err := http.PostForm(serverURL, url.Values{
		"connParams": {connParams},
	})
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	key := string(body)
	fmt.Printf("Your Portkey: %s\n", key)

	resp, err = http.PostForm((serverURL + "wait"), url.Values{
		"key": {key},
	})
	if err != nil {
		return err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	err = utils.Decode(string(body), remoteSignal)
	return err
}

func SignalExchangeWithKey(localSignal, remoteSignal *Signal, key string) error {
	connParams, err := utils.Encode(localSignal)
	if err != nil {
		return err
	}
	resp, err := http.PostForm((serverURL + "key"), url.Values{
		"key":        {key},
		"connParams": {connParams},
	})
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	err = utils.Decode(string(body), remoteSignal)
	return err
}
