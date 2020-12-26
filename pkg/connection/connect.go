package connection

import (
	"sync"

	"github.com/pion/quic"
	"github.com/pion/webrtc/v3"
	log "github.com/sirupsen/logrus"

	"github.com/sdslabs/portkey/pkg/benchmark"
	"github.com/sdslabs/portkey/pkg/session"
	"github.com/sdslabs/portkey/pkg/signal"
)

var wg sync.WaitGroup
var stunServers []string = []string{"stun:stun.l.google.com:19302"}

func Connect(key string, sendPath string, receive bool, receivePath string, doBenchmarking bool) {
	isOffer := (key == "")
	api := webrtc.NewAPI()
	iceOptions := webrtc.ICEGatherOptions{
		ICEServers: []webrtc.ICEServer{
			{URLs: stunServers},
		},
	}
	gatherer, err := api.NewICEGatherer(iceOptions)
	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("Constructing ICE transport...")
	ice := api.NewICETransport(gatherer)

	log.Infoln("Constructing Quic transport...")
	qt, err := api.NewQUICTransport(ice, nil)
	if err != nil {
		log.Fatal(err)
	}

	receiveErr := make(chan error)
	if receive {
		wg.Add(1)
		qt.OnBidirectionalStream(func(stream *quic.BidirectionalStream) {
			log.Infof("New stream received: streamid = %d\n", stream.StreamID())
			go session.ReadLoop(stream, receivePath, receiveErr, &wg)
		})
		log.Infoln("Deployed incoming stream handler")
	}

	gatherFinished := make(chan struct{})
	gatherer.OnLocalCandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			close(gatherFinished)
		}
	})

	log.Infoln("Gathering ICE candidates...")
	err = gatherer.Gather()
	if err != nil {
		log.Fatal(err)
	}

	<-gatherFinished

	iceCandidates, err := gatherer.GetLocalCandidates()
	if err != nil {
		log.Fatal(err)
	}

	iceParams, err := gatherer.GetLocalParameters()
	if err != nil {
		log.Fatal(err)
	}

	quicParams, err := qt.GetLocalParameters()
	if err != nil {
		log.Fatal(err)
	}

	s := signal.Signal{
		ICECandidates:  iceCandidates,
		ICEParameters:  iceParams,
		QuicParameters: quicParams,
	}

	remoteSignal := signal.Signal{}

	if isOffer {
		err = signal.SignalExchange(&s, &remoteSignal)
	} else {
		err = signal.SignalExchangeWithKey(&s, &remoteSignal, key)
	}
	if err != nil {
		log.WithError(err).Fatalln("Unable to exchange signal")
	}

	log.Infoln("ICE candidates exchange successful")

	iceRole := webrtc.ICERoleControlled
	if isOffer {
		iceRole = webrtc.ICERoleControlling
	}

	err = ice.SetRemoteCandidates(remoteSignal.ICECandidates)
	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("Starting ICE transport...")
	err = ice.Start(nil, remoteSignal.ICEParameters, &iceRole)
	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("Starting Quic transport...")
	err = qt.Start(remoteSignal.QuicParameters)
	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("------------Connection established------------")

	if doBenchmarking {
		if err = benchmark.StartTransfer(isOffer); err != nil {
			log.WithError(err).Errorln("Error in starting benchmarking")
		}
		defer func() {
			if err = benchmark.EndTransfer(isOffer); err != nil {
				log.WithError(err).Errorln("Error in ending benchmarking")
			}
		}()
	}

	sendErr := make(chan error)
	if sendPath != "" {
		stream, err := qt.CreateBidirectionalStream()
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("New stream created: streamid = %d\n", stream.StreamID())
		wg.Add(1)
		go session.WriteLoop(stream, sendPath, sendErr, &wg)
	}

	wg.Wait()

	log.Infoln("Closing Quic transport...")
	if err = qt.Stop(quic.TransportStopInfo{}); err != nil {
		log.Fatal(err)
	}

	log.Infoln("Closing ICE transport...")
	if err = ice.Stop(); err != nil {
		log.Fatal(err)
	}
}
