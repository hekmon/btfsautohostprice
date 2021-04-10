package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hekmon/btfsautohostprice/coinmarketcap"
)

func main() {
	market := coinmarketcap.New(APIKey)

	fiatIDs, creditCount, err := market.GetFiatCoinMarketCapIDMap(context.Background(), coinmarketcap.GetFiatCoinMarketCapIDMapParams{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Credit count:", creditCount)
	for id, data := range fiatIDs {
		if data.Symbol == "USD" {
			fmt.Printf("%s (%s) (%s) has ID %d\n", data.Name, data.Symbol, data.Sign, id)
		}
	}

	cryptoIDs, creditCount, err := market.GetCryptoCoinMarketCapIDMap(context.Background(), coinmarketcap.GetCryptoCoinMarketCapIDMapParams{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(cryptoIDs))
	fmt.Println("Credit count:", creditCount)
	for id, data := range cryptoIDs {
		if data.Symbol == "BTT" {
			fmt.Printf("%s (%s) has ID %d\n", data.Name, data.Symbol, id)
			fmt.Printf("%+v\n", data)
			if data.Platform != nil {
				fmt.Printf("%+v\n", data.Platform)
			}
		}
	}
}
