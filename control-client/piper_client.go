package controlclient

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/bytedance/sonic"
)

// ====================================================================================
// PIPER FILE LOCK OPERATIONS
// ====================================================================================

// AcquireFileLock attempts to acquire a lock on a file for processing
func (c *Client) AcquireFileLock(ctx context.Context, tenantID, datasetID, fileKey, lockedBy string, lockDurationSeconds int) (*PiperFileLock, error) {
	payload := map[string]interface{}{
		"tenant_id":             tenantID,
		"dataset_id":            datasetID,
		"file_key":              fileKey,
		"locked_by":             lockedBy,
		"lock_duration_seconds": lockDurationSeconds,
	}

	resp, err := c.doRequest(ctx, "POST", "/api/v1/piper/locks/files", payload)
	if err != nil {
		return nil, err
	}

	var lock PiperFileLock
	if err := c.parseResponse(resp, &lock); err != nil {
		return nil, err
	}

	return &lock, nil
}

// ReleaseFileLock releases a lock on a file
func (c *Client) ReleaseFileLock(ctx context.Context, tenantID, datasetID, fileKey, lockedBy string) error {
	payload := map[string]interface{}{
		"tenant_id":  tenantID,
		"dataset_id": datasetID,
		"file_key":   fileKey,
		"locked_by":  lockedBy,
	}

	body, err := sonic.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, "DELETE", "/api/v1/piper/locks/files", body)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// CheckFileLock checks if a file is currently locked
func (c *Client) CheckFileLock(ctx context.Context, tenantID, datasetID, fileKey string) (*PiperFileLock, error) {
	path := fmt.Sprintf("/api/v1/piper/locks/files/%s/%s/%s", tenantID, datasetID, url.PathEscape(fileKey))
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil // No lock exists
	}

	var lock PiperFileLock
	if err := c.parseResponse(resp, &lock); err != nil {
		return nil, err
	}

	return &lock, nil
}

// CleanupExpiredFileLocks removes all expired file locks
func (c *Client) CleanupExpiredFileLocks(ctx context.Context) (int, error) {
	resp, err := c.doRequest(ctx, "DELETE", "/api/v1/piper/locks/files/cleanup/expired", nil)
	if err != nil {
		return 0, err
	}

	var result CountResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return 0, err
	}

	return result.Count, nil
}

// CleanupStaleFileLocks removes file locks with stale heartbeats
func (c *Client) CleanupStaleFileLocks(ctx context.Context, staleThresholdSeconds int) (int, error) {
	path := fmt.Sprintf("/api/v1/piper/locks/files/cleanup/stale?stale_threshold_seconds=%d", staleThresholdSeconds)
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
// PIPER JOB RECORD OPERATIONS
// ====================================================================================

// CreatePiperJob creates a new job record
func (c *Client) CreatePiperJob(ctx context.Context, job *PiperJobRecord) (*PiperJobRecord, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/piper/jobs", job)
	if err != nil {
		return nil, err
	}

	var createdJob PiperJobRecord
	if err := c.parseResponse(resp, &createdJob); err != nil {
		return nil, err
	}

	return &createdJob, nil
}

// UpdatePiperJobStatus updates the status of a job
func (c *Client) UpdatePiperJobStatus(ctx context.Context, jobID, status, errorMessage, outputFile string, recordsProcessed int64) error {
	payload := map[string]interface{}{
		"status":            status,
		"error_message":     errorMessage,
		"output_file":       outputFile,
		"records_processed": recordsProcessed,
	}

	path := fmt.Sprintf("/api/v1/piper/jobs/%s/status", jobID)
	resp, err := c.doRequest(ctx, "PUT", path, payload)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// GetPiperJob retrieves a job by ID
func (c *Client) GetPiperJob(ctx context.Context, jobID string) (*PiperJobRecord, error) {
	path := fmt.Sprintf("/api/v1/piper/jobs/%s", jobID)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}

	var job PiperJobRecord
	if err := c.parseResponse(resp, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// GetPiperJobsByStatus retrieves jobs by status
func (c *Client) GetPiperJobsByStatus(ctx context.Context, status string, limit int) ([]PiperJobRecord, error) {
	path := "/api/v1/piper/jobs"
	params := url.Values{}
	if status != "" {
		params.Set("status", status)
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
		Items []PiperJobRecord `json:"items"`
		Total int              `json:"total"`
	}
	if err := c.parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// GetPiperJobsForTenant retrieves all jobs for a tenant
func (c *Client) GetPiperJobsForTenant(ctx context.Context, tenantID string, limit int) ([]PiperJobRecord, error) {
	path := fmt.Sprintf("/api/v1/piper/jobs/tenant/%s", tenantID)
	if limit > 0 {
		path += fmt.Sprintf("?limit=%d", limit)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Items []PiperJobRecord `json:"items"`
		Total int              `json:"total"`
	}
	if err := c.parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// CleanupOldPiperJobs deletes jobs older than specified days
func (c *Client) CleanupOldPiperJobs(ctx context.Context, olderThanDays int) (int, error) {
	path := fmt.Sprintf("/api/v1/piper/jobs/cleanup/old?older_than_days=%d", olderThanDays)
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
// PIPER PIPELINE CONFIGURATION CACHE OPERATIONS
// ====================================================================================

// CachePipelineConfiguration caches a pipeline configuration
func (c *Client) CachePipelineConfiguration(ctx context.Context, tenantID, datasetID string, configuration map[string]interface{}, ttlSeconds int) error {
	payload := map[string]interface{}{
		"tenant_id":     tenantID,
		"dataset_id":    datasetID,
		"configuration": configuration,
		"ttl_seconds":   ttlSeconds,
	}

	resp, err := c.doRequest(ctx, "POST", "/api/v1/piper/cache/pipelines", payload)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// GetCachedPipelineConfiguration retrieves a cached pipeline configuration
func (c *Client) GetCachedPipelineConfiguration(ctx context.Context, tenantID, datasetID string) (*PiperPipelineConfiguration, error) {
	path := fmt.Sprintf("/api/v1/piper/cache/pipelines/%s/%s", tenantID, datasetID)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}

	var config PiperPipelineConfiguration
	if err := c.parseResponse(resp, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// InvalidatePipelineConfiguration removes a cached pipeline configuration
func (c *Client) InvalidatePipelineConfiguration(ctx context.Context, tenantID, datasetID string) error {
	path := fmt.Sprintf("/api/v1/piper/cache/pipelines/%s/%s", tenantID, datasetID)
	resp, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// ListCachedPipelines lists all cached pipeline configurations
func (c *Client) ListCachedPipelines(ctx context.Context, limit int) ([]PiperPipelineConfiguration, error) {
	path := "/api/v1/piper/cache/pipelines"
	if limit > 0 {
		path += fmt.Sprintf("?limit=%d", limit)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Items []PiperPipelineConfiguration `json:"items"`
		Total int                          `json:"total"`
	}
	if err := c.parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// CleanupExpiredPipelineCache removes expired pipeline cache entries
func (c *Client) CleanupExpiredPipelineCache(ctx context.Context) (int, error) {
	resp, err := c.doRequest(ctx, "DELETE", "/api/v1/piper/cache/pipelines/cleanup/expired", nil)
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
// PIPER TENANT CACHE OPERATIONS
// ====================================================================================

// CacheTenant caches tenant information
func (c *Client) CacheTenant(ctx context.Context, tenantID string, tenantData map[string]interface{}, ttlSeconds int) error {
	payload := map[string]interface{}{
		"tenant_id":   tenantID,
		"tenant_data": tenantData,
		"ttl_seconds": ttlSeconds,
	}

	resp, err := c.doRequest(ctx, "POST", "/api/v1/piper/cache/tenants", payload)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// GetCachedTenants retrieves all cached tenants
func (c *Client) GetCachedTenants(ctx context.Context, limit int) ([]PiperTenantCache, error) {
	path := "/api/v1/piper/cache/tenants"
	if limit > 0 {
		path += fmt.Sprintf("?limit=%d", limit)
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Items []PiperTenantCache `json:"items"`
		Total int                `json:"total"`
	}
	if err := c.parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// InvalidateTenantCache removes cached tenant(s)
func (c *Client) InvalidateTenantCache(ctx context.Context, tenantID string) error {
	path := "/api/v1/piper/cache/tenants"
	if tenantID != "" {
		path += "?tenant_id=" + tenantID
	}

	resp, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// CleanupExpiredTenantCache removes expired tenant cache entries
func (c *Client) CleanupExpiredTenantCache(ctx context.Context) (int, error) {
	resp, err := c.doRequest(ctx, "DELETE", "/api/v1/piper/cache/tenants/cleanup/expired", nil)
	if err != nil {
		return 0, err
	}

	var result CountResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return 0, err
	}

	return result.Count, nil
}
