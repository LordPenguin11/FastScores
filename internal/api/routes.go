package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/ws", websocket.New(LiveScoreWebSocketHandler))
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Leagues endpoints
	app.Post("/leagues/sync", FetchLeaguesFromSportmonks)
	app.Get("/leagues", GetLeaguesHandler)
} 
