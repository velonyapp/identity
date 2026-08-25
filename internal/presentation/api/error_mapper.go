package api

import (
	"errors"

	v1 "github.com/velony-app/identity/gen/api/v1"
	"github.com/velony-app/identity/internal/application/command"
	"github.com/velony-app/identity/internal/application/common"
	"github.com/velony-app/identity/internal/domain/entity"
	"github.com/velony-app/identity/internal/domain/vo"

	kerrors "github.com/go-kratos/kratos/v3/errors"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, vo.ErrFullNameTooShort):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_FULL_NAME.String(),
			vo.ErrFullNameTooShort.Error(),
		)

	case errors.Is(err, vo.ErrFullNameTooLong):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_FULL_NAME.String(),
			vo.ErrFullNameTooLong.Error(),
		)

	case errors.Is(err, vo.ErrUsernameTooShort):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_USERNAME.String(),
			vo.ErrUsernameTooShort.Error(),
		)

	case errors.Is(err, vo.ErrUsernameTooLong):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_USERNAME.String(),
			vo.ErrUsernameTooLong.Error(),
		)

	case errors.Is(err, vo.ErrUsernameInvalidCharacter):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_USERNAME.String(),
			vo.ErrUsernameInvalidCharacter.Error(),
		)

	case errors.Is(err, vo.ErrInvalidEmail):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_EMAIL.String(),
			vo.ErrInvalidEmail.Error(),
		)

	case errors.Is(err, vo.ErrPasswordTooShort):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_PASSWORD.String(),
			vo.ErrPasswordTooShort.Error(),
		)

	case errors.Is(err, vo.ErrPasswordTooLong):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_PASSWORD.String(),
			vo.ErrPasswordTooLong.Error(),
		)

	case errors.Is(err, vo.ErrPasswordMissingUppercase):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_PASSWORD.String(),
			vo.ErrPasswordMissingUppercase.Error(),
		)

	case errors.Is(err, vo.ErrPasswordMissingLowercase):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_PASSWORD.String(),
			vo.ErrPasswordMissingLowercase.Error(),
		)

	case errors.Is(err, vo.ErrPasswordMissingNumber):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_PASSWORD.String(),
			vo.ErrPasswordMissingNumber.Error(),
		)

	case errors.Is(err, vo.ErrPasswordMissingSymbol):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_PASSWORD.String(),
			vo.ErrPasswordMissingSymbol.Error(),
		)

	case errors.Is(err, vo.ErrPasswordContainsWhitespace):
		return kerrors.BadRequest(
			v1.ErrorReason_INVALID_PASSWORD.String(),
			vo.ErrPasswordContainsWhitespace.Error(),
		)

	case errors.Is(err, entity.ErrUserDeleted):
		return kerrors.NotFound(
			v1.ErrorReason_USER_NOT_FOUND.String(),
			common.ErrUserNotFound.Error(),
		)

	case errors.Is(err, common.ErrUserNotFound):
		return kerrors.NotFound(
			v1.ErrorReason_USER_NOT_FOUND.String(),
			common.ErrUserNotFound.Error(),
		)

	case errors.Is(err, entity.ErrSessionExpired):
		return kerrors.Unauthorized(
			v1.ErrorReason_INVALID_REFRESH_TOKEN.String(),
			command.ErrInvalidRefreshToken.Error(),
		)

	case errors.Is(err, entity.ErrSessionRevoked):
		return kerrors.Unauthorized(
			v1.ErrorReason_INVALID_REFRESH_TOKEN.String(),
			command.ErrInvalidRefreshToken.Error(),
		)

	case errors.Is(err, command.ErrInvalidRefreshToken):
		return kerrors.Unauthorized(
			v1.ErrorReason_INVALID_REFRESH_TOKEN.String(),
			command.ErrInvalidRefreshToken.Error(),
		)

	default:
		return err // let this be i want to debug
		// return kerrors.InternalServer(
		// 	"",
		// 	"internal server error",
		// )
	}
}
