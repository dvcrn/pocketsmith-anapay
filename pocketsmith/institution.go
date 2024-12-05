package pocketsmith

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Institution struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	CurrencyCode string `json:"currency_code"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func (c *Client) CreateInstitution(userID int, title string, currencyCode string) (*Institution, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/institutions", userID)

	payload := struct {
		Title        string `json:"title"`
		CurrencyCode string `json:"currency_code"`
	}{
		Title:        title,
		CurrencyCode: currencyCode,
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

	// Read response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Response body: %s\n", string(bodyBytes))

	// Create new reader with the body bytes for subsequent decoding
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var institution Institution
	err = json.NewDecoder(resp.Body).Decode(&institution)
	if err != nil {
		return nil, err
	}

	return &institution, nil
}

func (c *Client) ListInstitutions(userID int) ([]Institution, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/institutions", userID)

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

	var institutions []Institution
	err = json.NewDecoder(resp.Body).Decode(&institutions)
	if err != nil {
		return nil, err
	}

	return institutions, nil
}

func (c *Client) FindInstitutionByName(userID int, name string) (*Institution, error) {
	institutions, err := c.ListInstitutions(userID)
	if err != nil {
		return nil, err
	}

	for _, institution := range institutions {
		if institution.Title == name {
			return &institution, nil
		}
	}

	return nil, ErrNotFound
}
