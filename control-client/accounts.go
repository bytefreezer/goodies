package controlclient

import (
	"context"
	"fmt"
)

// CreateAccount creates a new account
func (c *Client) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/accounts", req)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := c.parseResponse(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// GetAccount retrieves an account by ID
func (c *Client) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	path := fmt.Sprintf("/api/v1/accounts/%s", accountID)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := c.parseResponse(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// ListAccounts retrieves all accounts
func (c *Client) ListAccounts(ctx context.Context, limit int) ([]Account, error) {
	path := "/api/v1/accounts"
	if limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", path, limit)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Items []Account `json:"items"`
		Total int       `json:"total"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

// UpdateAccount updates an existing account
func (c *Client) UpdateAccount(ctx context.Context, accountID string, req UpdateAccountRequest) (*Account, error) {
	path := fmt.Sprintf("/api/v1/accounts/%s", accountID)
	resp, err := c.doRequest(ctx, "PUT", path, req)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := c.parseResponse(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// DeleteAccount deletes an account
func (c *Client) DeleteAccount(ctx context.Context, accountID string) error {
	path := fmt.Sprintf("/api/v1/accounts/%s", accountID)
	resp, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	var result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("delete failed: %s", result.Message)
	}

	return nil
}
