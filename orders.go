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

type Order struct {
	BuyerInfo struct {
		Email string `json:"Email"`
	} `json:"BuyerInfo"`
	AmazonOrderId                string    `json:"AmazonOrderId"`
	EarliestShipDate             time.Time `json:"EarliestShipDate"`
	SalesChannel                 string    `json:"SalesChannel"`
	OrderStatus                  string    `json:"OrderStatus"`
	NumberOfItemsShipped         int       `json:"NumberOfItemsShipped"`
	OrderType                    string    `json:"OrderType"`
	IsPremiumOrder               bool      `json:"IsPremiumOrder"`
	IsPrime                      bool      `json:"IsPrime"`
	FulfillmentChannel           string    `json:"FulfillmentChannel"`
	NumberOfItemsUnshipped       int       `json:"NumberOfItemsUnshipped"`
	HasRegulatedItems            bool      `json:"HasRegulatedItems"`
	IsReplacementOrder           string    `json:"IsReplacementOrder"`
	IsSoldByAB                   bool      `json:"IsSoldByAB"`
	LatestShipDate               time.Time `json:"LatestShipDate"`
	ShipServiceLevel             string    `json:"ShipServiceLevel"`
	IsISPU                       bool      `json:"IsISPU"`
	MarketplaceId                string    `json:"MarketplaceId"`
	PurchaseDate                 time.Time `json:"PurchaseDate"`
	IsAccessPointOrder           bool      `json:"IsAccessPointOrder"`
	SellerOrderId                string    `json:"SellerOrderId"`
	PaymentMethod                string    `json:"PaymentMethod"`
	IsBusinessOrder              bool      `json:"IsBusinessOrder"`
	PaymentMethodDetails         []string  `json:"PaymentMethodDetails"`
	IsGlobalExpressEnabled       bool      `json:"IsGlobalExpressEnabled"`
	LastUpdateDate               time.Time `json:"LastUpdateDate"`
	ShipmentServiceLevelCategory string    `json:"ShipmentServiceLevelCategory"`
}

type GetOrdersResponse struct {
	Orders        []Order   `json:"Orders"`
	NextToken     string    `json:"NextToken"`
	CreatedBefore time.Time `json:"CreatedBefore"`
}

func (s *Client) paginateOrders(ctx context.Context, qs url.Values, limit int) (*GetOrdersResponse, error) {
	resp, err := s.getOrders(ctx, qs)
	if err != nil {
		return nil, err
	}

	if resp.NextToken != "" {
		nextPage := resp
		i := 0
		for {
			if i > limit && limit != -1 {
				break
			}
			i++

			v := url.Values{}

			resp.NextToken = strings.ReplaceAll(resp.NextToken, "+", "%2B")
			resp.NextToken = strings.ReplaceAll(resp.NextToken, "=", "%3D")
			v.Set("NextToken", resp.NextToken)
			nextPage, err = s.getOrders(ctx, v)
			if err != nil {
				if v, ok := err.(Error); ok {
					if len(v.Errors) > 0 && v.Errors[0].Code == "InvalidInput" {
						return resp, nil
					}
				}
				return nil, err
			}

			resp.Orders = append(resp.Orders, nextPage.Orders...)
			if nextPage.NextToken == "" {
				break
			}
		}
	}

	return resp, nil
}

func (s *Client) getOrders(ctx context.Context, qs url.Values) (*GetOrdersResponse, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     s.Marketplace.Endpoint,
		Path:     "/orders/v0/orders",
		RawQuery: qs.Encode(),
	}

	req := request{
		Method: http.MethodGet,
		URL:    &u,
		Token:  s.Token,
	}

	res, err := s.retry(ctx, req, retryOptions{
		retryLimit:    10,
		sleepDuration: 1 * time.Minute,
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Payload GetOrdersResponse `json:"payload"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error decoding spapi response: %w", err)
	}

	return &resp.Payload, nil
}

type GetOrdersRequest struct {
	CreatedAfter      time.Time
	CreatedBefore     time.Time
	LastUpdatedBefore time.Time
	LastUpdatedAfter  time.Time
	OrderStatuses     []string
}

func (s *Client) GetOrders(ctx context.Context, opts *GetOrdersRequest) ([]Order, error) {
	qs := url.Values{}

	if opts == nil {
		opts = &GetOrdersRequest{
			CreatedAfter: time.Now().Add(-24 * time.Hour),
		}
	}

	if !opts.CreatedAfter.IsZero() {
		qs.Set("CreatedAfter", opts.CreatedAfter.Format(time.RFC3339))
	}
	if !opts.CreatedBefore.IsZero() {
		qs.Set("CreatedBefore", opts.CreatedBefore.Format(time.RFC3339))
	}

	if !opts.LastUpdatedAfter.IsZero() {
		qs.Set("LastUpdatedAfter", opts.LastUpdatedAfter.Format(time.RFC3339))
	}

	if !opts.LastUpdatedBefore.IsZero() {
		qs.Set("LastUpdatedBefore", opts.LastUpdatedBefore.Format(time.RFC3339))
	}

	if len(opts.OrderStatuses) > 0 {
		qs.Set("OrderStatuses", strings.Join(opts.OrderStatuses, ","))
	}

	qs.Set("MarketplaceIds", s.Marketplace.ID)

	resp, err := s.paginateOrders(ctx, qs, -1)
	if err != nil {
		return nil, err
	}

	return resp.Orders, nil
}
