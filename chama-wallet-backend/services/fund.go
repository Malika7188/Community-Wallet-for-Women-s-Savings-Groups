package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// FundTestAccount uses the Stellar Friendbot to send test XLM to a new account
func FundTestAccount(address string) error {
	url := fmt.Sprintf("https://friendbot.stellar.org/?addr=%s", address)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to call friendbot: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("friendbot returned non-200 status: %s - %s", resp.Status, body)
	}

	return nil
}
