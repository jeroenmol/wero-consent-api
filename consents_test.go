package weroapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	weroapi "github.com/jeroenmol/wero-api"
	"github.com/jeroenmol/wero-api/epierrors"
	"github.com/jeroenmol/wero-api/models"
)

// newTestService creates a ConsentsService backed by a fake HTTP server.
// The server is closed automatically when the test ends.
func newTestService(t *testing.T, handler http.HandlerFunc) *weroapi.ConsentsService {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return weroapi.New(srv.URL, "test-token")
}

// mustEncode marshals v to JSON and writes it to w; panics on error.
func mustEncode(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Errorf("failed to encode response: %v", err)
	}
}

// assertAuthHeader verifies the request carries the expected Bearer token.
func assertAuthHeader(t *testing.T, r *http.Request) {
	t.Helper()
	auth := r.Header.Get("Authorization")
	if auth != "Bearer test-token" {
		t.Errorf("Authorization = %q, want %q", auth, "Bearer test-token")
	}
}

// -------------------------------------------------------------------------
// 1. RequestConsent
// -------------------------------------------------------------------------

func TestRequestConsent_Success(t *testing.T) {
	want := models.ConsentResource{
		ID: "5f7c8ca7-581a-43cb-ba27-08409c1cfb4b",
		Status: models.ConsentStatus{
			Value: "Requested",
			At:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Meta: models.ConsentMeta{
			CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		OrderReferenceID: "ord-123",
	}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/consents" {
			t.Errorf("path = %q, want /api/consents", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		mustEncode(t, w, want)
	})

	got, err := svc.RequestConsent(context.Background(), models.RequestConsentRequest{
		PaymentPlan: models.RequestedPaymentPlan{
			Type:   "SingleImmediate",
			Amount: &models.Money{EuroCents: 1999},
		},
		OrderReferenceID: "ord-123",
		Acceptor: models.RequestedAcceptor{
			ID:          "12345672",
			Name:        "Test Merchant",
			MCC:         "1234",
			CountryCode: "BE",
			Shop:        models.AcceptorShop{ID: "s1", Name: "Shop"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != want.ID {
		t.Errorf("ID = %q, want %q", got.ID, want.ID)
	}
	if got.OrderReferenceID != want.OrderReferenceID {
		t.Errorf("OrderReferenceID = %q, want %q", got.OrderReferenceID, want.OrderReferenceID)
	}
	if got.Status.Value != want.Status.Value {
		t.Errorf("Status.Value = %q, want %q", got.Status.Value, want.Status.Value)
	}
}

func TestRequestConsent_ConstraintViolation(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		mustEncode(t, w, epierrors.ConstraintViolationProblem{
			Type:   "/problem/constraint-violation",
			Title:  "Constraint Violation",
			Status: 400,
			Violations: []epierrors.ConstraintViolation{
				{Field: ".paymentPlan.amount", Message: "must not be null"},
			},
		})
	})

	_, err := svc.RequestConsent(context.Background(), models.RequestConsentRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "request consent") {
		t.Errorf("error %q missing 'request consent' prefix", err.Error())
	}
}

func TestRequestConsent_NetworkError(t *testing.T) {
	// Close the server immediately so all connections fail.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	_, err := svc.RequestConsent(context.Background(), models.RequestConsentRequest{})
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// 2. ListConsents
// -------------------------------------------------------------------------

func TestListConsents_Success(t *testing.T) {
	want := models.ConsentCollectionResource{}
	want.Embedded.Consents = []models.ConsentResource{
		{ID: "consent-1", Status: models.ConsentStatus{Value: "Consented"}},
		{ID: "consent-2", Status: models.ConsentStatus{Value: "Requested"}},
	}
	want.Page = &models.Page{TotalElements: 2, TotalPages: 1, Size: 2}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/api/consents" {
			t.Errorf("path = %q, want /api/consents", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, want)
	})

	got, err := svc.ListConsents(context.Background(), models.ListConsentsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Embedded.Consents) != 2 {
		t.Errorf("len(consents) = %d, want 2", len(got.Embedded.Consents))
	}
	if got.Embedded.Consents[0].ID != "consent-1" {
		t.Errorf("consents[0].ID = %q, want %q", got.Embedded.Consents[0].ID, "consent-1")
	}
}

func TestListConsents_QueryParamsForwarded(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("findBy") != "orderReferenceId" {
			t.Errorf("findBy = %q, want %q", q.Get("findBy"), "orderReferenceId")
		}
		if q.Get("orderReferenceId") != "ord-abc" {
			t.Errorf("orderReferenceId = %q, want %q", q.Get("orderReferenceId"), "ord-abc")
		}
		if q.Get("paymentPlan") != "SingleImmediate" {
			t.Errorf("paymentPlan = %q, want %q", q.Get("paymentPlan"), "SingleImmediate")
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, models.ConsentCollectionResource{})
	})

	_, err := svc.ListConsents(context.Background(), models.ListConsentsParams{
		FindBy:           "orderReferenceId",
		OrderReferenceID: "ord-abc",
		PaymentPlan:      "SingleImmediate",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListConsents_ServerError(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		mustEncode(t, w, epierrors.Problem{
			Type:   "/problem/internal-server-error",
			Title:  "Internal Server Error",
			Status: 500,
		})
	})

	_, err := svc.ListConsents(context.Background(), models.ListConsentsParams{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListConsents_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	_, err := svc.ListConsents(context.Background(), models.ListConsentsParams{})
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// 3. GetConsent
// -------------------------------------------------------------------------

func TestGetConsent_Success(t *testing.T) {
	consentID := "abc-123"
	want := models.ConsentResource{
		ID:     consentID,
		Status: models.ConsentStatus{Value: "Consented"},
	}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if !strings.Contains(r.URL.Path, consentID) {
			t.Errorf("path %q does not contain consentID %q", r.URL.Path, consentID)
		}
		if r.URL.Path != "/api/consents/"+consentID {
			t.Errorf("path = %q, want /api/consents/%s", r.URL.Path, consentID)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, want)
	})

	got, err := svc.GetConsent(context.Background(), consentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != consentID {
		t.Errorf("ID = %q, want %q", got.ID, consentID)
	}
}

func TestGetConsent_NotFound(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusNotFound)
		mustEncode(t, w, epierrors.Problem{
			Type:   "/problem/not-found",
			Title:  "Not Found",
			Status: 404,
			Detail: "consent not found",
		})
	})

	_, err := svc.GetConsent(context.Background(), "nonexistent-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get consent") {
		t.Errorf("error %q missing 'get consent' prefix", err.Error())
	}
}

func TestGetConsent_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	_, err := svc.GetConsent(context.Background(), "some-id")
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// 4. GetConsentDetails
// -------------------------------------------------------------------------

func TestGetConsentDetails_Success(t *testing.T) {
	consentID := "detail-consent-id"
	want := models.ConsentDetailsResource{
		Description:  "Order details",
		OrderPageURL: "https://example.com/orders/123",
		Gift:         true,
		BillingAddress: &models.Address{
			RecipientName: "John Doe",
			Line1:         "Main Street 1",
			City:          "Brussels",
			PostalCode:    "1000",
			CountryCode:   "BE",
		},
	}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if !strings.Contains(r.URL.Path, consentID) {
			t.Errorf("path %q does not contain consentID %q", r.URL.Path, consentID)
		}
		expectedPath := "/api/consents/" + consentID + "/details"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, want)
	})

	got, err := svc.GetConsentDetails(context.Background(), consentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Description != want.Description {
		t.Errorf("Description = %q, want %q", got.Description, want.Description)
	}
	if got.BillingAddress == nil {
		t.Fatal("BillingAddress is nil")
	}
	if got.BillingAddress.City != want.BillingAddress.City {
		t.Errorf("BillingAddress.City = %q, want %q", got.BillingAddress.City, want.BillingAddress.City)
	}
}

func TestGetConsentDetails_ServerError(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		mustEncode(t, w, epierrors.Problem{Title: "Not Found", Status: 404})
	})

	_, err := svc.GetConsentDetails(context.Background(), "missing-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// -------------------------------------------------------------------------
// 5. UpdateConsentDetails
// -------------------------------------------------------------------------

func TestUpdateConsentDetails_Success(t *testing.T) {
	consentID := "update-detail-id"
	want := models.ConsentDetailsResource{
		Description: "Updated description",
		Contact:     &models.ConsentContact{EmailAddress: "user@example.com"},
	}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		expectedPath := "/api/consents/" + consentID + "/details"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, want)
	})

	got, err := svc.UpdateConsentDetails(context.Background(), consentID, models.SetEasyCheckoutProfileRequest{
		ShippingAddress: models.Address{
			RecipientName: "Jane Doe",
			Line1:         "123 Street",
			City:          "Ghent",
			PostalCode:    "9000",
			CountryCode:   "BE",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Description != want.Description {
		t.Errorf("Description = %q, want %q", got.Description, want.Description)
	}
}

func TestUpdateConsentDetails_ConstraintViolation(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		mustEncode(t, w, epierrors.ConstraintViolationProblem{
			Type:   "/problem/constraint-violation",
			Title:  "Constraint Violation",
			Status: 400,
			Violations: []epierrors.ConstraintViolation{
				{Field: ".shippingAddress.city", Message: "must not be blank"},
			},
		})
	})

	_, err := svc.UpdateConsentDetails(context.Background(), "some-id", models.SetEasyCheckoutProfileRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "update consent details") {
		t.Errorf("error %q missing 'update consent details' prefix", err.Error())
	}
}

func TestUpdateConsentDetails_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	_, err := svc.UpdateConsentDetails(context.Background(), "some-id", models.SetEasyCheckoutProfileRequest{})
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// 6. GetPaymentRules
// -------------------------------------------------------------------------

func TestGetPaymentRules_Success(t *testing.T) {
	consentID := "payment-rules-id"
	want := models.ConsentPaymentRulesResource{
		PaymentRules: []models.ConsentPaymentRule{
			{
				Type: "Window",
				Constraint: &models.ConsentPaymentRuleConstraint{
					Type:      "MaxAmount",
					MaxAmount: models.Money{EuroCents: 5000},
				},
			},
		},
	}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		expectedPath := "/api/consents/" + consentID + "/payment-rules"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, want)
	})

	got, err := svc.GetPaymentRules(context.Background(), consentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.PaymentRules) != 1 {
		t.Fatalf("len(PaymentRules) = %d, want 1", len(got.PaymentRules))
	}
	if got.PaymentRules[0].Type != "Window" {
		t.Errorf("PaymentRules[0].Type = %q, want %q", got.PaymentRules[0].Type, "Window")
	}
}

func TestGetPaymentRules_ServerError(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		mustEncode(t, w, epierrors.Problem{Title: "Not Found", Status: 404})
	})

	_, err := svc.GetPaymentRules(context.Background(), "missing-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get payment rules") {
		t.Errorf("error %q missing 'get payment rules' prefix", err.Error())
	}
}

// -------------------------------------------------------------------------
// 7. GetCaptureRules
// -------------------------------------------------------------------------

func TestGetCaptureRules_Success(t *testing.T) {
	consentID := "capture-rules-id"
	want := models.CaptureRulesResource{
		CaptureRules: models.CaptureRules{
			MaxDurationAfterAuthorization: 3600,
			MultiCapturesAllowed:          true,
		},
	}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		expectedPath := "/api/consents/" + consentID + "/capture-rules"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, want)
	})

	got, err := svc.GetCaptureRules(context.Background(), consentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.CaptureRules.MaxDurationAfterAuthorization != want.CaptureRules.MaxDurationAfterAuthorization {
		t.Errorf("MaxDurationAfterAuthorization = %d, want %d",
			got.CaptureRules.MaxDurationAfterAuthorization,
			want.CaptureRules.MaxDurationAfterAuthorization)
	}
	if !got.CaptureRules.MultiCapturesAllowed {
		t.Error("MultiCapturesAllowed = false, want true")
	}
}

func TestGetCaptureRules_ServerError(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		mustEncode(t, w, epierrors.Problem{Title: "Not Found", Status: 404})
	})

	_, err := svc.GetCaptureRules(context.Background(), "missing-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetCaptureRules_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	_, err := svc.GetCaptureRules(context.Background(), "some-id")
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// 8. GetPaymentPlan
// -------------------------------------------------------------------------

func TestGetPaymentPlan_Success(t *testing.T) {
	consentID := "payment-plan-id"
	amount := models.Money{EuroCents: 9999}
	want := models.PaymentPlanResource{
		Type:   "SingleImmediate",
		Amount: &amount,
	}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		expectedPath := "/api/consents/" + consentID + "/payment-plan"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, want)
	})

	got, err := svc.GetPaymentPlan(context.Background(), consentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Type != want.Type {
		t.Errorf("Type = %q, want %q", got.Type, want.Type)
	}
	if got.Amount == nil {
		t.Fatal("Amount is nil")
	}
	if got.Amount.EuroCents != amount.EuroCents {
		t.Errorf("Amount.EuroCents = %d, want %d", got.Amount.EuroCents, amount.EuroCents)
	}
}

func TestGetPaymentPlan_ServerError(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		mustEncode(t, w, epierrors.Problem{Title: "Not Found", Status: 404})
	})

	_, err := svc.GetPaymentPlan(context.Background(), "missing-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get payment plan") {
		t.Errorf("error %q missing 'get payment plan' prefix", err.Error())
	}
}

func TestGetPaymentPlan_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	_, err := svc.GetPaymentPlan(context.Background(), "some-id")
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// 9. GetCurrentToken
// -------------------------------------------------------------------------

func TestGetCurrentToken_Success(t *testing.T) {
	consentID := "token-consent-id"
	want := models.ConsentTokenResource{
		Token: "eyJhbGciOiJSUzI1NiJ9.example-token",
	}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		expectedPath := "/api/consents/" + consentID + "/current-token"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, want)
	})

	got, err := svc.GetCurrentToken(context.Background(), consentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Token != want.Token {
		t.Errorf("Token = %q, want %q", got.Token, want.Token)
	}
}

func TestGetCurrentToken_ServerError(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		mustEncode(t, w, epierrors.Problem{Title: "Not Found", Status: 404})
	})

	_, err := svc.GetCurrentToken(context.Background(), "missing-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get current token") {
		t.Errorf("error %q missing 'get current token' prefix", err.Error())
	}
}

func TestGetCurrentToken_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	_, err := svc.GetCurrentToken(context.Background(), "some-id")
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// 10. AddCurrentToken
// -------------------------------------------------------------------------

func TestAddCurrentToken_Returns200WithBody(t *testing.T) {
	consentID := "add-token-consent-id"
	want := models.ConsentTokenResource{
		Token: "new-token-value",
	}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		expectedPath := "/api/consents/" + consentID + "/current-token"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, want)
	})

	got, err := svc.AddCurrentToken(context.Background(), consentID, models.AddConsentTokenRequest{
		Token: "new-token-value",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("got nil, want token resource")
	}
	if got.Token != want.Token {
		t.Errorf("Token = %q, want %q", got.Token, want.Token)
	}
}

func TestAddCurrentToken_Returns204NilNil(t *testing.T) {
	consentID := "add-token-204-id"

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		expectedPath := "/api/consents/" + consentID + "/current-token"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	got, err := svc.AddCurrentToken(context.Background(), consentID, models.AddConsentTokenRequest{
		Token: "some-token",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("got non-nil resource on 204: %+v", got)
	}
}

func TestAddCurrentToken_ServerError(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		mustEncode(t, w, epierrors.ConstraintViolationProblem{
			Type:   "/problem/constraint-violation",
			Title:  "Constraint Violation",
			Status: 400,
			Violations: []epierrors.ConstraintViolation{
				{Field: ".token", Message: "must not be blank"},
			},
		})
	})

	_, err := svc.AddCurrentToken(context.Background(), "some-id", models.AddConsentTokenRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "add consent token") {
		t.Errorf("error %q missing 'add consent token' prefix", err.Error())
	}
}

func TestAddCurrentToken_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	_, err := svc.AddCurrentToken(context.Background(), "some-id", models.AddConsentTokenRequest{Token: "t"})
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// 11. GetShopFront
// -------------------------------------------------------------------------

func TestGetShopFront_Success(t *testing.T) {
	consentID := "shop-front-consent-id"
	want := models.ShopFrontResource{
		ShopName:                   "My Shop",
		ShopLogoFallbackIdentifier: "shop-logo-id",
		WebsiteURL:                 "https://myshop.example.com",
	}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		expectedPath := "/api/consents/" + consentID + "/shop-front"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, want)
	})

	got, err := svc.GetShopFront(context.Background(), consentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ShopName != want.ShopName {
		t.Errorf("ShopName = %q, want %q", got.ShopName, want.ShopName)
	}
	if got.WebsiteURL != want.WebsiteURL {
		t.Errorf("WebsiteURL = %q, want %q", got.WebsiteURL, want.WebsiteURL)
	}
}

func TestGetShopFront_ServerError(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		mustEncode(t, w, epierrors.Problem{Title: "Not Found", Status: 404})
	})

	_, err := svc.GetShopFront(context.Background(), "missing-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get shop front") {
		t.Errorf("error %q missing 'get shop front' prefix", err.Error())
	}
}

func TestGetShopFront_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	_, err := svc.GetShopFront(context.Background(), "some-id")
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// 12. DenyConsent
// -------------------------------------------------------------------------

func TestDenyConsent_Success_WithRequest(t *testing.T) {
	consentID := "deny-consent-id"

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		expectedPath := "/api/consents/" + consentID + "/denial"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := svc.DenyConsent(context.Background(), consentID, &models.DenyConsentRequest{
		ConsumerDeviceTelemetry: &models.ConsumerDeviceTelemetry{
			DeviceIPAddress: "192.168.1.1",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDenyConsent_Success_NilRequest_EmptyBody(t *testing.T) {
	consentID := "deny-nil-req-id"

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		// Verify no body was sent when req is nil.
		if r.ContentLength > 0 {
			t.Errorf("ContentLength = %d, want 0 for nil request", r.ContentLength)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := svc.DenyConsent(context.Background(), consentID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDenyConsent_ServerError(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		mustEncode(t, w, epierrors.Problem{Title: "Conflict", Status: 409, Detail: "consent already denied"})
	})

	err := svc.DenyConsent(context.Background(), "some-id", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "deny consent") {
		t.Errorf("error %q missing 'deny consent' prefix", err.Error())
	}
}

func TestDenyConsent_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	err := svc.DenyConsent(context.Background(), "some-id", nil)
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// 13. RevokeConsent
// -------------------------------------------------------------------------

func TestRevokeConsent_Success(t *testing.T) {
	consentID := "revoke-consent-id"

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		expectedPath := "/api/consents/" + consentID + "/revocation"
		if r.URL.Path != expectedPath {
			t.Errorf("path = %q, want %q", r.URL.Path, expectedPath)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := svc.RevokeConsent(context.Background(), consentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRevokeConsent_ServerError(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		mustEncode(t, w, epierrors.Problem{Title: "Conflict", Status: 409, Detail: "consent already revoked"})
	})

	err := svc.RevokeConsent(context.Background(), "already-revoked-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "revoke consent") {
		t.Errorf("error %q missing 'revoke consent' prefix", err.Error())
	}
}

func TestRevokeConsent_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	err := svc.RevokeConsent(context.Background(), "some-id")
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// 14. CreateOfflineQRCodeConsent
// -------------------------------------------------------------------------

func TestCreateOfflineQRCodeConsent_Success(t *testing.T) {
	want := models.ConsentResource{
		ID:     "qr-consent-id",
		Status: models.ConsentStatus{Value: "Requested"},
	}

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		assertAuthHeader(t, r)
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/consents/offline-qr-code/create" {
			t.Errorf("path = %q, want /api/consents/offline-qr-code/create", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/hal+json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		mustEncode(t, w, want)
	})

	got, err := svc.CreateOfflineQRCodeConsent(context.Background(), models.OfflineQRCodeRequest{
		Type:             "SingleImmediate",
		Amount:           models.Money{EuroCents: 500},
		OrderReferenceID: "ord-qr-123",
		Description:      "QR payment",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != want.ID {
		t.Errorf("ID = %q, want %q", got.ID, want.ID)
	}
	if got.Status.Value != want.Status.Value {
		t.Errorf("Status.Value = %q, want %q", got.Status.Value, want.Status.Value)
	}
}

func TestCreateOfflineQRCodeConsent_ConstraintViolation(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		mustEncode(t, w, epierrors.ConstraintViolationProblem{
			Type:   "/problem/constraint-violation",
			Title:  "Constraint Violation",
			Status: 400,
			Violations: []epierrors.ConstraintViolation{
				{Field: ".amount", Message: "must be positive"},
				{Field: ".orderReferenceId", Message: "must not be blank"},
			},
		})
	})

	_, err := svc.CreateOfflineQRCodeConsent(context.Background(), models.OfflineQRCodeRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "create offline qr code consent") {
		t.Errorf("error %q missing 'create offline qr code consent' prefix", err.Error())
	}
}

func TestCreateOfflineQRCodeConsent_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	svc := weroapi.New(srv.URL, "test-token")
	_, err := svc.CreateOfflineQRCodeConsent(context.Background(), models.OfflineQRCodeRequest{
		Type:             "SingleImmediate",
		Amount:           models.Money{EuroCents: 100},
		OrderReferenceID: "ord-123",
	})
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

// -------------------------------------------------------------------------
// Cross-cutting: bearer token is always sent
// -------------------------------------------------------------------------

func TestBearerToken_SentOnEveryRequest(t *testing.T) {
	// Spot-check using GetConsent — the assertAuthHeader helper in each test
	// verifies "Bearer test-token" is present. This test explicitly verifies
	// the custom token value "my-secret-token" is forwarded correctly.
	const customToken = "my-secret-token"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer "+customToken {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer "+customToken)
		}
		w.WriteHeader(http.StatusOK)
		mustEncode(t, w, models.ConsentResource{ID: "x"})
	}))
	t.Cleanup(srv.Close)

	svc := weroapi.New(srv.URL, customToken)
	_, err := svc.GetConsent(context.Background(), "x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
