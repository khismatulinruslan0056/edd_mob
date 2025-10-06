package enrichment

import (
	"Effective_Mobile/internal/logger"
	"Effective_Mobile/internal/model"
	"encoding/json"
	"fmt"
)

func EnrichAge(user *model.User) error {
	const op = "service.enrichment.enrichAge"
	url := ageEnrichURL + user.Name
	logger.Debug("%s: enriching age from %s", op, url)

	body, err := FetchBody(url, "Age")
	if err != nil {
		return err
	}

	var userAge model.UserAge
	err = json.Unmarshal(body, &userAge)
	if err != nil {
		logger.Error("%s: failed to unmarshal age data: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	user.Age = &userAge.Age
	logger.Info("%s: enriched age: %d", op, userAge.Age)
	return nil
}
