# Consent API Client

## Example
```go
// ConsentsService provides access to the /api/consents endpoints.
type ConsentsService struct {
	client *resty.Client
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
		return nil, fmt.Errorf("add consent token: failed making request: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("add consent token: %w", &errResp)
	}
	if resp.StatusCode() == 204 {
		return nil, nil
	}
	return &result, nil
}
```