package spapi

type Marketplace struct {
	ID       string
	Endpoint string
}

var (
	MarketplaceUS = Marketplace{
		ID:       "ATVPDKIKX0DER",
		Endpoint: EndpointNorthAmerica,
	}
	MarketplaceCA = Marketplace{
		ID:       "A2EUQ1WTGCTBG2",
		Endpoint: EndpointNorthAmerica,
	}
	MarketplaceMX = Marketplace{
		ID:       "A1AM78C64UM0Y8",
		Endpoint: EndpointNorthAmerica,
	}
	MarketplaceGB = Marketplace{
		ID:       "A1F83G8C2ARO7P",
		Endpoint: EndpointEurope,
	}
	MarketplaceDE = Marketplace{
		ID:       "A1PA6795UKMFR9",
		Endpoint: EndpointEurope,
	}
	MarketplaceFR = Marketplace{
		ID:       "A13V1IB3VIYZZH",
		Endpoint: EndpointEurope,
	}
	MarketplaceIT = Marketplace{
		ID:       "APJ6JRA9NG5V4",
		Endpoint: EndpointEurope,
	}
	MarketplaceES = Marketplace{
		ID:       "A1RKKUPIHCS9HS",
		Endpoint: EndpointEurope,
	}
	MarketplaceNL = Marketplace{
		ID:       "A1805IZSGTT6HS",
		Endpoint: EndpointEurope,
	}
	MarketplaceSE = Marketplace{
		ID:       "A2NODRKZP88ZB9",
		Endpoint: EndpointEurope,
	}
	MarketplaceTR = Marketplace{
		ID:       "A33AVAJ2PDY3EV",
		Endpoint: EndpointEurope,
	}
	MarketplaceAE = Marketplace{
		ID:       "A2VIGQ35RCS4UG",
		Endpoint: EndpointFarEast,
	}
	MarketplaceJP = Marketplace{
		ID:       "A1VC38T7YXB528",
		Endpoint: EndpointFarEast,
	}
)

var MarketplaceMap = map[string]Marketplace{
	"US": MarketplaceUS,
	"CA": MarketplaceCA,
	"MX": MarketplaceMX,
	"GB": MarketplaceGB,
	"DE": MarketplaceDE,
	"FR": MarketplaceFR,
	"IT": MarketplaceIT,
	"ES": MarketplaceES,
	"NL": MarketplaceNL,
	"SE": MarketplaceSE,
	"TR": MarketplaceTR,
	"AE": MarketplaceAE,
	"JP": MarketplaceJP,
}
