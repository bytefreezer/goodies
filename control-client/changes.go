// Licensed under Elastic License 2.0
// See LICENSE.txt for details

package controlclient

import (
	"context"
	"fmt"
	"time"
)

// Change category constants
const (
	ChangeCategoryTenants         = "tenants"
	ChangeCategoryDatasets        = "datasets"
	ChangeCategoryProxyConfig     = "proxy_config"
	ChangeCategoryTransformations = "transformations"
	ChangeCategoryAccount         = "account"
)

// ChangeStatus represents a single change category's state
type ChangeStatus struct {
	Hash        string    `json:"hash"`
	LastChanged time.Time `json:"last_changed"`
}

// ChangesResponse represents the GET /api/v1/changes response
type ChangesResponse struct {
	Tenants         *ChangeStatus `json:"tenants,omitempty"`
	Datasets        *ChangeStatus `json:"datasets,omitempty"`
	ProxyConfig     *ChangeStatus `json:"proxy_config,omitempty"`
	Transformations *ChangeStatus `json:"transformations,omitempty"`
	Account         *ChangeStatus `json:"account,omitempty"`
}

// GetChanges calls GET /api/v1/changes with optional account_id scoping
func (c *Client) GetChanges(ctx context.Context, accountID string) (*ChangesResponse, error) {
	path := "/api/v1/changes"
	if accountID != "" {
		path = fmt.Sprintf("%s?account_id=%s", path, accountID)
	}
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var result ChangesResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
