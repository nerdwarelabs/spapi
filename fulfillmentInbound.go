package spapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ItemEligibilityPreviewResponse struct {
	ASIN                 string `json:"asin"`
	Program              string `json:"program"`
	MarketplaceId        string `json:"marketplaceId"`
	IsEligibleForProgram bool   `json:"isEligibleForProgram"`
}

type FulfillmentInboundProgram string

var (
	FulfillmentInboundProgramInbound   FulfillmentInboundProgram = "INBOUND"
	FulfillmentInboundProgramCommingle FulfillmentInboundProgram = "COMMINGLING"
)

func (s *Client) GetItemEligibilityPreview(ctx context.Context, asin string, program FulfillmentInboundProgram) (*ItemEligibilityPreviewResponse, error) {
	qs := url.Values{}
	qs.Add("asin", asin)
	qs.Add("program", string(program))
	qs.Add("marketplaceIds", s.Marketplace.ID)

	u := url.URL{
		Scheme:   "https",
		Host:     s.Marketplace.Endpoint,
		Path:     "/fba/inbound/v1/eligibility/itemPreview",
		RawQuery: qs.Encode(),
	}

	req := request{
		Method: http.MethodGet,
		URL:    &u,
		Token:  s.Token,
	}

	res, err := s.retry(ctx, req, retryOptions{
		retryLimit:    10,
		sleepDuration: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Payload ItemEligibilityPreviewResponse `json:"payload"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error decoding spapi response: %w", err)
	}

	return &resp.Payload, nil
}

type GetPrepInstructionsResponse struct {
	ASINPrepInstructionsList []struct {
		ASIN                string   `json:"ASIN"`
		BarcodeInstruction  string   `json:"BarcodeInstruction"`
		PrepGuidance        string   `json:"PrepGuidance"`
		PrepInstructionList []string `json:"PrepInstructionList"`
	} `json:"AsinPrepInstructionsList"`

	InvalidASINList []struct {
		ASIN        string `json:"ASIN"`
		ErrorReason string `json:"ErrorReason"`
	} `json:"InvalidASINList"`
}

type PrepInstruction struct {
	PrepInstructionList []string
	BarcodeInstruction  string
}

func (s *Client) GetItemPrepInstructions(ctx context.Context, shipTo string, asins []string) (*GetPrepInstructionsResponse, error) {
	qs := url.Values{}
	qs.Add("ASINList", strings.Join(asins, ","))
	qs.Add("ShipToCountryCode", shipTo)

	u := url.URL{
		Scheme:   "https",
		Host:     s.Marketplace.Endpoint,
		Path:     "/fba/inbound/v0/prepInstructions",
		RawQuery: qs.Encode(),
	}

	req := request{
		Method: http.MethodGet,
		URL:    &u,
		Token:  s.Token,
	}

	res, err := s.retry(ctx, req, retryOptions{
		retryLimit:    10,
		sleepDuration: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Payload GetPrepInstructionsResponse `json:"payload"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error decoding spapi response: %w", err)
	}

	return &resp.Payload, nil
}
