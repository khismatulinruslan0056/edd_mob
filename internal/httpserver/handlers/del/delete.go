package del

import (
	"Effective_Mobile/internal/httpserver/handlers"
	"Effective_Mobile/internal/httpserver/handlers/dto"
	"Effective_Mobile/internal/logger"
	"encoding/json"
	"fmt"
	"net/http"
)

type Deleter interface {
	Delete(id int) error
}

// @Summary Удалить пользователя
// @Description Удаляет пользователя по ID
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /people/{id} [delete]
func New(deleter Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "httpserver.handlers.del.new"

		logger.Debug("%s: incoming %s request on %s", op, r.Method, r.URL.Path)

		if http.MethodDelete != r.Method {
			logger.Error("%s: method not allowed: %s", op, r.Method)
			handlers.WriteError(w, http.StatusMethodNotAllowed, fmt.Sprintf("%s: %s", op, handlers.MethodNotAllowed))
			return
		}

		id, err := handlers.GetID(r.URL.Path)
		if err != nil {
			logger.Error("%s: failed to extract ID from URL: %v", op, err)
			handlers.WriteError(w, http.StatusBadRequest, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}
		logger.Debug("%s: extracted id: %d", op, id)

		err = deleter.Delete(id)
		if err != nil {
			logger.Error("%s: failed to delete user with id %d: %v", op, id, err)
			handlers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}
		logger.Info("%s: successfully deleted user with id %d", op, id)

		response := dto.Response{
			ID:      id,
			Message: "user deleted",
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
