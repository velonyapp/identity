package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/query"
)

var _ query.User = (*userQuery)(nil)

type userQuery struct {
	db *sql.DB
}

func NewUserQuery(db *sql.DB) query.User {
	return &userQuery{db: db}
}

func (q *userQuery) GetUser(ctx context.Context, userID string) (*common.User, error) {
	const query = `
		SELECT
			users.id,
			users.username,
			users.full_name,
			users.email,
			avatars.key
		FROM users
		LEFT JOIN avatars
			ON avatars.id = users.avatar_id
		WHERE users.id = ?
		LIMIT 1
	`

	row := executor(ctx, q.db).QueryRowContext(ctx, query, userID)

	user, err := scanQueryUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (q *userQuery) BatchGetUser(ctx context.Context, userIDs []string) ([]*common.User, error) {
	if len(userIDs) == 0 {
		return []*common.User{}, nil
	}

	placeholders := make([]string, len(userIDs))
	args := make([]any, len(userIDs))

	for i, userID := range userIDs {
		placeholders[i] = "?"
		args[i] = userID
	}

	query := `
		SELECT
			users.id,
			users.username,
			users.full_name,
			users.email,
			avatars.key
		FROM users
		LEFT JOIN avatars
			ON avatars.id = users.avatar_id
		WHERE users.id IN (` + strings.Join(placeholders, ", ") + `)
	`

	rows, err := executor(ctx, q.db).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*common.User, 0, len(userIDs))

	for rows.Next() {
		user, err := scanQueryUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

type queryUserScanner interface {
	Scan(dest ...any) error
}

func scanQueryUser(scanner queryUserScanner) (*common.User, error) {
	var (
		id        string
		username  string
		fullName  string
		email     sql.NullString
		avatarKey sql.NullString
	)

	if err := scanner.Scan(
		&id,
		&username,
		&fullName,
		&email,
		&avatarKey,
	); err != nil {
		return nil, err
	}

	user := &common.User{
		ID:       id,
		Username: username,
		FullName: fullName,
	}

	if email.Valid {
		user.Email = &email.String
	}
	if avatarKey.Valid {
		user.AvatarKey = &avatarKey.String
	}

	return user, nil
}
