package handlers

import (
	"Effective_Mobile/internal/logger"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	InvalidPath      = errors.New("invalid path")
	InvalidUserID    = errors.New("invalid user id")
	UserNotFound     = errors.New("user not found")
	MethodNotAllowed = errors.New("method not allowed")
)

func GetID(path string) (int, error) {
	const op = "httpserver.handlers.del.getID"

	parts := strings.Split(path, "/")
	logger.Debug("%s: split path '%s' into %v", op, path, parts)

	if len(parts) != 3 || parts[1] != "people" {
		logger.Error("%s: invalid path format", op)
		return -1, fmt.Errorf("%s: %w", op, InvalidPath)
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		logger.Error("%s: failed to convert id '%s': %v", op, parts[2], err)
		return -1, fmt.Errorf("%s: %w", op, InvalidUserID)
	}

	logger.Debug("%s: extracted id: %d", op, id)
	return id, nil
}
