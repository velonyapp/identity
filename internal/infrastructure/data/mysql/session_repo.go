package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ repo.Session = (*sessionRepo)(nil)

type sessionRepo struct {
	db  *sql.DB
	log *slog.Logger
}

func NewSessionRepo(
	db *sql.DB,
	log *slog.Logger,
) repo.Session {
	return &sessionRepo{
		db:  db,
		log: log,
	}
}

type repoSessionScanner interface {
	Scan(dest ...any) error
}

func (r *sessionRepo) FindByTokenHash(ctx context.Context, tokenHash vo.SessionTokenHash) (*entity.Session, error) {
	const query = `
		SELECT
			id,
			user_id,
			token_hash,
			create_time,
			expire_time,
			revoke_time,
			version
		FROM sessions
		WHERE token_hash = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, tokenHash.Bytes())

	session, err := scanRepoSession(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		r.log.ErrorContext(ctx,
			"failed to find session by token hash",
			slog.Any("error", err),
		)
		return nil, err
	}

	return session, nil
}

func (r *sessionRepo) Save(ctx context.Context, session *entity.Session) error {
	if !session.IsDirty() {
		return nil
	}

	var revokeTime any
	if session.RevokeTime() != nil {
		revokeTime = *session.RevokeTime()
	}

	if session.IsNew() {
		const query = `
			INSERT INTO sessions (
				id,
				user_id,
				token_hash,
				create_time,
				expire_time,
				revoke_time,
				version
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`

		if _, err := executor(ctx, r.db).ExecContext(ctx, query,
			session.ID().String(),
			session.UserID().String(),
			session.TokenHash().Bytes(),
			session.CreateTime(),
			session.ExpireTime(),
			revokeTime,
			session.Version(),
		); err != nil {
			r.log.ErrorContext(ctx,
				"failed to insert session",
				slog.String("session_id", session.ID().String()),
				slog.String("user_id", session.UserID().String()),
				slog.Any("error", err),
			)
			return err
		}
	} else {
		const query = `
		    UPDATE sessions
		    SET
		        token_hash = ?,
		        expire_time = ?,
		        revoke_time = ?,
		        version = ?
		    WHERE id = ?
		        AND version = ?
		`

		result, err := executor(ctx, r.db).ExecContext(ctx, query,
			session.TokenHash().Bytes(),
			session.ExpireTime(),
			revokeTime,
			session.Version(),
			session.ID().String(),
			session.Version()-1,
		)
		if err != nil {
			r.log.ErrorContext(ctx,
				"failed to update session",
				slog.String("session_id", session.ID().String()),
				slog.String("user_id", session.UserID().String()),
				slog.Any("error", err),
			)
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			r.log.ErrorContext(ctx,
				"failed to get affected rows after updating session",
				slog.String("session_id", session.ID().String()),
				slog.Any("error", err),
			)
			return err
		}
		if rowsAffected == 0 {
			return repo.ErrSessionVersionConflict
		}
	}

	session.MarkPersisted()

	return nil
}

func scanRepoSession(scanner repoSessionScanner) (*entity.Session, error) {
	var (
		idRaw         string
		userIDRaw     string
		tokenHashRaw  []byte
		createTimeRaw time.Time
		expireTimeRaw time.Time
		revokeTimeRaw sql.NullTime
		versionRaw    uint32
	)

	if err := scanner.Scan(
		&idRaw,
		&userIDRaw,
		&tokenHashRaw,
		&createTimeRaw,
		&expireTimeRaw,
		&revokeTimeRaw,
		&versionRaw,
	); err != nil {
		return nil, err
	}

	id, _ := vo.NewSessionID(idRaw)
	userID, _ := vo.NewUserID(userIDRaw)
	tokenHash, _ := vo.NewSessionTokenHash(tokenHashRaw)
	createTime := createTimeRaw
	expireTime := expireTimeRaw
	var revokeTime *time.Time
	if revokeTimeRaw.Valid {
		value := revokeTimeRaw.Time
		revokeTime = &value
	}
	version := versionRaw

	return entity.ReconstituteSession(
		id,
		userID,
		tokenHash,
		createTime,
		expireTime,
		revokeTime,
		version,
	), nil
}
