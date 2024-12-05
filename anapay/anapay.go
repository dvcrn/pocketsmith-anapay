package anapay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

// LoginResponse represents the response from the login endpoint
type LoginResponse struct {
	EmailAuthenticated string `json:"emailAuthenticated"`
	LoginAuthID        string `json:"loginAuthId"`
	AccessToken        string `json:"accessToken"`
	TokenType          string `json:"tokenType"`
	ExpiresIn          string `json:"expiresIn"`
	RefreshToken       string `json:"refreshToken"`
	Scope              string `json:"scope"`
}

// AccountInfo represents the account information
type CreditCardInfo struct {
	PaymentSourceID        string      `json:"paymentSourceId"`
	Digit4byte             string      `json:"digit4byte"`
	Brand                  string      `json:"brand"`
	ChargeAvailability     string      `json:"chargeAvailability"`
	PaymentAvailability    string      `json:"paymentAvailability"`
	SpecificCreditCardType interface{} `json:"specificCreditCardType"`
}

type PointInfo struct {
	PointNumber           string `json:"pointNumber"`
	PointPriorityUseSetup string `json:"pointPriorityUseSetup"`
}

type NFCRegisterStatus struct {
	ApplePayIDRegisterStatus         string `json:"applePayIdRegisterStatus"`
	ApplePayVisaTouchRegisterStatus  string `json:"applePayVisaTouchRegisterStatus"`
	GooglePayIDRegisterStatus        string `json:"googlePayIdRegisterStatus"`
	GooglePayVisaTouchRegisterStatus string `json:"googlePayVisaTouchRegisterStatus"`
}

type ServiceRegisterStatus struct {
	HousePrepaid   string `json:"housePrepaid"`
	BrandPrepaid   string `json:"brandPrepaid"`
	CrankUp        string `json:"crankUp"`
	LinkSettlement string `json:"linkSettlement"`
	EKyc           string `json:"eKyc"`
	EKycResult     string `json:"eKycResult"`
}

type AccountInfo struct {
	AllianceID            string                `json:"allianceId"`
	ReferenceNumber       string                `json:"referenceNumber"`
	AccountStatus         string                `json:"accountStatus"`
	Balance               string                `json:"balance"`
	MainPaymentSourceID   string                `json:"mainPaymentSourceId"`
	CreditCardInfo        []CreditCardInfo      `json:"creditCardInfo"`
	BankPayInfo           []interface{}         `json:"bankPayInfo"`
	PointInfo             PointInfo             `json:"pointInfo"`
	NFCRegisterStatus     NFCRegisterStatus     `json:"nfcRegisterStatus"`
	ServiceRegisterStatus ServiceRegisterStatus `json:"serviceRegisterStatus"`
	BankpayFirstAuthFlag  string                `json:"bankpayFirstAuthFlag"`
}

// Transaction represents a single transaction
// example:
//
//	 {
//	  "saleDatetime": "20241223011207",
//	  "settlementType": "",
//	  "dealType": "05",
//	  "delKbn": "01",
//	  "descriptionType": "3009",
//	  "shopName": "",
//	  "amount": "5000",
//	  "walletSettlementNo": "23123401120741221241",
//	  "walletSettlementSubNo": "01",
//	  "pointConversionAmount": ""
//	}
type Transaction struct {
	SaleDatetime          string `json:"saleDatetime"`
	SettlementType        string `json:"settlementType"`
	DealType              string `json:"dealType"`
	DelKbn                string `json:"delKbn"`
	DescriptionType       string `json:"descriptionType"`
	ShopName              string `json:"shopName"`
	Amount                string `json:"amount"`
	WalletSettlementNo    string `json:"walletSettlementNo"`
	WalletSettlementSubNo string `json:"walletSettlementSubNo"`
	PointConversionAmount string `json:"pointConversionAmount"`
}

// TransactionResponse represents the response from the transactions endpoint
type TransactionResponse struct {
	History []Transaction `json:"history"`
}

// Login authenticates with the ANA service
func Login(anaWalletID, deviceID string) (*LoginResponse, error) {
	url := "https://teikei1.api.mkpst.com/ana/accounts/login"
	headers := map[string]string{
		"Host":            "teikei1.api.mkpst.com",
		"Accept":          "application/json",
		"User-Agent":      "ANAMileage/4.31.0 (jp.co.ana.anamile; build:4; iOS 18.1.0) Alamofire/5.9.1",
		"Accept-Language": "ja-JP;q=1.0, en-AU;q=0.9, de-JP;q=0.8",
		"Content-Type":    "application/json",
	}

	body := map[string]string{
		"anaWalletId": anaWalletID,
		"deviceId":    deviceID,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &loginResp, nil
}

// GetAccounts retrieves account information
func GetAccounts(accessToken string) (*AccountInfo, error) {
	url := "https://teikei1.api.mkpst.com/accounts?balanceReferenceFlag=1&nfcStatusReferenceFlag=1"
	headers := map[string]string{
		"Host":            "teikei1.api.mkpst.com",
		"Accept":          "application/json",
		"User-Agent":      "ANAMileage/4.31.0 (jp.co.ana.anamile; build:4; iOS 18.1.0) Alamofire/5.9.1",
		"Authorization":   "Bearer " + accessToken,
		"Accept-Language": "ja-JP;q=1.0, en-AU;q=0.9, de-JP;q=0.8",
		"Content-Type":    "application/json",
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// print body
	fmt.Println(string(respBody))

	var accountInfo AccountInfo
	if err := json.Unmarshal(respBody, &accountInfo); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &accountInfo, nil
}

// GetTransactions retrieves transaction history
func GetTransactions(accessToken string, pageNumber, pageSize int) ([]Transaction, error) {
	if pageNumber == 0 {
		pageNumber = 1
	}
	if pageSize == 0 {
		pageSize = 999
	}

	url := fmt.Sprintf("https://teikei1.api.mkpst.com/salesDetails?pageSize=%d&pageNumber=%d&historyType=&settlementType=",
		pageSize, pageNumber)

	headers := map[string]string{
		"Host":            "teikei1.api.mkpst.com",
		"Accept":          "application/json",
		"User-Agent":      "ANAMileage/4.31.0 (jp.co.ana.anamile; build:4; iOS 18.1.0) Alamofire/5.9.1",
		"Authorization":   "Bearer " + accessToken,
		"Accept-Language": "ja-JP;q=1.0, en-AU;q=0.9, de-JP;q=0.8",
		"Content-Type":    "application/json",
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var transResp TransactionResponse
	if err := json.Unmarshal(respBody, &transResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return transResp.History, nil
}
