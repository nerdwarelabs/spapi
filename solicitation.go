package spapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func (s *Client) CreateProductReviewAndSellerFeedbackSolicitation(ctx context.Context, orderId string) error {
	query := url.Values{}
	query.Add("marketplaceIds", s.Marketplace.ID)

	u := url.URL{
		Scheme:   "https",
		Host:     s.Marketplace.Endpoint,
		Path:     fmt.Sprintf("/solicitations/v1/orders/%s/solicitations/productReviewAndSellerFeedback", orderId),
		RawQuery: query.Encode(),
	}

	req := request{
		Method: http.MethodPost,
		URL:    &u,
		Token:  s.Token,
	}

	_, err := s.retry(ctx, req, retryOptions{
		retryLimit:    10,
		sleepDuration: 1 * time.Second,
	})
	if err != nil {
		return err
	}

	return nil
}
