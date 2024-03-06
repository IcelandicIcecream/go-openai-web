package ai

import (
	"context"
	"fmt"
	"icelandicicecream/openai-go/model"
	"log"

	"github.com/carlmjohnson/requests"
	"github.com/sashabaranov/go-openai"
)

var (
	accountingUrl     = "http://localhost:8080"
	getAccountsPath   = "v2/accounts"
	getCurrenciesPath = "v2/currencies"
	getCategoriesPath = "v2/categories"
)

var GetAccountsDefinition = openai.FunctionDefinition{
	Name:        "GetAccounts",
	Description: "Get user's chart of accounts",
}

var GetAccountsTool = openai.Tool{
	Type:     openai.ToolTypeFunction,
	Function: &GetAccountsDefinition,
}

func GetAccounts(orgSchema string) (accounts []string, err error) {
	ctx := context.Background()

	var response struct {
		Message string          `json:"message"`
		Payload []model.Account `json:"payload"`
	}

	err = requests.URL(fmt.Sprintf("%s/%s", accountingUrl, getAccountsPath)).ToJSON(&response).Header("X-Org-Id", orgSchema).Fetch(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, account := range response.Payload {
		accounts = append(accounts, account.Name)
	}

	return accounts, nil
}

func GetCurrencies(orgSchema string) ([]model.Currency, error) {
	ctx := context.Background()

	var response struct {
		Message string           `json:"message"`
		Payload []model.Currency `json:"payload"`
	}

	err := requests.URL(fmt.Sprintf("%s/%s", accountingUrl, getCurrenciesPath)).Header("X-Org-Id", orgSchema).ToJSON(&response).Fetch(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return response.Payload, nil
}

func GetCategories(orgSchema string) ([]model.Category, error) {
	ctx := context.Background()

	var response struct {
		Message string           `json:"message"`
		Payload []model.Category `json:"payload"`
	}

	err := requests.URL(fmt.Sprintf("%s/%s", accountingUrl, getCategoriesPath)).Header("X-Org-Id", orgSchema).ToJSON(&response).Fetch(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return response.Payload, nil
}
