package providers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BillingUsageProvider struct {
	apiKey string
}

func NewBillingUsageProvider(apiKey string) *BillingUsageProvider {
	return &BillingUsageProvider{
		apiKey: apiKey,
	}
}

func (p *BillingUsageProvider) IncrementUsage(customerID string) error {
	data := url.Values{}
	data.Set("event_name", "events")
	data.Set("payload[value]", "56610")
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
