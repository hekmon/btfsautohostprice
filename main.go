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
	// Config env var names
	apiKeyEnvVarName     = "COINMARKETCAP_APIKEY"
	BTFSTargetEnvVarName = "BTFS_TARGET"
	amountEnvVarName     = "TERABYTEMONTH_USD"
	// CMC IDs
	USDID = 2781
	BTTID = 3718
)

var (
	// Config
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
	if tmp := os.Getenv(amountEnvVarName); tmp != "" {
		var err error
		if amount, err = strconv.ParseFloat(tmp, 64); err != nil {
			fmt.Printf("amount can not be converted to float64: %s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("%s env var must be set\n", amountEnvVarName)
		os.Exit(1)
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
		fmt.Printf("can not find BTT id %d within the %d returned quotes:\n%+v\n", BTTID, len(quotes.Quote), quotes.Quote)
		os.Exit(2)
	}
	hostPrice := bttquote.Price / 3
	fmt.Printf("%0.2f USD is worth %f BTT at %v: with the 3x network redundancy, a host price must be %f BTT for a user to store 1TB/month on the network at this price\n",
		amount, bttquote.Price, bttquote.LastUpdated, hostPrice)
	// Update host
	if err = updateHostPrice(hostPrice); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println("host pricing updated")
}

func updateHostPrice(tokens float64) (err error) {
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/storage/announce?host-storage-price=%d", btfsTarget, BTFSPriceConvert(tokens)),
		"application/json", nil) // mimiq web UI request
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
