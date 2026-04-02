package mocks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/myorg/worker-service/internal/domain"
	"github.com/myorg/worker-service/internal/repository"
)

// Ensure interface satisfaction at compile time.
var _ repository.JobRepository = (*MockJobRepository)(nil)

// MockJobRepository is an in-memory, thread-safe test double for JobRepository.
type MockJobRepository struct {
	mu   sync.Mutex
	jobs map[string]*domain.Job // keyed by event_id

	// Hooks for injecting errors in specific tests.
	CreateErr        error
	FindByEventIDErr error
	UpdateStatusErr  error
	IncrementRetryErr error
}

// NewMockJobRepository constructs an empty mock repository.
func NewMockJobRepository() *MockJobRepository {
	return &MockJobRepository{jobs: make(map[string]*domain.Job)}
}

func (m *MockJobRepository) Create(ctx context.Context, job *domain.Job) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().UTC()
	job.CreatedAt = now
	job.UpdatedAt = now

	// Store a copy to prevent external mutation.
	copy := *job
	m.jobs[job.EventID] = &copy
	return nil
}

func (m *MockJobRepository) FindByEventID(ctx context.Context, eventID string) (*domain.Job, error) {
	if m.FindByEventIDErr != nil {
		return nil, m.FindByEventIDErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	j, ok := m.jobs[eventID]
	if !ok {
		return nil, nil
	}
	copy := *j
	return &copy, nil
}

func (m *MockJobRepository) UpdateStatus(
	ctx context.Context,
	jobID string,
	status domain.JobStatus,
	txHashes []string,
	errMsg string,
) error {
	if m.UpdateStatusErr != nil {
		return m.UpdateStatusErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, j := range m.jobs {
		if j.ID == jobID {
			j.Status = status
			j.TxHashes = txHashes
			j.ErrorMessage = errMsg
			j.UpdatedAt = time.Now().UTC()
			return nil
		}
	}
	return fmt.Errorf("mock: job not found: %s", jobID)
}

func (m *MockJobRepository) IncrementRetry(ctx context.Context, jobID string, errMsg string) error {
	if m.IncrementRetryErr != nil {
		return m.IncrementRetryErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, j := range m.jobs {
		if j.ID == jobID {
			j.RetryCount++
			j.ErrorMessage = errMsg
			j.UpdatedAt = time.Now().UTC()
			return nil
		}
	}
	return fmt.Errorf("mock: job not found: %s", jobID)
}

// All returns every stored job (useful for assertions).
func (m *MockJobRepository) All() []*domain.Job {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]*domain.Job, 0, len(m.jobs))
	for _, j := range m.jobs {
		copy := *j
		result = append(result, &copy)
	}
	return result
}

// JobByEventID returns the stored job for an event_id, or nil.
func (m *MockJobRepository) JobByEventID(eventID string) *domain.Job {
	m.mu.Lock()
	defer m.mu.Unlock()
	j := m.jobs[eventID]
	if j == nil {
		return nil
	}
	copy := *j
	return &copy
}
