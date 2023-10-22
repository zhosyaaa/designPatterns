package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CurrencyClient struct {
	APIKey        string
	exchangeRates map[string]interface{}
}

func NewCurrencyClient(apiKey string) *CurrencyClient {
	return &CurrencyClient{
		APIKey: apiKey,
	}
}

func (c *CurrencyClient) GetExchangeRates() (map[string]interface{}, error) {
	API_URL := "https://api.apilayer.com/exchangerates_data/latest" // URL API для получения данных о курсах валют
	client := &http.Client{}
	req, err := http.NewRequest("GET", API_URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Error: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	fmt.Println("data: ", data)
	exchangeRates := data["rates"].(map[string]interface{})
	fmt.Println("exchangeRates: ", exchangeRates)
	c.exchangeRates = exchangeRates
	return exchangeRates, nil
}
