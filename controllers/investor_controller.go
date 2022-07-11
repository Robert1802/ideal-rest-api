package controllers

import (
	"context"
	"ideal-rest-api/configs"
	"ideal-rest-api/models"
	"ideal-rest-api/responses"
	s "ideal-rest-api/services"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var investorCollection *mongo.Collection = configs.GetCollection(configs.DB, "investors")

func CreateInvestor(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var investor models.Investor
	defer cancel()
	var validate = validator.New()

	err := c.BodyParser(&investor)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if validationErr := validate.Struct(&investor); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if len(investor.CPF) != 11 {
		return c.Status(http.StatusFound).JSON(&fiber.Map{"data": "CPF must have 11 characters"})
	}

	newInvestor := models.Investor{
		CPF:    investor.CPF,
		Name:   investor.Name,
		Email:  investor.Email,
		Assets: investor.Assets,
	}

	check, _ := investorCollection.FindOne(ctx, bson.M{"cpf": investor.CPF}).DecodeBytes()

	if check != nil {
		return c.Status(http.StatusFound).JSON(&fiber.Map{"data": "CPF already in use"})
	}

	_, err = investorCollection.InsertOne(ctx, newInvestor)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	return c.Status(http.StatusCreated).JSON(newInvestor)

}

func GetAInvestor(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cpf := c.Params("cpf")
	var investor models.Investor
	defer cancel()

	err := investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&investor)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	assets := s.AssetListToAsset(investor.Assets)

	response, err := s.LatestPrices(assets)
	if err != nil {
		return err
	}

	_, _ = investorCollection.UpdateOne(ctx, bson.M{"cpf": cpf}, bson.M{"$set": bson.M{"assets": response}})

	err = investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&investor)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(investor)
}

func EditAInvestor(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cpf := c.Params("cpf")
	var investor models.Investor
	defer cancel()
	var validate = validator.New()

	if err := c.BodyParser(&investor); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if validationErr := validate.Struct(&investor); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(validationErr.Error())
	}

	update := bson.M{"name": investor.Name, "email": investor.Email}

	result, err := investorCollection.UpdateOne(ctx, bson.M{"cpf": cpf}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	var updatedUser models.UserInfo
	if result.MatchedCount == 1 {
		err := investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&updatedUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(err.Error())
		}
	}

	return c.Status(http.StatusOK).JSON(updatedUser)
}

func DeleteAInvestor(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cpf := c.Params("cpf")
	defer cancel()

	result, err := investorCollection.DeleteOne(ctx, bson.M{"cpf": cpf})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.InvestorResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Investor with specified CPF " + cpf + " not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.InvestorResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Investor successfully deleted!"}},
	)
}

func GetAllInvestors(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var investors []models.UserInfo
	defer cancel()

	results, err := investorCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.UserInfo
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(err.Error())
		}

		investors = append(investors, singleUser)
	}

	return c.Status(http.StatusOK).JSON(investors)
}

// func InsertAssetOnInvestor(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	cpf := c.Params("cpf")
// 	var assets = new(models.Asset)
// 	defer cancel()

// 	err := c.BodyParser(&assets)
// 	if err != nil {
// 		return err
// 	}

// 	var symbol interface{}

// 	// Loop only adds new assets and get their current prices
// 	iter := quote.List(assets.Symbol)
// 	for iter.Next() {
// 		symbol = nil
// 		_ = investorCollection.FindOne(ctx, bson.M{"cpf": cpf, "assets.symbol": iter.Quote().Symbol}).Decode(&symbol)
// 		q := iter.Quote()
// 		item := models.AssetList{
// 			Symbol: q.Symbol,
// 			Price:  q.RegularMarketPrice,
// 		}

// 		if symbol == nil {
// 			_, err := investorCollection.UpdateOne(ctx, bson.M{"cpf": cpf}, bson.M{"$push": bson.M{"assets": item}})
// 			if err != nil {
// 				return c.Status(http.StatusInternalServerError).JSON(err.Error())
// 			}
// 		}

// 	}

// 	var investor models.Investor

// 	err = investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&investor)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(err.Error())
// 	}

// 	totalAssets := s.AssetListToAsset(investor.Assets)

// 	response, err := s.LatestPrices(totalAssets)
// 	if err != nil {
// 		return err
// 	}

// 	_, _ = investorCollection.UpdateOne(ctx, bson.M{"cpf": cpf}, bson.M{"$set": bson.M{"assets": response}})

// 	err = investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&investor)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Status(http.StatusOK).JSON(investor)

// }

// func RemoveAssetOfInvestor(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	cpf := c.Params("cpf")
// 	var assets = new(models.Asset)
// 	defer cancel()

// 	err := c.BodyParser(&assets)
// 	if err != nil {
// 		return err
// 	}

// 	var symbol interface{}

// 	// Loop only adds new assets and get their current prices
// 	iter := quote.List(assets.Symbol)
// 	for iter.Next() {
// 		symbol = nil
// 		_ = investorCollection.FindOne(ctx, bson.M{"cpf": cpf, "assets.symbol": iter.Quote().Symbol}).Decode(&symbol)
// 		q := iter.Quote()
// 		item := models.AssetList{
// 			Symbol: q.Symbol,
// 			Price:  q.RegularMarketPrice,
// 		}

// 		if symbol != nil {
// 			_, err := investorCollection.UpdateOne(ctx, bson.M{"cpf": cpf}, bson.M{"$pull": bson.M{"assets": item}})
// 			if err != nil {
// 				return c.Status(http.StatusInternalServerError).JSON(err.Error())
// 			}
// 		}

// 	}

// 	var investor models.Investor

// 	err = investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&investor)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(err.Error())
// 	}

// 	totalAssets := s.AssetListToAsset(investor.Assets)

// 	response, err := s.LatestPrices(totalAssets)
// 	if err != nil {
// 		return err
// 	}

// 	_, _ = investorCollection.UpdateOne(ctx, bson.M{"cpf": cpf}, bson.M{"$set": bson.M{"assets": response}})

// 	err = investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&investor)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Status(http.StatusOK).JSON(investor)

// }

// func SortAssets(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	cpf := c.Params("cpf")
// 	sortType := c.Params("type")     // name, price or list
// 	ascendingType := c.Params("asc") // ascendent = 0 or descendent = 1

// 	var investor = new(models.Investor)
// 	var assetList []models.AssetList

// 	err := investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&investor)
// 	if err != nil {
// 		return err
// 	}

// 	switch sortType {
// 	case "name":
// 		assetList = s.SortByName(*investor, ascendingType)
// 	case "price":
// 		assetList = s.SortByPrice(*investor, ascendingType)
// 	case "list":
// 		var assetsOrder = new(models.Asset)
// 		err = c.BodyParser(&assetsOrder)
// 		if err != nil {
// 			return err
// 		}
// 		assetList = s.SortByList(*investor, *assetsOrder, ascendingType)
// 	default:
// 		return c.Status(http.StatusNotFound).JSON(&fiber.Map{"data": "Sorting type " + sortType + " does not exist"})
// 	}

// 	assets := s.AssetListToAsset(assetList)

// 	response, err := s.LatestPrices(assets)
// 	if err != nil {
// 		return err
// 	}

// 	_, _ = investorCollection.UpdateOne(ctx, bson.M{"cpf": cpf}, bson.M{"$set": bson.M{"assets": response}})

// 	err = investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&investor)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Status(http.StatusOK).JSON(investor)

// }

// func GetAssetPrice(c *fiber.Ctx) error {

// 	var assets = new(models.Asset)

// 	err := c.BodyParser(&assets)
// 	if err != nil {
// 		return err
// 	}

// 	response, err := s.LatestPrices(*assets)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Status(http.StatusOK).JSON(response)
// }
