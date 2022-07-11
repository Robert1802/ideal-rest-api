package controllers

import (
	"context"
	"ideal-rest-api/models"
	s "ideal-rest-api/services"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/piquette/finance-go/quote"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertAssetOnInvestor(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cpf := c.Params("cpf")
	var assets = new(models.Asset)
	defer cancel()

	err := c.BodyParser(&assets)
	if err != nil {
		return err
	}

	var symbol interface{}

	// Loop only adds new assets and get their current prices
	iter := quote.List(assets.Symbol)
	for iter.Next() {
		symbol = nil
		_ = investorCollection.FindOne(ctx, bson.M{"cpf": cpf, "assets.symbol": iter.Quote().Symbol}).Decode(&symbol)
		q := iter.Quote()
		item := models.AssetList{
			Symbol: q.Symbol,
			Price:  q.RegularMarketPrice,
		}

		if symbol == nil {
			_, err := investorCollection.UpdateOne(ctx, bson.M{"cpf": cpf}, bson.M{"$push": bson.M{"assets": item}})
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(err.Error())
			}
		}

	}

	var investor models.Investor

	err = investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&investor)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	totalAssets := s.AssetListToAsset(investor.Assets)

	response, err := s.LatestPrices(totalAssets)
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

func RemoveAssetOfInvestor(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cpf := c.Params("cpf")
	var assets = new(models.Asset)
	defer cancel()

	err := c.BodyParser(&assets)
	if err != nil {
		return err
	}

	var symbol interface{}

	iter := quote.List(assets.Symbol)
	for iter.Next() {
		symbol = nil
		_ = investorCollection.FindOne(ctx, bson.M{"cpf": cpf, "assets.symbol": iter.Quote().Symbol}).Decode(&symbol)
		q := iter.Quote()
		item := models.AssetList{
			Symbol: q.Symbol,
			Price:  q.RegularMarketPrice,
		}

		if symbol != nil {
			_, err := investorCollection.UpdateOne(ctx, bson.M{"cpf": cpf}, bson.M{"$pull": bson.M{"assets": item}})
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(err.Error())
			}
		}

	}

	var investor models.Investor

	err = investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&investor)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	totalAssets := s.AssetListToAsset(investor.Assets)

	response, err := s.LatestPrices(totalAssets)
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

func SortAssets(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cpf := c.Params("cpf")
	sortType := c.Params("type")     // name, price or list
	ascendingType := c.Params("asc") // ascendent = 0 or descendent = 1

	var investor = new(models.Investor)
	var assetList []models.AssetList

	err := investorCollection.FindOne(ctx, bson.M{"cpf": cpf}).Decode(&investor)
	if err != nil {
		return err
	}

	switch sortType {
	case "name":
		assetList = s.SortByName(*investor, ascendingType)
	case "price":
		assetList = s.SortByPrice(*investor, ascendingType)
	case "list":
		var assetsOrder = new(models.Asset)
		err = c.BodyParser(&assetsOrder)
		if err != nil {
			return err
		}
		assetList = s.SortByList(*investor, *assetsOrder, ascendingType)
	default:
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{"data": "Sorting type " + sortType + " does not exist"})
	}

	assets := s.AssetListToAsset(assetList)

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

func GetAssetPrice(c *fiber.Ctx) error {

	var assets = new(models.Asset)

	err := c.BodyParser(&assets)
	if err != nil {
		return err
	}

	response, err := s.LatestPrices(*assets)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(response)
}
