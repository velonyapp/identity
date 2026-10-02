package mysql

import (
	"context"
	"database/sql"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ repo.Avatar = (*avatarRepo)(nil)

type avatarRepo struct {
	db *sql.DB
}

func NewAvatarRepo(
	db *sql.DB,
) repo.Avatar {
	return &avatarRepo{
		db: db,
	}
}

type repoAvatarScanner interface {
	Scan(dest ...any) error
}

func (repo *avatarRepo) FindByID(ctx context.Context, avatarID vo.AvatarID) (*entity.Avatar, error) {
	const query = `
		SELECT
			avatars.id,
			avatars.user_id,
			avatars.key
		FROM avatars
		WHERE avatars.id = ?
		LIMIT 1
	`

	row := executor(ctx, repo.db).QueryRowContext(ctx, query, avatarID.Value())

	avatar, err := scanRepoAvatar(row)
	if err != nil {
		return nil, err
	}

	return avatar, nil
}

func (repo *avatarRepo) Save(ctx context.Context, avatar *entity.Avatar) error {
	if avatar.IsDeleted() {
		const query = `
			DELETE FROM avatars
			WHERE id = ?
		`

		if _, err := executor(ctx, repo.db).ExecContext(ctx, query, avatar.ID().Value()); err != nil {
			return err
		}

		return nil
	}

	const query = `
		INSERT INTO avatars (
			id,
			user_id,
			key
		)
		VALUES (?, ?, ?)
	`

	if _, err := executor(ctx, repo.db).ExecContext(ctx, query,
		avatar.ID().Value(),
		avatar.UserID().Value(),
		avatar.Key().Value(),
	); err != nil {
		return err
	}

	return nil
}

func scanRepoAvatar(scanner repoAvatarScanner) (*entity.Avatar, error) {
	var (
		id     string
		userID string
		key    string
	)

	if err := scanner.Scan(
		&id,
		&userID,
		&key,
	); err != nil {
		return nil, err
	}

	keyVO, _ := vo.NewAvatarKey(key)

	return entity.ReconstituteAvatar(
		vo.NewAvatarID(id),
		vo.NewUserID(userID),
		keyVO,
		nil,
	), nil
}
