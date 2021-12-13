package field

import (
	"io"
	"os"
	"path"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	log "github.com/sirupsen/logrus"

	"github.com/natesales/bunnyfms/internal/driverstation"
)

var (
	autoDuration    time.Duration
	teleopDuration  time.Duration
	endgameDuration time.Duration

	autoTimer    *time.Timer
	teleopTimer  *time.Timer
	endgameTimer *time.Timer

	autoStartedAt    time.Time
	teleopStartedAt  time.Time
	endgameStartedAt time.Time
)

var (
	gameSounds            bool
	matchState, matchName string
)

const (
	stateIdle    = "Idle"
	stateAuto    = "Auto"
	stateTeleop  = "Teleop"
	stateEndGame = "Endgame"
)

// playSound plays a game sound file
func playSound(file string) {
	if !gameSounds {
		log.Warnf("Game sounds disabled, not playing %s", file)
		return
	}

	f, err := os.Open(path.Join("sounds/", file))
	if err != nil {
		log.Warn(err)
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		log.Warn(err)
	}

	// TODO: This panics if sound is already going
	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		log.Warn(err)
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		log.Warn(err)
	}
}

// Setup creates a new field setup (once per event)
func Setup(auto, teleop, endGame string, sounds bool) error {
	// Parse durations
	var err error
	autoDuration, err = time.ParseDuration(auto)
	if err != nil {
		return err
	}
	teleopDuration, err = time.ParseDuration(teleop)
	if err != nil {
		return err
	}
	endgameDuration, err = time.ParseDuration(endGame)
	if err != nil {
		return err
	}

	matchState = stateIdle
	gameSounds = sounds

	log.Infof("Configuring FMS with auto: %s, teleop: %s, endgame: %s, sounds: %v", autoDuration, teleopDuration, endgameDuration, sounds)

	return nil
}

func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return "-"
	}
	return time.Unix(0, 0).UTC().Add(d.Round(time.Second)).Format("4:05")
}

// State gets the game state
func State() map[string]interface{} {
	now := time.Now()

	o := map[string]interface{}{
		"name":      matchName,
		"state":     matchState,
		"alliances": TeamNumbers(),
		"ds":        driverstation.ConnectionStats(),
	}

	if matchState == "Idle" {
		o["auto_timer"] = formatDuration(autoDuration)
		o["teleop_timer"] = formatDuration(teleopDuration)
		o["endgame_timer"] = formatDuration(endgameDuration)
		o["current_timer"] = "0:00"
	} else {
		o["auto_timer"] = formatDuration(autoDuration - now.Sub(autoStartedAt))
		o["teleop_timer"] = formatDuration(teleopDuration - now.Sub(teleopStartedAt))
		o["endgame_timer"] = formatDuration(endgameDuration - now.Sub(endgameStartedAt))

		if matchState == "Auto" {
			o["current_timer"] = formatDuration(autoDuration - now.Sub(autoStartedAt))
		} else {
			o["current_timer"] = formatDuration(teleopDuration - now.Sub(teleopStartedAt))
		}
	}

	return o
}

// Start starts a match
func Start() {
	go func() {
		log.Infof("Match %s: starting auto", matchName)
		go playSound("auto.mp3")
		driverstation.StartAuto()
		matchState = stateAuto
		autoStartedAt = time.Now()
		autoTimer = time.NewTimer(autoDuration)
		<-autoTimer.C

		log.Infof("Match %s: starting teleop", matchName)
		go playSound("teleop.mp3")
		driverstation.StartTeleop()
		matchState = stateTeleop
		teleopStartedAt = time.Now()
		teleopTimer = time.NewTimer(teleopDuration - endgameDuration)
		<-teleopTimer.C

		log.Infof("Match %s: starting endgame", matchName)
		go playSound("endgame.mp3")
		matchState = stateEndGame
		driverstation.StopMatch()
		endgameStartedAt = time.Now()
		endgameTimer = time.NewTimer(endgameDuration)
		<-endgameTimer.C

		log.Infof("Match %s: finished", matchName)
		go playSound("end.mp3")
		matchState = stateIdle
	}()
}

// Stop stops a match
func Stop() {
	log.Infof("Match %s: aborting", matchName)
	go playSound("abort.mp3")
	matchState = "Idle"
	for _, timer := range []*time.Timer{autoTimer, teleopTimer, endgameTimer} {
		if timer != nil {
			timer.Stop()
		}
	}
}

// PlayAllSounds plays all game sounds to test audio levels
func PlayAllSounds() {
	playSound("auto.mp3")
	playSound("teleop.mp3")
	playSound("endgame.mp3")
	playSound("end.mp3")
	playSound("abort.mp3")
}

// UpdateTeamNumbers updates all alliance station team numbers
func UpdateTeamNumbers(alliances map[string]int) {
	if driverstation.AllianceStations == nil {
		driverstation.AllianceStations = map[string]*driverstation.AllianceStation{}
	}

	for position, team := range alliances {
		if driverstation.AllianceStations[position] == nil {
			driverstation.AllianceStations[position] = &driverstation.AllianceStation{Team: team}
		} else {
			driverstation.AllianceStations[position].Team = team
		}
	}
}

// TeamNumbers gets a map of alliance station position to team number
func TeamNumbers() map[string]int {
	var o = make(map[string]int, len(driverstation.AllianceStations))
	for position, allianceStation := range driverstation.AllianceStations {
		o[position] = allianceStation.Team
	}
	return o
}

// UpdateMatchName sets the match name
func UpdateMatchName(n string) {
	log.Infof("Updating match name to %s", n)
	matchName = n
}

// ResetAlliances clears all alliance stations
func ResetAlliances() {
	log.Info("Resetting alliances")
	driverstation.CloseAll()
	driverstation.AllianceStations = map[string]*driverstation.AllianceStation{}
}
