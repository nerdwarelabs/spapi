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

type BrowseClassification struct {
	DisplayName      string `json:"displayName"`
	ClassificationId string `json:"classificationId"`
}

type CatalogItemContributorRole struct {
	DisplayName string `json:"displayName"`
	Value       string `json:"value"`
}
type CatalogItemContributor struct {
	Role  CatalogItemContributorRole `json:"role"`
	Value string                     `json:"value"`
}

type CatalogItemSummary struct {
	MarketplaceId           string                   `json:"marketplaceId"`
	AdultProduct            bool                     `json:"adultProduct"`
	Autographed             bool                     `json:"autographed"`
	Brand                   string                   `json:"brand"`
	BrowseClassificiations  []BrowseClassification   `json:"browseClassifications"`
	Color                   string                   `json:"color"`
	Contributors            []CatalogItemContributor `json:"contributors"`
	ItemClassification      string                   `json:"itemClassification"`
	ItemName                string                   `json:"itemName"`
	Manufacturer            string                   `json:"manufacturer"`
	Memorabilia             bool                     `json:"memorabilia"`
	ModelNumber             string                   `json:"modelNumber"`
	PackageQuantity         int                      `json:"packageQuantity"`
	PartNumber              string                   `json:"partNumber"`
	ReleaseDate             string                   `json:"releaseDate"`
	Size                    string                   `json:"size"`
	Style                   string                   `json:"style"`
	TradeInEligible         bool                     `json:"tradeInEligible"`
	WebsiteDisplayGroup     string                   `json:"websiteDisplayGroup"`
	WebsiteDisplayGroupName string                   `json:"websiteDisplayGroupName"`
}

type CatalogItemRelationshipVariationTheme struct {
	Attributes []string `json:"attributes"`
	Theme      string   `json:"theme"`
}

type CatalogItemRelationship struct {
	ChildASINs     []string                              `json:"childAsins"`
	ParentASINs    string                                `json:"parentAsins"`
	VariationTheme CatalogItemRelationshipVariationTheme `json:"variationTheme"`
	Type           string                                `json:"type"`
}

type CatalogItemRelationships struct {
	MarketplaceId string                    `json:"marketplaceId"`
	Relationships []CatalogItemRelationship `json:"relationships"`
}

type CatalogItemImage struct {
	Height  int    `json:"height"`
	Variant string `json:"variant"`
	URL     string `json:"link"`
	Width   int    `json:"width"`
}

type CatalogItemImages struct {
	MarketplaceId string             `json:"marketplaceId"`
	Images        []CatalogItemImage `json:"images"`
}

type CatalogItemIdentifiersIdentifier struct {
	Type  string `json:"identifierType"`
	Value string `json:"identifier"`
}

type CatalogItemIdentifiers struct {
	MarketplaceId string                             `json:"marketplaceId"`
	Identifiers   []CatalogItemIdentifiersIdentifier `json:"identifiers"`
}

type ClassificationRank struct {
	ClassificiationID string `json:"classificationId"`
	Title             string `json:"title"`
	Rank              int    `json:"rank"`
	Link              string `json:"link"`
}

type DisplayGroupRank struct {
	WebsiteDisplayGroup string `json:"websiteDisplayGroup"`
	Title               string `json:"title"`
	Rank                int    `json:"rank"`
	Link                string `json:"link"`
}

type SalesRank struct {
	MarketplaceId        string               `json:"marketplaceId"`
	ClassificiationRanks []ClassificationRank `json:"classificationRanks"`
	DisplayGroupRanks    []DisplayGroupRank   `json:"displayGroupRanks"`
}

// spapi searchCatalogItem response
type CatalogItem struct {
	ASIN          string                   `json:"asin"`
	Summaries     []CatalogItemSummary     `json:"summaries"`
	Identifiers   []CatalogItemIdentifiers `json:"identifiers"`
	SalesRanks    []SalesRank              `json:"salesRanks"`
	Relationships CatalogItemRelationships `json:"relationships"`
	Images        []CatalogItemImages      `json:"images"`
}

type Refinement struct {
	NumberOfResults int    `json:"numberOfResults"`
	BrandName       string `json:"brandName"`
}

type Refinements struct {
	Brands []Refinement `json:"brands"`
}

type Pagination struct {
	NextToken     string `json:"nextToken"`
	PreviousToken string `json:"previousToken"`
}

type SearchCatalogItemsResponse struct {
	NumberOfResults int           `json:"numberOfResults"`
	Items           []CatalogItem `json:"items"`
	Pagination      Pagination    `json:"pagination"`
	Refinements     Refinements   `json:"refinements"`
}

var (
	IdentifierTypeASIN   = "ASIN"
	IdentifierTypeEAN    = "EAN"
	IdentifierTypeUPC    = "UPC"
	IdentifierTypeSKU    = "SKU"
	IdenfiiterTypeISBN   = "ISBN"
	IdentifierTypeJAN    = "JAN"
	IdentifierTypeGTIN   = "GTIN"
	IdentifierTypeMINSAN = "MINSAN"
)

var (
	IncludedDataAttributes  = "attributes"
	IncludedDataDimensions  = "dimensions"
	IncludedDataIdentifiers = "identifiers"
	IncludedDataImages      = "images"
	IncludedDataProductType = "productTypes"
	IncludedDataRelations   = "relationships"
	IncludedDataSummaries   = "summaries"
	IncludedDataSalesRank   = "salesRanks"
)

// SearchCatalogItems
func (s *Client) SearchCatalogItemsByIdentifer(ctx context.Context, includedData []string, identifers []string, identiferType string) (*SearchCatalogItemsResponse, error) {
	qs := url.Values{}
	qs.Set("identifiers", strings.Join(identifers, ","))
	qs.Set("identifiersType", identiferType)
	qs.Set("marketplaceIds", s.Marketplace.ID)
	qs.Set("pageSize", "20")
	qs.Set("includedData", strings.Join(includedData, ","))

	return s.paginate(ctx, qs, -1)
}

func (s *Client) SearchCatalogItemsByKeyword(ctx context.Context, includedData []string, keywords []string, pageLimit int) (*SearchCatalogItemsResponse, error) {
	qs := url.Values{}
	for i, keyword := range keywords {
		keywords[i] = url.QueryEscape(keyword)
	}
	qs.Set("keywords", strings.Join(keywords, ","))
	qs.Set("marketplaceIds", s.Marketplace.ID)
	qs.Set("pageSize", "20")
	qs.Set("includedData", strings.Join(includedData, ","))

	return s.paginate(ctx, qs, pageLimit)
}

func (s *Client) SearchCatalogItemsByKeywordAndBrand(ctx context.Context, includedData []string, keywords, brandNames []string, pageLimit int) (*SearchCatalogItemsResponse, error) {
	qs := url.Values{}
	for i, brandName := range brandNames {
		brandNames[i] = url.QueryEscape(brandName)
	}
	for i, keyword := range keywords {
		keywords[i] = url.QueryEscape(keyword)
	}
	qs.Set("keywords", strings.Join(keywords, ","))
	qs.Set("brandNames", strings.Join(brandNames, ","))
	qs.Set("marketplaceIds", s.Marketplace.ID)
	qs.Set("pageSize", "20")
	qs.Set("includedData", strings.Join(includedData, ","))

	return s.paginate(ctx, qs, pageLimit)
}

func (s *Client) paginate(ctx context.Context, qs url.Values, limit int) (*SearchCatalogItemsResponse, error) {
	resp, err := s.searchCatalogItems(ctx, qs)
	if err != nil {
		return nil, err
	}

	if resp.Pagination.NextToken != "" {
		nextPage := resp
		i := 0
		for {
			if i > limit && limit != -1 {
				break
			}
			i++
			qs.Set("pageToken", nextPage.Pagination.NextToken)
			nextPage, err = s.searchCatalogItems(ctx, qs)
			if err != nil {
				if v, ok := err.(Error); ok {
					if len(v.Errors) > 0 && v.Errors[0].Code == "InvalidInput" {
						return resp, nil
					}
				}
				return nil, err
			}

			resp.Items = append(resp.Items, nextPage.Items...)
			if nextPage.Pagination.NextToken == "" {
				break
			}
		}
	}

	return resp, nil
}

func (s *Client) searchCatalogItems(ctx context.Context, qs url.Values) (*SearchCatalogItemsResponse, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     s.Marketplace.Endpoint,
		Path:     "/catalog/2022-04-01/items",
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

	var resp SearchCatalogItemsResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error decoding spapi response: %w", err)
	}

	return &resp, nil
}
