package api

import (
	"errors"

	apiv1 "github.com/velonyapp/identity/gen/api/v1"
	applicationcommand "github.com/velonyapp/identity/internal/application/command"
	applicationcommon "github.com/velonyapp/identity/internal/application/common"
	domainentity "github.com/velonyapp/identity/internal/domain/entity"
	domainservice "github.com/velonyapp/identity/internal/domain/service"
	domainvo "github.com/velonyapp/identity/internal/domain/vo"

	kerrors "github.com/go-kratos/kratos/v3/errors"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domainvo.ErrFullNameTooShort):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_FULL_NAME.String(),
			domainvo.ErrFullNameTooShort.Error(),
		)

	case errors.Is(err, domainvo.ErrFullNameTooLong):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_FULL_NAME.String(),
			domainvo.ErrFullNameTooLong.Error(),
		)

	case errors.Is(err, domainvo.ErrUsernameTooShort):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_USERNAME.String(),
			domainvo.ErrUsernameTooShort.Error(),
		)

	case errors.Is(err, domainvo.ErrUsernameTooLong):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_USERNAME.String(),
			domainvo.ErrUsernameTooLong.Error(),
		)

	case errors.Is(err, domainvo.ErrUsernameInvalidCharacter):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_USERNAME.String(),
			domainvo.ErrUsernameInvalidCharacter.Error(),
		)

	case errors.Is(err, domainvo.ErrInvalidEmail):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_EMAIL.String(),
			domainvo.ErrInvalidEmail.Error(),
		)

	case errors.Is(err, domainvo.ErrAvatarKeyEmpty):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_AVATAR_KEY.String(),
			domainvo.ErrAvatarKeyEmpty.Error(),
		)

	case errors.Is(err, domainvo.ErrAvatarKeyInvalidFormat):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_AVATAR_KEY.String(),
			domainvo.ErrAvatarKeyInvalidFormat.Error(),
		)

	case errors.Is(err, domainvo.ErrAvatarKeyInvalidUTF8):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_AVATAR_KEY.String(),
			domainvo.ErrAvatarKeyInvalidUTF8.Error(),
		)

	case errors.Is(err, domainvo.ErrAvatarKeyTooLong):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_AVATAR_KEY.String(),
			domainvo.ErrAvatarKeyTooLong.Error(),
		)

	case errors.Is(err, domainvo.ErrPasswordTooShort):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_PASSWORD.String(),
			domainvo.ErrPasswordTooShort.Error(),
		)

	case errors.Is(err, domainvo.ErrPasswordTooLong):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_PASSWORD.String(),
			domainvo.ErrPasswordTooLong.Error(),
		)

	case errors.Is(err, domainvo.ErrPasswordMissingUppercase):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_PASSWORD.String(),
			domainvo.ErrPasswordMissingUppercase.Error(),
		)

	case errors.Is(err, domainvo.ErrPasswordMissingLowercase):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_PASSWORD.String(),
			domainvo.ErrPasswordMissingLowercase.Error(),
		)

	case errors.Is(err, domainvo.ErrPasswordMissingNumber):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_PASSWORD.String(),
			domainvo.ErrPasswordMissingNumber.Error(),
		)

	case errors.Is(err, domainvo.ErrPasswordMissingSymbol):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_PASSWORD.String(),
			domainvo.ErrPasswordMissingSymbol.Error(),
		)

	case errors.Is(err, domainvo.ErrPasswordContainsWhitespace):
		return kerrors.BadRequest(
			apiv1.ErrorReason_INVALID_PASSWORD.String(),
			domainvo.ErrPasswordContainsWhitespace.Error(),
		)

	case errors.Is(err, domainentity.ErrUserDeleted):
		return kerrors.NotFound(
			apiv1.ErrorReason_USER_NOT_FOUND.String(),
			applicationcommon.ErrUserNotFound.Error(),
		)

	case errors.Is(err, applicationcommon.ErrUserNotFound):
		return kerrors.NotFound(
			apiv1.ErrorReason_USER_NOT_FOUND.String(),
			applicationcommon.ErrUserNotFound.Error(),
		)

	case errors.Is(err, domainentity.ErrSessionExpired):
		return kerrors.Unauthorized(
			apiv1.ErrorReason_INVALID_REFRESH_TOKEN.String(),
			applicationcommand.ErrInvalidRefreshToken.Error(),
		)

	case errors.Is(err, domainentity.ErrSessionRevoked):
		return kerrors.Unauthorized(
			apiv1.ErrorReason_INVALID_REFRESH_TOKEN.String(),
			applicationcommand.ErrInvalidRefreshToken.Error(),
		)

	case errors.Is(err, applicationcommand.ErrInvalidRefreshToken):
		return kerrors.Unauthorized(
			apiv1.ErrorReason_INVALID_REFRESH_TOKEN.String(),
			applicationcommand.ErrInvalidRefreshToken.Error(),
		)

	case errors.Is(err, applicationcommand.ErrInvalidCredentials):
		return kerrors.Unauthorized(
			apiv1.ErrorReason_INVALID_CREDENTIALS.String(),
			applicationcommand.ErrInvalidCredentials.Error(),
		)

	case errors.Is(err, domainservice.ErrUsernameAlreadyExists):
		return kerrors.Conflict(
			apiv1.ErrorReason_USERNAME_ALREADY_EXISTS.String(),
			domainservice.ErrUsernameAlreadyExists.Error(),
		)

	default:
		return err // let this be i want to debug
		// return kerrors.InternalServer(
		// 	"",
		// 	"internal server error",
		// )
	}
}
