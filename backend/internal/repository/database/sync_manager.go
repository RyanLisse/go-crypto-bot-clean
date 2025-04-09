package database

import (
	"context"
	"log"
	"sync"
	"time"
)

// SyncManager manages synchronization with the cloud database
type SyncManager struct {
	repo         Repository
	syncInterval time.Duration
	stopCh       chan struct{}
	wg           sync.WaitGroup
	mu           sync.Mutex
	isRunning    bool
}

// NewSyncManager creates a new synchronization manager
func NewSyncManager(repo Repository, syncInterval time.Duration) *SyncManager {
	return &SyncManager{
		repo:         repo,
		syncInterval: syncInterval,
		stopCh:       make(chan struct{}),
	}
}

// Start starts the synchronization manager
func (m *SyncManager) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return
	}

	// Check if the repository supports synchronization
	if !m.repo.SupportsSynchronization() {
		log.Println("Repository does not support synchronization, sync manager will not start")
		return
	}

	m.isRunning = true
	m.wg.Add(1)

	go m.syncLoop()
}

// Stop stops the synchronization manager
func (m *SyncManager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return
	}

	close(m.stopCh)
	m.wg.Wait()
	m.isRunning = false
	m.stopCh = make(chan struct{})
}

// SyncNow triggers an immediate synchronization
func (m *SyncManager) SyncNow(ctx context.Context) error {
	if !m.repo.SupportsSynchronization() {
		return nil
	}

	return m.repo.Synchronize(ctx)
}

// syncLoop runs the synchronization loop
func (m *SyncManager) syncLoop() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.syncInterval)
	defer ticker.Stop()

	// Perform initial sync
	ctx := context.Background()
	if err := m.repo.Synchronize(ctx); err != nil {
		log.Printf("Initial synchronization failed: %v", err)
	} else {
		timestamp, _ := m.repo.GetLastSyncTimestamp(ctx)
		log.Printf("Initial synchronization completed at %v", timestamp)
	}

	for {
		select {
		case <-ticker.C:
			ctx := context.Background()
			if err := m.repo.Synchronize(ctx); err != nil {
				log.Printf("Synchronization failed: %v", err)
			} else {
				timestamp, _ := m.repo.GetLastSyncTimestamp(ctx)
				log.Printf("Synchronization completed at %v", timestamp)
			}
		case <-m.stopCh:
			// Perform final sync before stopping
			ctx := context.Background()
			if err := m.repo.Synchronize(ctx); err != nil {
				log.Printf("Final synchronization failed: %v", err)
			} else {
				timestamp, _ := m.repo.GetLastSyncTimestamp(ctx)
				log.Printf("Final synchronization completed at %v", timestamp)
			}
			return
		}
	}
}
