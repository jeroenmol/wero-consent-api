package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Event name constants matching the "name" discriminator values in the OpenAPI spec.
const (
	EventNameCaptureStatusChanged                 = "CaptureStatusChanged"
	EventNameRefundStatusChanged                  = "RefundStatusChanged"
	EventNameChargebackStatusChanged              = "ChargebackStatusChanged"
	EventNameMoneyTransferStatusChanged           = "MoneyTransferStatusChanged"
	EventNameAccountAuthorizedPersonAccessRevoked = "AccountAuthorizedPersonAccessRevoked"
)

// SettlementMethod represents the settlement method used by a consumer PSP.
type SettlementMethod string

const (
	SettlementMethodSCTInst             SettlementMethod = "SCTInst"
	SettlementMethodOnUs                SettlementMethod = "OnUs"
	SettlementMethodFinancialGuarantees SettlementMethod = "FinancialGuarantees"
)

// ConsumerPspEvent is the sealed union of all incoming notification event types.
// The discriminator field is "name".
type ConsumerPspEvent interface {
	EventName() string
	sealed()
}

// NotificationRequest is the top-level request body for POST /api/notifications.
// UnmarshalJSON peeks at "name" and dispatches to the right concrete event type.
type NotificationRequest struct {
	Event ConsumerPspEvent
}

var eventRegistry = map[string]func() ConsumerPspEvent{
	EventNameCaptureStatusChanged:                 func() ConsumerPspEvent { return &CaptureStatusChanged{} },
	EventNameRefundStatusChanged:                  func() ConsumerPspEvent { return &RefundStatusChanged{} },
	EventNameChargebackStatusChanged:              func() ConsumerPspEvent { return &ChargebackStatusChanged{} },
	EventNameMoneyTransferStatusChanged:           func() ConsumerPspEvent { return &MoneyTransferStatusChanged{} },
	EventNameAccountAuthorizedPersonAccessRevoked: func() ConsumerPspEvent { return &AccountAuthorizedPersonAccessRevoked{} },
}

func (e *NotificationRequest) UnmarshalJSON(data []byte) error {
	var peek struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	factory, ok := eventRegistry[peek.Name]
	if !ok {
		return fmt.Errorf("unknown event name: %q", peek.Name)
	}
	v := factory()
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	e.Event = v
	return nil
}

// --- Capture status union ---

// CaptureStatus is the sealed union for capture status variants (discriminator: "value").
type CaptureStatus interface {
	StatusValue() string
	captureStatusSealed()
}

// CaptureStatusEnvelope decodes the nested capture "status" field.
type CaptureStatusEnvelope struct {
	Status CaptureStatus
}

var captureStatusRegistry = map[string]func() CaptureStatus{
	"Failed":   func() CaptureStatus { return &CaptureStatusFailed{} },
	"Pending":  func() CaptureStatus { return &CaptureStatusPending{} },
	"Rejected": func() CaptureStatus { return &CaptureStatusRejected{} },
	"Settled":  func() CaptureStatus { return &CaptureStatusSettled{} },
}

func (e *CaptureStatusEnvelope) UnmarshalJSON(data []byte) error {
	var peek struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	factory, ok := captureStatusRegistry[peek.Value]
	if !ok {
		return fmt.Errorf("unknown capture status: %q", peek.Value)
	}
	v := factory()
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	e.Status = v
	return nil
}

type CaptureStatusFailed struct {
	Value           string    `json:"value" validate:"required,eq=Failed"`
	At              time.Time `json:"at" validate:"required"`
	TechnicalReason string    `json:"technicalReason" validate:"required"`
}

func (d *CaptureStatusFailed) StatusValue() string  { return "Failed" }
func (d *CaptureStatusFailed) captureStatusSealed() {}

type CaptureStatusPending struct {
	Value                string    `json:"value" validate:"required,eq=Pending"`
	At                   time.Time `json:"at" validate:"required"`
	ReasonCode           string    `json:"reasonCode" validate:"required"`
	Reason               string    `json:"reason" validate:"required"`
	ReinitiationRequired bool      `json:"reinitiationRequired"`
}

func (d *CaptureStatusPending) StatusValue() string  { return "Pending" }
func (d *CaptureStatusPending) captureStatusSealed() {}

type CaptureStatusRejected struct {
	Value      string    `json:"value" validate:"required,eq=Rejected"`
	At         time.Time `json:"at" validate:"required"`
	ReasonCode string    `json:"reasonCode" validate:"required"`
	Reason     string    `json:"reason" validate:"required"`
}

func (d *CaptureStatusRejected) StatusValue() string  { return "Rejected" }
func (d *CaptureStatusRejected) captureStatusSealed() {}

type CaptureStatusSettled struct {
	Value            string    `json:"value" validate:"required,eq=Settled"`
	At               time.Time `json:"at" validate:"required"`
	OriginatorPspRef string    `json:"originatorPspReference" validate:"required"`
	SettlementMethod SettlementMethod `json:"settlementMethod" validate:"required,oneof=SCTInst OnUs FinancialGuarantees"`
}

func (d *CaptureStatusSettled) StatusValue() string  { return "Settled" }
func (d *CaptureStatusSettled) captureStatusSealed() {}

// CaptureStatusChanged is the event for CaptureStatusChanged events.
type CaptureStatusChanged struct {
	Name            string                 `json:"name" validate:"required,eq=CaptureStatusChanged"`
	ID              string                 `json:"id" validate:"required,uuid"`
	RootAggregateID string                 `json:"rootAggregateId" validate:"required,uuid"`
	CaptureID       string                 `json:"captureId" validate:"required,uuid"`
	Status          *CaptureStatusEnvelope `json:"status" validate:"required"`
}

func (d *CaptureStatusChanged) EventName() string { return EventNameCaptureStatusChanged }
func (d *CaptureStatusChanged) sealed()           {}

// --- Refund status union ---

// RefundStatus is the sealed union for refund status variants (discriminator: "value").
type RefundStatus interface {
	StatusValue() string
	refundStatusSealed()
}

// RefundStatusEnvelope decodes the nested refund "status" field.
type RefundStatusEnvelope struct {
	Status RefundStatus
}

var refundStatusRegistry = map[string]func() RefundStatus{
	"Failed":   func() RefundStatus { return &RefundStatusFailed{} },
	"Rejected": func() RefundStatus { return &RefundStatusRejected{} },
	"Settled":  func() RefundStatus { return &RefundStatusSettled{} },
}

func (e *RefundStatusEnvelope) UnmarshalJSON(data []byte) error {
	var peek struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	factory, ok := refundStatusRegistry[peek.Value]
	if !ok {
		return fmt.Errorf("unknown refund status: %q", peek.Value)
	}
	v := factory()
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	e.Status = v
	return nil
}

type RefundStatusFailed struct {
	Value           string    `json:"value" validate:"required,eq=Failed"`
	At              time.Time `json:"at" validate:"required"`
	TechnicalReason string    `json:"technicalReason" validate:"required"`
}

func (d *RefundStatusFailed) StatusValue() string { return "Failed" }
func (d *RefundStatusFailed) refundStatusSealed() {}

type RefundStatusRejected struct {
	Value      string    `json:"value" validate:"required,eq=Rejected"`
	At         time.Time `json:"at" validate:"required"`
	ReasonCode string    `json:"reasonCode" validate:"required"`
	Reason     string    `json:"reason" validate:"required"`
}

func (d *RefundStatusRejected) StatusValue() string { return "Rejected" }
func (d *RefundStatusRejected) refundStatusSealed() {}

type RefundStatusSettled struct {
	Value string    `json:"value" validate:"required,eq=Settled"`
	At    time.Time `json:"at" validate:"required"`
}

func (d *RefundStatusSettled) StatusValue() string { return "Settled" }
func (d *RefundStatusSettled) refundStatusSealed() {}

// RefundStatusChanged is the event for RefundStatusChanged events.
type RefundStatusChanged struct {
	Name           string                `json:"name" validate:"required,eq=RefundStatusChanged"`
	ID             string                `json:"id" validate:"required,uuid"`
	RefundID       string                `json:"refundId" validate:"required,uuid"`
	Status         *RefundStatusEnvelope `json:"status" validate:"required"`
	RemittanceInfo string                `json:"remittanceInfo"`
	Amount         Money                 `json:"amount" validate:"required"`
}

func (d *RefundStatusChanged) EventName() string { return EventNameRefundStatusChanged }
func (d *RefundStatusChanged) sealed()           {}

// ChargebackStatusChanged is the event for ChargebackStatusChanged events.
// Status is always "Settled"; Stage defaults to "ChargebackRepresented".
type ChargebackStatusChanged struct {
	Name           string    `json:"name" validate:"required,eq=ChargebackStatusChanged"`
	ID             string    `json:"id" validate:"required,uuid"`
	DisputeID      string    `json:"disputeId" validate:"required,uuid"`
	Status         string    `json:"status" validate:"omitempty,eq=Settled"`
	RemittanceInfo string    `json:"remittanceInfo" validate:"required"`
	Amount         Money     `json:"amount" validate:"required"`
	At             time.Time `json:"at" validate:"required"`
	Stage          string    `json:"stage" validate:"omitempty,eq=ChargebackRepresented"`
}

func (d *ChargebackStatusChanged) EventName() string { return EventNameChargebackStatusChanged }
func (d *ChargebackStatusChanged) sealed()           {}

// --- Money transfer status union ---

// MoneyTransferStatus is the sealed union for money transfer status variants (discriminator: "value").
type MoneyTransferStatus interface {
	StatusValue() string
	mtStatusSealed()
}

// MoneyTransferStatusEnvelope decodes the nested money transfer "status" field.
type MoneyTransferStatusEnvelope struct {
	Status MoneyTransferStatus
}

var mtStatusRegistry = map[string]func() MoneyTransferStatus{
	"Accepted": func() MoneyTransferStatus { return &MTStatusAccepted{} },
	"Failed":   func() MoneyTransferStatus { return &MTStatusFailed{} },
	"Settled":  func() MoneyTransferStatus { return &MTStatusSettled{} },
	"Rejected": func() MoneyTransferStatus { return &MTStatusRejected{} },
}

func (e *MoneyTransferStatusEnvelope) UnmarshalJSON(data []byte) error {
	var peek struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return err
	}
	factory, ok := mtStatusRegistry[peek.Value]
	if !ok {
		return fmt.Errorf("unknown money transfer status: %q", peek.Value)
	}
	v := factory()
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	e.Status = v
	return nil
}

type MTStatusAccepted struct {
	Value string    `json:"value" validate:"required,eq=Accepted"`
	At    time.Time `json:"at" validate:"required"`
}

func (d *MTStatusAccepted) StatusValue() string { return "Accepted" }
func (d *MTStatusAccepted) mtStatusSealed()     {}

type MTStatusFailed struct {
	Value         string    `json:"value" validate:"required,eq=Failed"`
	At            time.Time `json:"at" validate:"required"`
	FailureCode   string    `json:"failureCode" validate:"required,len=4"`
	FailureReason string    `json:"failureReason" validate:"omitempty,min=1,max=105"`
}

func (d *MTStatusFailed) StatusValue() string { return "Failed" }
func (d *MTStatusFailed) mtStatusSealed()     {}

type MTStatusSettled struct {
	Value            string    `json:"value" validate:"required,eq=Settled"`
	At               time.Time `json:"at" validate:"required"`
	OriginatorPspRef string    `json:"originatorPspReference" validate:"required,max=35"`
	SettlementMethod SettlementMethod `json:"settlementMethod" validate:"omitempty,oneof=SCTInst OnUs"`
}

func (d *MTStatusSettled) StatusValue() string { return "Settled" }
func (d *MTStatusSettled) mtStatusSealed()     {}

type MTStatusRejected struct {
	Value            string    `json:"value" validate:"required,eq=Rejected"`
	At               time.Time `json:"at" validate:"required"`
	RejectionReason  string    `json:"rejectionReason" validate:"required"`
	RejectionMessage string    `json:"rejectionMessage" validate:"omitempty,min=1,max=140"`
}

func (d *MTStatusRejected) StatusValue() string { return "Rejected" }
func (d *MTStatusRejected) mtStatusSealed()     {}

// MoneyTransferStatusChanged is the event for MoneyTransferStatusChanged events.
type MoneyTransferStatusChanged struct {
	Name       string                       `json:"name" validate:"required,eq=MoneyTransferStatusChanged"`
	ID         string                       `json:"id" validate:"required,uuid"`
	ResourceID string                       `json:"resourceId" validate:"required,uuid"`
	Status     *MoneyTransferStatusEnvelope `json:"status" validate:"required"`
}

func (d *MoneyTransferStatusChanged) EventName() string {
	return EventNameMoneyTransferStatusChanged
}
func (d *MoneyTransferStatusChanged) sealed() {}

// AccountAuthorizedPersonAccessRevoked is the event for AccountAuthorizedPersonAccessRevoked events.
type AccountAuthorizedPersonAccessRevoked struct {
	Name                    string `json:"name" validate:"required,eq=AccountAuthorizedPersonAccessRevoked"`
	ID                      string `json:"id" validate:"required,uuid"`
	ExternalConsumerID      string `json:"externalConsumerId" validate:"required"`
	ExternalPaymentSourceID string `json:"externalPaymentSourceId" validate:"required"`
}

func (d *AccountAuthorizedPersonAccessRevoked) EventName() string {
	return EventNameAccountAuthorizedPersonAccessRevoked
}
func (d *AccountAuthorizedPersonAccessRevoked) sealed() {}
