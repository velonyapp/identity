package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/velonyapp/identity/internal/application/domainevent"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ repo.Session = (*sessionRepo)(nil)

type sessionRepo struct {
	db         *sql.DB
	dispatcher *domainevent.Dispatcher
}

func NewSessionRepo(
	db *sql.DB,
	dispatcher *domainevent.Dispatcher,
) repo.Session {
	return &sessionRepo{
		db:         db,
		dispatcher: dispatcher,
	}
}

type sessionScanner interface {
	Scan(dest ...any) error
}

func (repo *sessionRepo) FindByTokenHash(ctx context.Context, tokenHash vo.SessionTokenHash) (*entity.Session, error) {
	const query = `
		SELECT
			id,
			user_id,
			token_hash,
			expire_time,
			revoke_time
		FROM sessions
		WHERE token_hash = ?
		LIMIT 1
	`

	row := executor(ctx, repo.db).QueryRowContext(ctx, query, tokenHash.Value())

	session, err := scanSession(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return session, nil
}

func (repo *sessionRepo) Save(ctx context.Context, session *entity.Session) error {
	const query = `
		INSERT INTO sessions (
			id,
			user_id,
			token_hash,
			expire_time,
			revoke_time
		)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			token_hash = ?,
			expire_time = ?,
			revoke_time = ?
	`

	var revokeTime any
	if session.RevokeTime() != nil {
		revokeTime = *session.RevokeTime()
	}

	if _, err := executor(ctx, repo.db).ExecContext(ctx, query,
		session.ID().Value(),
		session.UserID().Value(),
		session.TokenHash().Value(),
		session.ExpireTime(),
		revokeTime,

		session.TokenHash().Value(),
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

func scanSession(scanner sessionScanner) (*entity.Session, error) {
	var (
		id         string
		userID     string
		tokenHash  []byte
		expireTime time.Time
		revokeTime sql.NullTime
	)

	if err := scanner.Scan(
		&id,
		&userID,
		&tokenHash,
		&expireTime,
		&revokeTime,
	); err != nil {
		return nil, err
	}

	tokenHashVO, _ := vo.NewSessionTokenHash(tokenHash)

	var revokeTimeVO *time.Time
	if revokeTime.Valid {
		value := revokeTime.Time
		revokeTimeVO = &value
	}

	return entity.ReconstituteSession(
		vo.NewSessionID(id),
		vo.NewUserID(userID),
		tokenHashVO,
		expireTime,
		revokeTimeVO,
	), nil
}
