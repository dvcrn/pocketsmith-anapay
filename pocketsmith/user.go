package pocketsmith

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID                       int    `json:"id"`
	Login                    string `json:"login"`
	Name                     string `json:"name"`
	Email                    string `json:"email"`
	AvatarURL                string `json:"avatar_url"`
	BetaUser                 bool   `json:"beta_user"`
	TimeZone                 string `json:"time_zone"`
	WeekStartDay             int    `json:"week_start_day"`
	IsReviewingTransactions  bool   `json:"is_reviewing_transactions"`
	BaseCurrencyCode         string `json:"base_currency_code"`
	AlwaysShowBaseCurrency   bool   `json:"always_show_base_currency"`
	UsingMultipleCurrencies  bool   `json:"using_multiple_currencies"`
	AvailableAccounts        int    `json:"available_accounts"`
	AvailableBudgets         int    `json:"available_budgets"`
	ForecastLastUpdatedAt    string `json:"forecast_last_updated_at"`
	ForecastLastAccessedAt   string `json:"forecast_last_accessed_at"`
	ForecastStartDate        string `json:"forecast_start_date"`
	ForecastEndDate          string `json:"forecast_end_date"`
	ForecastDeferRecalculate bool   `json:"forecast_defer_recalculate"`
	ForecastNeedsRecalculate bool   `json:"forecast_needs_recalculate"`
	LastLoggedInAt           string `json:"last_logged_in_at"`
	LastActivityAt           string `json:"last_activity_at"`
	CreatedAt                string `json:"created_at"`
	UpdatedAt                string `json:"updated_at"`
}

func (c *Client) GetCurrentUser() (*User, error) {
	url := "https://api.pocketsmith.com/v2/me"
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

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
