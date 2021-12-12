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
	eventName       = flag.String("event-name", "Offseason Event", "Event name")
	noDriveStations = flag.Bool("no-ds", false, "Disable drive station communication")
	noSounds        = flag.Bool("no-sounds", false, "Disable game sounds")
)

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)

	if err := field.Setup(*autoDuration, *teleOpDuration, *endgameDuration, *eventName, !*noSounds); err != nil {
		log.Fatal(err)
	}

	if !*noDriveStations {
		driverstation.StartComms()
	} else {
		log.Warn("-no-ds flag set, not enabling driver station communication")
	}

	log.Printf("Starting HTTP server on %s", *listenAddr)
	api.Serve(*listenAddr)
}
