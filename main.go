package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/hekmon/btfsautohostprice/coinmarketcap"
)

const (
	apiKeyEnvVarName     = "COINMARKETCAP_APIKEY"
	BTFSTargetEnvVarName = "BTFS_TARGET"
	amountEnvVarName     = "TERABYTEMONTH_USD"
	USDID                = 2781
	BTTID                = 3718
)

var (
	coinmarketcapAPIKey string
	btfsTarget          = "http://127.0.0.1:5001"
	amount              float64
)

func main() {
	// Get env
	if coinmarketcapAPIKey = os.Getenv(apiKeyEnvVarName); coinmarketcapAPIKey == "" {
		fmt.Printf("coinmarketcap API Key must be provided thru %s env var\n", apiKeyEnvVarName)
		os.Exit(1)
	}
	if tmp := os.Getenv(BTFSTargetEnvVarName); tmp != "" {
		if _, err := url.Parse(tmp); err != nil {
			fmt.Printf("%s value is invalid: %s\n", BTFSTargetEnvVarName, err)
			os.Exit(1)
		}
		btfsTarget = tmp
	} else {
		fmt.Printf("no custom BTFS API target set thru %s env var, using default: %s\n", BTFSTargetEnvVarName, btfsTarget)
	}
	if tmp := os.Getenv(amountEnvVarName); tmp == "" {
		fmt.Printf("%s env var must be set\n", amountEnvVarName)
		os.Exit(1)
	} else {
		var err error
		if amount, err = strconv.ParseFloat(tmp, 64); err != nil {
			fmt.Printf("amount can not be converted to float64: %s\n", err)
			os.Exit(1)
		}
	}

	// Get pricing
	market := coinmarketcap.New(APIKey)
	quotes, _, err := market.PriceConversion(context.Background(), coinmarketcap.PriceConversionParams{
		Amount:     amount,
		ID:         USDID,
		ConvertIDs: []int{BTTID},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	bttquote, found := quotes.Quote[BTTID]
	if !found {
		fmt.Printf("can not find BTT id %d within the %d returned quotes\n", BTTID, len(quotes.Quote))
		os.Exit(2)
	}
	fmt.Printf("%f $ is worth of %f BTT at %v\n", amount, bttquote.Price, bttquote.LastUpdated)
	fmt.Printf("as a user push 3 times the amount of data on the network for redundancy, if we want a user to be able to store 1TB for a month for %f$\n", amount)
	fmt.Printf("we need to set the host price for 1TB/month at %f BTT, this amount of BTT is equivalent to the value %d for the BTFS API\n", bttquote.Price/3, BTFSPriceConvert(bttquote.Price/3))

	// Update host
	if err = UpdateHostPrice(bttquote.Price / 3); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println("host pricing updated")
}

func UpdateHostPrice(tokens float64) (err error) {
	url := fmt.Sprintf("%s/api/v1/storage/announce?host-storage-price=%d", btfsTarget, BTFSPriceConvert(tokens))
	resp, err := http.Post(url, "application/json", nil) // mimiq web UI request
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return
}

func BTFSPriceConvert(tokens float64) (param int) {
	return int(math.Round(1e6 * tokens / 30 / 1024))
}
