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

func TestClientCallOptionsCompressOnlyWhenNodeAdvertisesIt(t *testing.T) {
	legacy := &Client{}
	if opts := legacy.CallOptions(); len(opts) != 0 {
		t.Fatalf("node without grpc_gzip must receive uncompressed requests, got %d options", len(opts))
	}
	current := &Client{}
	current.SetHandshake("1.6.2", []string{CapabilityGzipRequests})
	if opts := current.CallOptions(); len(opts) != 1 {
		t.Fatalf("node with grpc_gzip must get the gzip call option, got %d options", len(opts))
	}
}
