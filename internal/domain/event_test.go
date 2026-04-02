package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/myorg/worker-service/internal/domain"
)

func TestEventType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		et    domain.EventType
		valid bool
	}{
		{"CREATE_DID", domain.EventTypeCreateDID, true},
		{"ADD_KEY", domain.EventTypeAddKey, true},
		{"ADD_CLAIM", domain.EventTypeAddClaim, true},
		{"STORE_DATA", domain.EventTypeStoreData, true},
		{"unknown", domain.EventType("UNKNOWN"), false},
		{"empty", domain.EventType(""), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.et.IsValid(); got != tc.valid {
				t.Errorf("EventType(%q).IsValid() = %v, want %v", tc.et, got, tc.valid)
			}
		})
	}
}

func TestEvent_Validate(t *testing.T) {
	validPayload := json.RawMessage(`{"did":"did:example:123"}`)

	tests := []struct {
		name    string
		event   domain.Event
		wantErr bool
	}{
		{
			name: "valid event",
			event: domain.Event{
				ID:        "evt-001",
				Type:      domain.EventTypeCreateDID,
				Payload:   validPayload,
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing id",
			event: domain.Event{
				Type:    domain.EventTypeCreateDID,
				Payload: validPayload,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			event: domain.Event{
				ID:      "evt-002",
				Type:    domain.EventType("BADTYPE"),
				Payload: validPayload,
			},
			wantErr: true,
		},
		{
			name: "nil payload",
			event: domain.Event{
				ID:   "evt-003",
				Type: domain.EventTypeAddKey,
			},
			wantErr: true,
		},
		{
			name: "empty payload",
			event: domain.Event{
				ID:      "evt-004",
				Type:    domain.EventTypeAddClaim,
				Payload: json.RawMessage(""),
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.event.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Event.Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
