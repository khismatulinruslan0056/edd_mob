package pg

import (
	"Effective_Mobile/internal/config"
	"Effective_Mobile/internal/logger"
	"Effective_Mobile/internal/model"
	"Effective_Mobile/internal/storage"
	"Effective_Mobile/lib/null"
	"database/sql"
	"fmt"
	pq "github.com/lib/pq"
	"strconv"
	"strings"
)

type Storage struct {
	db *sql.DB
}

func New(cfDSN *config.DsnPG) (*Storage, error) {
	const op = "storage.pg.new"

	dsn := DSN(cfDSN)
	logger.Info("%s: connecting to PostgreSQL", op)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("%s: failed to open DB: %v", op, err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("%s: successfully connected to PostgreSQL", op)
	return &Storage{db: db}, nil
}

func DSN(dsn *config.DsnPG) string {
	return fmt.Sprintf(
		"host=%v port=%v user=%s password=%s dbname=%s sslmode=disable",
		dsn.Host,
		dsn.Port,
		dsn.User,
		dsn.Password,
		dsn.Name,
	)
}

func (s *Storage) Add(user model.User) (int, error) {
	const op = "storage.pg.add"
	//tx, err := s.db.BeginTx(ctx, nil)
	//if err != nil {
	//	return -1, fmt.Errorf("%s: %w", op, err)
	//}
	//defer func() {
	//	if err != nil {
	//		tx.Rollback()
	//	}
	//}()

	args, columns, placeHolders := prepareQuery(user)

	query := fmt.Sprintf("INSERT INTO people (%s) VALUES (%s) RETURNING id",
		strings.Join(columns, ", "),
		strings.Join(placeHolders, ", "),
	)
	//res, err := tx.ExecContext(ctx, query, args)
	var id int
	err := s.db.QueryRow(query, args...).Scan(&id)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			logger.Error("%s: user already exists: %v", op, err)
			return -1, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		logger.Error("%s: insert failed: %v", op, err)
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	//if err = tx.Commit(); err != nil {
	//	return -1, fmt.Errorf("%s: %w", op, err)
	//}

	logger.Debug("%s: user added with ID %d", op, id)
	return id, nil
}

func (s *Storage) List(params *ListParam) ([]*model.User, error) {
	const op = "storage.pg.list"

	args, columns, placeHolders := prepareQuery(params.User)

	sb := strings.Builder{}
	sb.WriteString("SELECT * FROM people")

	if len(columns) > 0 {
		sb.WriteString(" WHERE ")
		for i, column := range columns {
			sb.WriteString(column)
			sb.WriteString(" = ")
			sb.WriteString(placeHolders[i])
			if i < len(columns)-1 {
				sb.WriteString(" AND ")
			}
		}
	}

	if params.Limit > 0 {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("LIMIT $%d OFFSET $%d", len(columns)+1, len(columns)+2))
		args = append(args, params.Limit, params.Offset)
	}
s.
	rows, err := s.db.Query(sb.String(), args...)
	if err != nil {
		logger.Error("%s: list query failed: %v", op, err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("%s: rows close failed: %v", op, err)
		} else {
			logger.Debug("%s: rows closed successfully", op)
		}
	}()

	users := make([]*model.User, 0)

	for rows.Next() {
		var (
			patronymic  sql.NullString
			age         sql.NullInt64
			gender      sql.NullString
			nationality sql.NullString
		)

		user := &model.User{}
		if err = rows.Scan(&user.ID, &user.Name, &user.Surname,
			&patronymic, &gender, &age, &nationality); err != nil {
			logger.Error("%s: scan failed: %v", op, err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		user.Patronymic = null.SqlNullStringValid(patronymic)
		user.Gender = null.SqlNullStringValid(gender)
		user.Nationality = null.SqlNullStringValid(nationality)
		user.Age = null.SqlNullInt64Valid(age)

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		logger.Error("%s: rows error: %v", op, err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Debug("%s: found %d users", op, len(users))
	return users, err
}

func (s *Storage) Delete(id int) error {
	const op = "storage.pg.del"

	logger.Debug("%s: deleting user with ID %d", op, id)

	query := `DELETE FROM people WHERE id = $1`
	res, err := s.db.Exec(query, id)
	if err != nil {
		logger.Error("%s: delete failed: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		logger.Error("%s: failed to get affected rows: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}
	if affected == 0 {
		logger.Debug("%s: user with ID %d not found", op, id)
		return storage.ErrUserNotFound
	}

	logger.Debug("%s: user with ID %d deleted", op, id)
	return nil
}

func (s *Storage) Update(id int, user *model.User) error {
	const op = "storage.pg.update"

	args, columns, placeHolders := prepareQuery(*user)
	if len(columns) == 0 {
		logger.Error("%s: nothing to update", op)
		return fmt.Errorf("%s: %w", op, storage.ErrNothingUpdate)
	}

	sb := strings.Builder{}

	sb.WriteString("UPDATE people SET ")
	for i, column := range columns {
		sb.WriteString(column)
		sb.WriteString(" = ")
		sb.WriteString(placeHolders[i])
		if i != len(columns)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString(" WHERE id = $")
	sb.WriteString(strconv.Itoa(len(columns) + 1))
	args = append(args, id)

	res, err := s.db.Exec(sb.String(), args...)
	if err != nil {
		logger.Error("%s: update failed: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		logger.Error("%s: failed to get affected rows: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}
	if affected == 0 {
		logger.Debug("%s: user with ID %d not found for update", op, id)
		return storage.ErrUserNotFound
	}

	logger.Debug("%s: updated user with ID %d", op, id)
	return nil
}

func (s *Storage) Close() error {
	const op = "storage.pg.close"
	logger.Info("%s: closing DB connection", op)

	if err := s.db.Close(); err != nil {
		logger.Error("%s: DB connection close failed: %v", op, err)
		return err
	}

	return nil
}

type ListParam struct {
	User   model.User
	Limit  int
	Offset int
}

func prepareElemForQuery(args []interface{}, columns []string, placeHolders []string,
	index *int, arg interface{}, column string) ([]interface{}, []string, []string) {

	placeholder := "$" + strconv.Itoa(*index)
	logger.Debug("prepareElemForQuery: column=%s, placeholder=%s, arg=%v, index=%d", column, placeholder, arg, *index)

	logger.Debug("prepareElemForQuery: column=%s, placeholder=%s, arg=%v, index=%d", column, placeholder, arg, *index)

	args = append(args, arg)
	columns = append(columns, column)
	placeHolders = append(placeHolders, placeholder)
	*index++
	return args, columns, placeHolders
}

func prepareQuery(user model.User) ([]interface{}, []string, []string) {
	args := make([]interface{}, 0)
	columns := make([]string, 0)
	placeHolders := make([]string, 0)
	index := 1

	if user.ID > 0 {
		column := "id"
		arg := user.ID
		args, columns, placeHolders = prepareElemForQuery(args, columns, placeHolders, &index, arg, column)
	}

	if user.Name != "" {
		column := "name"
		arg := user.Name
		args, columns, placeHolders = prepareElemForQuery(args, columns, placeHolders, &index, arg, column)
	}

	if user.Surname != "" {
		column := "surname"
		arg := user.Surname
		args, columns, placeHolders = prepareElemForQuery(args, columns, placeHolders, &index, arg, column)
	}

	if user.Patronymic != nil {
		column := "patronymic"
		arg := user.Patronymic
		args, columns, placeHolders = prepareElemForQuery(args, columns, placeHolders, &index, arg, column)
	}

	if user.Age != nil {
		column := "age"
		arg := user.Age
		args, columns, placeHolders = prepareElemForQuery(args, columns, placeHolders, &index, arg, column)

	}

	if user.Gender != nil {
		column := "gender"
		arg := user.Gender
		args, columns, placeHolders = prepareElemForQuery(args, columns, placeHolders, &index, arg, column)

	}

	if user.Nationality != nil {
		column := "nationality"
		arg := user.Nationality
		args, columns, placeHolders = prepareElemForQuery(args, columns, placeHolders, &index, arg, column)
	}

	return args, columns, placeHolders
}
