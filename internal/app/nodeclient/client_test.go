package nodeclient

import (
	"testing"
	"time"
)

func TestClientHandshakeIsDueOnlyAfterInterval(t *testing.T) {
	client := &Client{}
	client.SetHandshake("1.7.1", []string{"ip_scope_lookup"})
	if client.HandshakeDue(time.Minute) {
		t.Fatal("fresh handshake must not be due")
	}

	client.handshakeAt = time.Now().Add(-2 * time.Minute)
	if !client.HandshakeDue(time.Minute) {
		t.Fatal("handshake older than the interval must be due")
	}

	client.DeferHandshake()
	if client.HandshakeDue(time.Minute) {
		t.Fatal("deferred handshake must not be due again before the interval")
	}
}
