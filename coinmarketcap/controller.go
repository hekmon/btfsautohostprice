package coinmarketcap

import (
	"net/http"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

func New(apiKey string) *Controller {
	return &Controller{
		apiKey: apiKey,
		client: cleanhttp.DefaultClient(),
	}
}

type Controller struct {
	apiKey string
	client *http.Client
}
