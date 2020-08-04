package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type TimeLog struct {
	startTime time.Time
	endTime   time.Time
}

var timeFormat = time.UnixDate
var offererLog TimeLog
var receiverLog TimeLog
var offererOn bool = false
var receiverOn bool = false

func startTransfer(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.WithError(err).Infoln("Bad request on /start")
		w.WriteHeader(400)
		return
	}
	isOffer := (req.FormValue("isOffer") == "true")
	startTime, err := time.Parse(timeFormat, req.FormValue("time"))
	if err != nil || (isOffer && offererOn) || (!isOffer && receiverOn) {
		if err != nil {
			log.WithError(err).Infoln("Bad request on /start")
		} else {
			log.Infoln("Bad request on /start")
		}
		w.WriteHeader(400)
		return
	}
	if isOffer {
		offererLog.startTime = startTime
		offererOn = true
		log.Infof("Offerer starts at %v\n", startTime)
	} else {
		receiverLog.startTime = startTime
		receiverOn = true
		log.Infof("Receiver starts at %v\n", startTime)
	}
	w.WriteHeader(200)
}

func endTransfer(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.WithError(err).Infoln("Bad request on /end")
		w.WriteHeader(400)
		return
	}
	isOffer := (req.FormValue("isOffer") == "true")
	endTime, err := time.Parse(timeFormat, req.FormValue("time"))
	if err != nil || (isOffer && !offererOn) || (!isOffer && !receiverOn) {
		if err != nil {
			log.WithError(err).Infoln("Bad request on /end")
		} else {
			log.Infoln("Bad request on /end")
		}
		w.WriteHeader(400)
		return
	}
	if isOffer {
		offererLog.endTime = endTime
		offererOn = false
		log.Infof("Offerer ends at %v\n", endTime)
	} else {
		receiverLog.endTime = endTime
		receiverOn = false
		log.Infof("Receiver ends at %v\n", endTime)
	}
	if !offererOn && !receiverOn {
		printDuration()
	}
	w.WriteHeader(200)
}

func printDuration() {
	transferStart := offererLog.startTime
	transferEnd := offererLog.endTime
	if offererLog.startTime.After(receiverLog.startTime) {
		transferStart = receiverLog.startTime
	}
	if offererLog.endTime.Before(receiverLog.endTime) {
		transferEnd = receiverLog.endTime
	}
	transferDuration := transferEnd.Sub(transferStart)
	log.Infof("Transfer duration : %v\n\n\n", transferDuration)
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		log.Info("env variable PORT not found. Defaulting to PORT=8080\n\n")
		port = "8080"
	}
	return ":" + port
}

func main() {
	http.HandleFunc("/start", startTransfer)
	http.HandleFunc("/end", endTransfer)
	log.Fatal(http.ListenAndServe(getPort(), nil))
}
