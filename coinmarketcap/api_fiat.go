package coinmarketcap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/go-querystring/query"
)

// GetFiatCoinMarketCapIDMapParams contains query parameters for https://coinmarketcap.com/api/documentation/v1/#operation/getV1FiatMap
type GetFiatCoinMarketCapIDMapParams struct {
	Start         int    `url:"start,omitempty"`          // >= 1, default: 1, Optionally offset the start (1-based index) of the paginated list of items to return.
	Limit         int    `url:"limit,omitempty"`          // [ 1 .. 5000 ], Optionally specify the number of results to return. Use this parameter and the "start" parameter to determine your own pagination size.
	Sort          string `url:"sort,omitempty"`           // default: "id", valid values: ["name", "id"], What field to sort the list by.
	IncludeMetals bool   `url:"include_metals,omitempty"` // default: false, Pass true to include precious metals.
}

func (cmcidmp GetFiatCoinMarketCapIDMapParams) validate() error {
	if cmcidmp.Start < 0 {
		return errors.New("\"Start\" must be >= 1")
	}
	if cmcidmp.Limit < 0 || cmcidmp.Limit > 5000 {
		return errors.New("\"Limit\" must be [1..5000]")
	}
	switch cmcidmp.Sort {
	case "":
	case "name":
	case "id":
	default:
		return errors.New("\"Sort\" valid values are \"name\" \"id\"")
	}
	return nil
}

// GetFiatCoinMarketCapIDMap maps https://coinmarketcap.com/api/documentation/v1/#operation/getV1FiatMap
func (c *Controller) GetFiatCoinMarketCapIDMap(ctx context.Context, params GetFiatCoinMarketCapIDMapParams) (ids FiatCoinMarketCapIDMap, creditCount int, err error) {
	if err = params.validate(); err != nil {
		err = fmt.Errorf("query parameters invalid: %w", err)
		return
	}
	paramsValues, err := query.Values(params)
	if err != nil {
		err = fmt.Errorf("failed to convert params to url values: %w", err)
		return
	}
	if creditCount, err = c.request(ctx, "GET", "/v1/fiat/map", paramsValues, nil, &ids); err != nil {
		err = fmt.Errorf("requesting endpoint failed: %w", err)
		return
	}
	return
}

type FiatCoinMarketCapIDMap map[int]FiatCoinMarketCapIDMapItem

func (cmcidm *FiatCoinMarketCapIDMap) UnmarshalJSON(data []byte) error {
	var tmp []struct {
		ID int `json:"id"`
		FiatCoinMarketCapIDMapItem
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshall into temporary structure: %w", err)
	}
	*cmcidm = make(FiatCoinMarketCapIDMap, len(tmp))
	for _, item := range tmp {
		(*cmcidm)[item.ID] = item.FiatCoinMarketCapIDMapItem
	}
	return nil
}

type FiatCoinMarketCapIDMapItem struct {
	Name   string `json:"name"`
	Sign   string `json:"sign"`
	Symbol string `json:"symbol"`
}
