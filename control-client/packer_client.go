// Licensed under Elastic License 2.0
// See LICENSE.txt for details

package controlclient

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ====================================================================================
// PACKER TENANT LOCK OPERATIONS
// ====================================================================================

// AcquireTenantLock attempts to acquire a lock on a tenant for packer processing
func (c *Client) AcquireTenantLock(ctx context.Context, tenantID, lockedBy string, lockDurationSeconds int) (*PackerTenantLock, error) {
	payload := map[string]interface{}{
		"tenant_id":             tenantID,
		"locked_by":             lockedBy,
		"lock_duration_seconds": lockDurationSeconds,
	}

	resp, err := c.doRequest(ctx, "POST", "/api/v1/packer/locks/tenants", payload)
	if err != nil {
		return nil, err
	}

	var lock PackerTenantLock
	if err := c.parseResponse(resp, &lock); err != nil {
		return nil, err
	}

	return &lock, nil
}

// ReleaseTenantLock releases a tenant lock
func (c *Client) ReleaseTenantLock(ctx context.Context, tenantID, lockedBy string) error {
	path := fmt.Sprintf("/api/v1/packer/locks/tenants/%s?locked_by=%s", tenantID, url.QueryEscape(lockedBy))
	resp, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// UpdateTenantLockHeartbeat updates the heartbeat for a tenant lock
func (c *Client) UpdateTenantLockHeartbeat(ctx context.Context, tenantID, lockedBy string) error {
	payload := map[string]interface{}{
		"locked_by": lockedBy,
	}

	path := fmt.Sprintf("/api/v1/packer/locks/tenants/%s/heartbeat", tenantID)
	resp, err := c.doRequest(ctx, "PUT", path, payload)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// CheckTenantLock checks if a tenant is currently locked
func (c *Client) CheckTenantLock(ctx context.Context, tenantID string) (*PackerTenantLock, error) {
	path := fmt.Sprintf("/api/v1/packer/locks/tenants/%s", tenantID)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil // No lock exists
	}

	var lock PackerTenantLock
	if err := c.parseResponse(resp, &lock); err != nil {
		return nil, err
	}

	return &lock, nil
}

// CleanupExpiredTenantLocks removes all expired tenant locks
func (c *Client) CleanupExpiredTenantLocks(ctx context.Context) (int, error) {
	resp, err := c.doRequest(ctx, "DELETE", "/api/v1/packer/locks/tenants/cleanup/expired", nil)
	if err != nil {
		return 0, err
	}

	var result CountResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return 0, err
	}

	return result.Count, nil
}

// ClearAllTenantLocks removes all tenant locks (use with caution)
func (c *Client) ClearAllTenantLocks(ctx context.Context) (int, error) {
	resp, err := c.doRequest(ctx, "DELETE", "/api/v1/packer/locks/tenants/cleanup/all", nil)
	if err != nil {
		return 0, err
	}

	var result CountResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return 0, err
	}

	return result.Count, nil
}

// CleanupStaleTenantLocks removes tenant locks with stale heartbeats
func (c *Client) CleanupStaleTenantLocks(ctx context.Context, staleThresholdSeconds int) (int, error) {
	path := fmt.Sprintf("/api/v1/packer/locks/tenants/cleanup/stale?stale_threshold_seconds=%d", staleThresholdSeconds)
	resp, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return 0, err
	}

	var result CountResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return 0, err
	}

	return result.Count, nil
}

// ====================================================================================
// PACKER PARQUET METADATA OPERATIONS
// ====================================================================================

// UpsertParquetFileMetadata inserts or updates parquet file metadata
func (c *Client) UpsertParquetFileMetadata(ctx context.Context, metadata *PackerParquetFileMetadata) (*PackerParquetFileMetadata, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/packer/metadata/files", metadata)
	if err != nil {
		return nil, err
	}

	var result PackerParquetFileMetadata
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetParquetFileMetadataByPartition retrieves metadata for files in a specific partition
func (c *Client) GetParquetFileMetadataByPartition(ctx context.Context, tenantID, datasetID, partitionPath string, limit int) ([]PackerParquetFileMetadata, error) {
	path := fmt.Sprintf("/api/v1/packer/metadata/files/%s/%s", tenantID, datasetID)
	params := url.Values{}
	if partitionPath != "" {
		params.Set("partition_path", partitionPath)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Items []PackerParquetFileMetadata `json:"items"`
		Total int                         `json:"total"`
	}
	if err := c.parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// GetAllParquetFileMetadata retrieves all metadata for a tenant/dataset
func (c *Client) GetAllParquetFileMetadata(ctx context.Context, tenantID, datasetID string, limit int) ([]PackerParquetFileMetadata, error) {
	path := fmt.Sprintf("/api/v1/packer/metadata/files/%s/%s/all", tenantID, datasetID)
	if limit > 0 {
		path += fmt.Sprintf("?limit=%d", limit)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Items []PackerParquetFileMetadata `json:"items"`
		Total int                         `json:"total"`
	}
	if err := c.parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// DeleteParquetFileMetadata deletes metadata for a specific file or partition
func (c *Client) DeleteParquetFileMetadata(ctx context.Context, tenantID, datasetID, filePath string) error {
	path := fmt.Sprintf("/api/v1/packer/metadata/files/%s/%s", tenantID, datasetID)
	if filePath != "" {
		path += "?file_path=" + url.QueryEscape(filePath)
	}

	resp, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// CleanupOrphanedParquetMetadata removes metadata for files that no longer exist in S3
func (c *Client) CleanupOrphanedParquetMetadata(ctx context.Context, tenantID, datasetID string) (int, error) {
	payload := map[string]interface{}{
		"tenant_id":  tenantID,
		"dataset_id": datasetID,
	}

	resp, err := c.doRequest(ctx, "POST", "/api/v1/packer/metadata/files/cleanup/orphaned", payload)
	if err != nil {
		return 0, err
	}

	var result CountResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return 0, err
	}

	return result.Count, nil
}

// CleanupExpiredParquetMetadata removes expired parquet metadata
func (c *Client) CleanupExpiredParquetMetadata(ctx context.Context) (int, error) {
	resp, err := c.doRequest(ctx, "DELETE", "/api/v1/packer/metadata/files/cleanup/expired", nil)
	if err != nil {
		return 0, err
	}

	var result CountResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return 0, err
	}

	return result.Count, nil
}

// ====================================================================================
// PACKER METADATA GENERATION STATUS OPERATIONS
// ====================================================================================

// UpdateMetadataGenerationStatus updates the metadata generation status for a partition
func (c *Client) UpdateMetadataGenerationStatus(ctx context.Context, status *PackerMetadataGenerationStatus) error {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/packer/metadata/generation/status", status)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// GetMetadataGenerationStatus retrieves the metadata generation status for a partition
func (c *Client) GetMetadataGenerationStatus(ctx context.Context, tenantID, datasetID, partitionPath string) (*PackerMetadataGenerationStatus, error) {
	path := fmt.Sprintf("/api/v1/packer/metadata/generation/status/%s/%s/%s", tenantID, datasetID, url.PathEscape(partitionPath))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}

	var status PackerMetadataGenerationStatus
	if err := c.parseResponse(resp, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// ====================================================================================
// PACKER METADATA SUMMARY OPERATIONS
// ====================================================================================

// GetParquetMetadataSummary retrieves aggregated metadata summary for a partition
func (c *Client) GetParquetMetadataSummary(ctx context.Context, tenantID, datasetID, partitionPath string) (*PackerParquetMetadataSummary, error) {
	path := fmt.Sprintf("/api/v1/packer/metadata/summary/%s/%s/%s", tenantID, datasetID, url.PathEscape(partitionPath))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}

	var summary PackerParquetMetadataSummary
	if err := c.parseResponse(resp, &summary); err != nil {
		return nil, err
	}

	return &summary, nil
}
