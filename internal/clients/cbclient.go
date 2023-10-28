package clients

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"sync"
)

type CurrencyClient struct {
}

func NewCurrencyClient() *CurrencyClient {
	return &CurrencyClient{}
}

var (
	currencyData map[string]float64
	mu           sync.Mutex
)

func (c *CurrencyClient) GetExchangeRates() (map[string]float64, error) {
	mu.Lock()
	defer mu.Unlock()
	API_URL := "https://open.er-api.com/v6/latest/USD"
	resp, err := http.Get(API_URL)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching data")
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading data")
		return nil, err
	}
	var data map[string]float64
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error().Err(err).Msg("Error parsing data")
		return nil, err
	}
	currencyData = data
	return currencyData, nil
}
