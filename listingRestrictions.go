package spapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ListingRestrictionsReasonLink struct {
	Resource string `json:"resource"`
	Verb     string `json:"verb"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

type ListingRestrictionReason struct {
	Message string                          `json:"message"`
	Code    string                          `json:"reasonCode"` // APPROVAL_REQUIRED, ASIN_NOT_FOUND, NOT_ELIGIBLE
	Links   []ListingRestrictionsReasonLink `json:"links"`
}

type ListingRestriction struct {
	MarketplaceId string                     `json:"marketplaceId"`
	ConditionType string                     `json:"conditionType"` // new_new
	Reasons       []ListingRestrictionReason `json:"reasons"`
}

type ListingRestrictionsCondition string

var (
	ListingRestrictionsConditionNew ListingRestrictionsCondition = "new_new"
)

type ListingRestrictionsResponse struct {
	Restrictions []ListingRestriction `json:"restrictions"`
}

func (s *Client) GetListingRestrictions(ctx context.Context, asin string, condition ListingRestrictionsCondition) (*ListingRestrictionsResponse, error) {
	qs := url.Values{}
	qs.Add("asin", asin)
	qs.Add("conditionType", string(condition))
	qs.Add("sellerId", s.SellerID)
	qs.Add("marketplaceIds", s.Marketplace.ID)

	u := url.URL{
		Scheme:   "https",
		Host:     s.Marketplace.Endpoint,
		Path:     "/listings/2021-08-01/restrictions",
		RawQuery: qs.Encode(),
	}

	req := request{
		Method: http.MethodGet,
		URL:    &u,
		Token:  s.Token,
	}

	res, err := s.retry(ctx, req, retryOptions{
		retryLimit:    10,
		sleepDuration: 500 * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}

	var resp ListingRestrictionsResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error decoding spapi response: %w", err)
	}

	return &resp, nil
}
