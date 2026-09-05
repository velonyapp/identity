package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/velonyapp/identity/internal/application/domainevent"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type UserRepo struct {
	db         *sql.DB
	dispatcher *domainevent.Dispatcher
}

func NewUserRepo(
	db *sql.DB,
	dispatcher *domainevent.Dispatcher,
) repo.User {
	return &UserRepo{
		db:         db,
		dispatcher: dispatcher,
	}
}

type userScanner interface {
	Scan(dest ...any) error
}

func (repo *UserRepo) FindByID(ctx context.Context, userID vo.UserID) (*entity.User, error) {
	const query = `
		SELECT
			users.id,
			users.username,
			users.email,
			users.full_name,
			users.avatar_key,
			users.create_time,
			users.update_time,
			local_auth_strategies.password_hash
		FROM users
		LEFT JOIN local_auth_strategies
			ON local_auth_strategies.user_id = users.id
		WHERE users.id = ?
		LIMIT 1
	`

	row := executor(ctx, repo.db).QueryRowContext(ctx, query, userID.Value())

	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (repo *UserRepo) FindByIDs(ctx context.Context, userIDs []vo.UserID) ([]*entity.User, error) {
	if len(userIDs) == 0 {
		return []*entity.User{}, nil
	}

	placeholders := make([]string, len(userIDs))
	args := make([]any, len(userIDs))

	for i, userID := range userIDs {
		placeholders[i] = "?"
		args[i] = userID.Value()
	}

	query := `
		SELECT
			users.id,
			users.username,
			users.email,
			users.full_name,
			users.avatar_key,
			users.create_time,
			users.update_time,
			local_auth_strategies.password_hash
		FROM users
		LEFT JOIN local_auth_strategies
			ON local_auth_strategies.user_id = users.id
		WHERE users.id IN (` + strings.Join(placeholders, ", ") + `)
	`

	rows, err := executor(ctx, repo.db).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*entity.User, 0, len(userIDs))

	for rows.Next() {
		user, err := scanUser(rows)
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

func (repo *UserRepo) FindByUsername(ctx context.Context, username vo.Username) (*entity.User, error) {
	const query = `
		SELECT
			users.id,
			users.username,
			users.email,
			users.full_name,
			users.avatar_key,
			users.create_time,
			users.update_time,
			local_auth_strategies.password_hash
		FROM users
		LEFT JOIN local_auth_strategies
			ON local_auth_strategies.user_id = users.id
		WHERE users.username = ?
		LIMIT 1
	`

	row := executor(ctx, repo.db).QueryRowContext(ctx, query, username.Value())

	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (repo *UserRepo) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	const query = `
		SELECT
			users.id,
			users.username,
			users.email,
			users.full_name,
			users.avatar_key,
			users.create_time,
			users.update_time,
			local_auth_strategies.password_hash
		FROM users
		LEFT JOIN local_auth_strategies
			ON local_auth_strategies.user_id = users.id
		WHERE users.email = ?
		LIMIT 1
	`

	row := executor(ctx, repo.db).QueryRowContext(ctx, query, email.Value())

	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (repo *UserRepo) Save(ctx context.Context, user *entity.User) error {
	if user.DeleteTime != nil {
		const query = `
			DELETE FROM users
			WHERE id = ?
		`

		if _, err := executor(ctx, repo.db).ExecContext(ctx, query, user.ID.Value()); err != nil {
			return err
		}
	} else {
		const query = `
			INSERT INTO users (
				id,
				username,
				email,
				full_name,
				avatar_key,
				create_time,
				update_time
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				username = ?,
				email = ?,
				full_name = ?,
				avatar_key = ?,
				update_time = ?
		`

		var email any
		if user.Email != nil {
			email = user.Email.Value()
		}

		var avatarKey any
		if user.AvatarKey != nil {
			avatarKey = user.AvatarKey.Value()
		}

		if _, err := executor(ctx, repo.db).ExecContext(ctx, query,
			user.ID.Value(),
			user.Username.Value(),
			email,
			user.FullName.Value(),
			avatarKey,
			user.CreateTime.Value(),
			user.UpdateTime.Value(),

			user.Username.Value(),
			email,
			user.FullName.Value(),
			avatarKey,
			user.UpdateTime.Value(),
		); err != nil {
			return err
		}

		if user.LocalAuthStrategy != nil {
			const localAuthStrategyQuery = `
				INSERT INTO local_auth_strategies (
					user_id,
					password_hash
				)
				VALUES (?, ?)
				ON DUPLICATE KEY UPDATE
					password_hash = ?
			`

			if _, err := executor(ctx, repo.db).ExecContext(ctx, localAuthStrategyQuery,
				user.LocalAuthStrategy.UserID.Value(),
				user.LocalAuthStrategy.PasswordHash.Value(),

				user.LocalAuthStrategy.PasswordHash.Value(),
			); err != nil {
				return err
			}
		} else {
			const localAuthStrategyQuery = `
				DELETE FROM local_auth_strategies
				WHERE user_id = ?
			`

			if _, err := executor(ctx, repo.db).ExecContext(
				ctx,
				localAuthStrategyQuery,
				user.ID.Value(),
			); err != nil {
				return err
			}
		}
	}

	for _, domainEvent := range user.PullEvents() {
		if err := repo.dispatcher.Dispatch(ctx, domainEvent); err != nil {
			return err
		}
	}

	return nil
}

func scanUser(scanner userScanner) (*entity.User, error) {
	var (
		id           string
		username     string
		email        sql.NullString
		fullName     string
		avatarKey    sql.NullString
		createTime   time.Time
		updateTime   time.Time
		passwordHash sql.NullString
	)

	if err := scanner.Scan(
		&id,
		&username,
		&email,
		&fullName,
		&avatarKey,
		&createTime,
		&updateTime,
		&passwordHash,
	); err != nil {
		return nil, err
	}

	usernameVO, _ := vo.NewUsername(username)

	var emailVO *vo.Email
	if email.Valid {
		value, _ := vo.NewEmail(email.String)
		emailVO = &value
	}

	fullNameVO, _ := vo.NewFullName(fullName)

	var avatarKeyVO *vo.StorageKey
	if avatarKey.Valid {
		value, _ := vo.NewStorageKey(avatarKey.String)
		avatarKeyVO = &value
	}

	var localAuthStrategy *entity.LocalAuthStrategy
	if passwordHash.Valid {
		passwordHashVO, _ := vo.NewPasswordHash(passwordHash.String)

		localAuthStrategy = &entity.LocalAuthStrategy{
			UserID:       vo.NewUserID(id),
			PasswordHash: passwordHashVO,
		}
	}

	return &entity.User{
		ID:                vo.NewUserID(id),
		Username:          usernameVO,
		Email:             emailVO,
		FullName:          fullNameVO,
		AvatarKey:         avatarKeyVO,
		CreateTime:        vo.NewTime(createTime),
		UpdateTime:        vo.NewTime(updateTime),
		LocalAuthStrategy: localAuthStrategy,
	}, nil
}
