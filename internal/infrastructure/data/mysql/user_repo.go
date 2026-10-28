package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
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

type repoUserScanner interface {
	Scan(dest ...any) error
}

func (r *userRepo) FindByID(ctx context.Context, userID vo.UserID) (*entity.User, error) {
	const query = `
		SELECT
			users.id,
			users.username,
			users.email,
			users.avatar_id,
			users.create_time,
			users.version,

			email_change_requests.email,
			email_change_requests.token,
			email_change_requests.create_time,
			email_change_requests.expire_time,

			avatar_change_requests.avatar_id,
			avatar_change_requests.create_time,
			avatar_change_requests.expire_time,

			local_auth_strategies.password_hash
		FROM users
		LEFT JOIN email_change_requests
			ON email_change_requests.user_id = users.id
		LEFT JOIN avatar_change_requests
			ON avatar_change_requests.user_id = users.id
		LEFT JOIN local_auth_strategies
			ON local_auth_strategies.user_id = users.id
		WHERE users.id = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, userID.String())

	user, err := scanRepoUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username vo.Username) (*entity.User, error) {
	const query = `
		SELECT
			users.id,
			users.username,
			users.email,
			users.avatar_id,
			users.create_time,
			users.version,

			email_change_requests.email,
			email_change_requests.token,
			email_change_requests.create_time,
			email_change_requests.expire_time,

			avatar_change_requests.avatar_id,
			avatar_change_requests.create_time,
			avatar_change_requests.expire_time,

			local_auth_strategies.password_hash
		FROM users
		LEFT JOIN email_change_requests
			ON email_change_requests.user_id = users.id
		LEFT JOIN avatar_change_requests
			ON avatar_change_requests.user_id = users.id
		LEFT JOIN local_auth_strategies
			ON local_auth_strategies.user_id = users.id
		WHERE users.username = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, username.String())

	user, err := scanRepoUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	const query = `
		SELECT
			users.id,
			users.username,
			users.email,
			users.avatar_id,
			users.create_time,
			users.version,

			email_change_requests.email,
			email_change_requests.token,
			email_change_requests.create_time,
			email_change_requests.expire_time,

			avatar_change_requests.avatar_id,
			avatar_change_requests.create_time,
			avatar_change_requests.expire_time,

			local_auth_strategies.password_hash
		FROM users
		LEFT JOIN email_change_requests
			ON email_change_requests.user_id = users.id
		LEFT JOIN avatar_change_requests
			ON avatar_change_requests.user_id = users.id
		LEFT JOIN local_auth_strategies
			ON local_auth_strategies.user_id = users.id
		WHERE users.email = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, email.String())

	user, err := scanRepoUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *userRepo) Save(ctx context.Context, user *entity.User) error {
	if user.IsDeleted() {
		if !user.IsNew() {
			const query = `
				DELETE FROM users
				WHERE id = ?
			`

			_, err := executor(ctx, r.db).ExecContext(ctx, query, user.ID().String())
			if err != nil {
				return err
			}
		}
	} else {
		var email any
		if user.Email() != nil {
			email = user.Email().String()
		}
		var avatarID any
		if user.AvatarID() != nil {
			avatarID = user.AvatarID().String()
		}

		if user.IsNew() {
			const query = `
				INSERT INTO users (
					id,
					username,
					email,
					avatar_id,
					create_time,
					version
				)
				VALUES (?, ?, ?, ?, ?, ?)
			`

			if _, err := executor(ctx, r.db).ExecContext(ctx, query,
				user.ID().String(),
				user.Username().String(),
				email,
				avatarID,
				user.CreateTime(),
				user.Version(),
			); err != nil {
				var mysqlErr *mysql.MySQLError
				if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
					switch {
					case strings.Contains(mysqlErr.Message, "uq_users_username"):
						return repo.ErrUsernameAlreadyExists
					case strings.Contains(mysqlErr.Message, "uq_users_email"):
						return repo.ErrEmailAlreadyExists
					}
				}

				return err
			}
		} else {
			const query = `
			    UPDATE users
			    SET
			        username = ?,
			        email = ?,
			        avatar_id = ?,
			        version = version + 1
			    WHERE id = ?
			        AND version = ?
			`

			result, err := executor(ctx, r.db).ExecContext(ctx, query,
				user.Username().String(),
				email,
				avatarID,
				user.ID().String(),
				user.Version(),
			)
			if err != nil {
				return err
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if rowsAffected == 0 {
				return repo.ErrUserVersionConflict
			}
		}

		if err := r.saveEmailRequestChange(ctx, user); err != nil {
			return err
		}
		if err := r.saveAvatarRequestChange(ctx, user); err != nil {
			return err
		}
		if err := r.saveLocalAuthStrategy(ctx, user); err != nil {
			return err
		}
	}

	for _, domainEvent := range user.PullEvents() {
		if err := r.dispatcher.Dispatch(ctx, domainEvent); err != nil {
			return err
		}
	}

	user.MarkPersisted()

	return nil
}

func (r *userRepo) saveEmailRequestChange(ctx context.Context, user *entity.User) error {
	if user.EmailChangeRequest() == nil {
		const query = `
			DELETE FROM email_change_requests
			WHERE user_id = ?
		`

		if _, err := executor(ctx, r.db).ExecContext(ctx, query,
			user.ID().String(),
		); err != nil {
			return err
		}

		return nil
	}

	const query = `
		INSERT INTO email_change_requests (
			user_id,
			email,
			token,
			create_time,
			expire_time
		)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			email = ?,
			token = ?,
			create_time = ?,
			expire_time = ?
	`

	emailChangeRequest := user.EmailChangeRequest()

	if _, err := executor(ctx, r.db).ExecContext(ctx, query,
		user.ID().String(),
		emailChangeRequest.Email().String(),
		emailChangeRequest.Token().String(),
		emailChangeRequest.CreateTime(),
		emailChangeRequest.ExpireTime(),

		emailChangeRequest.Email().String(),
		emailChangeRequest.Token().String(),
		emailChangeRequest.CreateTime(),
		emailChangeRequest.ExpireTime(),
	); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) saveAvatarRequestChange(ctx context.Context, user *entity.User) error {
	if user.AvatarChangeRequest() == nil {
		const query = `
			DELETE FROM avatar_change_requests
			WHERE user_id = ?
		`

		if _, err := executor(ctx, r.db).ExecContext(ctx, query,
			user.ID().String(),
		); err != nil {
			return err
		}

		return nil
	}

	const query = `
		INSERT INTO avatar_change_requests (
			user_id,
			avatar_id,
			create_time,
			expire_time
		)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			avatar_id = ?,
			create_time = ?,
			expire_time = ?
	`

	avatarChangeRequest := user.AvatarChangeRequest()

	if _, err := executor(ctx, r.db).ExecContext(ctx, query,
		user.ID().String(),
		avatarChangeRequest.AvatarID().String(),
		avatarChangeRequest.CreateTime(),
		avatarChangeRequest.ExpireTime(),

		avatarChangeRequest.AvatarID().String(),
		avatarChangeRequest.CreateTime(),
		avatarChangeRequest.ExpireTime(),
	); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) saveLocalAuthStrategy(ctx context.Context, user *entity.User) error {
	if user.LocalAuthStrategy() == nil {
		const query = `
				DELETE FROM local_auth_strategies
				WHERE user_id = ?
			`

		if _, err := executor(ctx, r.db).ExecContext(ctx, query,
			user.ID().String(),
		); err != nil {
			return err
		}

		return nil
	}

	const query = `
		INSERT INTO local_auth_strategies (
			user_id,
			password_hash
		)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
			password_hash = ?
	`

	if _, err := executor(ctx, r.db).ExecContext(ctx, query,
		user.ID().String(),
		user.LocalAuthStrategy().PasswordHash().Bytes(),

		user.LocalAuthStrategy().PasswordHash().Bytes(),
	); err != nil {
		return err
	}

	return nil
}

func scanRepoUser(scanner repoUserScanner) (*entity.User, error) {
	var (
		idRaw       string
		usernameRaw string
		emailRaw    sql.NullString
		avatarIDRaw sql.NullString
		createTime  time.Time
		version     uint32

		emailChangeRequestEmailRaw      sql.NullString
		emailChangeRequestTokenRaw      sql.NullString
		emailChangeRequestCreateTimeRaw sql.NullTime
		emailChangeRequestExpireTimeRaw sql.NullTime

		avatarChangeRequestAvatarIDRaw   sql.NullString
		avatarChangeRequestCreateTimeRaw sql.NullTime
		avatarChangeRequestExpireTimeRaw sql.NullTime

		passwordHashRaw []byte
	)

	if err := scanner.Scan(
		&idRaw,
		&usernameRaw,
		&emailRaw,
		&avatarIDRaw,
		&createTime,
		&version,

		&emailChangeRequestEmailRaw,
		&emailChangeRequestTokenRaw,
		&emailChangeRequestCreateTimeRaw,
		&emailChangeRequestExpireTimeRaw,

		&avatarChangeRequestAvatarIDRaw,
		&avatarChangeRequestCreateTimeRaw,
		&avatarChangeRequestExpireTimeRaw,

		&passwordHashRaw,
	); err != nil {
		return nil, err
	}

	id, _ := vo.NewUserID(idRaw)
	username, _ := vo.NewUsername(usernameRaw)
	var email *vo.Email
	if emailRaw.Valid {
		value, _ := vo.NewEmail(emailRaw.String)
		email = &value
	}
	var avatarID *vo.AvatarID
	if avatarIDRaw.Valid {
		value, _ := vo.NewAvatarID(avatarIDRaw.String)
		avatarID = &value
	}
	var emailChangeRequest *vo.EmailChangeRequest
	if emailChangeRequestEmailRaw.Valid {
		email, _ := vo.NewEmail(emailChangeRequestEmailRaw.String)
		token, _ := vo.NewEmailChangeRequestToken(emailChangeRequestTokenRaw.String)

		value := vo.NewEmailChangeRequest(
			email,
			token,
			emailChangeRequestCreateTimeRaw.Time,
			emailChangeRequestExpireTimeRaw.Time,
		)
		emailChangeRequest = &value
	}
	var avatarChangeRequest *vo.AvatarChangeRequest
	if avatarChangeRequestAvatarIDRaw.Valid {
		avatarID, _ := vo.NewAvatarID(avatarChangeRequestAvatarIDRaw.String)

		value := vo.NewAvatarChangeRequest(
			avatarID,
			avatarChangeRequestCreateTimeRaw.Time,
			avatarChangeRequestExpireTimeRaw.Time,
		)
		avatarChangeRequest = &value
	}
	var localAuthStrategy *vo.LocalAuthStrategy
	if passwordHashRaw != nil {
		passwordHash, _ := vo.NewPasswordHash(passwordHashRaw)

		value := vo.NewLocalAuthStrategy(passwordHash)
		localAuthStrategy = &value
	}

	return entity.ReconstituteUser(
		id,
		username,
		email,
		avatarID,
		createTime,
		nil,
		emailChangeRequest,
		avatarChangeRequest,
		localAuthStrategy,
		version,
	), nil
}
