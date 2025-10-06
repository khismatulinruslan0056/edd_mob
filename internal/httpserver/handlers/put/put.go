package put

import (
	"Effective_Mobile/internal/httpserver/handlers"
	"Effective_Mobile/internal/httpserver/handlers/dto"
	"Effective_Mobile/internal/httpserver/handlers/get"
	"Effective_Mobile/internal/logger"
	"Effective_Mobile/internal/model"
	"Effective_Mobile/internal/service/enrichment"
	"Effective_Mobile/internal/storage/pg"
	"encoding/json"
	"fmt"
	"net/http"
)

type Putter interface {
	Update(id int, user *model.User) error
}

// @Summary Обновить пользователя
// @Description Обновляет данные пользователя по ID
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param user body dto.UserRequest true "Обновлённая информация о пользователе"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /people/{id} [put]
func New(putter Putter, getter get.Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "httpserver.handlers.put.new"

		logger.Debug("%s: incoming %s request on %s", op, r.Method, r.URL.Path)

		if http.MethodPut != r.Method {
			logger.Error("%s: method not allowed: %s", op, r.Method)
			handlers.WriteError(w, http.StatusMethodNotAllowed, fmt.Sprintf("%s: %s", op, handlers.MethodNotAllowed))
			return
		}

		if r.Body != nil {
			defer func() {
				if cerr := r.Body.Close(); cerr != nil {
					logger.Error("%s: request body not closed: %s", op, cerr)
				} else {
					logger.Debug("%s: request body closed", op)
				}
			}()
		}

		id, err := handlers.GetID(r.URL.Path)
		if err != nil {
			logger.Error("%s: failed to parse id: %v", op, err)
			handlers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}
		logger.Debug("%s: extracted id: %d", op, id)

		var req dto.UserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("%s: failed to decode request body: %v", op, err)
			handlers.WriteError(w, http.StatusBadRequest, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}

		user := model.User{
			Name:       req.Name,
			Surname:    req.Surname,
			Patronymic: req.Patronymic,
		}

		logger.Debug("%s: decoded user: %+v", op, user)

		ok, err := existName(id, user.Name, getter)
		if err != nil {
			logger.Error("%s: failed to check existing name: %v", op, err)
			handlers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}

		if !ok {
			logger.Debug("%s: name changed, enriching...", op)
			err = enrichment.Enrich(&user)
			if err != nil {
				logger.Error("%s: enrichment failed: %v", op, err)
				handlers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("%s: %s", op, err.Error()))
				return
			}
			logger.Debug("%s: enriched user: %+v", op, user)
		} else {
			logger.Debug("%s: name unchanged, enrichment skipped", op)
		}

		err = putter.Update(id, &user)
		if err != nil {
			logger.Error("%s: failed to update user: %v", op, err)
			handlers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}
		logger.Info("%s: user %d updated", op, id)

		response := dto.Response{
			ID:      id,
			Message: "user updated",
		}

		responseJson, err := json.Marshal(&response)
		if err != nil {
			logger.Error("%s: failed to marshal response: %v", op, err)
			handlers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJson)
		logger.Debug("%s: response written: %s", op, string(responseJson))
	}
}

func existName(id int, name string, getter get.Getter) (bool, error) {
	const op = "httpserver.handlers.put.existName"

	users, err := getter.List(&pg.ListParam{User: model.User{ID: id}})
	if err != nil {
		logger.Error("%s: failed to get user by id %d: %v", op, id, err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if len(users) == 0 {
		logger.Error("%s: user with id %d not found", op, id)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	match := users[0].Name == name
	logger.Debug("%s: existing name: %s, incoming name: %s, match: %t", op, users[0].Name, name, match)

	return match, nil
}
