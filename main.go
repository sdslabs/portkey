package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

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
	key := flag.String("key", "", "key for connection")
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

	qt.OnBidirectionalStream(func(stream *quic.BidirectionalStream) {
		fmt.Printf("New stream %d\n", stream.StreamID())

		go ReadLoop(stream)

		go WriteLoop(stream, isOffer)
	})

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

	stream, err := qt.CreateBidirectionalStream()
	if err != nil {
		panic(err)
	}
	fmt.Println("\n\n------------Connection established------------")
	go ReadLoop(stream)
	go WriteLoop(stream, isOffer)

	select {}
}

func ReadLoop(s *quic.BidirectionalStream) {
	for {
		buffer := make([]byte, messageSize)
		params, err := s.ReadInto(buffer)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Message from stream '%d': %s\n", s.StreamID(), string(buffer[:params.Amount]))
	}
}

func WriteLoop(s *quic.BidirectionalStream, isOffer bool) {
	i := 0
	for range time.NewTicker(1 * time.Second).C {
		if isOffer {
			message := "so this shit works? " + strconv.Itoa(i)
			i++
			fmt.Printf("Sending %s \n", message)

			data := quic.StreamWriteParameters{
				Data: []byte(message),
			}
			err := s.Write(data)
			if err != nil {
				panic(err)
			}
		} else {
			message := "i can send too " + strconv.Itoa(i)
			i++
			fmt.Printf("Sending %s \n", message)

			data := quic.StreamWriteParameters{
				Data: []byte(message),
			}
			err := s.Write(data)
			if err != nil {
				panic(err)
			}
		}
	}
}
