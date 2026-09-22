package mysql

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type EmailChangeRequestRepo struct {
	db *sql.DB
}

func NewEmailChangeRequestRepo(
	db *sql.DB,
) repo.EmailChangeRequest {
	return &EmailChangeRequestRepo{
		db: db,
	}
}

type emailChangeRequestScanner interface {
	Scan(dest ...any) error
}

func (repo *EmailChangeRequestRepo) FindByToken(ctx context.Context, token vo.RequestToken) (*entity.EmailChangeRequest, error) {
	const query = `
		SELECT
			id,
			user_id,
			new_email,
			create_time,
			expire_time,
			confirm_time
		FROM email_change_requests
		WHERE token_hash = ?
		LIMIT 1
	`

	tokenHash := sha256.Sum256([]byte(token.Value()))

	row := executor(ctx, repo.db).QueryRowContext(ctx, query, tokenHash[:])

	request, err := scanEmailChangeRequest(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return request, nil
}

func (repo *EmailChangeRequestRepo) Save(ctx context.Context, request *entity.EmailChangeRequest) error {
	var confirmTime any
	if request.ConfirmTime != nil {
		confirmTime = request.ConfirmTime.Value()
	}

	if request.PendingToken != nil {
		const query = `
			INSERT INTO email_change_requests (
				id,
				token_hash,
				user_id,
				new_email,
				create_time,
				expire_time,
				confirm_time
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				token_hash = ?,
				user_id = ?,
				new_email = ?,
				create_time = ?,
				expire_time = ?,
				confirm_time = ?
		`

		tokenHash := sha256.Sum256([]byte(request.PendingToken.Value()))

		_, err := executor(ctx, repo.db).ExecContext(
			ctx,
			query,
			request.ID.Value(),
			tokenHash[:],
			request.UserID.Value(),
			request.NewEmail.Value(),
			request.CreateTime.Value(),
			request.ExpireTime.Value(),
			confirmTime,

			tokenHash[:],
			request.UserID.Value(),
			request.NewEmail.Value(),
			request.CreateTime.Value(),
			request.ExpireTime.Value(),
			confirmTime,
		)

		return err
	}

	const query = `
		UPDATE email_change_requests
		SET
			user_id = ?,
			new_email = ?,
			create_time = ?,
			expire_time = ?,
			confirm_time = ?
		WHERE id = ?
	`

	_, err := executor(ctx, repo.db).ExecContext(
		ctx,
		query,
		request.UserID.Value(),
		request.NewEmail.Value(),
		request.CreateTime.Value(),
		request.ExpireTime.Value(),
		confirmTime,
		request.ID.Value(),
	)

	return err
}

func scanEmailChangeRequest(scanner emailChangeRequestScanner) (*entity.EmailChangeRequest, error) {
	var (
		id          string
		userID      string
		newEmail    string
		createTime  time.Time
		expireTime  time.Time
		confirmTime sql.NullTime
	)

	if err := scanner.Scan(
		&id,
		&userID,
		&newEmail,
		&createTime,
		&expireTime,
		&confirmTime,
	); err != nil {
		return nil, err
	}

	emailVO, _ := vo.NewEmail(newEmail)

	var confirmTimeVO *vo.Time
	if confirmTime.Valid {
		value := vo.NewTime(confirmTime.Time)
		confirmTimeVO = &value
	}

	return &entity.EmailChangeRequest{
		ID:          vo.NewEmailChangeRequestID(id),
		UserID:      vo.NewUserID(userID),
		NewEmail:    emailVO,
		CreateTime:  vo.NewTime(createTime),
		ExpireTime:  vo.NewTime(expireTime),
		ConfirmTime: confirmTimeVO,
	}, nil
}
