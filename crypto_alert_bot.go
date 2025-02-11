package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

// API & Telegram settings
const (
	apiURL        = "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd"
	telegramToken = "YOUR_TELEGRAM_BOT_TOKEN" // Replace with your Telegram bot token
	chatID        = "YOUR_CHAT_ID"            // Replace with your Telegram chat ID
	alertThresholdBTC = 45000.0  // Set threshold for BTC alerts
	alertThresholdETH = 3000.0   // Set threshold for ETH alerts
)

// CryptoPrice stores price data
type CryptoPrice struct {
	Bitcoin  struct{ USD float64 `json:"usd"` } `json:"bitcoin"`
	Ethereum struct{ USD float64 `json:"usd"` } `json:"ethereum"`
}

// fetchCryptoPrices gets real-time prices from CoinGecko
func fetchCryptoPrices() (*CryptoPrice, error) {
	client := resty.New()
	resp, err := client.R().
		SetResult(&CryptoPrice{}).
		Get(apiURL)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*CryptoPrice), nil
}

// sendTelegramAlert sends a message to Telegram
func sendTelegramAlert(message string) error {
	telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramToken)
	client := resty.New()
	_, err := client.R().
		SetQueryParams(map[string]string{
			"chat_id": chatID,
			"text":    message,
		}).
		Get(telegramURL)

	return err
}

func main() {
	fmt.Println("Crypto Price Alert Bot Running...")

	for {
		// Fetch prices
		prices, err := fetchCryptoPrices()
		if err != nil {
			log.Println("Error fetching crypto prices:", err)
			continue
		}

		// Check for BTC price alert
		if prices.Bitcoin.USD >= alertThresholdBTC {
			msg := fmt.Sprintf("🚀 Bitcoin Alert! BTC has reached $%.2f", prices.Bitcoin.USD)
			_ = sendTelegramAlert(msg)
			fmt.Println(msg)
		}

		// Check for ETH price alert
		if prices.Ethereum.USD >= alertThresholdETH {
			msg := fmt.Sprintf("🚀 Ethereum Alert! ETH has reached $%.2f", prices.Ethereum.USD)
			_ = sendTelegramAlert(msg)
			fmt.Println(msg)
		}

		// Wait for 60 seconds before checking again
		time.Sleep(60 * time.Second)
	}
}
