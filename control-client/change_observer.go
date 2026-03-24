// Licensed under Elastic License 2.0
// See LICENSE.txt for details

package controlclient

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ChangeObserver polls /api/v1/changes and dispatches callbacks when changes are detected
type ChangeObserver struct {
	client      *Client
	accountID   string
	interval    time.Duration
	callbacks   map[string][]func()
	localHashes map[string]string
	mu          sync.RWMutex
	stopChan    chan struct{}
	running     bool
	logFunc     func(format string, args ...interface{})
}

// NewChangeObserver creates a new change observer
func NewChangeObserver(client *Client, accountID string, interval time.Duration) *ChangeObserver {
	return &ChangeObserver{
		client:      client,
		accountID:   accountID,
		interval:    interval,
		callbacks:   make(map[string][]func()),
		localHashes: make(map[string]string),
		stopChan:    make(chan struct{}),
	}
}

// SetLogFunc sets a custom log function. If not set, logs are discarded.
func (o *ChangeObserver) SetLogFunc(f func(format string, args ...interface{})) {
	o.logFunc = f
}

func (o *ChangeObserver) log(format string, args ...interface{}) {
	if o.logFunc != nil {
		o.logFunc(format, args...)
	}
}

// OnChange registers a callback for a change category.
// The callback is called when the hash for that category changes.
func (o *ChangeObserver) OnChange(category string, callback func()) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.callbacks[category] = append(o.callbacks[category], callback)
}

// Start begins polling for changes in a background goroutine
func (o *ChangeObserver) Start() {
	o.mu.Lock()
	if o.running {
		o.mu.Unlock()
		return
	}
	o.running = true
	o.mu.Unlock()

	go o.pollLoop()
}

// Stop stops the change observer
func (o *ChangeObserver) Stop() {
	o.mu.Lock()
	defer o.mu.Unlock()
	if !o.running {
		return
	}
	o.running = false
	close(o.stopChan)
}

func (o *ChangeObserver) pollLoop() {
	ticker := time.NewTicker(o.interval)
	defer ticker.Stop()

	// Do an initial poll to seed local hashes (no callbacks on first poll)
	o.seedHashes()

	for {
		select {
		case <-o.stopChan:
			return
		case <-ticker.C:
			o.poll()
		}
	}
}

// seedHashes does an initial poll to populate local hashes without firing callbacks
func (o *ChangeObserver) seedHashes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := o.client.GetChanges(ctx, o.accountID)
	if err != nil {
		o.log("change observer: initial seed failed (will retry): %v", err)
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	o.applyHash(resp.Tenants, ChangeCategoryTenants)
	o.applyHash(resp.Datasets, ChangeCategoryDatasets)
	o.applyHash(resp.ProxyConfig, ChangeCategoryProxyConfig)
	o.applyHash(resp.Transformations, ChangeCategoryTransformations)
	o.applyHash(resp.Account, ChangeCategoryAccount)
}

func (o *ChangeObserver) applyHash(status *ChangeStatus, category string) {
	if status != nil {
		o.localHashes[category] = status.Hash
	}
}

func (o *ChangeObserver) poll() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := o.client.GetChanges(ctx, o.accountID)
	if err != nil {
		// Graceful degradation: old control versions return 404
		// Only log occasionally to avoid spam
		return
	}

	o.mu.Lock()
	changed := o.detectChanges(resp)
	o.mu.Unlock()

	// Fire callbacks outside the lock
	for _, category := range changed {
		o.mu.RLock()
		cbs := o.callbacks[category]
		o.mu.RUnlock()
		for _, cb := range cbs {
			go cb()
		}
	}
}

// detectChanges compares response hashes with local state, updates local state,
// and returns a list of categories that changed. Must be called with o.mu held.
func (o *ChangeObserver) detectChanges(resp *ChangesResponse) []string {
	var changed []string

	changed = append(changed, o.checkCategory(resp.Tenants, ChangeCategoryTenants)...)
	changed = append(changed, o.checkCategory(resp.Datasets, ChangeCategoryDatasets)...)
	changed = append(changed, o.checkCategory(resp.ProxyConfig, ChangeCategoryProxyConfig)...)
	changed = append(changed, o.checkCategory(resp.Transformations, ChangeCategoryTransformations)...)
	changed = append(changed, o.checkCategory(resp.Account, ChangeCategoryAccount)...)

	return changed
}

func (o *ChangeObserver) checkCategory(status *ChangeStatus, category string) []string {
	if status == nil {
		return nil
	}
	old := o.localHashes[category]
	if old != status.Hash {
		o.localHashes[category] = status.Hash
		if old != "" { // Don't fire on first discovery
			o.log("change observer: %s changed (hash %s → %s)", category, old, status.Hash)
			return []string{category}
		}
	}
	return nil
}

// String returns a description of the observer for logging
func (o *ChangeObserver) String() string {
	scope := "global"
	if o.accountID != "" {
		scope = fmt.Sprintf("account:%s", o.accountID)
	}
	return fmt.Sprintf("ChangeObserver(interval=%v, scope=%s)", o.interval, scope)
}
