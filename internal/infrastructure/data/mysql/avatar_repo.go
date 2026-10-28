package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ repo.Avatar = (*avatarRepo)(nil)

type avatarRepo struct {
	db *sql.DB
}

func NewAvatarRepo(db *sql.DB) repo.Avatar {
	return &avatarRepo{
		db: db,
	}
}

type repoAvatarScanner interface {
	Scan(dest ...any) error
}

func (r *avatarRepo) FindByID(ctx context.Context, avatarID vo.AvatarID) (*entity.Avatar, error) {
	const query = `
		SELECT
			id,
			user_id,
			object_key
		FROM avatars
		WHERE id = ?
		LIMIT 1
	`

	row := executor(ctx, r.db).QueryRowContext(ctx, query, avatarID.String())

	avatar, err := scanRepoAvatar(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return avatar, nil
}

func (r *avatarRepo) Create(ctx context.Context, avatar *entity.Avatar) error {
	const query = `
		INSERT INTO avatars (
			id,
			user_id,
			object_key
		)
		VALUES (?, ?, ?)
	`

	if _, err := executor(ctx, r.db).ExecContext(ctx, query,
		avatar.ID().String(),
		avatar.UserID().String(),
		avatar.Key().String(),
	); err != nil {
		return err
	}

	return nil
}

func (r *avatarRepo) Delete(ctx context.Context, avatarID vo.AvatarID) error {
	const query = `
		DELETE FROM avatars
		WHERE id = ?
	`

	if _, err := executor(ctx, r.db).ExecContext(ctx, query,
		avatarID.String(),
	); err != nil {
		return err
	}

	return nil
}

func scanRepoAvatar(scanner repoAvatarScanner) (*entity.Avatar, error) {
	var (
		idRaw     string
		userIDRaw string
		keyRaw    string
	)

	if err := scanner.Scan(
		&idRaw,
		&userIDRaw,
		&keyRaw,
	); err != nil {
		return nil, err
	}

	id, _ := vo.NewAvatarID(idRaw)
	userID, _ := vo.NewUserID(userIDRaw)
	key, _ := vo.NewAvatarKey(keyRaw)

	return entity.ReconstituteAvatar(
		id,
		userID,
		key,
	), nil
}
