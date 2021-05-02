package datasources

import (
	"encoding/json"
	"errors"
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"gobacktrader/events"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	avBaseURL = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol="
	avURLtail = "{STOCK}&outputsize=full&apikey={API_KEY}"
)

// AlphaVantageResponse defines the json response from
// alphavantage on daily adjusted prices.
// See the docs at https://www.alphavantage.co/documentation/
type AlphaVantageResponse struct {
	MetaData struct {
		Information   string `json:"1. Information"`
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
		OutputSize    string `json:"4. Output Size"`
		TimeZone      string `json:"5. Time Zone"`
	} `json:"Meta Data"`
	TimeSeriesDaily map[string]struct {
		Open             string `json:"1. open"`
		High             string `json:"2. high"`
		Low              string `json:"3. low"`
		Close            string `json:"4. close"`
		AdjustedClose    string `json:"5. adjusted close"`
		Volume           string `json:"6. volume"`
		DividendAmount   string `json:"7. dividend amount"`
		SplitCoefficient string `json:"8. split coefficient"`
	} `json:"Time Series (Daily)"`
}

// AlphaVantageQuery defines the query details when scraping data from alphavantage.
type AlphaVantageQuery struct {
	Query
}

// NewAlphaVantageQuery returns a new instance of AlphaVantageQuery.
func NewAlphaVantageQuery(targetAsset asset.IAssetReadOnly, startDate time.Time, endDate time.Time) AlphaVantageQuery {
	return AlphaVantageQuery{
		Query: Query{
			targetAsset: targetAsset,
			startDate:   startDate,
			endDate:     endDate,
			apiKey:      "demo",
		},
	}
}

// GetURL returns the formatted query URL.
func (q AlphaVantageQuery) GetURL() string {
	replacements := map[string]string{
		"{STOCK}":   q.GetTicker(),
		"{API_KEY}": q.apiKey,
	}

	urlTail := btutil.ReplaceStrings(avURLtail, replacements)
	return avBaseURL + urlTail
}

// Run returns the query reponse from alphavantage.
func (q AlphaVantageQuery) Run() (AlphaVantageResponse, error) {
	var alphaVantageResponse AlphaVantageResponse

	response, err := http.Get(q.GetURL())
	if err != nil {
		return alphaVantageResponse, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return alphaVantageResponse, err
	}

	err = json.Unmarshal(data, &alphaVantageResponse)
	return alphaVantageResponse, err
}

// GenerateEvents returns all price events from the fmpcloud response.
func (q AlphaVantageQuery) GenerateEvents() ([]events.AssetPriceEvent, error) {
	var priceEvents []events.AssetPriceEvent

	targetAsset, ok := q.GetAsset().(asset.IAssetWriteOnly)
	if !ok {
		return priceEvents, errors.New("Unable to cast to IAssetWriteOnly")
	}

	avResponse, err := q.Run()
	if err != nil {
		return priceEvents, err
	}

	dateLayout := "2006-01-02"
	for datestr, item := range avResponse.TimeSeriesDaily {
		eventTime, err := time.Parse(dateLayout, datestr)
		if err != nil {
			return priceEvents, err
		}

		include := true
		if eventTime.Before(q.startDate) || eventTime.After(q.endDate) {
			include = false
		}

		if include {
			closestr := item.AdjustedClose
			close, err := strconv.ParseFloat(closestr, 64)
			if err != nil {
				return priceEvents, err
			}
			price := asset.Price{Float64: close, Valid: true}
			assetPriceEvent := events.NewAssetPriceEvent(targetAsset, eventTime, price)
			priceEvents = append(priceEvents, assetPriceEvent)
		}
	}

	return priceEvents, nil
}
