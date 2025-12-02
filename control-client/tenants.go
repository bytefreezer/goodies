// Licensed under Elastic License 2.0
// See LICENSE.txt for details

package controlclient

import (
	"context"
	"fmt"
)

// CreateTenant creates a new tenant for an account
func (c *Client) CreateTenant(ctx context.Context, accountID string, req CreateTenantRequest) (*Tenant, error) {
	path := fmt.Sprintf("/api/v1/accounts/%s/tenants", accountID)
	resp, err := c.doRequest(ctx, "POST", path, req)
	if err != nil {
		return nil, err
	}

	var tenant Tenant
	if err := c.parseResponse(resp, &tenant); err != nil {
		return nil, err
	}

	return &tenant, nil
}

// GetTenant retrieves a tenant by ID
func (c *Client) GetTenant(ctx context.Context, accountID, tenantID string) (*Tenant, error) {
	path := fmt.Sprintf("/api/v1/accounts/%s/tenants/%s", accountID, tenantID)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var tenant Tenant
	if err := c.parseResponse(resp, &tenant); err != nil {
		return nil, err
	}

	return &tenant, nil
}

// ListTenants retrieves all tenants for an account
func (c *Client) ListTenants(ctx context.Context, accountID string, limit int) ([]Tenant, error) {
	path := fmt.Sprintf("/api/v1/accounts/%s/tenants", accountID)
	if limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", path, limit)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Items []Tenant `json:"items"`
		Total int      `json:"total"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

// UpdateTenant updates an existing tenant
func (c *Client) UpdateTenant(ctx context.Context, accountID, tenantID string, req UpdateTenantRequest) (*Tenant, error) {
	path := fmt.Sprintf("/api/v1/accounts/%s/tenants/%s", accountID, tenantID)
	resp, err := c.doRequest(ctx, "PUT", path, req)
	if err != nil {
		return nil, err
	}

	var tenant Tenant
	if err := c.parseResponse(resp, &tenant); err != nil {
		return nil, err
	}

	return &tenant, nil
}

// DeleteTenant deletes a tenant
func (c *Client) DeleteTenant(ctx context.Context, accountID, tenantID string) error {
	path := fmt.Sprintf("/api/v1/accounts/%s/tenants/%s", accountID, tenantID)
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
