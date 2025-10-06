package enrichment

import (
	"Effective_Mobile/internal/logger"
	"Effective_Mobile/internal/model"
	"encoding/json"
	"fmt"
)

func EnrichNationality(user *model.User) error {
	const op = "service.enrichment.enrichNationality"
	url := nationalityEnrichURL + user.Name
	logger.Debug("%s: enriching nationality from %s", op, url)

	body, err := FetchBody(url, "Nationality")
	if err != nil {
		return err
	}

	var userNationality model.UserNationality
	err = json.Unmarshal(body, &userNationality)
	if err != nil {
		logger.Error("%s: failed to unmarshal nationality data: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(userNationality.Countries) > 0 {
		user.Nationality = &userNationality.Countries[0].CountryID // выбираем наиболее вероятную национальность
		logger.Info("%s: enriched nationality: %s", op, userNationality.Countries[0].CountryID)
	} else {
		logger.Info("%s: no nationality data found", op)
	}

	return nil
}
