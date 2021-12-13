package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/natesales/bunnyfms/internal/api"
	"github.com/natesales/bunnyfms/internal/driverstation"
	"github.com/natesales/bunnyfms/internal/field"
)

var (
	adminListenAddr  = flag.String("admin", "localhost:8080", "Admin listen address")
	viewerListenAddr = flag.String("viewer", ":8081", "Viewer listen address")
	autoDuration     = flag.String("auto-duration", "10s", "Auto duration")
	teleOpDuration   = flag.String("teleop-duration", "2m20s", "Teleop duration")
	endgameDuration  = flag.String("endgame-duration", "30s", "Endgame duration")
	noDriveStations  = flag.Bool("no-ds", false, "Disable drive station communication")
	noSounds         = flag.Bool("no-sounds", false, "Disable game sounds")
)

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)

	if err := field.Setup(*autoDuration, *teleOpDuration, *endgameDuration, !*noSounds); err != nil {
		log.Fatal(err)
	}

	if !*noDriveStations {
		driverstation.StartComms()
	} else {
		log.Warn("-no-ds flag set, not enabling driver station communication")
	}

	api.Serve(*adminListenAddr, *viewerListenAddr)
}
