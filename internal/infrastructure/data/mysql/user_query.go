package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/query"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ query.User = (*userQuery)(nil)

type userQuery struct {
	db *sql.DB
}

func NewUserQuery(db *sql.DB) query.User {
	return &userQuery{
		db: db,
	}
}

type queryUserScanner interface {
	Scan(dest ...any) error
}

func (r *userQuery) Find(ctx context.Context, userID vo.UserID) (*common.UserResult, error) {
	const query = `
		SELECT
			users.id,
			users.username,
			users.email,
			users.create_time,

			avatars.object_key
		FROM users
		LEFT JOIN avatars
			ON avatars.id = users.avatar_id
		WHERE users.id = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, userID.String())

	user, err := scanQueryUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *userQuery) BatchFind(ctx context.Context, userIDs []vo.UserID) (map[vo.UserID]common.UserResult, error) {
	results := make(map[vo.UserID]common.UserResult, len(userIDs))
	if len(userIDs) == 0 {
		return results, nil
	}

	placeholders := make([]string, len(userIDs))
	args := make([]any, len(userIDs))

	for i, userID := range userIDs {
		placeholders[i] = "?"
		args[i] = userID.String()
	}

	query := `
		SELECT
			users.id,
			users.username,
			users.email,
			users.create_time,

			avatars.object_key
		FROM users
		LEFT JOIN avatars
			ON avatars.id = users.avatar_id
		WHERE users.id IN (` + strings.Join(placeholders, ", ") + `)
	`

	rows, err := executor(ctx, r.db).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(map[string]common.UserResult, len(userIDs))

	for rows.Next() {
		user, err := scanQueryUser(rows)
		if err != nil {
			return nil, err
		}

		users[user.ID] = user
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, userID := range userIDs {
		if user, ok := users[userID.String()]; ok {
			results[userID] = user
		}
	}

	return results, nil
}

func scanQueryUser(scanner queryUserScanner) (common.UserResult, error) {
	var (
		id           string
		username     string
		emailRaw     sql.NullString
		createTime   time.Time
		avatarKeyRaw sql.NullString
	)

	if err := scanner.Scan(
		&id,
		&username,
		&emailRaw,
		&createTime,
		&avatarKeyRaw,
	); err != nil {
		return common.UserResult{}, err
	}

	var email *string
	if emailRaw.Valid {
		value := emailRaw.String
		email = &value
	}
	var avatarKey *string
	if avatarKeyRaw.Valid {
		value := avatarKeyRaw.String
		avatarKey = &value
	}

	return common.UserResult{
		ID:         id,
		Username:   username,
		Email:      email,
		AvatarKey:  avatarKey,
		CreateTime: createTime,
	}, nil
}
