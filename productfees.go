package spapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type getMyFeesEstimateRequest []getMyFeesEstimateItem

type priceToEstimateFees struct {
	ListingPrice Money `json:"ListingPrice"`
}
type feesEstimateRequest struct {
	MarketplaceID               string              `json:"MarketplaceId"`
	Identifier                  string              `json:"Identifier"`
	IsAmazonFulfilled           bool                `json:"IsAmazonFulfilled"`
	PriceToEstimateFees         priceToEstimateFees `json:"PriceToEstimateFees"`
	OptionalFulfillmentPrograms []string            `json:"OptionalFulfillmentPrograms,omitempty"`
}

type getMyFeesEstimateItem struct {
	IdType              string              `json:"IdType"`
	IdValue             string              `json:"IdValue"` // ASIN or SellerSKU
	FeesEstimateRequest feesEstimateRequest `json:"FeesEstimateRequest"`
}

type Money struct {
	CurrencyCode string  `json:"CurrencyCode"`
	Amount       float64 `json:"Amount"`
}

type GetMyFeesResponseItem struct {
	Status                 string `json:"Status"`
	FeesEstimateIdentifier struct {
		MarketplaceId         string `json:"MarketplaceId"`
		IdType                string `json:"IdType"`
		SellerId              string `json:"SellerId"`
		SellerInputIdentifier string `json:"SellerInputIdentifier"`
		IsAmazonFulfilled     bool   `json:"IsAmazonFulfilled"`
		IdValue               string `json:"IdValue"`
		PriceToEstimateFees   struct {
			ListingPrice struct {
				CurrencyCode string  `json:"CurrencyCode"`
				Amount       float64 `json:"Amount"`
			} `json:"ListingPrice"`
		} `json:"PriceToEstimateFees"`
	} `json:"FeesEstimateIdentifier"`
	FeesEstimate struct {
		TimeOfFeesEstimation string `json:"TimeOfFeesEstimation"`
		TotalFeesEstimate    struct {
			CurrencyCode string  `json:"CurrencyCode"`
			Amount       float64 `json:"Amount"`
		} `json:"TotalFeesEstimate"`
		FeeDetailList []struct {
			FeeType               string `json:"FeeType"`
			FeeAmount             Money  `json:"FeeAmount"`
			FinalFee              Money  `json:"FinalFee"`
			FeePromotion          Money  `json:"FeePromotion"`
			IncludedFeeDetailList []struct {
				FeeType      string `json:"FeeType"`
				FeeAmount    Money  `json:"FeeAmount"`
				FinalFee     Money  `json:"FinalFee"`
				FeePromotion Money  `json:"FeePromotion"`
			} `json:"IncludedFeeDetailList"`
		} `json:"FeeDetailList"`
	} `json:"FeesEstimate"`
}

type GetMyFeesItem struct {
	ASIN     string
	Price    float64
	Currency string
}

// GetProductFees From Amazon
func (s *Client) GetProductFees(ctx context.Context, items []GetMyFeesItem) ([]GetMyFeesResponseItem, error) {
	u := url.URL{
		Scheme: "https",
		Host:   s.Marketplace.Endpoint,
		Path:   "/products/fees/v0/feesEstimate",
	}

	reqItems := getMyFeesEstimateRequest{}
	for _, item := range items {
		reqItems = append(reqItems, getMyFeesEstimateItem{
			IdType:  "ASIN",
			IdValue: item.ASIN,
			FeesEstimateRequest: feesEstimateRequest{
				MarketplaceID:     s.Marketplace.ID,
				Identifier:        item.ASIN,
				IsAmazonFulfilled: true,
				PriceToEstimateFees: priceToEstimateFees{
					ListingPrice: Money{
						CurrencyCode: item.Currency,
						Amount:       item.Price,
					},
				},
			},
		})
	}

	body, err := json.Marshal(reqItems)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	req := request{
		Method: http.MethodPost,
		URL:    &u,
		Body:   body,
		Token:  s.Token,
	}

	res, err := s.retry(ctx, req, retryOptions{
		retryLimit:    10,
		sleepDuration: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	var resp []GetMyFeesResponseItem
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error decoding spapi response: %w", err)
	}

	return resp, nil
}
