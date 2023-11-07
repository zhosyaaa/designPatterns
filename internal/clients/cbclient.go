package clients

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"sync"
)

type ExchangeRates struct {
	Result            string             `json:"result"`
	Provider          string             `json:"provider"`
	Documentation     string             `json:"documentation"`
	TermsOfUse        string             `json:"terms_of_use"`
	TimeLastUpdateUTC string             `json:"time_last_update_utc"`
	TimeNextUpdateUTC string             `json:"time_next_update_utc"`
	BaseCode          string             `json:"base_code"`
	Rates             map[string]float64 `json:"rates"`
}

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
	var exchangeRates ExchangeRates
	err = json.Unmarshal([]byte(body), &exchangeRates)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing JSON")
		return nil, err
	}
	currencyData = exchangeRates.Rates
	return currencyData, nil
}
