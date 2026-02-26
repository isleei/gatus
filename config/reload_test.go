package config

import "testing"

func TestRequestImmediateReload(t *testing.T) {
	// Drain any pending signal from prior tests
	select {
	case <-immediateReloadRequests:
	default:
	}

	RequestImmediateReload()
	select {
	case <-ImmediateReloadRequests():
		// success
	default:
		t.Fatal("expected a signal on the reload channel, got none")
	}
}

func TestRequestImmediateReload_Coalescing(t *testing.T) {
	// Drain any pending signal from prior tests
	select {
	case <-immediateReloadRequests:
	default:
	}

	// Send twice to fill the buffered channel (capacity 1) and verify it doesn't block
	RequestImmediateReload()
	RequestImmediateReload()

	// Should only be able to receive once
	select {
	case <-ImmediateReloadRequests():
	default:
		t.Fatal("expected at least one signal")
	}
	select {
	case <-ImmediateReloadRequests():
		t.Fatal("expected channel to be empty after one receive")
	default:
		// success
	}
}

func TestImmediateReloadRequests_ReturnsReadableChannel(t *testing.T) {
	ch := ImmediateReloadRequests()
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}
}
