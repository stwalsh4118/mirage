package vault

import (
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	// ErrCircuitOpen is returned when the circuit breaker is open
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// CircuitBreaker implements the circuit breaker pattern for Vault operations
type CircuitBreaker struct {
	state            CircuitState
	failureCount     int
	lastFailure      time.Time
	successCount     int
	failureThreshold int
	timeout          time.Duration
	successThreshold int
	mu               sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker with default settings
func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{
		state:            CircuitClosed,
		failureThreshold: DefaultFailureThreshold,
		timeout:          DefaultCircuitTimeout,
		successThreshold: DefaultSuccessThreshold,
	}
}

// Call executes the given function if the circuit breaker allows it
func (cb *CircuitBreaker) Call(fn func() error) error {
	if !cb.canExecute() {
		return ErrCircuitOpen
	}

	err := fn()
	cb.recordResult(err)
	return err
}

// canExecute checks if the circuit breaker allows execution
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		// Check if enough time has passed to try half-open
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = CircuitHalfOpen
			cb.successCount = 0
			log.Info().Msg("Vault circuit breaker transitioning to half-open")
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	}
	return false
}

// recordResult updates the circuit breaker state based on operation result
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailure = time.Now()
		cb.successCount = 0

		if cb.state == CircuitClosed && cb.failureCount >= cb.failureThreshold {
			cb.state = CircuitOpen
			log.Warn().
				Int("failure_count", cb.failureCount).
				Int("threshold", cb.failureThreshold).
				Msg("Vault circuit breaker opened due to failures")
		} else if cb.state == CircuitHalfOpen {
			// Failed in half-open, go back to open
			cb.state = CircuitOpen
			log.Warn().Msg("Vault circuit breaker reopened after half-open failure")
		}
	} else {
		cb.successCount++

		if cb.state == CircuitHalfOpen && cb.successCount >= cb.successThreshold {
			// Recovered successfully, close the circuit
			cb.state = CircuitClosed
			cb.failureCount = 0
			cb.successCount = 0
			log.Info().Msg("Vault circuit breaker closed after successful recovery")
		} else if cb.state == CircuitClosed && cb.failureCount > 0 {
			// Reset failure count on success in closed state
			cb.failureCount = 0
		}
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetMetrics returns current metrics for the circuit breaker
func (cb *CircuitBreaker) GetMetrics() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"state":         cb.state.String(),
		"failure_count": cb.failureCount,
		"success_count": cb.successCount,
		"last_failure":  cb.lastFailure,
		"time_in_state": time.Since(cb.lastFailure),
	}
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = CircuitClosed
	cb.failureCount = 0
	cb.successCount = 0
	log.Info().Msg("Vault circuit breaker manually reset")
}
