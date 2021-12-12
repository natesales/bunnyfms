package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/natesales/bunnyfms/internal/api"
	"github.com/natesales/bunnyfms/internal/driverstation"
	"github.com/natesales/bunnyfms/internal/field"
)

var (
	listenAddr      = flag.String("listen", ":8080", "HTTP listen address")
	autoDuration    = flag.String("auto-duration", "10s", "Auto duration")
	teleOpDuration  = flag.String("teleop-duration", "2m20s", "Telop duration")
	endgameDuration = flag.String("endgame-duration", "10s", "Endgame duration")
	eventName       = flag.String("event-name", "", "Event name")
)

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)

	if err := field.Setup(*autoDuration, *teleOpDuration, *endgameDuration, *eventName); err != nil {
		log.Fatal(err)
	}

	driverstation.Start()

	log.Printf("Starting HTTP server on %s", *listenAddr)
	api.Serve(*listenAddr)
}
