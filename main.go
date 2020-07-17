package main

import (
	"flag"
	"fmt"

	"github.com/pion/quic"
	"github.com/pion/webrtc/v2"
)

const messageSize = 100

var stunServers []string = []string{"stun:stun.l.google.com:19302"}

type Signal struct {
	ICECandidates  []webrtc.ICECandidate `json:"iceCandidates"`
	ICEParameters  webrtc.ICEParameters  `json:"iceParameters"`
	QuicParameters webrtc.QUICParameters `json:"quicParameters"`
}

func main() {
	local := flag.Bool("local", false, "for testing on localhost server")
	key := flag.String("k", "", "key for connection")
	receive := flag.Bool("r", false, "set if want to receive files")
	sendPath := flag.String("s", "", "absolute path of file/directory to send")
	receivePath := flag.String("rpath", "", "absolute path of directory to save files in, pwd by default")

	flag.Parse()
	if *local {
		serverURL = "http://localhost:8080/"
	}
	isOffer := (*key == "")
	api := webrtc.NewAPI()
	iceOptions := webrtc.ICEGatherOptions{
		ICEServers: []webrtc.ICEServer{
			{URLs: stunServers},
		},
	}
	gatherer, err := api.NewICEGatherer(iceOptions)
	if err != nil {
		panic(err)
	}

	ice := api.NewICETransport(gatherer)

	qt, err := api.NewQUICTransport(ice, nil)
	if err != nil {
		panic(err)
	}

	err = gatherer.Gather()
	if err != nil {
		panic(err)
	}

	iceCandidates, err := gatherer.GetLocalCandidates()
	if err != nil {
		panic(err)
	}

	iceParams, err := gatherer.GetLocalParameters()
	if err != nil {
		panic(err)
	}

	quicParams, err := qt.GetLocalParameters()
	if err != nil {
		panic(err)
	}

	if *receive {
		qt.OnBidirectionalStream(func(stream *quic.BidirectionalStream) {
			fmt.Printf("New stream %d\n", stream.StreamID())
			go ReadLoop(stream, *receivePath) 
		})
	}

	s := Signal{
		ICECandidates:  iceCandidates,
		ICEParameters:  iceParams,
		QuicParameters: quicParams,
	}

	remoteSignal := Signal{}

	if isOffer {
		signalExchange(&s, &remoteSignal)
	} else {
		signalExchangeWithKey(&s, &remoteSignal, key)
	}
	iceRole := webrtc.ICERoleControlled
	if isOffer {
		iceRole = webrtc.ICERoleControlling
	}

	err = ice.SetRemoteCandidates(remoteSignal.ICECandidates)
	if err != nil {
		panic(err)
	}

	err = ice.Start(nil, remoteSignal.ICEParameters, &iceRole)
	if err != nil {
		panic(err)
	}

	err = qt.Start(remoteSignal.QuicParameters)
	if err != nil {
		panic(err)
	}

	if *sendPath != "" {
		stream, err := qt.CreateBidirectionalStream()
		if err != nil {
			panic(err)
		}
		fmt.Println("\n\n------------Connection established------------")
		go WriteLoop(stream, *sendPath)
	}

	select {}
}


