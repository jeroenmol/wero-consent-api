package weroapi

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/jeroenmol/wero-api/epierrors"
	"github.com/jeroenmol/wero-api/models"
)

// ConsentsService provides access to the /api/consents endpoints.
type ConsentsService struct {
	client *resty.Client
}

// New creates a ConsentsService backed by a resty client pointed at baseURL.
// The client uses Bearer token authentication.
func New(baseURL, bearerToken string) *ConsentsService {
	client := resty.New().
		SetBaseURL(baseURL).
		SetAuthToken(bearerToken).
		SetHeader("Accept", "application/hal+json; charset=UTF-8")
	return &ConsentsService{client: client}
}

// RequestConsent creates a new consent via POST /api/consents.
func (s *ConsentsService) RequestConsent(ctx context.Context, req models.RequestConsentRequest) (*models.ConsentResource, error) {
	var result models.ConsentResource
	var errResp epierrors.ConstraintViolationProblem
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&result).
		SetError(&errResp).
		Post("/api/consents")
	if err != nil {
		return nil, fmt.Errorf("request consent: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("request consent: %w", &errResp)
	}
	return &result, nil
}

// ListConsents retrieves a paginated list of consents via GET /api/consents.
func (s *ConsentsService) ListConsents(ctx context.Context, params models.ListConsentsParams) (*models.ConsentCollectionResource, error) {
	var result models.ConsentCollectionResource
	var errResp epierrors.Problem
	resp, err := s.client.R().
		SetContext(ctx).
		SetQueryParams(listConsentsQueryParams(params)).
		SetResult(&result).
		SetError(&errResp).
		Get("/api/consents")
	if err != nil {
		return nil, fmt.Errorf("list consents: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("list consents: %w", &errResp)
	}
	return &result, nil
}

// GetConsent retrieves a single consent by ID via GET /api/consents/{consentId}.
func (s *ConsentsService) GetConsent(ctx context.Context, consentID string) (*models.ConsentResource, error) {
	var result models.ConsentResource
	var errResp epierrors.Problem
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("consentId", consentID).
		SetResult(&result).
		SetError(&errResp).
		Get("/api/consents/{consentId}")
	if err != nil {
		return nil, fmt.Errorf("get consent: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get consent: %w", &errResp)
	}
	return &result, nil
}

// GetConsentDetails retrieves enriched details for a consent via GET /api/consents/{consentId}/details.
func (s *ConsentsService) GetConsentDetails(ctx context.Context, consentID string) (*models.ConsentDetailsResource, error) {
	var result models.ConsentDetailsResource
	var errResp epierrors.Problem
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("consentId", consentID).
		SetResult(&result).
		SetError(&errResp).
		Get("/api/consents/{consentId}/details")
	if err != nil {
		return nil, fmt.Errorf("get consent details: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get consent details: %w", &errResp)
	}
	return &result, nil
}

// UpdateConsentDetails sets the EasyCheckout profile for a consent via PUT /api/consents/{consentId}/details.
//
// Deprecated: use the standard consent flow instead.
func (s *ConsentsService) UpdateConsentDetails(ctx context.Context, consentID string, req models.SetEasyCheckoutProfileRequest) (*models.ConsentDetailsResource, error) {
	var result models.ConsentDetailsResource
	var errResp epierrors.ConstraintViolationProblem
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("consentId", consentID).
		SetBody(req).
		SetResult(&result).
		SetError(&errResp).
		Put("/api/consents/{consentId}/details")
	if err != nil {
		return nil, fmt.Errorf("update consent details: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("update consent details: %w", &errResp)
	}
	return &result, nil
}

// GetPaymentRules retrieves the payment rules for a consent via GET /api/consents/{consentId}/payment-rules.
func (s *ConsentsService) GetPaymentRules(ctx context.Context, consentID string) (*models.ConsentPaymentRulesResource, error) {
	var result models.ConsentPaymentRulesResource
	var errResp epierrors.Problem
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("consentId", consentID).
		SetResult(&result).
		SetError(&errResp).
		Get("/api/consents/{consentId}/payment-rules")
	if err != nil {
		return nil, fmt.Errorf("get payment rules: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get payment rules: %w", &errResp)
	}
	return &result, nil
}

// GetCaptureRules retrieves the capture rules for a consent via GET /api/consents/{consentId}/capture-rules.
func (s *ConsentsService) GetCaptureRules(ctx context.Context, consentID string) (*models.CaptureRulesResource, error) {
	var result models.CaptureRulesResource
	var errResp epierrors.Problem
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("consentId", consentID).
		SetResult(&result).
		SetError(&errResp).
		Get("/api/consents/{consentId}/capture-rules")
	if err != nil {
		return nil, fmt.Errorf("get capture rules: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get capture rules: %w", &errResp)
	}
	return &result, nil
}

// GetPaymentPlan retrieves the payment plan for a consent via GET /api/consents/{consentId}/payment-plan.
func (s *ConsentsService) GetPaymentPlan(ctx context.Context, consentID string) (*models.PaymentPlanResource, error) {
	var result models.PaymentPlanResource
	var errResp epierrors.Problem
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("consentId", consentID).
		SetResult(&result).
		SetError(&errResp).
		Get("/api/consents/{consentId}/payment-plan")
	if err != nil {
		return nil, fmt.Errorf("get payment plan: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get payment plan: %w", &errResp)
	}
	return &result, nil
}

// GetCurrentToken retrieves the current consent token via GET /api/consents/{consentId}/current-token.
func (s *ConsentsService) GetCurrentToken(ctx context.Context, consentID string) (*models.ConsentTokenResource, error) {
	var result models.ConsentTokenResource
	var errResp epierrors.Problem
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("consentId", consentID).
		SetResult(&result).
		SetError(&errResp).
		Get("/api/consents/{consentId}/current-token")
	if err != nil {
		return nil, fmt.Errorf("get current token: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get current token: %w", &errResp)
	}
	return &result, nil
}

// AddCurrentToken uploads a new consent token, invalidating the previous one [Wallet].
// Returns the stored token resource (200) or nil when the server responds 204.
func (s *ConsentsService) AddCurrentToken(ctx context.Context, consentID string, req models.AddConsentTokenRequest) (*models.ConsentTokenResource, error) {
	var result models.ConsentTokenResource
	var errResp epierrors.ConstraintViolationProblem
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("consentId", consentID).
		SetBody(req).
		SetResult(&result).
		SetError(&errResp).
		Put("/api/consents/{consentId}/current-token")
	if err != nil {
		return nil, fmt.Errorf("add consent token: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("add consent token: %w", &errResp)
	}
	if resp.StatusCode() == 204 {
		return nil, nil
	}
	return &result, nil
}

// GetShopFront retrieves the shop front resource for a consent via GET /api/consents/{consentId}/shop-front.
func (s *ConsentsService) GetShopFront(ctx context.Context, consentID string) (*models.ShopFrontResource, error) {
	var result models.ShopFrontResource
	var errResp epierrors.Problem
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("consentId", consentID).
		SetResult(&result).
		SetError(&errResp).
		Get("/api/consents/{consentId}/shop-front")
	if err != nil {
		return nil, fmt.Errorf("get shop front: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get shop front: %w", &errResp)
	}
	return &result, nil
}

// DenyConsent denies a consent via PUT /api/consents/{consentId}/denial.
// req is optional; pass nil if no telemetry is available.
func (s *ConsentsService) DenyConsent(ctx context.Context, consentID string, req *models.DenyConsentRequest) error {
	var errResp epierrors.Problem
	r := s.client.R().
		SetContext(ctx).
		SetPathParam("consentId", consentID).
		SetError(&errResp)
	if req != nil {
		r = r.SetBody(req)
	}
	resp, err := r.Put("/api/consents/{consentId}/denial")
	if err != nil {
		return fmt.Errorf("deny consent: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("deny consent: %w", &errResp)
	}
	return nil
}

// RevokeConsent revokes a consent via PUT /api/consents/{consentId}/revocation.
func (s *ConsentsService) RevokeConsent(ctx context.Context, consentID string) error {
	var errResp epierrors.Problem
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("consentId", consentID).
		SetError(&errResp).
		Put("/api/consents/{consentId}/revocation")
	if err != nil {
		return fmt.Errorf("revoke consent: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("revoke consent: %w", &errResp)
	}
	return nil
}

// CreateOfflineQRCodeConsent creates an offline QR code consent via POST /api/consents/offline-qr-code/create.
//
// Deprecated: use the standard consent flow instead.
func (s *ConsentsService) CreateOfflineQRCodeConsent(ctx context.Context, req models.OfflineQRCodeRequest) (*models.ConsentResource, error) {
	var result models.ConsentResource
	var errResp epierrors.ConstraintViolationProblem
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&result).
		SetError(&errResp).
		Post("/api/consents/offline-qr-code/create")
	if err != nil {
		return nil, fmt.Errorf("create offline qr code consent: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("create offline qr code consent: %w", &errResp)
	}
	return &result, nil
}

func listConsentsQueryParams(p models.ListConsentsParams) map[string]string {
	m := make(map[string]string)
	if p.FindBy != "" {
		m["findBy"] = p.FindBy
	}
	if p.OrderReferenceID != "" {
		m["orderReferenceId"] = p.OrderReferenceID
	}
	if p.ShortID != "" {
		m["shortId"] = p.ShortID
	}
	if p.PaymentMeansID != "" {
		m["paymentMeansId"] = p.PaymentMeansID
	}
	if p.PaymentPlan != "" {
		m["paymentPlan"] = p.PaymentPlan
	}
	return m
}
