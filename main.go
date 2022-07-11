package main

import (
	"ideal-rest-api/configs"
	"ideal-rest-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//Executa database
	configs.ConnectDB()

	//Rotas
	routes.InvestorRoute(app)

	app.Listen(":6000")
}
