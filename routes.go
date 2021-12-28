package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Setup(app *fiber.App) {
	app.Get("/api/dashboard", monitor.New())
	app.Get("/api/bugun", bugunkiYemekler)
	app.Get("/api/yarin", yarinkiYemekler)
	app.Get("/api/tum", tumYemekler)

}
