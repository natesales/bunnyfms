package field

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	log "github.com/sirupsen/logrus"
)

var (
	autoDuration    time.Duration
	teleopDuration  time.Duration
	endgameDuration time.Duration

	autoStartedAt    time.Time
	teleopStartedAt  time.Time
	endgameStartedAt time.Time
)

var matchState string

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
func Setup(auto, teleop, endGame string) error {
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

	return nil
}

func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return "-"
	}
	return fmt.Sprintf("%+v", d.Round(time.Second))
}

// State gets the game state
func State() map[string]interface{} {
	now := time.Now()

	if matchState == "Idle" {
		return map[string]interface{}{
			"state":         matchState,
			"auto_timer":    formatDuration(autoDuration),
			"teleop_timer":  formatDuration(teleopDuration),
			"endgame_timer": formatDuration(endgameDuration),
		}
	} else {
		return map[string]interface{}{
			"state":         matchState,
			"auto_timer":    formatDuration(autoDuration - now.Sub(autoStartedAt)),
			"teleop_timer":  formatDuration(teleopDuration - now.Sub(teleopStartedAt)),
			"endgame_timer": formatDuration(endgameDuration - now.Sub(endgameStartedAt)),
		}
	}
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
