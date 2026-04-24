package kafka

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type RegisterOTPEvent struct {
	Email       string    `json:"email"`
	Firstname   *string   `json:"firstname,omitempty"`
	Lastname    *string   `json:"lastname,omitempty"`
	Otp         string    `json:"otp"`
	RequestedAt time.Time `json:"requestedAt"`
}

func (event RegisterOTPEvent) Validate() error {
	if strings.TrimSpace(event.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(event.Otp) == "" {
		return fmt.Errorf("otp is required")
	}
	return nil
}

func (event RegisterOTPEvent) RecipientName() string {
	first := strings.TrimSpace(stringOrEmpty(event.Firstname))
	last := strings.TrimSpace(stringOrEmpty(event.Lastname))
	full := strings.TrimSpace(first + " " + last)
	if full != "" {
		return full
	}
	return event.Email
}

func (event RegisterOTPEvent) Marshal() ([]byte, error) {
	if err := event.Validate(); err != nil {
		return nil, err
	}
	if event.RequestedAt.IsZero() {
		event.RequestedAt = time.Now().UTC()
	}
	return json.Marshal(event)
}

func ParseRegisterOTPEvent(raw []byte) (RegisterOTPEvent, error) {
	var event RegisterOTPEvent
	if err := json.Unmarshal(raw, &event); err != nil {
		return RegisterOTPEvent{}, err
	}
	if err := event.Validate(); err != nil {
		return RegisterOTPEvent{}, err
	}
	return event, nil
}

func ParseBrokers(brokersCSV string) []string {
	parts := strings.Split(brokersCSV, ",")
	brokers := make([]string, 0, len(parts))
	for _, part := range parts {
		broker := strings.TrimSpace(part)
		if broker == "" {
			continue
		}
		brokers = append(brokers, broker)
	}
	return brokers
}

func stringOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
