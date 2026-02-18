package config

var immediateReloadRequests = make(chan struct{}, 1)

// RequestImmediateReload requests a configuration reload as soon as possible.
// The request is coalesced if one is already pending.
func RequestImmediateReload() {
	select {
	case immediateReloadRequests <- struct{}{}:
	default:
	}
}

// ImmediateReloadRequests returns a channel that emits when an immediate reload is requested.
func ImmediateReloadRequests() <-chan struct{} {
	return immediateReloadRequests
}
