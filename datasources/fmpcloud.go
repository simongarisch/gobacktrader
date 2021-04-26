package datasources

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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
	stock     string
	startDate time.Time
	endDate   time.Time
	apiKey    string
}

// NewFmpCloudQuery returns a new instance of FmpCloudQuery.
func NewFmpCloudQuery(stock string, startDate time.Time, endDate time.Time) FmpCloudQuery {
	return FmpCloudQuery{
		stock:     stock,
		startDate: startDate,
		endDate:   endDate,
		apiKey:    "demo",
	}
}

// GetStock returns the query stock code.
func (q FmpCloudQuery) GetStock() string {
	return q.stock
}

// GetStartDate returns the query start date.
func (q FmpCloudQuery) GetStartDate() time.Time {
	return q.startDate
}

// GetEndDate returns the query end date.
func (q FmpCloudQuery) GetEndDate() time.Time {
	return q.endDate
}

// GetAPIKey gets the query api key.
func (q FmpCloudQuery) GetAPIKey() string {
	return q.apiKey
}

// SetAPIKey sets the query api key.
func (q *FmpCloudQuery) SetAPIKey(apiKey string) *FmpCloudQuery {
	q.apiKey = apiKey
	return q
}

// GetURL returns the formatted query URL.
func (q FmpCloudQuery) GetURL() string {
	replacements := map[string]string{
		"{STOCK}":      q.stock,
		"{START_DATE}": q.startDate.Format("2006-01-02"),
		"{END_DATE}":   q.endDate.Format("2006-01-02"),
		"{API_KEY}":    q.apiKey,
	}
	urlTail := fmpURLTail
	for oldString, newString := range replacements {
		urlTail = strings.ReplaceAll(urlTail, oldString, newString)
	}
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
