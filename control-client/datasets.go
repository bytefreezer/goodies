// Licensed under Elastic License 2.0
// See LICENSE.txt for details

package controlclient

import (
	"context"
	"fmt"
)

// ListDatasets retrieves all datasets for a tenant
func (c *Client) ListDatasets(ctx context.Context, tenantID string, limit int) ([]Dataset, error) {
	// Need to determine the account ID from tenant first, or change the API path
	// For now, using the direct tenant path which should work
	path := fmt.Sprintf("/api/v1/tenants/%s/datasets", tenantID)
	if limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", path, limit)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Items []Dataset `json:"items"`
		Total int       `json:"total"`
	}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

// GetDataset retrieves a dataset by ID
func (c *Client) GetDataset(ctx context.Context, tenantID, datasetID string) (*Dataset, error) {
	path := fmt.Sprintf("/api/v1/tenants/%s/datasets/%s", tenantID, datasetID)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var dataset Dataset
	if err := c.parseResponse(resp, &dataset); err != nil {
		return nil, err
	}

	return &dataset, nil
}
