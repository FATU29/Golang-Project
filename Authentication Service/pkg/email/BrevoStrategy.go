package email

import (
	"Authentication_Service/internal/config"
	"Authentication_Service/pkg/email/model"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/gopkg/util/logger"
)

type BrevoEmail struct {
	Cfg *config.Config
}

func (b *BrevoEmail) SendEmail(to model.To, subject string, body string) error {
	apiKey := b.Cfg.BrevoApiKey
	apiUrl := b.Cfg.BrevoEndpointMail

	const htmlTpl = `
    <!DOCTYPE html>
    <html>
    <body>
       <p><strong style="color:blue; font-size:20px;">{{.body}}</strong></p>
    </body>
    </html>
    `

	// 1. Parse the template
	t, err := template.New("email").Parse(htmlTpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	var tplBuffer bytes.Buffer

	// 2. Inject data into the template
	data := map[string]string{
		"body": body,
	}

	if err := t.Execute(&tplBuffer, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// 3. Prepare the payload
	payload := model.EmailPayload{
		Sender: model.Sender{
			Name:  b.Cfg.SenderName,
			Email: b.Cfg.SenderMail,
		},
		To:          []model.To{to},
		Subject:     subject,
		HtmlContent: tplBuffer.String(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	// 4. Create the HTTP Request
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", apiKey)

	// 5. Send the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("network error while sending mail: %v", err)
	}
	defer resp.Body.Close()

	logger.Info("response Status:", *resp)

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errorMsg := string(bodyBytes)
		logger.Errorf("Brevo Error Detail: %s", errorMsg)
		fmt.Printf("Brevo Error Detail: %s\n", errorMsg)
		return fmt.Errorf("brevo api error: %s", errorMsg)
	}

	// 6. Check the response status
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil // Success
	}

	// Optional: Read response body here to see why Brevo rejected it
	return fmt.Errorf("send mail failed with status code: %d", resp.StatusCode)
}
