package services

import (
	"encoding/json"
	"ideal-rest-api/models"
	"sort"

	"github.com/piquette/finance-go/quote"
)

func LatestPrices(assets models.Asset) ([]models.AssetList, error) {
	response := []models.AssetList{}

	var myMap []string
	data, err := json.Marshal(assets.Symbol)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &myMap)

	iter := quote.List(myMap)
	for iter.Next() {
		q := iter.Quote()
		item := models.AssetList{
			Symbol: q.Symbol,
			Price:  q.RegularMarketPrice,
		}
		response = append(response, item)
	}
	if iter.Err() != nil {
		return nil, iter.Err()
	}

	return response, nil
}

func SortByName(investor models.Investor, asc string) []models.AssetList {
	assetsList := investor.Assets

	// ascendent
	if asc == "0" {
		sort.Slice(assetsList, func(i, j int) bool {
			return assetsList[i].Symbol < assetsList[j].Symbol
		})
	}
	// descendent
	if asc == "1" {
		sort.Slice(assetsList, func(i, j int) bool {
			return assetsList[i].Symbol > assetsList[j].Symbol
		})
	}

	return assetsList
}

func SortByPrice(investor models.Investor, asc string) []models.AssetList {
	assetsList := investor.Assets

	// ascendent
	if asc == "0" {
		sort.Slice(assetsList, func(i, j int) bool {
			return assetsList[i].Price < assetsList[j].Price
		})
	}
	// descendent
	if asc == "1" {
		sort.Slice(assetsList, func(i, j int) bool {
			return assetsList[i].Price > assetsList[j].Price
		})
	}

	return assetsList
}

func SortByList(investor models.Investor, assetsOrder models.Asset, asc string) []models.AssetList {
	assetsList := investor.Assets

	var newAssets []models.AssetList
	for _, asor := range assetsOrder.Symbol {
		for _, asli := range assetsList {
			if asli.Symbol == asor {
				newAssets = append(newAssets, asli)
			}
		}
	}

	for _, asli := range assetsList {
		j := 0
		for _, asor := range assetsOrder.Symbol {
			if asor != asli.Symbol {
				j++
			}
			if j == len(assetsOrder.Symbol) {
				newAssets = append(newAssets, asli)
			}
		}
	}

	return newAssets
}

func AssetListToAsset(assetList []models.AssetList) models.Asset {
	var assets models.Asset
	for _, as := range assetList {
		assets.Symbol = append(assets.Symbol, as.Symbol)
	}

	return assets
}
