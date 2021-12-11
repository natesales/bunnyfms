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
	Message string `json:"message"`
	Arg     string `json:"arg"`
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
				// field.Stop()
				// TODO
			case "ds_reconnect":
				log.Debug("Reconnecting to driver stations")
				driverstation.Reset()
			case "estop":
				log.Debugf("Estopping %s", msg.Arg)
				// TODO
			}
		}
	}))
}

// Serve starts the API server
func Serve(listenAddr string) {
	if app == nil {
		log.Debug("Registering application handlers")
		register()
	}
	log.Fatal(app.Listen(listenAddr))
}
