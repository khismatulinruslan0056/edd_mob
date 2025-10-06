package enrichment

import (
	"Effective_Mobile/internal/logger"
	"Effective_Mobile/internal/model"
	"encoding/json"
	"fmt"
)

func EnrichGender(user *model.User) error {
	const op = "service.enrichment.enrichGender"
	url := genderEnrichURL + user.Name
	logger.Debug("%s: enriching gender from %s", op, url)

	body, err := FetchBody(url, "Gender")
	if err != nil {
		return err
	}

	var userGender model.UserGender
	err = json.Unmarshal(body, &userGender)
	if err != nil {
		logger.Error("%s: failed to unmarshal gender data: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	user.Gender = &userGender.Gender
	logger.Info("%s: enriched gender: %s", op, userGender.Gender)
	return nil
}
