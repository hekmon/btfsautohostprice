package coinmarketcap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// GetCryptoCoinMarketCapIDMapParams contains query parameters for https://coinmarketcap.com/api/documentation/v1/#operation/getV1CryptocurrencyMap
type GetCryptoCoinMarketCapIDMapParams struct {
	ListingStatus string `url:"listing_status,omitempty"` // default: "active", Only active cryptocurrencies are returned by default. Pass inactive to get a list of cryptocurrencies that are no longer active. Pass untracked to get a list of cryptocurrencies that are listed but do not yet meet methodology requirements to have tracked markets available. You may pass one or more comma-separated values.
	Start         int    `url:"start,omitempty"`          // >= 1, default: 1, Optionally offset the start (1-based index) of the paginated list of items to return.
	Limit         int    `url:"limit,omitempty"`          // [ 1 .. 5000 ], Optionally specify the number of results to return. Use this parameter and the "start" parameter to determine your own pagination size.
	Sort          string `url:"sort,omitempty"`           // default: "id", valid values: ["cmc_rank", "id"], What field to sort the list by.
	Symbol        string `url:"symbol,omitempty"`         // Optionally pass a comma-separated list of cryptocurrency symbols to return CoinMarketCap IDs for. If this option is passed, other options will be ignored.
	Aux           string `url:"aux,omitempty"`            // default: "platform,first_historical_data,last_historical_data,is_active", Optionally specify a comma-separated list of supplemental data fields to return. Pass "platform,first_historical_data,last_historical_data,is_active,status" to include all auxiliary fields.
}

func (cmcidmp GetCryptoCoinMarketCapIDMapParams) validate() error {
	if cmcidmp.ListingStatus != "" {
		for index, value := range strings.Split(cmcidmp.ListingStatus, ",") {
			switch value {
			case "active":
			case "inactive":
			case "untracked":
			default:
				return fmt.Errorf("listing status invalid at index %d: %s", index, value)
			}
		}
	}
	if cmcidmp.Start < 0 {
		return errors.New("\"Start\" must be >= 1")
	}
	if cmcidmp.Limit < 0 || cmcidmp.Limit > 5000 {
		return errors.New("\"Limit\" must be [1..5000]")
	}
	switch cmcidmp.Sort {
	case "":
	case "cmc_rank":
	case "id":
	default:
		return errors.New("\"Sort\" valid values are \"cmc_rank\" \"id\"")
	}
	if cmcidmp.Aux != "" {
		for index, value := range strings.Split(cmcidmp.Aux, ",") {
			switch value {
			case "platform":
			case "first_historical_data":
			case "last_historical_data":
			case "is_active":
			case "status":
			default:
				return fmt.Errorf("aux date field invalid at index %d: %s", index, value)
			}
		}
	}
	return nil
}

// GetCryptoCoinMarketCapIDMap maps https://coinmarketcap.com/api/documentation/v1/#operation/getV1CryptocurrencyMap
func (c *Controller) GetCryptoCoinMarketCapIDMap(ctx context.Context, params GetCryptoCoinMarketCapIDMapParams) (ids CryptoCoinMarketCapIDMap, creditCount int, err error) {
	if err = params.validate(); err != nil {
		err = fmt.Errorf("query parameters invalid: %w", err)
		return
	}
	paramsValues, err := query.Values(params)
	if err != nil {
		err = fmt.Errorf("failed to convert params to url values: %w", err)
		return
	}
	if creditCount, err = c.request(ctx, "GET", "/v1/cryptocurrency/map", paramsValues, nil, &ids); err != nil {
		err = fmt.Errorf("requesting endpoint failed: %w", err)
		return
	}
	return
}

type CryptoCoinMarketCapIDMap map[int]CryptoCoinMarketCapIDMapItem

func (cmcidm *CryptoCoinMarketCapIDMap) UnmarshalJSON(data []byte) error {
	var tmp []struct {
		ID       int `json:"id"`
		IsActive int `json:"is_active"`
		CryptoCoinMarketCapIDMapItem
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshall into temporary structure: %w", err)
	}
	*cmcidm = make(CryptoCoinMarketCapIDMap, len(tmp))
	for _, item := range tmp {
		if item.IsActive == 1 {
			item.CryptoCoinMarketCapIDMapItem.Active = true
		}
		(*cmcidm)[item.ID] = item.CryptoCoinMarketCapIDMapItem
	}
	return nil
}

type CryptoCoinMarketCapIDMapItem struct {
	Name                string                            `json:"name"`
	Symbol              string                            `json:"symbol"`
	Slug                string                            `json:"slug"`
	Active              bool                              `json:"active"`
	FirstHistoricalData time.Time                         `json:"first_historical_data"`
	LastHistoricalData  time.Time                         `json:"last_historical_data"`
	Platform            *CryptoCoinMarketCapIDMapPlatform `json:"platform"`
}

type CryptoCoinMarketCapIDMapPlatform struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Symbol       string `jaon:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}
