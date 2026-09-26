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

var _ repo.User = (*userRepo)(nil)

type userRepo struct {
	db         *sql.DB
	dispatcher *domainevent.Dispatcher
}

func NewUserRepo(
	db *sql.DB,
	dispatcher *domainevent.Dispatcher,
) repo.User {
	return &userRepo{
		db:         db,
		dispatcher: dispatcher,
	}
}

type userScanner interface {
	Scan(dest ...any) error
}

func (repo *userRepo) FindByID(ctx context.Context, userID vo.UserID) (*entity.User, error) {
	const query = `
		SELECT
			users.id,
			users.username,
			users.full_name,
			users.email,
			users.avatar_key,
			users.create_time,
			users.update_time,

			email_change_requests.value,
			email_change_requests.time,
			email_change_requests.expire_time,

			avatar_change_requests.value,
			avatar_change_requests.time,
			avatar_change_requests.expire_time,

			local_auth_strategies.password_hash,

			google_auth_strategies.sub
		FROM users
		LEFT JOIN email_change_requests
			ON email_change_requests.user_id = users.id
		LEFT JOIN avatar_change_requests
			ON avatar_change_requests.user_id = users.id
		LEFT JOIN local_auth_strategies
			ON local_auth_strategies.user_id = users.id
		LEFT JOIN google_auth_strategies
			ON google_auth_strategies.user_id = users.id
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

func (repo *userRepo) FindByIDs(ctx context.Context, userIDs []vo.UserID) ([]*entity.User, error) {
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
			users.full_name,
			users.email,
			users.avatar_key,
			users.create_time,
			users.update_time,

			email_change_requests.value,
			email_change_requests.time,
			email_change_requests.expire_time,

			avatar_change_requests.value,
			avatar_change_requests.time,
			avatar_change_requests.expire_time,

			local_auth_strategies.password_hash,

			google_auth_strategies.sub
		FROM users
		LEFT JOIN email_change_requests
			ON email_change_requests.user_id = users.id
		LEFT JOIN avatar_change_requests
			ON avatar_change_requests.user_id = users.id
		LEFT JOIN local_auth_strategies
			ON local_auth_strategies.user_id = users.id
		LEFT JOIN google_auth_strategies
			ON google_auth_strategies.user_id = users.id
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

func (repo *userRepo) FindByUsername(ctx context.Context, username vo.Username) (*entity.User, error) {
	const query = `
		SELECT
			users.id,
			users.username,
			users.full_name,
			users.email,
			users.avatar_key,
			users.create_time,
			users.update_time,

			email_change_requests.value,
			email_change_requests.time,
			email_change_requests.expire_time,

			avatar_change_requests.value,
			avatar_change_requests.time,
			avatar_change_requests.expire_time,

			local_auth_strategies.password_hash,

			google_auth_strategies.sub
		FROM users
		LEFT JOIN email_change_requests
			ON email_change_requests.user_id = users.id
		LEFT JOIN avatar_change_requests
			ON avatar_change_requests.user_id = users.id
		LEFT JOIN local_auth_strategies
			ON local_auth_strategies.user_id = users.id
		LEFT JOIN google_auth_strategies
			ON google_auth_strategies.user_id = users.id
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

func (repo *userRepo) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	const query = `
		SELECT
			users.id,
			users.username,
			users.full_name,
			users.email,
			users.avatar_key,
			users.create_time,
			users.update_time,

			email_change_requests.value,
			email_change_requests.time,
			email_change_requests.expire_time,

			avatar_change_requests.value,
			avatar_change_requests.time,
			avatar_change_requests.expire_time,

			local_auth_strategies.password_hash,

			google_auth_strategies.sub
		FROM users
		LEFT JOIN email_change_requests
			ON email_change_requests.user_id = users.id
		LEFT JOIN avatar_change_requests
			ON avatar_change_requests.user_id = users.id
		LEFT JOIN local_auth_strategies
			ON local_auth_strategies.user_id = users.id
		LEFT JOIN google_auth_strategies
			ON google_auth_strategies.user_id = users.id
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

func (repo *userRepo) Save(ctx context.Context, user *entity.User) error {
	if user.DeleteTime() != nil {
		const query = `
			DELETE FROM users
			WHERE id = ?
		`

		if _, err := executor(ctx, repo.db).ExecContext(ctx, query, user.ID().Value()); err != nil {
			return err
		}
	} else {
		const query = `
			INSERT INTO users (
				id,
				username,
				full_name,
				email,
				avatar_key,
				create_time,
				update_time
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				username = ?,
				full_name = ?,
				email = ?,
				avatar_key = ?,
				update_time = ?
		`

		var email any
		if user.Email() != nil {
			email = user.Email().Value()
		}

		var avatarKey any
		if user.AvatarKey() != nil {
			avatarKey = user.AvatarKey().String()
		}

		if _, err := executor(ctx, repo.db).ExecContext(ctx, query,
			user.ID().Value(),
			user.Username().Value(),
			user.FullName().Value(),
			email,
			avatarKey,
			user.CreateTime(),
			user.UpdateTime(),

			user.Username().Value(),
			user.FullName().Value(),
			email,
			avatarKey,
			user.UpdateTime(),
		); err != nil {
			return err
		}

		if user.EmailChangeRequest() != nil {
			const emailChangeRequestQuery = `
				INSERT INTO email_change_requests (
					user_id,
					value,
					time,
					expire_time
				)
				VALUES (?, ?, ?, ?)
				ON DUPLICATE KEY UPDATE
					value = ?,
					time = ?,
					expire_time = ?
			`

			emailChangeRequest := user.EmailChangeRequest()

			if _, err := executor(ctx, repo.db).ExecContext(ctx, emailChangeRequestQuery,
				user.ID().Value(),
				emailChangeRequest.Value().Value(),
				emailChangeRequest.Time(),
				emailChangeRequest.ExpireTime(),

				emailChangeRequest.Value().Value(),
				emailChangeRequest.Time(),
				emailChangeRequest.ExpireTime(),
			); err != nil {
				return err
			}
		} else {
			const emailChangeRequestQuery = `
				DELETE FROM email_change_requests
				WHERE user_id = ?
			`

			if _, err := executor(ctx, repo.db).ExecContext(
				ctx,
				emailChangeRequestQuery,
				user.ID().Value(),
			); err != nil {
				return err
			}
		}

		if user.AvatarChangeRequest() != nil {
			const avatarChangeRequestQuery = `
				INSERT INTO avatar_change_requests (
					user_id,
					value,
					time,
					expire_time
				)
				VALUES (?, ?, ?, ?)
				ON DUPLICATE KEY UPDATE
					value = ?,
					time = ?,
					expire_time = ?
			`

			avatarChangeRequest := user.AvatarChangeRequest()

			if _, err := executor(ctx, repo.db).ExecContext(ctx, avatarChangeRequestQuery,
				user.ID().Value(),
				avatarChangeRequest.Value().String(),
				avatarChangeRequest.Time(),
				avatarChangeRequest.ExpireTime(),

				avatarChangeRequest.Value().String(),
				avatarChangeRequest.Time(),
				avatarChangeRequest.ExpireTime(),
			); err != nil {
				return err
			}
		} else {
			const avatarChangeRequestQuery = `
				DELETE FROM avatar_change_requests
				WHERE user_id = ?
			`

			if _, err := executor(ctx, repo.db).ExecContext(
				ctx,
				avatarChangeRequestQuery,
				user.ID().Value(),
			); err != nil {
				return err
			}
		}

		if user.LocalAuthStrategy() != nil {
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
				user.ID().Value(),
				user.LocalAuthStrategy().PasswordHash().Value(),

				user.LocalAuthStrategy().PasswordHash().Value(),
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
				user.ID().Value(),
			); err != nil {
				return err
			}
		}

		if user.GoogleAuthStrategy() != nil {
			const googleAuthStrategyQuery = `
				INSERT INTO google_auth_strategies (
					user_id,
					sub
				)
				VALUES (?, ?)
				ON DUPLICATE KEY UPDATE
					sub = ?
			`

			if _, err := executor(ctx, repo.db).ExecContext(ctx, googleAuthStrategyQuery,
				user.ID().Value(),
				user.GoogleAuthStrategy().Sub.Value(),

				user.GoogleAuthStrategy().Sub.Value(),
			); err != nil {
				return err
			}
		} else {
			const googleAuthStrategyQuery = `
				DELETE FROM google_auth_strategies
				WHERE user_id = ?
			`

			if _, err := executor(ctx, repo.db).ExecContext(
				ctx,
				googleAuthStrategyQuery,
				user.ID().Value(),
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
		id         string
		username   string
		fullName   string
		email      sql.NullString
		avatarKey  sql.NullString
		createTime time.Time
		updateTime time.Time

		emailChangeRequestValue      sql.NullString
		emailChangeRequestTime       sql.NullTime
		emailChangeRequestExpireTime sql.NullTime

		avatarChangeRequestValue      sql.NullString
		avatarChangeRequestTime       sql.NullTime
		avatarChangeRequestExpireTime sql.NullTime

		passwordHash sql.NullString

		googleSub sql.NullString
	)

	if err := scanner.Scan(
		&id,
		&username,
		&fullName,
		&email,
		&avatarKey,
		&createTime,
		&updateTime,

		&emailChangeRequestValue,
		&emailChangeRequestTime,
		&emailChangeRequestExpireTime,

		&avatarChangeRequestValue,
		&avatarChangeRequestTime,
		&avatarChangeRequestExpireTime,

		&passwordHash,

		&googleSub,
	); err != nil {
		return nil, err
	}

	usernameVO, _ := vo.NewUsername(username)
	fullNameVO, _ := vo.NewFullName(fullName)

	var emailVO *vo.Email
	if email.Valid {
		value, _ := vo.NewEmail(email.String)
		emailVO = &value
	}

	var avatarKeyVO *vo.AvatarKey
	if avatarKey.Valid {
		value, _ := vo.NewAvatarKey(avatarKey.String)
		avatarKeyVO = &value
	}

	var emailChangeRequest *vo.EmailChangeRequest
	if emailChangeRequestValue.Valid {
		newEmail, _ := vo.NewEmail(emailChangeRequestValue.String)

		value := vo.NewEmailChangeRequest(
			newEmail,
			emailChangeRequestTime.Time,
			emailChangeRequestExpireTime.Time,
		)
		emailChangeRequest = &value
	}

	var avatarChangeRequest *vo.AvatarChangeRequest
	if avatarChangeRequestValue.Valid {
		newAvatarKey, _ := vo.NewAvatarKey(avatarChangeRequestValue.String)

		value := vo.NewAvatarChangeRequest(
			newAvatarKey,
			avatarChangeRequestTime.Time,
			avatarChangeRequestExpireTime.Time,
		)
		avatarChangeRequest = &value
	}

	var localAuthStrategy *vo.LocalAuthStrategy
	if passwordHash.Valid {
		passwordHashVO, _ := vo.NewPasswordHash(passwordHash.String)

		value := vo.NewLocalAuthStrategy(passwordHashVO)
		localAuthStrategy = &value
	}

	var googleAuthStrategy *vo.GoogleAuthStrategy
	if googleSub.Valid {
		googleSubVO, _ := vo.NewGoogleSub(googleSub.String)

		value := vo.NewGoogleAuthStrategy(googleSubVO)
		googleAuthStrategy = &value
	}

	return entity.ReconstituteUser(
		vo.NewUserID(id),
		usernameVO,
		fullNameVO,
		emailVO,
		avatarKeyVO,
		createTime,
		updateTime,
		nil,
		emailChangeRequest,
		avatarChangeRequest,
		localAuthStrategy,
		googleAuthStrategy,
	), nil
}
