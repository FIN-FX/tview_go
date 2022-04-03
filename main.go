package main

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"log"
	syslog2 "log/syslog"
	"strconv"
	"strings"
	"tview_go/config"
	"tview_go/lib"
)

const logProgramName string = "tview_go_main"

var client *binance.Client

func initLogger() {
	syslog, err := syslog2.New(syslog2.LOG_INFO, logProgramName)
	if err != nil {
		log.Fatal(err)
	} else {
		log.SetOutput(syslog)
	}
}

func getTickers() []*binance.SymbolPrice {
	prices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return prices
}

func getBalances() (map[string]float64, binance.Balance) {
	var usdtB binance.Balance
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil, usdtB
	}
	bb := make(map[string]float64, len(account.Balances))
	for _, v := range account.Balances {
		free, _ := strconv.ParseFloat(v.Free, 64)
		lock, _ := strconv.ParseFloat(v.Locked, 64)
		bb[v.Asset] = free + lock
		if v.Asset == "USDT" {
			usdtB = v
		}
	}
	return bb, usdtB
}

func main() {
	initLogger()
	log.Println("Start")
	client = binance.NewClient(config.ApiKey, config.SecretKey)
	tickers := getTickers()
	balances, usdtB := getBalances()
	var sum float64
	for _, ticker := range tickers {
		asset, after, ok := strings.Cut(ticker.Symbol, "USDT")
		if ok && after == "" {
			item, found := balances[asset]
			if !found {
				continue
			}
			if item == 0 {
				continue
			}
			lastPrice, _ := strconv.ParseFloat(ticker.Price, 64)
			b := item * lastPrice
			sum += b
		}
	}
	sum = sum + balances["USDT"]
	lib.SendMessage("ALL: " + strconv.FormatFloat(sum, 'f', 2, 64))
	lib.SendMessage("FREE: " + usdtB.Free)
}
