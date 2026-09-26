package mysql

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/velonyapp/identity/internal/application/domainevent"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ repo.Session = (*SessionRepo)(nil)

type SessionRepo struct {
	db         *sql.DB
	dispatcher *domainevent.Dispatcher
}

func NewSessionRepo(
	db *sql.DB,
	dispatcher *domainevent.Dispatcher,
) repo.Session {
	return &SessionRepo{
		db:         db,
		dispatcher: dispatcher,
	}
}

type sessionScanner interface {
	Scan(dest ...any) error
}

func (repo *SessionRepo) FindByToken(ctx context.Context, token vo.SessionToken) (*entity.Session, error) {
	const query = `
		SELECT
			id,
			user_id,
			expire_time,
			revoke_time
		FROM sessions
		WHERE token = ?
		LIMIT 1
	`

	tokenHash := sha256.Sum256([]byte(token.Value()))

	row := executor(ctx, repo.db).QueryRowContext(ctx, query, tokenHash[:])

	session, err := scanSession(row, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return session, nil
}

func (repo *SessionRepo) Save(ctx context.Context, session *entity.Session) error {
	const query = `
		INSERT INTO sessions (
			id,
			user_id,
			token,
			expire_time,
			revoke_time
		)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			user_id = ?,
			token = ?,
			expire_time = ?,
			revoke_time = ?
	`

	var revokeTime any
	if session.RevokeTime() != nil {
		revokeTime = *session.RevokeTime()
	}

	tokenHash := sha256.Sum256([]byte(session.Token().Value()))

	if _, err := executor(ctx, repo.db).ExecContext(ctx, query,
		session.ID().Value(),
		session.UserID().Value(),
		tokenHash[:],
		session.ExpireTime(),
		revokeTime,

		session.UserID().Value(),
		tokenHash[:],
		session.ExpireTime(),
		revokeTime,
	); err != nil {
		return err
	}

	for _, domainEvent := range session.PullEvents() {
		if err := repo.dispatcher.Dispatch(ctx, domainEvent); err != nil {
			return err
		}
	}

	return nil
}

func scanSession(scanner sessionScanner, token vo.SessionToken) (*entity.Session, error) {
	var (
		id         string
		userID     string
		expireTime time.Time
		revokeTime sql.NullTime
	)

	if err := scanner.Scan(
		&id,
		&userID,
		&expireTime,
		&revokeTime,
	); err != nil {
		return nil, err
	}

	var revokeTimeVO *time.Time
	if revokeTime.Valid {
		value := revokeTime.Time
		revokeTimeVO = &value
	}

	return entity.ReconstituteSession(
		vo.NewSessionID(id),
		vo.NewUserID(userID),
		token,
		expireTime,
		revokeTimeVO,
	), nil
}
