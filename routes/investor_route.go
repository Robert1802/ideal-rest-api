package routes

import (
	"ideal-rest-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func InvestorRoute(app *fiber.App) {

	// Investor
	// Create Investor
	app.Post("/investor", controllers.CreateInvestor)
	// Get Investor By CPF
	app.Get("/investor/:cpf", controllers.GetAInvestor)
	// Get All Investors
	app.Get("/investors", controllers.GetAllInvestors)
	// Edit Investor
	app.Put("/investor/:cpf", controllers.EditAInvestor)
	// Delete Investor
	app.Delete("/investor/:cpf", controllers.DeleteAInvestor)

	// Asset
	// Add Asset to Investor List
	app.Post("/investor/asset/:cpf", controllers.InsertAssetOnInvestor)
	// Order Assets
	app.Post("/investor/asset/order/:type/:asc/:cpf", controllers.SortAssets)

	// Get Asset By Symbol
	app.Post("/investor/assets", controllers.GetAssetPrice)

}
