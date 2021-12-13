package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	log "github.com/sirupsen/logrus"

	"github.com/natesales/bunnyfms/internal/driverstation"
	"github.com/natesales/bunnyfms/internal/field"
)

var app *fiber.App

type message struct {
	Message         string         `json:"message"`
	AllianceStation string         `json:"alliance_station"`
	Alliances       map[string]int `json:"alliances"`
	Name            string         `json:"name"`
}

func register() {
	app = fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Static("/", "static/")

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			var msg message
			if err := c.ReadJSON(&msg); err != nil {
				log.Println("read:", err)
				break
			}

			switch msg.Message {
			case "ping":
				if err := c.WriteJSON(field.State()); err != nil {
					log.Println("write:", err)
				}
			case "start":
				log.Debug("Starting match")
				field.Start()
			case "stop":
				log.Debug("Stopping match")
				field.Stop()
			case "ds_reconnect":
				log.Debug("Reconnecting to driver stations")
				driverstation.ResetComms()
			case "estop":
				log.Debugf("Estopping %s", msg.AllianceStation)
				driverstation.Estop(msg.AllianceStation)
			case "test_sounds":
				log.Debug("Playing all sounds")
				field.PlayAllSounds()
			case "update_alliances":
				log.Debugf("Updating alliances to %+v", msg.Alliances)
				field.UpdateTeamNumbers(msg.Alliances)
			case "match_name":
				log.Debugf("Updating match name to %+v", msg.Name)
				field.UpdateMatchName(msg.Name)
			case "reset_alliances":
				log.Debug("Resetting alliances")
				field.ResetAlliances()
			}
		}
	}))
}

// Serve starts the API server
func Serve(listenAddr string) {
	if app == nil {
		register()
	}
	log.Fatal(app.Listen(listenAddr))
}
