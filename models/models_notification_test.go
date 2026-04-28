package models_test

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/jeroenmol/wero-api/models"
)

// --- NotificationRequest dispatch ---

func TestNotificationRequest_DispatchesToCorrectType(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantName string
		wantType any
	}{
		{
			name: "CaptureStatusChanged",
			body: `{
				"name": "CaptureStatusChanged",
				"id": "c0204631-b190-425d-bfa6-aa053b88f382",
				"rootAggregateId": "d8bcb97b-3524-4eb1-bde3-419fbe0a9e2e",
				"captureId": "f2f84a80-2c26-458e-bcfc-f2b5a92228c9",
				"status": {"value": "Settled","at": "2023-05-03T10:38:30Z","originatorPspReference": "1234","settlementMethod": "SCTInst"}
			}`,
			wantName: models.EventNameCaptureStatusChanged,
			wantType: &models.CaptureStatusChanged{},
		},
		{
			name: "RefundStatusChanged",
			body: `{
				"name": "RefundStatusChanged",
				"id": "c0204631-b190-425d-bfa6-aa053b88f382",
				"refundId": "d8bcb97b-3524-4eb1-bde3-419fbe0a9e2e",
				"amount": {"euroCents": 4900},
				"status": {"value": "Settled","at": "2023-05-03T10:38:30Z"}
			}`,
			wantName: models.EventNameRefundStatusChanged,
			wantType: &models.RefundStatusChanged{},
		},
		{
			name: "ChargebackStatusChanged",
			body: `{
				"name": "ChargebackStatusChanged",
				"id": "c0204631-b190-425d-bfa6-aa053b88f382",
				"disputeId": "d8bcb97b-3524-4eb1-bde3-419fbe0a9e2e",
				"status": "Settled",
				"remittanceInfo": "info",
				"amount": {"euroCents": 4900},
				"at": "2023-05-03T10:38:30Z"
			}`,
			wantName: models.EventNameChargebackStatusChanged,
			wantType: &models.ChargebackStatusChanged{},
		},
		{
			name: "MoneyTransferStatusChanged",
			body: `{
				"name": "MoneyTransferStatusChanged",
				"id": "4ff0cddc-75f2-4d4c-9126-e47134386480",
				"resourceId": "f2f84a80-2c26-458e-bcfc-f2b5a92228c9",
				"status": {"value": "Accepted","at": "2023-05-03T10:38:30Z"}
			}`,
			wantName: models.EventNameMoneyTransferStatusChanged,
			wantType: &models.MoneyTransferStatusChanged{},
		},
		{
			name: "AccountAuthorizedPersonAccessRevoked",
			body: `{
				"name": "AccountAuthorizedPersonAccessRevoked",
				"id": "4ff0cddc-75f2-4d4c-9126-e47134386480",
				"externalConsumerId": "consumer-1",
				"externalPaymentSourceId": "source-1"
			}`,
			wantName: models.EventNameAccountAuthorizedPersonAccessRevoked,
			wantType: &models.AccountAuthorizedPersonAccessRevoked{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var req models.NotificationRequest
			if err := json.Unmarshal([]byte(tc.body), &req); err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}
			if req.Event == nil {
				t.Fatal("Event is nil after unmarshal")
			}
			if req.Event.EventName() != tc.wantName {
				t.Errorf("EventName() = %q, want %q", req.Event.EventName(), tc.wantName)
			}
			switch req.Event.(type) {
			case *models.CaptureStatusChanged:
				if _, ok := tc.wantType.(*models.CaptureStatusChanged); !ok {
					t.Errorf("wrong concrete type: got CaptureStatusChanged")
				}
			case *models.RefundStatusChanged:
				if _, ok := tc.wantType.(*models.RefundStatusChanged); !ok {
					t.Errorf("wrong concrete type: got RefundStatusChanged")
				}
			case *models.ChargebackStatusChanged:
				if _, ok := tc.wantType.(*models.ChargebackStatusChanged); !ok {
					t.Errorf("wrong concrete type: got ChargebackStatusChanged")
				}
			case *models.MoneyTransferStatusChanged:
				if _, ok := tc.wantType.(*models.MoneyTransferStatusChanged); !ok {
					t.Errorf("wrong concrete type: got MoneyTransferStatusChanged")
				}
			case *models.AccountAuthorizedPersonAccessRevoked:
				if _, ok := tc.wantType.(*models.AccountAuthorizedPersonAccessRevoked); !ok {
					t.Errorf("wrong concrete type: got AccountAuthorizedPersonAccessRevoked")
				}
			default:
				t.Errorf("unexpected concrete type: %T", req.Event)
			}
		})
	}
}

func TestNotificationRequest_RejectsUnknownEventName(t *testing.T) {
	body := `{"name": "UnknownEvent", "id": "c0204631-b190-425d-bfa6-aa053b88f382"}`
	var req models.NotificationRequest
	if err := json.Unmarshal([]byte(body), &req); err == nil {
		t.Fatal("expected error for unknown event name, got nil")
	}
}

func TestNotificationRequest_RejectsMissingName(t *testing.T) {
	body := `{"id": "c0204631-b190-425d-bfa6-aa053b88f382"}`
	var req models.NotificationRequest
	if err := json.Unmarshal([]byte(body), &req); err == nil {
		t.Fatal("expected error for missing name field, got nil")
	}
}

// --- Capture status sub-union dispatch ---

func TestCaptureStatusEnvelope_DispatchesToSettled(t *testing.T) {
	body := `{"value": "Settled","at": "2023-05-03T10:38:30Z","originatorPspReference": "REF-123","settlementMethod": "SCTInst"}`
	var env models.CaptureStatusEnvelope
	if err := json.Unmarshal([]byte(body), &env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	settled, ok := env.Status.(*models.CaptureStatusSettled)
	if !ok {
		t.Fatalf("expected *CaptureStatusSettled, got %T", env.Status)
	}
	if settled.OriginatorPspRef != "REF-123" {
		t.Errorf("OriginatorPspRef = %q, want %q", settled.OriginatorPspRef, "REF-123")
	}
	if settled.SettlementMethod != "SCTInst" {
		t.Errorf("SettlementMethod = %q, want %q", settled.SettlementMethod, "SCTInst")
	}
}

func TestCaptureStatusEnvelope_DispatchesToPending(t *testing.T) {
	body := `{"value": "Pending","at": "2023-05-03T10:38:30Z","reasonCode": "AC01","reason": "Invalid account","reinitiationRequired": true}`
	var env models.CaptureStatusEnvelope
	if err := json.Unmarshal([]byte(body), &env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pending, ok := env.Status.(*models.CaptureStatusPending)
	if !ok {
		t.Fatalf("expected *CaptureStatusPending, got %T", env.Status)
	}
	if pending.ReasonCode != "AC01" {
		t.Errorf("ReasonCode = %q, want %q", pending.ReasonCode, "AC01")
	}
	if !pending.ReinitiationRequired {
		t.Error("ReinitiationRequired = false, want true")
	}
}

func TestCaptureStatusEnvelope_RejectsUnknownValue(t *testing.T) {
	body := `{"value": "Unknown","at": "2023-05-03T10:38:30Z"}`
	var env models.CaptureStatusEnvelope
	if err := json.Unmarshal([]byte(body), &env); err == nil {
		t.Fatal("expected error for unknown status value, got nil")
	}
}

func TestMoneyTransferStatusEnvelope_DispatchesToRejected(t *testing.T) {
	body := `{"value": "Rejected","at": "2023-05-03T10:38:30Z","rejectionReason": "InsufficientFunds","rejectionMessage": "Not enough balance"}`
	var env models.MoneyTransferStatusEnvelope
	if err := json.Unmarshal([]byte(body), &env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rejected, ok := env.Status.(*models.MTStatusRejected)
	if !ok {
		t.Fatalf("expected *MTStatusRejected, got %T", env.Status)
	}
	if rejected.RejectionReason != "InsufficientFunds" {
		t.Errorf("RejectionReason = %q, want %q", rejected.RejectionReason, "InsufficientFunds")
	}
}

// --- Validation ---

func TestCaptureStatusChanged_ValidationPassesForValidInput(t *testing.T) {
	d := &models.CaptureStatusChanged{
		Name:            models.EventNameCaptureStatusChanged,
		ID:              "c0204631-b190-425d-bfa6-aa053b88f382",
		RootAggregateID: "d8bcb97b-3524-4eb1-bde3-419fbe0a9e2e",
		CaptureID:       "f2f84a80-2c26-458e-bcfc-f2b5a92228c9",
		Status:          &models.CaptureStatusEnvelope{},
	}
	v := validator.New()
	if err := v.Struct(d); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestCaptureStatusChanged_ValidationFailsForMissingCaptureID(t *testing.T) {
	d := &models.CaptureStatusChanged{
		Name:            models.EventNameCaptureStatusChanged,
		ID:              "c0204631-b190-425d-bfa6-aa053b88f382",
		RootAggregateID: "d8bcb97b-3524-4eb1-bde3-419fbe0a9e2e",
		// CaptureID omitted
	}
	v := validator.New()
	if err := v.Struct(d); err == nil {
		t.Error("expected validation error for missing CaptureID, got nil")
	}
}

func TestCaptureStatusChanged_ValidationFailsForInvalidUUID(t *testing.T) {
	d := &models.CaptureStatusChanged{
		Name:            models.EventNameCaptureStatusChanged,
		ID:              "not-a-uuid",
		RootAggregateID: "d8bcb97b-3524-4eb1-bde3-419fbe0a9e2e",
		CaptureID:       "f2f84a80-2c26-458e-bcfc-f2b5a92228c9",
		Status:          &models.CaptureStatusEnvelope{},
	}
	v := validator.New()
	if err := v.Struct(d); err == nil {
		t.Error("expected validation error for invalid UUID in ID, got nil")
	}
}

func TestMTStatusFailed_ValidationFailsForFailureCodeNotExactlyFourChars(t *testing.T) {
	tests := []struct {
		name        string
		failureCode string
	}{
		{"too short", "ABC"},
		{"too long", "ABCDE"},
		{"empty", ""},
	}
	v := validator.New()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &models.MTStatusFailed{
				Value:       "Failed",
				FailureCode: tc.failureCode,
			}
			if err := v.Struct(d); err == nil {
				t.Errorf("expected validation error for failureCode %q, got nil", tc.failureCode)
			}
		})
	}
}

func TestMTStatusSettled_ValidationFailsForOriginatorRefExceedingMaxLength(t *testing.T) {
	d := &models.MTStatusSettled{
		Value:            "Settled",
		OriginatorPspRef: "THIS-IS-A-REFERENCE-THAT-IS-WAY-TOO-LONG-FOR-THE-SPEC",
	}
	v := validator.New()
	if err := v.Struct(d); err == nil {
		t.Error("expected validation error for originatorPspReference exceeding max=35, got nil")
	}
}
