package field

import (
	"fmt"
	"time"
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
	return fmt.Sprintf("%02.0f:%02.0f", d.Minutes(), d.Seconds())
}

// State gets the game state
func State() map[string]interface{} {
	now := time.Now()
	return map[string]interface{}{
		"state":         matchState,
		"auto_timer":    formatDuration(autoDuration - now.Sub(autoStartedAt)),
		"teleop_timer":  formatDuration(teleopDuration - now.Sub(teleopStartedAt)),
		"endgame_timer": formatDuration(endgameDuration - now.Sub(endgameStartedAt)),
	}
}

// Start starts a match
func Start() {
	go func() {
		matchState = stateAuto
		autoStartedAt = time.Now()
		time.Sleep(autoDuration)

		matchState = stateTeleop
		teleopStartedAt = time.Now()
		time.Sleep(teleopDuration - endgameDuration)

		matchState = stateEndGame
		endgameStartedAt = time.Now()
		time.Sleep(endgameDuration)

		matchState = stateIdle
	}()
}
