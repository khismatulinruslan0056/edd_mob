package get

import (
	"Effective_Mobile/internal/httpserver/handlers"
	"Effective_Mobile/internal/httpserver/handlers/dto"
	"Effective_Mobile/internal/logger"
	"Effective_Mobile/internal/model"
	"Effective_Mobile/internal/storage/pg"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Getter interface {
	List(params *pg.ListParam) ([]*model.User, error)
}

// @Summary Получить список людей
// @Description Возвращает список людей с поддержкой фильтрации по полям и пагинации
// @Tags people
// @Accept json
// @Produce json
// @Param id query int false "ID пользователя"
// @Param name query string false "Имя"
// @Param surname query string false "Фамилия"
// @Param patronymic query string false "Отчество"
// @Param age query int false "Возраст"
// @Param gender query string false "Пол"
// @Param nationality query string false "Национальность"
// @Param limit query int false "Максимальное количество записей"
// @Param offset query int false "Смещение (offset) для пагинации"
// @Success 200 {array} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /people [get]
func New(getter Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "httpserver.handlers.get.new"

		logger.Info("%s: received request %s", op, r.URL.RawQuery)

		if http.MethodGet != r.Method {
			logger.Error("%s: method not allowed: %s", op, r.Method)
			handlers.WriteError(w, http.StatusMethodNotAllowed, fmt.Sprintf("%s: %s", op, handlers.MethodNotAllowed))
			return
		}

		rows := r.URL.Query()

		params, err := getParams(rows)
		if err != nil {
			logger.Error("%s: invalid query params: %v", op, err)
			handlers.WriteError(w, http.StatusBadRequest, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}
		logger.Debug("%s: parsed params: %+v", op, params)

		users, err := getter.List(params)
		if err != nil {
			logger.Error("%s: failed to list users: %v", op, err)
			handlers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}
		logger.Info("%s: successfully retrieved %d users", op, len(users))

		dtoUsers := make([]*dto.UserResponse, len(users))

		for i, user := range users {
			dtoUsers[i] = toDTO(user)
		}

		response, err := json.Marshal(dtoUsers)
		if err != nil {
			logger.Error("%s: failed to marshal users: %v", op, err)
			handlers.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("%s: %s", op, err.Error()))
			return
		}
		logger.Debug("%s: response payload: %s", op, response)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func getParams(rows url.Values) (*pg.ListParam, error) {
	const op = "httpserver.handlers.get.getParams"

	params := &pg.ListParam{}
	user := &model.User{}

	if limitStr := rows.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			logger.Error("%s: invalid limit: %v", op, err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		params.Limit = limit
		logger.Debug("%s: parsed limit: %d", op, limit)
	}

	if offsetStr := rows.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			logger.Error("%s: invalid offset: %v", op, err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		params.Offset = offset
		logger.Debug("%s: parsed offset: %d", op, offset)
	}

	if IDStr := rows.Get("id"); IDStr != "" {
		ID, err := strconv.Atoi(IDStr)
		if err != nil {
			logger.Error("%s: invalid id: %v", op, err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		user.ID = ID
		logger.Debug("%s: parsed id: %d", op, ID)
	}

	if name := rows.Get("name"); name != "" {
		user.Name = name
		logger.Debug("%s: parsed name: %s", op, name)
	}

	if surname := rows.Get("surname"); surname != "" {
		user.Surname = surname
		logger.Debug("%s: parsed surname: %s", op, surname)
	}

	if patronymic := rows.Get("patronymic"); patronymic != "" {
		user.Patronymic = &patronymic
		logger.Debug("%s: parsed patronymic: %s", op, patronymic)
	}

	if ageStr := rows.Get("age"); ageStr != "" {
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			logger.Error("%s: invalid age: %v", op, err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		user.Age = &age
		logger.Debug("%s: parsed age: %d", op, age)
	}

	if gender := rows.Get("gender"); gender != "" {
		user.Gender = &gender
		logger.Debug("%s: parsed gender: %s", op, gender)
	}

	if nationality := rows.Get("nationality"); nationality != "" {
		user.Nationality = &nationality
		logger.Debug("%s: parsed nationality: %s", op, nationality)
	}

	params.User = *user
	logger.Debug("%s: constructed ListParam: %+v", op, params)
	return params, nil
}

func toDTO(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Age:         user.Age,
		Gender:      user.Gender,
		Nationality: user.Nationality,
	}
}
