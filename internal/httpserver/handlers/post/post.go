package post

import (
	"Effective_Mobile/internal/httpserver/handlers"
	"Effective_Mobile/internal/httpserver/handlers/dto"
	"Effective_Mobile/internal/logger"
	"Effective_Mobile/internal/model"
	"Effective_Mobile/internal/service/enrichment"
	"encoding/json"
	"fmt"
	"net/http"
)

type Poster interface {
	Add(user model.User) (int, error)
}

// @Summary Добавить нового пользователя
// @Description Создаёт нового пользователя и возвращает его ID
// @Tags people
// @Accept json
// @Produce json
// @Param user body dto.UserRequest true "Информация о пользователе"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /people [post]
func New(poster Poster) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "httpserver.handlers.post.new"

		logger.Debug("%s: incoming %s request on %s", op, r.Method, r.URL.Path)

		if http.MethodPost != r.Method {
			logger.Error("%s: method not allowed: %s", op, r.Method)
			http.Error(w, fmt.Sprintf("%s: %s", op, handlers.MethodNotAllowed), http.StatusMethodNotAllowed)
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

		err := enrichment.Enrich(&user)
		if err != nil {
			logger.Error("%s: enrichment failed: %v", op, err)
			handlers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}

		logger.Debug("%s: enriched user: %+v", op, user)

		id, err := poster.Add(user)
		if err != nil {
			logger.Error("%s: failed to add user: %v", op, err)
			handlers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}
		logger.Info("%s: user added with id %d", op, id)

		response := dto.Response{
			ID:      id,
			Message: "user added",
		}

		responseJson, err := json.Marshal(&response)
		if err != nil {
			logger.Error("%s: failed to marshal response: %v", op, err)
			handlers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseJson)
		logger.Debug("%s: response written: %s", op, string(responseJson))
	}
}
