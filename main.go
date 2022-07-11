package main

import (
	"ideal-rest-api/configs"
	"ideal-rest-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//Execute database
	configs.ConnectDB()

	//Routes
	routes.InvestorRoute(app)
	routes.AssetsRoute(app)

	app.Listen(":6000")
}
