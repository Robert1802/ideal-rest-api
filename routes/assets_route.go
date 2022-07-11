package routes

import (
	"ideal-rest-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func AssetsRoute(app *fiber.App) {

	// Asset
	// Get Asset By Symbol
	app.Post("/assets", controllers.GetAssetPrice)
	// Add Asset to Investor List
	app.Post("/asset/:cpf", controllers.InsertAssetOnInvestor)
	// Remove Asset to Investor List
	app.Post("/asset/:cpf/remove", controllers.RemoveAssetOfInvestor)
	// Order Assets
	app.Post("/asset/order/:type/:asc/:cpf", controllers.SortAssets)

}
