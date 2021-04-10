package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hekmon/btfsautohostprice/coinmarketcap"
)

func main() {
	market := coinmarketcap.New(APIKey)
	ids, creditCount, err := market.GetFiatCoinMarketCapIDMap(context.Background(), coinmarketcap.GetFiatCoinMarketCapIDMapParams{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Credit count:", creditCount)
	for id, data := range ids {
		if data.Symbol == "USD" {
			fmt.Printf("%s (%s) (%s) has ID %d\n", data.Name, data.Symbol, data.Sign, id)
		}
	}
}
