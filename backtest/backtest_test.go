package backtest

import (
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"gobacktrader/compliance"
	"gobacktrader/datasources"
	"gobacktrader/events"
	"gobacktrader/trade"
	"os"
	"testing"
	"time"
)

func TestRegisterPortfolio(t *testing.T) {
	backtest := NewBacktest(myStrategy{})

	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Fatalf("Error in asset init - %s", err)
	}

	if backtest.HasPortfolio(portfolio) {
		t.Error("Backtest should not have portfolio registered.")
	}
	if backtest.HasAsset(stock) {
		t.Error("Backtest should not have asset registered.")
	}

	err1 = backtest.RegisterPortfolio(portfolio)
	err2 = backtest.RegisterAsset(stock)
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Fatalf("Error in backtest.Register - %s", err)
	}

	if !backtest.HasPortfolio(portfolio) {
		t.Error("Backtest should have portfolio registered.")
	}
	if !backtest.HasAsset(stock) {
		t.Error("Backtest should have asset registered.")
	}

	// we should be able to register the same asset and portfolio
	// again without issue
	err1 = backtest.RegisterPortfolio(portfolio)
	err2 = backtest.RegisterAsset(stock)
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in backtest.Register - %s", err)
	}

	// try to register a different portfolio and asset with the same codes
	portfolio2, err := asset.NewPortfolio("XXX", "USD")
	if err != nil {
		t.Errorf("Error in NewPortfolio - %s", err)
	}

	stock2, err := asset.NewStock("ZZB AU", "USD")
	if err != nil {
		t.Errorf("Error in NewStock - %s", err)
	}

	err1 = backtest.RegisterPortfolio(portfolio2)
	err2 = backtest.RegisterAsset(stock2)
	errStr1 := btutil.GetErrorString(err1)
	errStr2 := btutil.GetErrorString(err2)
	if errStr1 != "portfolio code 'XXX' is already in use and needs to be unique" {
		t.Errorf("Unexpected error string '%s'", errStr1)
	}
	if errStr2 != "asset ticker 'ZZB AU' is already in use and needs to be unique" {
		t.Errorf("Unexpected error string '%s'", errStr2)
	}
}

func TestBacktestBasicStrategy(t *testing.T) {
	// initialise our portfolio and assets
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	cash, err3 := asset.NewCash("AUD")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Fatalf("Error in asset creation - %s", err)
	}

	// create our trading strategy
	generateTrades := func() ([]*trade.Trade, error) {
		var trades []*trade.Trade
		stockPrice := stock.GetPrice()
		if stockPrice.Valid {
			if stockPrice.Float64 <= 2 {
				newTrade := trade.NewTrade(portfolio, stock, 100)
				trades = append(trades, newTrade)
			}
		}
		return trades, nil
	}

	strategy := NewStrategy(generateTrades)

	// seed the portfolio
	portfolio.Transfer(cash, 1000)

	// create our backtest instance, register assets
	backtest := NewBacktest(strategy)
	backtest.RegisterPortfolio(portfolio)
	backtest.RegisterAsset(stock)
	backtest.RegisterAsset(cash)

	// create events to process
	t1 := time.Date(2021, time.March, 13, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2021, time.March, 14, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2021, time.March, 15, 0, 0, 0, 0, time.UTC)

	e1 := events.NewAssetPriceEvent(stock, t1, asset.Price{Float64: 2.10, Valid: true})
	e2 := events.NewAssetPriceEvent(stock, t2, asset.Price{Float64: 2.00, Valid: true})
	e3 := events.NewAssetPriceEvent(stock, t3, asset.Price{Float64: 2.50, Valid: true})

	for _, event := range []events.IEvent{&e1, &e2, &e3} {
		backtest.AddEvent(event)
	}

	// we'd previously set our strategy to buy 100 shares whenever the stock price went <= $2
	// the backtest should start with a stock position of zero at t1, buy 100 shares at t2
	// and hold at t3.
	err := backtest.Run()
	if err != nil {
		t.Fatalf("Error in backtest.Run() - %s", err)
	}

	// check the portfolio units and value
	stockPosition := portfolio.GetUnits(stock)
	cashPosition := portfolio.GetUnits(cash)
	if stockPosition != 100 {
		t.Errorf("Unexpected stock position - wanted 100 got %0.2f", stockPosition)
	}
	if cashPosition != 800 {
		t.Errorf("Unexpected cash position - wanted 800, got %0.2f", cashPosition)
	}
	portfolioValue, err := portfolio.GetValue()
	if err != nil {
		t.Errorf("Error in portfolio.GetValue() - %s", err)
	}
	if !portfolioValue.Valid {
		t.Error("Expecting a valid portfolio value")
	}
	if portfolioValue.Float64 != 1050 { // 800 cash + 100 * 2.50 in stock
		t.Errorf("Unexpected portfolio value - wanted 1050, got %0.2f", portfolioValue.Float64)
	}

	// check the portfolio history
	portfolioHistory := portfolio.GetHistory()
	psnap1, _ := portfolioHistory[t1]
	psnap2, _ := portfolioHistory[t2]
	psnap3, _ := portfolioHistory[t3]

	if !psnap1.GetTime().Equal(t1) {
		t.Error("Unexpected time for psnap1")
	}
	if psnap1.GetValue().Float64 != 1000 {
		t.Errorf("Unexpected value for psnap1 - wanted 1000, got %0.2f", psnap1.GetValue().Float64)
	}

	if !psnap2.GetTime().Equal(t2) {
		t.Error("Unexpected time for psnap2")
	}
	if psnap2.GetValue().Float64 != 1000 {
		t.Errorf("Unexpected value for psnap2 - wanted 1000, got %0.2f", psnap2.GetValue().Float64)
	}

	if !psnap3.GetTime().Equal(t3) {
		t.Error("Unexpected time for psnap3")
	}
	if psnap3.GetValue().Float64 != 1050 {
		t.Errorf("Unexpected value for psnap3 - wanted 1050, got %0.2f", psnap3.GetValue().Float64)
	}

	// check the stock price history
	stockHistory := stock.GetHistory()
	snap1, _ := stockHistory[t1]
	snap2, _ := stockHistory[t2]
	snap3, _ := stockHistory[t3]
	if !snap1.GetTime().Equal(t1) {
		t.Error("snap1 - unexpected time.")
	}
	if snap1.GetPrice().Float64 != 2.1 {
		t.Error("snap1 - unexpected price.")
	}

	if !snap2.GetTime().Equal(t2) {
		t.Error("snap2 - unexpected time.")
	}
	if snap2.GetPrice().Float64 != 2.0 {
		t.Error("snap2 - unexpected price.")
	}

	if !snap3.GetTime().Equal(t3) {
		t.Error("snap3 - unexpected time.")
	}
	if snap3.GetPrice().Float64 != 2.5 {
		t.Error("snap3 - unexpected price.")
	}

	// check the cash history
	cashHistory := cash.GetHistory()
	snap1, _ = cashHistory[t1]
	snap2, _ = cashHistory[t2]
	snap3, _ = cashHistory[t3]
	if !snap1.GetTime().Equal(t1) {
		t.Error("snap1 - unexpected time.")
	}
	if snap1.GetPrice().Float64 != 1.0 {
		t.Error("snap1 - unexpected price.")
	}

	if !snap2.GetTime().Equal(t2) {
		t.Error("snap2 - unexpected time.")
	}
	if snap2.GetPrice().Float64 != 1.0 {
		t.Error("snap2 - unexpected price.")
	}

	if !snap3.GetTime().Equal(t3) {
		t.Error("snap3 - unexpected time.")
	}
	if snap3.GetPrice().Float64 != 1.0 {
		t.Error("snap3 - unexpected price.")
	}

	// check the snapshotTimes
	snapshotTimes := backtest.GetSnapshotTimes()
	if len(snapshotTimes) != 3 {
		t.Fatal("Expecting 3 times to be returned")
	}
	if !snapshotTimes[0].Equal(t1) {
		t.Error("Unexpected first time")
	}
	if !snapshotTimes[1].Equal(t2) {
		t.Error("Unexpected second time")
	}
	if !snapshotTimes[2].Equal(t3) {
		t.Error("Unexpected third time")
	}

	// write this history to csv
	err = backtest.HistoryToCsv("test")
	if err != nil {
		t.Fatalf("Error in HistoryToCsv - %s", err)
	}

	err = os.Remove("test.csv")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAaplTrading(t *testing.T) {
	// initialise our portfolio and assets
	portfolio, err1 := asset.NewPortfolio("MY_ACCOUNT", "USD")
	aapl, err2 := asset.NewStock("AAPL", "USD")
	usd, err3 := asset.NewCash("USD")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Fatalf("Error in asset creation - %s", err)
	}

	// transfer 1M USD to the portfolio
	portfolio.Transfer(usd, 1e6)

	// create a compliance rule where we cannot hold more than 500 shares of AAPL stock
	stockLimit := compliance.NewUnitLimit(aapl, 500)
	portfolio.AddComplianceRule(stockLimit)

	// create our trading strategy to buy 100 shares of AAPL on each cycle
	generateTrades := func() ([]*trade.Trade, error) {
		var trades []*trade.Trade
		price := aapl.GetPrice()
		if price.Valid {
			newTrade := trade.NewTrade(portfolio, aapl, 100)
			trades = append(trades, newTrade)
		}

		return trades, nil
	}

	// create a strategy instance
	strategy := NewStrategy(generateTrades)

	// create our backtest instance, register assets
	backtest := NewBacktest(strategy)
	backtest.RegisterPortfolio(portfolio)
	backtest.RegisterAsset(aapl)
	backtest.RegisterAsset(usd)

	// get events for the backtest
	startDate := btutil.Date(2021, 4, 1)
	endDate := btutil.Date(2021, 4, 30)
	dataQuery := datasources.NewFmpCloudQuery(aapl, startDate, endDate)
	events, err := dataQuery.GenerateEvents()
	if err != nil {
		t.Fatalf("Error in GenerageEvents - %s", err)
	}

	backtest.AddEvents(events)

	// run the backtest
	err = backtest.Run()
	if err != nil {
		t.Fatalf("Error in backtest.Run() - %s", err)
	}

	// check the portfolio holdings
	// we should have kept buying aapl until
	// we hit a compliance cap at 500 shares
	aaplPosition := portfolio.GetUnits(aapl)
	if aaplPosition != 500 {
		t.Errorf("Unexpected aapl position - wanted 500, got %0.2f", aaplPosition)
	}

	// write this history to csv
	err = backtest.HistoryToCsv("test")
	if err != nil {
		t.Fatalf("Error in HistoryToCsv - %s", err)
	}

	err = os.Remove("test.csv")
	if err != nil {
		t.Fatal(err)
	}
}
