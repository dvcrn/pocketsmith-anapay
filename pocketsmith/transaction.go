package pocketsmith

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateTransaction struct {
	Payee        string   `json:"payee"`
	Amount       float64  `json:"amount"`
	Date         string   `json:"date"`
	IsTransfer   bool     `json:"is_transfer"`
	Labels       []string `json:"labels,omitempty"`
	CategoryID   int      `json:"category_id,omitempty"`
	Note         string   `json:"note,omitempty"`
	Memo         string   `json:"memo,omitempty"`
	ChequeNumber string   `json:"cheque_number,omitempty"`
	NeedsReview  bool     `json:"needs_review"`
}

type Category struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Colour          *string    `json:"colour"`
	IsTransfer      bool       `json:"is_transfer"`
	IsBill          bool       `json:"is_bill"`
	RefundBehaviour *string    `json:"refund_behaviour"`
	Children        []Category `json:"children"`
	ParentID        *int       `json:"parent_id"`
	RollUp          bool       `json:"roll_up"`
	CreatedAt       string     `json:"created_at"`
	UpdatedAt       string     `json:"updated_at"`
}

type Transaction struct {
	ID                   int64              `json:"id"`
	Payee                string             `json:"payee"`
	OriginalPayee        string             `json:"original_payee"`
	Date                 string             `json:"date"`
	UploadSource         string             `json:"upload_source"`
	Category             Category           `json:"category"`
	ClosingBalance       float64            `json:"closing_balance"`
	ChequeNumber         *string            `json:"cheque_number"`
	Memo                 *string            `json:"memo"`
	Amount               float64            `json:"amount"`
	AmountInBaseCurrency float64            `json:"amount_in_base_currency"`
	Type                 string             `json:"type"`
	IsTransfer           bool               `json:"is_transfer"`
	NeedsReview          bool               `json:"needs_review"`
	Status               string             `json:"status"`
	Note                 *string            `json:"note"`
	Labels               []string           `json:"labels"`
	TransactionAccount   TransactionAccount `json:"transaction_account"`
	CreatedAt            string             `json:"created_at"`
	UpdatedAt            string             `json:"updated_at"`
}

func (c *Client) AddTransaction(accountID int, transaction *CreateTransaction) (*CreateTransaction, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transaction_accounts/%d/transactions", accountID)

	payload, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
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

	// {"id":1140680818,"payee":"チャージ","original_payee":"チャージ","date":"2024-12-04","upload_source":"api","category":null,"closing_balance":2346.0,"cheque_number":null,"memo":null,"amount":5000.0,"amount_in_base_currency":33.21,"type":"credit","is_transfer":true,"needs_review":false,"status":"posted","note":"チャージ","labels":[],"transaction_account":{"id":3466150,"account_id":3370081,"name":"ANA Pay","latest_feed_name":null,"number":null,"type":"credits","offline":true,"is_net_worth":false,"currency_code":"jpy","current_balance":0.0,"current_balance_in_base_currency":0.0,"current_balance_exchange_rate":null,"current_balance_date":"2024-12-05","current_balance_source":"closing_balance_as_today","data_feeds_balance_type":"balance","safe_balance":null,"safe_balance_in_base_currency":null,"has_safe_balance_adjustment":false,"starting_balance":0.0,"starting_balance_date":"2024-12-05","institution":{"id":1266583,"title":"ANA Pay","currency_code":"jpy","colour":"#4C5BA5","logo_url":null,"favicon_data_uri":null,"created_at":"2024-12-05T05:00:24Z","updated_at":"2024-12-05T05:00:24Z"},"data_feeds_account_id":null,"data_feeds_connection_id":null,"created_at":"2024-12-05T05:00:24Z","updated_at":"2024-12-05T05:00:24Z"},"created_at":"2024-12-05T05:17:26Z","updated_at":"2024-12-05T05:17:26Z"}

	var createdTransaction CreateTransaction
	err = json.NewDecoder(resp.Body).Decode(&createdTransaction)
	if err != nil {
		return nil, err
	}

	return &createdTransaction, nil
}

func (c *Client) SearchTransactions(accountID int, startDate, endDate, search string) ([]*Transaction, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transaction_accounts/%d/transactions", accountID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if startDate != "" {
		q.Add("start_date", startDate)
	}
	if endDate != "" {
		q.Add("end_date", endDate)
	}
	if search != "" {
		q.Add("search", search)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-Developer-Key", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var transactions []*Transaction
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (c *Client) ListTransactions(accountID int, startDate, endDate, updatedSince string, uncategorised int, transactionType string, needsReview int, search string, page int) ([]*Transaction, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transaction_accounts/%d/transactions", accountID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if startDate != "" {
		q.Add("start_date", startDate)
	}
	if endDate != "" {
		q.Add("end_date", endDate)
	}
	if updatedSince != "" {
		q.Add("updated_since", updatedSince)
	}
	if uncategorised > 0 {
		q.Add("uncategorised", fmt.Sprintf("%d", uncategorised))
	}
	if transactionType != "" {
		q.Add("type", transactionType)
	}
	if needsReview > 0 {
		q.Add("needs_review", fmt.Sprintf("%d", needsReview))
	}
	if search != "" {
		q.Add("search", search)
	}
	if page > 0 {
		q.Add("page", fmt.Sprintf("%d", page))
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-Developer-Key", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var transactions []*Transaction
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
