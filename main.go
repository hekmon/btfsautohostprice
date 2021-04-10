package main

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/hekmon/btfsautohostprice/coinmarketcap"
)

const (
	USDID = 2781
	BTTID = 3718
)

func main() {
	market := coinmarketcap.New(APIKey)

	// fiatIDs, creditCount, err := market.GetFiatCoinMarketCapIDMap(context.Background(), coinmarketcap.GetFiatCoinMarketCapIDMapParams{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Credit count:", creditCount)
	// for id, data := range fiatIDs {
	// 	if data.Symbol == "USD" {
	// 		fmt.Printf("%s (%s) (%s) has ID %d\n", data.Name, data.Symbol, data.Sign, id)
	// 	}
	// }

	// cryptoIDs, creditCount, err := market.GetCryptoCoinMarketCapIDMap(context.Background(), coinmarketcap.GetCryptoCoinMarketCapIDMapParams{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(len(cryptoIDs))
	// fmt.Println("Credit count:", creditCount)
	// for id, data := range cryptoIDs {
	// 	if data.Symbol == "BTT" {
	// 		fmt.Printf("%s (%s) has ID %d\n", data.Name, data.Symbol, id)
	// 		fmt.Printf("%+v\n", data)
	// 		if data.Platform != nil {
	// 			fmt.Printf("%+v\n", data.Platform)
	// 		}
	// 	}
	// }

	dollars := float64(10)
	quotes, _, err := market.PriceConversion(context.Background(), coinmarketcap.PriceConversionParams{
		Amount:     dollars,
		ID:         USDID,
		ConvertIDs: []int{BTTID},
	})
	if err != nil {
		log.Fatal(err)
	}
	bttquote, found := quotes.Quote[BTTID]
	if !found {
		log.Fatalf("can not find BTT id %d within the %d returned quotes", BTTID, len(quotes.Quote))
	}
	fmt.Printf("%f $ is worth of %f BTT at %v\n", dollars, bttquote.Price, bttquote.LastUpdated)
	fmt.Printf("as a user push 3 times the amount of data on the network for redundancy, if we want a user to be able to store 1TB for a month for %f$\n", dollars)
	fmt.Printf("we need to set the host price for 1TB/month at %f BTT, this amount of BTT is equivalent to the value %d for the BTFS API\n", bttquote.Price/3, BTFSPriceConvert(bttquote.Price/3))
}

func BTFSPriceConvert(tokens float64) (param int) {
	return int(math.Round(1e6 * tokens / 30 / 1024))
}
