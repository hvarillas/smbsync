package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendTelegramMessage(text string) error {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	
	if token == "" || chatID == "" {
		return fmt.Errorf("telegram credentials not configured")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	body := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	var respBody bytes.Buffer
	respBody.ReadFrom(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s\n%s", resp.Status, respBody.String())
	}

	return nil
}
