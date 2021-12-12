package field

import (
	"io"
	"os"
	"path"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	log "github.com/sirupsen/logrus"
)

type AllianceStation struct {
	Position string // R1, B1, etc
	Team     int    // Team number
}

var (
	autoDuration    time.Duration
	teleopDuration  time.Duration
	endgameDuration time.Duration

	autoStartedAt    time.Time
	teleopStartedAt  time.Time
	endgameStartedAt time.Time

	AllianceStations map[string]*AllianceStation
)

var matchState, matchName, eventName string

const (
	stateIdle    = "Idle"
	stateAuto    = "Auto"
	stateTeleop  = "Teleop"
	stateEndGame = "Endgame"
)

// playSound plays a game sound file
func playSound(file string) {
	f, err := os.Open(path.Join("sounds/", file))
	if err != nil {
		log.Warn(err)
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		log.Warn(err)
	}

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
func Setup(auto, teleop, endGame, event string) error {
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

	eventName = event
	matchState = stateIdle

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
		"name":       matchName,
		"state":      matchState,
		"alliances":  TeamNumbers(),
		"event_name": eventName,
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
		go playSound("auto.mp3")
		matchState = stateAuto
		autoStartedAt = time.Now()
		time.Sleep(autoDuration)

		go playSound("teleop.mp3")
		matchState = stateTeleop
		teleopStartedAt = time.Now()
		time.Sleep(teleopDuration - endgameDuration)

		matchState = stateEndGame
		endgameStartedAt = time.Now()
		time.Sleep(endgameDuration)

		go playSound("end.mp3")
		matchState = stateIdle
	}()
}

// PlayAllSounds plays all game sounds to test audio levels
func PlayAllSounds() {
	playSound("auto.mp3")
	playSound("teleop.mp3")
	playSound("end.mp3")
	playSound("abort.mp3")
}

// UpdateTeamNumbers updates all alliance station team numbers
func UpdateTeamNumbers(alliances map[string]int) {
	if AllianceStations == nil {
		AllianceStations = map[string]*AllianceStation{}
	}

	for position, team := range alliances {
		if AllianceStations[position] == nil {
			AllianceStations[position] = &AllianceStation{Team: team}
		} else {
			AllianceStations[position].Team = team
		}
	}
}

// TeamNumbers gets a map of alliance station position to team number
func TeamNumbers() map[string]int {
	var o = make(map[string]int, len(AllianceStations))
	for position, allianceStation := range AllianceStations {
		o[position] = allianceStation.Team
	}
	return o
}

// UpdateMatchName sets the match name
func UpdateMatchName(n string) {
	matchName = n
}

// ResetAlliances clears all alliance stations
func ResetAlliances() {
	AllianceStations = map[string]*AllianceStation{}
}
