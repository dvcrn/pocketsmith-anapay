package pocketsmith

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AccountType string

const (
	AccountTypeBank           AccountType = "bank"
	AccountTypeCredits        AccountType = "credits"
	AccountTypeCash           AccountType = "cash"
	AccountTypeLoans          AccountType = "loans"
	AccountTypeMortgage       AccountType = "mortgage"
	AccountTypeStocks         AccountType = "stocks"
	AccountTypeVehicle        AccountType = "vehicle"
	AccountTypeProperty       AccountType = "property"
	AccountTypeInsurance      AccountType = "insurance"
	AccountTypeOtherAsset     AccountType = "other_asset"
	AccountTypeOtherLiability AccountType = "other_liability"
)

type Scenario struct {
	ID                           int     `json:"id"`
	Title                        string  `json:"title"`
	Description                  string  `json:"description"`
	InterestRate                 float64 `json:"interest_rate"`
	InterestRateRepeatID         int     `json:"interest_rate_repeat_id"`
	Type                         string  `json:"type"`
	MinimumValue                 float64 `json:"minimum_value"`
	MaximumValue                 float64 `json:"maximum_value"`
	AchieveDate                  string  `json:"achieve_date"`
	StartingBalance              float64 `json:"starting_balance"`
	StartingBalanceDate          string  `json:"starting_balance_date"`
	ClosingBalance               float64 `json:"closing_balance"`
	ClosingBalanceDate           string  `json:"closing_balance_date"`
	CurrentBalance               float64 `json:"current_balance"`
	CurrentBalanceDate           string  `json:"current_balance_date"`
	CurrentBalanceInBaseCurrency float64 `json:"current_balance_in_base_currency"`
	CurrentBalanceExchangeRate   float64 `json:"current_balance_exchange_rate"`
	SafeBalance                  float64 `json:"safe_balance"`
	SafeBalanceInBaseCurrency    float64 `json:"safe_balance_in_base_currency"`
	CreatedAt                    string  `json:"created_at"`
	UpdatedAt                    string  `json:"updated_at"`
}

type TransactionAccount struct {
	ID                           int         `json:"id"`
	Name                         string      `json:"name"`
	Number                       string      `json:"number"`
	CurrentBalance               float64     `json:"current_balance"`
	CurrentBalanceDate           string      `json:"current_balance_date"`
	CurrentBalanceInBaseCurrency float64     `json:"current_balance_in_base_currency"`
	CurrentBalanceExchangeRate   float64     `json:"current_balance_exchange_rate"`
	SafeBalance                  float64     `json:"safe_balance"`
	SafeBalanceInBaseCurrency    float64     `json:"safe_balance_in_base_currency"`
	StartingBalance              float64     `json:"starting_balance"`
	StartingBalanceDate          string      `json:"starting_balance_date"`
	CreatedAt                    string      `json:"created_at"`
	UpdatedAt                    string      `json:"updated_at"`
	Institution                  Institution `json:"institution"`
	CurrencyCode                 string      `json:"currency_code"`
	Type                         AccountType `json:"type"`
}

type Account struct {
	ID                           int                  `json:"id"`
	Title                        string               `json:"title"`
	CurrencyCode                 string               `json:"currency_code"`
	Type                         AccountType          `json:"type"`
	IsNetWorth                   bool                 `json:"is_net_worth"`
	PrimaryTransactionAccount    TransactionAccount   `json:"primary_transaction_account"`
	PrimaryScenario              Scenario             `json:"primary_scenario"`
	TransactionAccounts          []TransactionAccount `json:"transaction_accounts"`
	Scenarios                    []Scenario           `json:"scenarios"`
	CreatedAt                    string               `json:"created_at"`
	UpdatedAt                    string               `json:"updated_at"`
	CurrentBalance               float64              `json:"current_balance"`
	CurrentBalanceDate           string               `json:"current_balance_date"`
	CurrentBalanceInBaseCurrency float64              `json:"current_balance_in_base_currency"`
	CurrentBalanceExchangeRate   float64              `json:"current_balance_exchange_rate"`
	SafeBalance                  float64              `json:"safe_balance"`
	SafeBalanceInBaseCurrency    float64              `json:"safe_balance_in_base_currency"`
}

func (c *Client) ListAccounts(userID int) ([]Account, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/accounts", userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-Developer-Key", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var accounts []Account
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (c *Client) ListTransactionAccounts(userID int) ([]TransactionAccount, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/transaction_accounts", userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-Developer-Key", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var transactionAccounts []TransactionAccount
	err = json.NewDecoder(resp.Body).Decode(&transactionAccounts)
	if err != nil {
		return nil, err
	}

	return transactionAccounts, nil
}

func (c *Client) CreateAccount(userID int, institutionID int, title string, currencyCode string, accountType AccountType) (*Account, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/accounts", userID)

	payload := struct {
		InstitutionID int    `json:"institution_id"`
		Title         string `json:"title"`
		CurrencyCode  string `json:"currency_code"`
		Type          string `json:"type"`
	}{
		InstitutionID: institutionID,
		Title:         title,
		CurrencyCode:  currencyCode,
		Type:          string(accountType),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-Developer-Key", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var account Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (c *Client) FindAccountByName(userID int, name string) (*Account, error) {
	accounts, err := c.ListAccounts(userID)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if account.Title == name {
			return &account, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateTransactionAccount(id int, institutionID int, startingBalance float64, startingBalanceDate string) (*TransactionAccount, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transaction_accounts/%d", id)

	payload := struct {
		InstitutionID       int     `json:"institution_id"`
		StartingBalance     float64 `json:"starting_balance"`
		StartingBalanceDate string  `json:"starting_balance_date"`
	}{
		InstitutionID:       institutionID,
		StartingBalance:     startingBalance,
		StartingBalanceDate: startingBalanceDate,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-Developer-Key", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var transactionAccount TransactionAccount
	err = json.NewDecoder(resp.Body).Decode(&transactionAccount)
	if err != nil {
		return nil, err
	}

	return &transactionAccount, nil
}
