package enrichment

import (
	"Effective_Mobile/internal/logger"
	"Effective_Mobile/internal/model"
	"fmt"
	"io"
	"net/http"
)

const (
	ageEnrichURL         = "https://api.agify.io/?name="
	genderEnrichURL      = "https://api.genderize.io/?name="
	nationalityEnrichURL = "https://api.nationalize.io/?name="
)

func Enrich(user *model.User) error {
	const op = "service.enrichment.enrich"
	logger.Info("%s: start enrichment for user: %s", op, user.Name)

	if err := EnrichGender(user); err != nil {
		logger.Error("%s: gender enrichment failed: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := EnrichAge(user); err != nil {
		logger.Error("%s: age enrichment failed: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := EnrichNationality(user); err != nil {
		logger.Error("%s: nationality enrichment failed: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("%s: enrichment complete for user: %s", op, user.Name)
	return nil
}

func FetchBody(url string, context string) ([]byte, error) {
	const op = "service.enrichment.FetchBody"
	logger.Debug("%s: sending GET request to %s", op, url)

	res, err := http.Get(url)
	if err != nil {
		logger.Error("%s: failed GET request for %s: %v", op, context, err)
		return nil, fmt.Errorf("%s: %w", op+context, err)
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("%s: %s returned status code %d", op, context, res.StatusCode)
		return nil, fmt.Errorf("%s: unexpected status code: %d", op+context, res.StatusCode)
	}

	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			logger.Error("%s: failed to close response body for %s: %v", op, context, cerr)
		} else {
			logger.Debug("%s: closed response body for %s", op, context)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("%s: failed to read response body for %s: %v", op, context, err)
		return nil, fmt.Errorf("%s: %w", op+context, err)
	}

	logger.Debug("%s: response body fetched for %s: %s", op, context, string(body))
	return body, nil
}
