package datasources

import (
	"encoding/json"
	"errors"
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"gobacktrader/events"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	fmpBaseURL = "https://fmpcloud.io/api/v3/historical-price-full/"
	fmpURLTail = "{STOCK}?from={START_DATE}&to={END_DATE}&apikey={API_KEY}"
)

// FmpCloudResponse defines the json response from fmpcloud.
// See the documentation at https://fmpcloud.io/documentation
type FmpCloudResponse struct {
	Symbol     string `json:"symbol"`
	Historical []struct {
		Date             string  `json:"date"`
		Open             float64 `json:"open"`
		High             float64 `json:"high"`
		Low              float64 `json:"low"`
		Close            float64 `json:"close"`
		AdjClose         float64 `json:"adjClose"`
		Volume           float64 `json:"volume"`
		UnadjustedVolume float64 `json:"unadjustedVolume"`
		Change           float64 `json:"change"`
		ChangePercent    float64 `json:"changePercent"`
		Vwap             float64 `json:"vwap"`
		Label            string  `json:"label"`
		ChangeOverTime   float64 `json:"changeOverTime"`
	} `json:"historical"`
}

// FmpCloudQuery defines the query details when scraping data from fmpcloud.io
type FmpCloudQuery struct {
	Query
}

// NewFmpCloudQuery returns a new instance of FmpCloudQuery.
func NewFmpCloudQuery(targetAsset asset.IAssetReadOnly, startDate time.Time, endDate time.Time) FmpCloudQuery {
	return FmpCloudQuery{
		Query: Query{
			targetAsset: targetAsset,
			startDate:   startDate,
			endDate:     endDate,
			apiKey:      "demo",
		},
	}
}

// GetURL returns the formatted query URL.
func (q FmpCloudQuery) GetURL() string {
	replacements := map[string]string{
		"{STOCK}":      q.GetTicker(),
		"{START_DATE}": q.startDate.Format("2006-01-02"),
		"{END_DATE}":   q.endDate.Format("2006-01-02"),
		"{API_KEY}":    q.apiKey,
	}
	urlTail := btutil.ReplaceStrings(fmpURLTail, replacements)
	return fmpBaseURL + urlTail
}

// Run returns the query response from fmpcloud.
func (q FmpCloudQuery) Run() (FmpCloudResponse, error) {
	var fmpCloudResponse FmpCloudResponse

	response, err := http.Get(q.GetURL())
	if err != nil {
		return fmpCloudResponse, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmpCloudResponse, err
	}

	err = json.Unmarshal(data, &fmpCloudResponse)
	return fmpCloudResponse, err
}

// GenerateEvents returns all price events from the fmpcloud response.
func (q FmpCloudQuery) GenerateEvents() ([]events.AssetPriceEvent, error) {
	var priceEvents []events.AssetPriceEvent

	targetAsset, ok := q.GetAsset().(asset.IAssetWriteOnly)
	if !ok {
		return priceEvents, errors.New("Unable to cast to IAssetWriteOnly")
	}

	fmpCloudResponse, err := q.Run()
	if err != nil {
		return priceEvents, err
	}

	dateLayout := "2006-01-02"
	for _, item := range fmpCloudResponse.Historical {
		datestr, close := item.Date, item.AdjClose
		eventTime, err := time.Parse(dateLayout, datestr)
		if err != nil {
			return priceEvents, err
		}

		price := asset.Price{Float64: close, Valid: true}
		assetPriceEvent := events.NewAssetPriceEvent(targetAsset, eventTime, price)
		priceEvents = append(priceEvents, assetPriceEvent)
	}

	return priceEvents, nil
}
