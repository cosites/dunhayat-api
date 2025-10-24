package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ZibalConfig struct {
	MerchantID string
	BaseURL    string
	Timeout    time.Duration
	APIToken   string
}

type ZibalClient struct {
	config     ZibalConfig
	httpClient *http.Client
}

type ZibalPaymentRequest struct {
	MerchantID        string             `json:"merchantId"`
	Amount            int                `json:"amount"`
	OrderID           string             `json:"orderId"`
	CallbackURL       string             `json:"callbackUrl"`
	Description       string             `json:"description,omitempty"`
	Mobile            string             `json:"mobile,omitempty"`
	MultiplexingInfos []MultiplexingInfo `json:"multiplexingInfos,omitempty"`
}

type MultiplexingInfo struct {
	SubMerchantID string `json:"subMerchantId"`
	Amount        int    `json:"amount"`
	Description   string `json:"description"`
}

type ZibalPaymentResponse struct {
	Result  int    `json:"result"`
	TrackID string `json:"trackId"`
	Message string `json:"message"`
}

type ZibalVerifyRequest struct {
	MerchantID string `json:"merchantId"`
	TrackID    string `json:"trackId"`
}

type ZibalVerifyResponse struct {
	Result           int    `json:"result"`
	Amount           int    `json:"amount"`
	OrderID          string `json:"orderId"`
	CardNumber       string `json:"cardNumber"`
	HashedCardNumber string `json:"hashedCardNumber"`
	Message          string `json:"message"`
}

type ZibalCallbackData struct {
	Success          bool   `json:"success"`
	Status           int    `json:"status"`
	TrackID          string `json:"trackId"`
	OrderID          string `json:"orderId"`
	Amount           int    `json:"amount"`
	CardNumber       string `json:"cardNumber"`
	HashedCardNumber string `json:"hashedCardNumber"`
}

func NewZibalClient(config ZibalConfig) *ZibalClient {
	return &ZibalClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (c *ZibalClient) CreatePaymentRequest(
	req ZibalPaymentRequest,
) (*ZibalPaymentResponse, error) {
	req.MerchantID = c.config.MerchantID

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to marshal payment request: %w", err,
		)
	}

	httpReq, err := http.NewRequest(
		"POST",
		c.config.BaseURL+"/request",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create HTTP request: %w", err,
		)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.config.APIToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.APIToken)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to send payment request: %w", err,
		)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read response body: %w", err,
		)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"payment request failed with status %d: %s",
			resp.StatusCode, string(body),
		)
	}

	var paymentResp ZibalPaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal payment response: %w", err,
		)
	}

	if paymentResp.Result != 100 {
		return nil, fmt.Errorf(
			"zibal payment request failed: %s (result: %d)",
			paymentResp.Message, paymentResp.Result,
		)
	}

	return &paymentResp, nil
}

func (c *ZibalClient) VerifyPayment(
	req ZibalVerifyRequest,
) (*ZibalVerifyResponse, error) {
	req.MerchantID = c.config.MerchantID

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to marshal verify request: %w", err,
		)
	}

	httpReq, err := http.NewRequest(
		"POST", c.config.BaseURL+"/verify", bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create HTTP request: %w", err,
		)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.config.APIToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.APIToken)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to send verify request: %w", err,
		)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read response body: %w", err,
		)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"verify request failed with status %d: %s",
			resp.StatusCode, string(body),
		)
	}

	var verifyResp ZibalVerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal verify response: %w", err,
		)
	}

	return &verifyResp, nil
}

func (c *ZibalClient) GetPaymentURL(trackID string) string {
	return fmt.Sprintf("%s/start/%s", c.config.BaseURL, trackID)
}
