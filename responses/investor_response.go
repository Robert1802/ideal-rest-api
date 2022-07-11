package responses

import "github.com/gofiber/fiber/v2"

type Investor struct {
	CPF    string      `json:"cpf,omitempty"`
	Name   string      `json:"name,omitempty"`
	Email  string      `json:"email,omitempty"`
	Assets []AssetList `json:"assets,omitempty"`
}

type AssetList struct {
	Symbol string  `json:"symbol,omitempty"`
	Price  float64 `json:"price,omitempty"`
}

type Asset struct {
	Symbol []string `json:"symbol,omitempty"`
}

type InvestorResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}

type UserInfo struct {
	CPF   string `json:"cpf,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
