package backtest

import (
	"fmt"
	"gobacktrader/asset"
	"gobacktrader/broker"
	"gobacktrader/btutil"
	"gobacktrader/events"
	"time"
)

// Backtest collects the assets we want to test and
// records history as backtest events are processed.
type Backtest struct {
	portfolios []*asset.Portfolio
	assets     []asset.IAssetReadOnly
	events     events.Events
	strategy   IStrategy
}

// NewBacktest returns a new Backtest instance.
func NewBacktest(strategy IStrategy) Backtest {
	return Backtest{strategy: strategy}
}

// codeRegistered checks if a code is registered either
// for an existing portfolio or asset.
// All codes should be unique.
func (backtest Backtest) codeRegistered(code string) bool {
	code = btutil.CleanString(code)
	// check if this exists as a portfolio code
	for _, portfolio := range backtest.portfolios {
		if code == btutil.CleanString(portfolio.GetCode()) {
			return true
		}
	}

	// or an asset ticker
	for _, asset := range backtest.assets {
		if code == btutil.CleanString(asset.GetTicker()) {
			return true
		}
	}

	return false
}

// RegisterPortfolio registers a portfolio within our backtest.
func (backtest *Backtest) RegisterPortfolio(p *asset.Portfolio) error {
	if backtest.HasPortfolio(p) {
		return nil // portfolio is already registered
	}

	// set the default executing broker if none is provided
	if p.GetBroker() == nil {
		executingBroker := broker.NewBroker(
			broker.NewNoCharges(),
			broker.NewFillAtLast(),
		)
		p.SetBroker(executingBroker)
	}

	portfolioCode := p.GetCode()
	if backtest.codeRegistered(portfolioCode) {
		return fmt.Errorf("portfolio code '%s' is already in use and needs to be unique", portfolioCode)
	}

	backtest.portfolios = append(backtest.portfolios, p)
	return nil
}

// RegisterAsset registers an asset within our backtest.
func (backtest *Backtest) RegisterAsset(a asset.IAssetReadOnly) error {
	if backtest.HasAsset(a) {
		return nil // asset is already registered
	}
	assetTicker := a.GetTicker()
	if backtest.codeRegistered(assetTicker) {
		return fmt.Errorf("asset ticker '%s' is already in use and needs to be unique", assetTicker)
	}

	backtest.assets = append(backtest.assets, a)
	return nil
}

// HasPortfolio returns true if the backtest has a specific portfolio registered, false otherwise.
func (backtest *Backtest) HasPortfolio(p *asset.Portfolio) bool {
	for _, registeredPortfolio := range backtest.portfolios {
		if registeredPortfolio == p {
			return true
		}
	}
	return false
}

// HasAsset returns true if the backtest has a specific asset registered, false otherwise.
func (backtest *Backtest) HasAsset(a asset.IAssetReadOnly) bool {
	for _, registeredAsset := range backtest.assets {
		if registeredAsset == a {
			return true
		}
	}
	return false
}

// AddEvent adds a new event to the backtest events collection.
func (backtest *Backtest) AddEvent(event events.IEvent) {
	backtest.events.Add(event)
}

// Run will execute our backtest.
func (backtest *Backtest) Run() error {
	for { // while we have events to process
		if backtest.events.IsEmpty() {
			break
		}

		// process events for the next time step
		var currentTime time.Time
		var eventsToProcess []events.IEvent
		eventsToProcess, err := backtest.events.FetchNextGroup()
		if err != nil {
			return err
		}

		currentTime = eventsToProcess[0].GetTime()
		for _, event := range eventsToProcess {
			if err := event.Process(); err != nil {
				return err
			}
		}

		// now that events have been processed for this time
		// we'll check to see if our strategy generates trades
		trades, err := backtest.strategy.GenerateTrades()
		if err != nil {
			return err
		}
		for _, trade := range trades {
			if _, err := trade.Execute(); err != nil {
				return err
			}
		}

		// once all events have been processed for this step
		// then take snapshots
		for _, asset := range backtest.assets {
			asset.TakeSnapshot(currentTime, asset)
		}
		for _, portfolio := range backtest.portfolios {
			portfolio.TakeSnapshot(currentTime)
		}
	}

	return nil
}
