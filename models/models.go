package models

import "time"

// --- Primitives / Shared ---

// Money represents a monetary amount in euro cents.
type Money struct {
	EuroCents int `json:"euroCents"`
}

// Link represents a hypermedia link.
type Link struct {
	Href string `json:"href"`
	Name string `json:"name,omitempty"`
}

// Page holds pagination metadata.
type Page struct {
	Size              int  `json:"size"`
	TotalElements     int  `json:"totalElements"`
	TotalPages        int  `json:"totalPages"`
	Number            int  `json:"number"`
	RequestedPageSize int  `json:"requestedPageSize"`
	LimitReached      bool `json:"limitReached"`
}

// --- Error types ---

// Problem represents a standard RFC 7807 problem detail.
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	CallID   string `json:"callId,omitempty"`
}

// ConstraintViolation is a single field violation within a ConstraintViolationProblem.
type ConstraintViolation struct {
	Field            string `json:"field"`
	Message          string `json:"message"`
	UserMessageKey   string `json:"userMessageKey,omitempty"`
	UserMessageValue string `json:"userMessageValue,omitempty"`
}

// ConstraintViolationProblem is returned when request validation fails.
type ConstraintViolationProblem struct {
	Type       string                `json:"type"`
	Title      string                `json:"title"`
	Status     int                   `json:"status"`
	CallID     string                `json:"callId,omitempty"`
	Violations []ConstraintViolation `json:"violations"`
}

// --- Consent status ---

// ConsentStatus holds the current status of a consent resource.
// The Value field determines which fields are populated.
type ConsentStatus struct {
	Value      string     `json:"value"`
	At         time.Time  `json:"at"`
	TimesOutAt *time.Time `json:"timesOutAt,omitempty"`
	Source     string     `json:"source,omitempty"`
}

// --- Payment plan resource types (response-side) ---

// PriceChangeHistoryItem records a historical price change for a subscription.
type PriceChangeHistoryItem struct {
	Price     Money     `json:"price"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// PaymentPlanLinks holds hypermedia links for a payment plan resource.
type PaymentPlanLinks struct {
	Self Link `json:"self"`
}

// PaymentPlanResource represents any payment plan variant returned by the API.
// The Type field determines which fields are populated.
type PaymentPlanResource struct {
	Type                    string                   `json:"type"`
	Amount                  *Money                   `json:"amount,omitempty"`
	UltimateBeneficiaryName string                   `json:"ultimateBeneficiaryName,omitempty"`
	UltimateBeneficiaryIBAN string                   `json:"ultimateBeneficiaryIban,omitempty"`
	Commercial              bool                     `json:"commercial,omitempty"`
	Description             string                   `json:"description,omitempty"`
	CaptureTrigger          string                   `json:"captureTrigger,omitempty"`
	AmountPaymentType       string                   `json:"amountPaymentType,omitempty"`
	MaxAuthToCaptureTime    int                      `json:"maxAuthToCaptureTime,omitempty"`
	MultiCapturesAllowed    bool                     `json:"multiCapturesAllowed,omitempty"`
	Reason                  string                   `json:"reason,omitempty"`
	MaxAmount               *Money                   `json:"maxAmount,omitempty"`
	PriceChangeHistory      []PriceChangeHistoryItem `json:"priceChangeHistory,omitempty"`
	StartAt                 *time.Time               `json:"startAt,omitempty"`
	FirstPaymentAt          *time.Time               `json:"firstPaymentAt,omitempty"`
	EndAt                   *time.Time               `json:"endAt,omitempty"`
	Repetition              string                   `json:"repetition,omitempty"`
	AmountOfFirstPayment    *Money                   `json:"amountOfFirstPayment,omitempty"`
	Links                   *PaymentPlanLinks        `json:"_links,omitempty"`
}

// --- Consent token resource ---

// ConsentTokenResource holds a consent token and its self link.
type ConsentTokenResource struct {
	Token string `json:"token"`
	Links struct {
		Self Link `json:"self"`
	} `json:"_links"`
}

// --- Payment rules resource ---

// ConsentPaymentRuleConstraint specifies the constraint for a payment rule window.
type ConsentPaymentRuleConstraint struct {
	Type      string `json:"type"`
	MaxAmount Money  `json:"maxAmount"`
}

// ConsentPaymentRuleCyclicModification describes a modification for specific cycles.
type ConsentPaymentRuleCyclicModification struct {
	CycleIndexes string                       `json:"cycleIndexes"`
	Constraint   ConsentPaymentRuleConstraint `json:"constraint"`
}

// ConsentPaymentRule represents any payment rule variant.
// The Type field determines which fields are populated.
type ConsentPaymentRule struct {
	Type           string                                 `json:"type"`
	Constraint     *ConsentPaymentRuleConstraint          `json:"constraint,omitempty"`
	StartOffset    string                                 `json:"startOffset,omitempty"`
	WindowDuration string                                 `json:"windowDuration,omitempty"`
	Modifications  []ConsentPaymentRuleCyclicModification `json:"modifications,omitempty"`
	Repetition     string                                 `json:"repetition,omitempty"`
	Duration       string                                 `json:"duration,omitempty"`
}

// ConsentPaymentRulesResource holds the set of payment rules for a consent.
type ConsentPaymentRulesResource struct {
	PaymentRules []ConsentPaymentRule `json:"paymentRules"`
	Links        struct {
		Self Link `json:"self"`
	} `json:"_links"`
}

// --- Capture rules resource ---

// CaptureRules holds the capture constraints for a consent.
type CaptureRules struct {
	MaxDurationAfterAuthorization int  `json:"maxDurationAfterAuthorization"`
	MultiCapturesAllowed          bool `json:"multiCapturesAllowed"`
}

// CaptureRulesResource wraps the capture rules and its self link.
type CaptureRulesResource struct {
	CaptureRules CaptureRules `json:"captureRules"`
	Links        struct {
		Self Link `json:"self"`
	} `json:"_links"`
}

// --- Address ---

// Address represents a postal address used in both request and response contexts.
type Address struct {
	CompanyName   string `json:"companyName,omitempty"`
	RecipientName string `json:"recipientName"`
	Line1         string `json:"line1"`
	Line2         string `json:"line2,omitempty"`
	City          string `json:"city"`
	PostalCode    string `json:"postalCode"`
	State         string `json:"state,omitempty"`
	CountryCode   string `json:"countryCode"`
}

// --- Consent details resource ---

// ConsentContact holds consumer contact information.
type ConsentContact struct {
	EmailAddress string `json:"emailAddress,omitempty"`
	PhoneNumber  string `json:"phoneNumber,omitempty"`
}

// ConsentDetailsResource holds enriched details for a consent.
type ConsentDetailsResource struct {
	Description     string          `json:"description,omitempty"`
	OrderPageURL    string          `json:"orderPageUrl,omitempty"`
	Gift            bool            `json:"gift,omitempty"`
	BillingAddress  *Address        `json:"billingAddress,omitempty"`
	ShippingAddress *Address        `json:"shippingAddress,omitempty"`
	ShippingCosts   *Money          `json:"shippingCosts,omitempty"`
	Contact         *ConsentContact `json:"contact,omitempty"`
	Links           struct {
		Self      Link  `json:"self"`
		OrderPage *Link `json:"orderPage,omitempty"`
	} `json:"_links"`
}

// --- Means of payment resource ---

// ConsentMeansOfPaymentResource holds consumer payment instrument info.
type ConsentMeansOfPaymentResource struct {
	ConsumerPspLogoURL string `json:"consumerPspLogoUrl"`
	MaskedIBAN         string `json:"maskedIban"`
}

// --- EasyCheckout & BusinessRules ---

// EasyCheckout configures which consumer contact fields are requested.
type EasyCheckout struct {
	IncludePhone bool `json:"includePhone"`
	IncludeEmail bool `json:"includeEmail"`
}

// DutchMigration holds routing advice for Dutch migration business rules.
type DutchMigration struct {
	RoutingAdviceID string `json:"routingAdviceId"`
}

// BusinessRules holds optional business-rule overrides for a consent.
type BusinessRules struct {
	DutchMigration *DutchMigration `json:"dutchMigration,omitempty"`
}

// --- Main consent resource ---

// AcceptorShop identifies the shop within an acceptor.
type AcceptorShop struct {
	ID                          string `json:"id"`
	Name                        string `json:"name"`
	ShopLogoFallbackIdentifier  string `json:"shopLogoFallbackIdentifier"`
}

// ConsentAcceptor identifies the merchant accepting the consent.
type ConsentAcceptor struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	MCC         string       `json:"mcc"`
	CountryCode string       `json:"countryCode"`
	PspID       string       `json:"pspId"`
	Shop        AcceptorShop `json:"shop"`
}

// ConsumerIdentifier identifies the consumer wallet.
type ConsumerIdentifier struct {
	WalletID string `json:"walletId"`
}

// ConsentMeta holds creation metadata for a consent.
type ConsentMeta struct {
	CreatedAt time.Time `json:"createdAt"`
}

// ConsentResourceLinks holds all hypermedia links for a consent resource.
type ConsentResourceLinks struct {
	Self         Link  `json:"self"`
	PaymentRules *Link `json:"paymentRules,omitempty"`
	CaptureRules *Link `json:"captureRules,omitempty"`
	PaymentPlan  *Link `json:"paymentPlan,omitempty"`
	CurrentToken *Link `json:"currentToken,omitempty"`
	Details      *Link `json:"details,omitempty"`
	LandingPage  *Link `json:"landingPage,omitempty"`
	ReturnPage   *Link `json:"returnPage,omitempty"`
	PrintLink    *Link `json:"printLink,omitempty"`
	PrintQRCode  *Link `json:"printQrCode,omitempty"`
	Denial       *Link `json:"denial,omitempty"`
	Revocation   *Link `json:"revocation,omitempty"`
	ShopLogo     *Link `json:"shopLogo,omitempty"`
	ShopInfo     *Link `json:"shopInfo,omitempty"`
}

// ConsentResourceEmbedded holds optionally embedded sub-resources for a consent.
type ConsentResourceEmbedded struct {
	PaymentPlan  *PaymentPlanResource           `json:"paymentPlan,omitempty"`
	CurrentToken *ConsentTokenResource          `json:"currentToken,omitempty"`
	PaymentRules *ConsentPaymentRulesResource   `json:"paymentRules,omitempty"`
	CaptureRules *CaptureRulesResource          `json:"captureRules,omitempty"`
	Details      *ConsentDetailsResource        `json:"details,omitempty"`
	PaymentMeans *ConsentMeansOfPaymentResource `json:"paymentMeans,omitempty"`
}

// ConsentResource is the primary consent representation returned by the API.
type ConsentResource struct {
	ID                  string               `json:"id"`
	ShortID             string               `json:"shortId,omitempty"`
	DurationUntilExpiry int                  `json:"durationUntilExpiry,omitempty"`
	ConsentedAt         *time.Time           `json:"consentedAt,omitempty"`
	ExpiresAt           *time.Time           `json:"expiresAt,omitempty"`
	Acceptor            ConsentAcceptor      `json:"acceptor"`
	ConsumerIdentifier  *ConsumerIdentifier  `json:"consumerIdentifier,omitempty"`
	OrderReferenceID    string               `json:"orderReferenceId,omitempty"`
	CallbackURL         string               `json:"callbackUrl,omitempty"`
	Status              ConsentStatus        `json:"status"`
	Meta                ConsentMeta          `json:"meta"`
	EasyCheckout        *EasyCheckout        `json:"easyCheckout,omitempty"`
	BusinessRules       *BusinessRules       `json:"businessRules,omitempty"`
	BeneficiaryName     string               `json:"beneficiaryName,omitempty"`
	Links               ConsentResourceLinks `json:"_links"`
	Embedded            *ConsentResourceEmbedded `json:"_embedded,omitempty"`
}

// --- Collection resources ---

// ConsentCollectionResource is returned by paginated consent list endpoints.
type ConsentCollectionResource struct {
	Embedded struct {
		Consents []ConsentResource `json:"consents"`
	} `json:"_embedded"`
	Links struct {
		Self  Link  `json:"self"`
		First *Link `json:"first,omitempty"`
		Prev  *Link `json:"prev,omitempty"`
		Next  *Link `json:"next,omitempty"`
		Last  *Link `json:"last,omitempty"`
	} `json:"_links"`
	Page *Page `json:"page,omitempty"`
}

// --- Request-side payment plans ---

// RequestedPaymentPlan describes the payment plan for a consent request.
// The Type field determines which fields are required.
type RequestedPaymentPlan struct {
	Type                    string     `json:"type"`
	Amount                  *Money     `json:"amount,omitempty"`
	UltimateBeneficiaryName string     `json:"ultimateBeneficiaryName,omitempty"`
	UltimateBeneficiaryIBAN string     `json:"ultimateBeneficiaryIban,omitempty"`
	Commercial              bool       `json:"commercial,omitempty"`
	Description             string     `json:"description,omitempty"`
	CaptureTrigger          string     `json:"captureTrigger,omitempty"`
	AmountPaymentType       string     `json:"amountPaymentType,omitempty"`
	MaxAuthToCaptureTime    int        `json:"maxAuthToCaptureTime,omitempty"`
	MultiCapturesAllowed    bool       `json:"multiCapturesAllowed,omitempty"`
	Reason                  string     `json:"reason,omitempty"`
	StartAt                 *time.Time `json:"startAt,omitempty"`
	FirstPaymentAt          *time.Time `json:"firstPaymentAt,omitempty"`
	EndAt                   *time.Time `json:"endAt,omitempty"`
	Repetition              string     `json:"repetition,omitempty"`
	AmountOfFirstPayment    *Money     `json:"amountOfFirstPayment,omitempty"`
}

// --- Request models ---

// RequestedAcceptor identifies the merchant in a consent request.
type RequestedAcceptor struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	MCC         string       `json:"mcc"`
	CountryCode string       `json:"countryCode"`
	Shop        AcceptorShop `json:"shop"`
}

// PrintOptions configures printed QR code output.
type PrintOptions struct {
	Lifetime int    `json:"lifetime,omitempty"`
	Format   string `json:"format,omitempty"`
}

// BasketItem represents a line item in a shopping basket.
type BasketItem struct {
	SKU         string  `json:"sku,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   Money   `json:"unitPrice"`
	Tax         *Money  `json:"tax,omitempty"`
}

// RequestedConsentDetails holds enriched order details for a consent request.
type RequestedConsentDetails struct {
	Description     string       `json:"description,omitempty"`
	OrderPageURL    string       `json:"orderPageUrl,omitempty"`
	Gift            bool         `json:"gift,omitempty"`
	BillingAddress  *Address     `json:"billingAddress,omitempty"`
	ShippingAddress *Address     `json:"shippingAddress,omitempty"`
	ShippingCosts   *Money       `json:"shippingCosts,omitempty"`
	BasketItems     []BasketItem `json:"basketItems,omitempty"`
}

// ConsumerDeviceTelemetry holds device fingerprint data for fraud prevention.
type ConsumerDeviceTelemetry struct {
	DeviceIPAddress         string     `json:"deviceIpAddress,omitempty"`
	DevicePort              int        `json:"devicePort,omitempty"`
	DeviceUserAgent         string     `json:"deviceUserAgent,omitempty"`
	DeviceModel             string     `json:"deviceModel,omitempty"`
	DeviceScreenResolution  string     `json:"deviceScreenResolution,omitempty"`
	DeviceCellularProviders string     `json:"deviceCellularProviders,omitempty"`
	DeviceLocation          string     `json:"deviceLocation,omitempty"`
	DeviceTimeZone          string     `json:"deviceTimeZone,omitempty"`
	DeviceID                string     `json:"deviceId,omitempty"`
	DeviceOS                string     `json:"deviceOs,omitempty"`
	DeviceLanguage          string     `json:"deviceLanguage,omitempty"`
	DeviceTime              *time.Time `json:"deviceTime,omitempty"`
}

// RequestConsentRequest is the body for POST /api/consents.
type RequestConsentRequest struct {
	PaymentPlan             RequestedPaymentPlan     `json:"paymentPlan"`
	OrderReferenceID        string                   `json:"orderReferenceId"`
	Acceptor                RequestedAcceptor        `json:"acceptor"`
	ConsumerIdentifier      *ConsumerIdentifier      `json:"consumerIdentifier,omitempty"`
	ReturnURL               string                   `json:"returnUrl,omitempty"`
	CallbackURL             string                   `json:"callbackUrl,omitempty"`
	PrintOptions            *PrintOptions            `json:"printOptions,omitempty"`
	Details                 *RequestedConsentDetails `json:"details,omitempty"`
	EasyCheckout            *EasyCheckout            `json:"easyCheckout,omitempty"`
	BusinessRules           *BusinessRules           `json:"businessRules,omitempty"`
	BeneficiaryName         string                   `json:"beneficiaryName,omitempty"`
	ShoppingDeviceTelemetry *ConsumerDeviceTelemetry `json:"shoppingDeviceTelemetry,omitempty"`
	InitiationType          string                   `json:"initiationType,omitempty"`
}

// ListConsentsParams are the query parameters for GET /api/consents.
type ListConsentsParams struct {
	FindBy           string `url:"findBy,omitempty"`
	OrderReferenceID string `url:"orderReferenceId,omitempty"`
	ShortID          string `url:"shortId,omitempty"`
	PaymentMeansID   string `url:"paymentMeansId,omitempty"`
	PaymentPlan      string `url:"paymentPlan,omitempty"`
}

// AddConsentTokenRequest is the body for PUT /api/consents/{consentId}/current-token.
type AddConsentTokenRequest struct {
	Token                   string                   `json:"token"`
	ConsumerDeviceTelemetry *ConsumerDeviceTelemetry `json:"consumerDeviceTelemetry,omitempty"`
}

// Contact holds consumer contact details.
type Contact struct {
	EmailAddress string `json:"emailAddress,omitempty"`
	PhoneNumber  string `json:"phoneNumber,omitempty"`
}

// SetEasyCheckoutProfileRequest is the body for PUT /api/consents/{consentId}/details.
type SetEasyCheckoutProfileRequest struct {
	ShippingAddress Address  `json:"shippingAddress"`
	BillingAddress  *Address `json:"billingAddress,omitempty"`
	Contact         *Contact `json:"contact,omitempty"`
}

// DenyConsentRequest is the optional body for PUT /api/consents/{consentId}/denial.
type DenyConsentRequest struct {
	ConsumerDeviceTelemetry *ConsumerDeviceTelemetry `json:"consumerDeviceTelemetry,omitempty"`
}

// ShopFrontResource is returned by GET /api/consents/{consentId}/shop-front.
type ShopFrontResource struct {
	WebsiteURL                 string `json:"websiteUrl,omitempty"`
	ShopName                   string `json:"shopName"`
	ShopLogoFallbackIdentifier string `json:"shopLogoFallbackIdentifier"`
	Links                      struct {
		Self     Link  `json:"self"`
		ShopLogo *Link `json:"shopLogo,omitempty"`
	} `json:"_links"`
}

// OfflineQRCodeRequest is the body for POST /api/consents/offline-qr-code/create.
//
// Deprecated: use the standard consent flow instead.
type OfflineQRCodeRequest struct {
	Type             string `json:"type"`
	Description      string `json:"description,omitempty"`
	Amount           Money  `json:"amount"`
	OrderReferenceID string `json:"orderReferenceId"`
	PayconiqID       string `json:"payconiqId,omitempty"`
	AcceptorQRCodeID string `json:"acceptorQRCodeId,omitempty"`
}
