package providers

import (
	"errors"
	"fmt"
	"io"
	domain_providers "ludiks/src/account/domain/providers"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/stripe/stripe-go/v82"
	cust "github.com/stripe/stripe-go/v82/customer"
	setupintent "github.com/stripe/stripe-go/v82/setupintent"
	subscription "github.com/stripe/stripe-go/v82/subscription"
)

type StripeProvider struct {
	apiKey string
}

func NewStripeProvider(apiKey string) domain_providers.StripeProvider {
	stripe.Key = apiKey
	return &StripeProvider{
		apiKey: apiKey,
	}
}

func (p *StripeProvider) autoDetectTaxIDType(country, taxIDValue string) string {
	if taxIDValue == "" {
		return ""
	}

	country = strings.ToUpper(country)

	switch country {
	case "FR", "DE", "IT", "ES", "NL", "BE", "AT", "IE", "PT", "FI", "SE", "DK", "PL", "CZ", "HU", "RO", "BG", "HR", "SI", "SK", "LT", "LV", "EE", "CY", "MT", "LU":
		return "eu_vat"
	case "GB":
		return "gb_vat"
	case "US":
		return "us_ein"
	case "CA":
		return "ca_gst"
	case "AU":
		return "au_abn"
	case "BR":
		return "br_cnpj"
	case "MX":
		return "mx_rfc"
	case "IN":
		return "in_gst"
	case "SG":
		return "sg_uen"
	case "TH":
		return "th_vat"
	case "TW":
		return "tw_vat"
	case "ZA":
		return "za_vat"
	case "CH":
		return "ch_vat"
	case "NO":
		return "no_vat"
	default:
		return ""
	}
}

func (p *StripeProvider) CreateCustomer(customerInfo domain_providers.CustomerInfo, organizationID string) (string, error) {
	custParams := &stripe.CustomerParams{
		Email: stripe.String(customerInfo.Email),
		Name:  stripe.String(customerInfo.CompanyName),
		Address: &stripe.AddressParams{
			Line1:      stripe.String(customerInfo.CompanyAddress.Line1),
			Line2:      stripe.String(customerInfo.CompanyAddress.Line2),
			City:       stripe.String(customerInfo.CompanyAddress.City),
			State:      stripe.String(customerInfo.CompanyAddress.State),
			PostalCode: stripe.String(customerInfo.CompanyAddress.PostalCode),
			Country:    stripe.String(customerInfo.CompanyAddress.Country),
		},
		Phone: stripe.String(customerInfo.ContactPhone),
		Metadata: map[string]string{
			"customer_type": "b2b",
			"contact_name":  customerInfo.ContactName,
		},
	}

	if customerInfo.TaxIDValue != "" {
		taxIDType := p.autoDetectTaxIDType(customerInfo.CompanyAddress.Country, customerInfo.TaxIDValue)
		if taxIDType != "" {
			custParams.TaxIDData = []*stripe.CustomerTaxIDDataParams{
				{
					Type:  stripe.String(taxIDType),
					Value: stripe.String(customerInfo.TaxIDValue),
				},
			}
		}
	}

	customer, err := cust.New(custParams)

	return customer.ID, err
}

func (p *StripeProvider) SetupPaymentMethod(
	customerID string,
	priceID *string,
	organizationID string,
) (*domain_providers.SetupPaymentMethodResult, error) {
	var subscriptionID string

	if priceID != nil && *priceID != "" {
		subParams := &stripe.SubscriptionParams{
			Customer: stripe.String(customerID),
			Items: []*stripe.SubscriptionItemsParams{
				{
					Price: stripe.String(*priceID),
				},
			},
			Metadata: map[string]string{
				"subscription_type": "b2b_pro",
			},
		}

		sub, err := subscription.New(subParams)
		if err != nil {
			return nil, fmt.Errorf("subscription creation failed: %v", err)
		}
		subscriptionID = sub.ID
	}

	setupIntentParams := &stripe.SetupIntentParams{
		Customer: stripe.String(customerID),
		Usage:    stripe.String("off_session"),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
			stripe.String("sepa_debit"),
		},
		Metadata: map[string]string{
			"setup_type":      "b2b_subscription",
			"organization_id": organizationID,
		},
	}

	setupIntent, err := setupintent.New(setupIntentParams)
	if err != nil {
		return nil, fmt.Errorf("setup intent creation failed: %v", err)
	}

	return &domain_providers.SetupPaymentMethodResult{
		ClientSecret:   setupIntent.ClientSecret,
		SubscriptionID: subscriptionID,
	}, nil
}

func (p *StripeProvider) GetOrganizationIDFromIntentPayment(setupIntentID string) (string, error) {
	params := &stripe.SetupIntentParams{}
	result, err := setupintent.Get(setupIntentID, params)
	if err != nil {
		return "", err
	}

	if result.Status == "succeeded" {
		return result.Metadata["organization_id"], nil
	} else {
		return "", errors.New("setup intent was not successful")
	}
}

func (p *StripeProvider) CancelSubscription(subscriptionID string, cancelAtPeriodEnd bool) (*domain_providers.CancelSubscriptionResult, error) {
	params := &stripe.SubscriptionParams{}

	if cancelAtPeriodEnd {
		params.CancelAtPeriodEnd = stripe.Bool(true)
	} else {
		params.CancelAtPeriodEnd = stripe.Bool(false)
	}

	sub, err := subscription.Update(subscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("subscription cancellation failed: %v", err)
	}

	var endDate int64
	if cancelAtPeriodEnd && sub.CancelAt != 0 {
		endDate = sub.CancelAt
	} else {
		endDate = time.Now().Unix()
	}

	return &domain_providers.CancelSubscriptionResult{
		SubscriptionID: sub.ID,
		CancelledAt:    time.Unix(sub.CanceledAt, 0).Format(time.RFC3339),
		EndDate:        time.Unix(endDate, 0).Format(time.RFC3339),
	}, nil
}

func (p *StripeProvider) GetSubscriptionDetails(subscriptionID string) (*domain_providers.SubscriptionDetails, error) {
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription details: %v", err)
	}

	return &domain_providers.SubscriptionDetails{
		ID:     sub.ID,
		Status: string(sub.Status),
	}, nil
}

func (p *StripeProvider) ReportUsage(customerID string) error {
	data := url.Values{}
	data.Set("event_name", "events")
	data.Set("payload[value]", fmt.Sprintf("%d", 1))
	data.Set("payload[stripe_customer_id]", customerID)
	data.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))

	req, err := http.NewRequest("POST", "https://api.stripe.com/v1/billing/meter_events", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
