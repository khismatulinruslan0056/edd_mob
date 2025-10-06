package null

import (
	"Effective_Mobile/internal/logger"
	"database/sql"
)

func SqlNullStringValid(param sql.NullString) *string {
	if param.Valid {
		logger.Debug("SqlNullStringValid: valid string found = %s", param.String)
		return &param.String
	}

	logger.Debug("SqlNullStringValid: string is NULL")
	return nil
}

func SqlNullInt64Valid(param sql.NullInt64) *int {
	if param.Valid {
		paramInt := int(param.Int64)
		logger.Debug("SqlNullInt64Valid: valid int64 found = %d", paramInt)
		return &paramInt
	}

	logger.Debug("SqlNullInt64Valid: int64 is NULL")
	return nil
}
