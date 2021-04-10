package coinmarketcap

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

type PriceConversionParams struct {
	Amount     float64   `url:"amount"`            // An amount of currency to convert. Example: 10.43
	ID         int       `url:"id,omitempty"`      // The CoinMarketCap currency ID of the base cryptocurrency or fiat to convert from. Example: 1
	Symbol     string    `url:"symbol,omitempty"`  // Alternatively the currency symbol of the base cryptocurrency or fiat to convert from. Example: "BTC". One "id" or "symbol" is required.
	Time       time.Time `url:"time,omitempty"`    // Optional time to reference historical pricing during conversion. If not passed, the current time will be used. If passed, we'll reference the closest historic values available for this conversion.
	Convert    string    `url:"convert,omitempty"` // Pass up to 120 comma-separated fiat or cryptocurrency symbols to convert the source amount to.
	ConvertIDs []int     `url:"-"`                 // Optionally calculate market quotes by CoinMarketCap ID instead of symbol. This option is identical to convert outside of ID format. Ex: convert_id=1,2781 would replace convert=BTC,USD in your query. This parameter cannot be used when convert is used.
}

func (pcp PriceConversionParams) validate() error {
	if pcp.Amount <= 0 {
		return errors.New("please enter a valid amount")
	}
	if (pcp.ID == 0 && pcp.Symbol == "") || (pcp.ID != 0 && pcp.Symbol != "") {
		return errors.New("please use either ID or Symbol")
	}
	if (pcp.Convert == "" && len(pcp.ConvertIDs) == 0) || (pcp.Convert != "" && len(pcp.ConvertIDs) != 0) {
		return errors.New("please use either Convert or ConvertIDs")
	}
	return nil
}

// PriceConversion maps to https://coinmarketcap.com/api/documentation/v1/#operation/getV1ToolsPriceconversion
func (c *Controller) PriceConversion(ctx context.Context, params PriceConversionParams) (quotes PriceConversionResponse, creditCount int, err error) {
	if err = params.validate(); err != nil {
		err = fmt.Errorf("query parameters invalid: %w", err)
		return
	}
	intermediate := struct {
		PriceConversionParams
		ConvertIDs string `url:"convert_id"`
	}{
		PriceConversionParams: params,
		ConvertIDs:            intSliceToStr(params.ConvertIDs),
	}
	paramsValues, err := query.Values(intermediate)
	if err != nil {
		err = fmt.Errorf("failed to convert params to url values: %w", err)
		return
	}
	if creditCount, err = c.request(ctx, "GET", "/v1/tools/price-conversion", paramsValues, nil, &quotes); err != nil {
		err = fmt.Errorf("requesting endpoint failed: %w", err)
		return
	}
	return
}

func intSliceToStr(input []int) string {
	tmp := make([]string, len(input))
	for index, integer := range input {
		tmp[index] = strconv.Itoa(integer)
	}
	return strings.Join(tmp, ",")
}

type PriceConversionResponse struct {
	Symbol      string                               `json:"symbol"`
	ID          int                                  `json:"id"`
	Name        string                               `json:"name"`
	Amount      float64                              `json:"amount"`
	LastUpdated time.Time                            `json:"last_updated"`
	Quote       map[int]PriceConversionResponseQuote `json:"quote"`
}

type PriceConversionResponseQuote struct {
	Price       float64   `json:"price"`
	LastUpdated time.Time `json:"last_updated"`
}
