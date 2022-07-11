package models

type Investor struct {
	CPF    string      `json:"cpf,omitempty" validate:"required"`
	Name   string      `json:"name,omitempty" validate:"required"`
	Email  string      `json:"email,omitempty" validate:"required"`
	Assets []AssetList `json:"assets,omitempty"`
}

type AssetList struct {
	Symbol string  `json:"symbol,omitempty"`
	Price  float64 `json:"price,omitempty"`
}

type Asset struct {
	Symbol []string `json:"symbol,omitempty"`
}

type UserInfo struct {
	CPF   string `json:"cpf,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
