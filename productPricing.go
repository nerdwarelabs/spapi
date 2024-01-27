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

type CompetitivePrice struct {
	BelongsToRequester bool   `json:"belongsToRequester"`
	Condition          string `json:"condition"`
	Subcondition       string `json:"subcondition"`
	Price              struct {
		LandedPrice  Money `json:"LandedPrice"`
		ListingPrice Money `json:"ListingPrice"`
		Shipping     Money `json:"Shipping"`
	} `json:"Price"`
	CompetitivePriceID string `json:"CompetitivePriceId"`
}

type CompetitiveOffers struct {
	Condition string `json:"condition"`
	Count     int    `json:"count"`
}

type GetCompetitivePricingForASINItem struct {
	ASIN    string `json:"ASIN"`
	Product struct {
		CompetitivePricing struct {
			CompetitivePrices     []CompetitivePrice  `json:"CompetitivePrices"`
			NumberOfOfferListings []CompetitiveOffers `json:"NumberOfOfferListings"`
		} `json:"CompetitivePricing"`
	} `json:"Product"`
	Status string `json:"status"`
}

func (c *Client) GetCompetitivePricingByASIN(ctx context.Context, asins []string) ([]*GetCompetitivePricingForASINItem, error) {
	qs := url.Values{}
	qs.Add("MarketplaceId", c.Marketplace.ID)
	qs.Add("Asins", strings.Join(asins, ","))
	qs.Add("ItemType", "Asin")

	u := url.URL{
		Scheme:   "https",
		Host:     c.Marketplace.Endpoint,
		Path:     "/products/pricing/v0/competitivePrice",
		RawQuery: qs.Encode(),
	}

	req := request{
		Method: http.MethodGet,
		URL:    &u,
		Token:  c.Token,
	}

	res, err := c.retry(ctx, req, retryOptions{
		retryLimit:    10,
		sleepDuration: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Payload []*GetCompetitivePricingForASINItem `json:"payload"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error decoding spapi response: %w", err)
	}

	return resp.Payload, nil
}
