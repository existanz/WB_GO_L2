package entities

import (
	"fmt"
	"time"
)

func ValidateEvent(event Event) error {
	if event.UserID == 0 {
		return fmt.Errorf("event's user_id is required")
	}

	if event.Title == "" {
		return fmt.Errorf("event's title is required")
	}

	if event.Date.IsZero() {
		return fmt.Errorf("event's date is required")
	}

	if event.Date.Before(time.Now()) {
		return fmt.Errorf("event's date should be in the future")
	}

	if event.Duration <= 0 {
		return fmt.Errorf("event's duration should be greater than 0")
	}

	return nil
}
