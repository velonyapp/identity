package middleware

import (
	"context"
	"errors"

	applicationcommand "github.com/velonyapp/identity/internal/application/command"
	applicationcommon "github.com/velonyapp/identity/internal/application/common"
	applicationport "github.com/velonyapp/identity/internal/application/port"
	domainentity "github.com/velonyapp/identity/internal/domain/entity"
	domainrepo "github.com/velonyapp/identity/internal/domain/repo"
	domainvo "github.com/velonyapp/identity/internal/domain/vo"
	presentationservice "github.com/velonyapp/identity/internal/presentation/service"

	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/middleware"
)

type ErrorMapper middleware.Middleware

func NewErrorMapper() ErrorMapper {
	return ErrorMapper(
		func(next middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req any) (any, error) {
				reply, err := next(ctx, req)
				if err == nil {
					return reply, nil
				}

				var kratosErr *kerrors.Error
				if errors.As(err, &kratosErr) && kratosErr.Code < 500 {
					return reply, err
				}

				switch {
				// Presentation Service
				case errors.Is(err, presentationservice.ErrInvalidUserResourceName):
					return reply, kerrors.BadRequest("",
						presentationservice.ErrInvalidUserResourceName.Error(),
					)
				case errors.Is(err, presentationservice.ErrInvalidUserBatchSize):
					return reply, kerrors.BadRequest("",
						presentationservice.ErrInvalidUserBatchSize.Error(),
					)
				case errors.Is(err, presentationservice.ErrUnauthenticated):
					return reply, kerrors.Unauthorized("",
						presentationservice.ErrUnauthenticated.Error(),
					)
				case errors.Is(err, presentationservice.ErrForbidden):
					return reply, kerrors.Forbidden("",
						presentationservice.ErrForbidden.Error(),
					)

				// Application Command
				case errors.Is(err, applicationcommand.ErrInvalidCredentials):
					return reply, kerrors.Unauthorized("",
						applicationcommand.ErrInvalidCredentials.Error(),
					)
				case errors.Is(err, applicationcommand.ErrInvalidRefreshToken):
					return reply, kerrors.Unauthorized("",
						applicationcommand.ErrInvalidRefreshToken.Error(),
					)

				// Application Common
				case errors.Is(err, applicationcommon.ErrUserNotFound):
					return reply, kerrors.NotFound("",
						applicationcommon.ErrUserNotFound.Error(),
					)

				// Application Port
				case errors.Is(err, applicationport.ErrPasswordMismatch):
					return reply, kerrors.Unauthorized("",
						applicationcommand.ErrInvalidCredentials.Error(),
					)
				case errors.Is(err, applicationport.ErrSessionTokenInvalid):
					return reply, kerrors.Unauthorized("",
						applicationcommand.ErrInvalidRefreshToken.Error(),
					)
				case errors.Is(err, applicationport.ErrEmailChangeRequestTokenInvalid):
					return reply, kerrors.BadRequest("",
						applicationport.ErrEmailChangeRequestTokenInvalid.Error(),
					)

				// Domain Entity
				case errors.Is(err, domainentity.ErrUserDeleted):
					return reply, kerrors.NotFound("",
						applicationcommon.ErrUserNotFound.Error(),
					)
				case errors.Is(err, domainentity.ErrSessionExpired):
					return reply, kerrors.Unauthorized("",
						applicationcommand.ErrInvalidRefreshToken.Error(),
					)
				case errors.Is(err, domainentity.ErrSessionRevoked):
					return reply, kerrors.Unauthorized("",
						applicationcommand.ErrInvalidRefreshToken.Error(),
					)
				case errors.Is(err, domainentity.ErrEmailAlreadySet):
					return reply, kerrors.BadRequest("",
						domainentity.ErrEmailAlreadySet.Error(),
					)
				case errors.Is(err, domainentity.ErrEmailChangeNotRequested):
					return reply, kerrors.BadRequest("",
						domainentity.ErrEmailChangeNotRequested.Error(),
					)
				case errors.Is(err, domainentity.ErrEmailChangeRequestExpired):
					return reply, kerrors.BadRequest("",
						domainentity.ErrEmailChangeRequestExpired.Error(),
					)
				case errors.Is(err, domainentity.ErrAvatarAlreadySet):
					return reply, kerrors.BadRequest("",
						domainentity.ErrAvatarAlreadySet.Error(),
					)
				case errors.Is(err, domainentity.ErrAvatarChangeNotRequested):
					return reply, kerrors.BadRequest("",
						domainentity.ErrAvatarChangeNotRequested.Error(),
					)
				case errors.Is(err, domainentity.ErrAvatarChangeRequestExpired):
					return reply, kerrors.BadRequest("",
						domainentity.ErrAvatarChangeRequestExpired.Error(),
					)

				// Domain Repository
				case errors.Is(err, domainrepo.ErrUsernameAlreadyExists):
					return reply, kerrors.Conflict("",
						domainrepo.ErrUsernameAlreadyExists.Error(),
					)
				case errors.Is(err, domainrepo.ErrEmailAlreadyExists):
					return reply, kerrors.Conflict("",
						domainrepo.ErrEmailAlreadyExists.Error(),
					)

				// Domain Value Object
				case errors.Is(err, domainvo.ErrUsernameTooShort):
					return reply, kerrors.BadRequest("",
						domainvo.ErrUsernameTooShort.Error(),
					)
				case errors.Is(err, domainvo.ErrUsernameTooLong):
					return reply, kerrors.BadRequest("",
						domainvo.ErrUsernameTooLong.Error(),
					)
				case errors.Is(err, domainvo.ErrUsernameInvalidCharacter):
					return reply, kerrors.BadRequest("",
						domainvo.ErrUsernameInvalidCharacter.Error(),
					)
				case errors.Is(err, domainvo.ErrInvalidEmail):
					return reply, kerrors.BadRequest("",
						domainvo.ErrInvalidEmail.Error(),
					)
				case errors.Is(err, domainvo.ErrPasswordTooShort):
					return reply, kerrors.BadRequest("",
						domainvo.ErrPasswordTooShort.Error(),
					)
				case errors.Is(err, domainvo.ErrPasswordTooLong):
					return reply, kerrors.BadRequest("",
						domainvo.ErrPasswordTooLong.Error(),
					)
				case errors.Is(err, domainvo.ErrPasswordMissingUppercase):
					return reply, kerrors.BadRequest("",
						domainvo.ErrPasswordMissingUppercase.Error(),
					)
				case errors.Is(err, domainvo.ErrPasswordMissingLowercase):
					return reply, kerrors.BadRequest("",
						domainvo.ErrPasswordMissingLowercase.Error(),
					)
				case errors.Is(err, domainvo.ErrPasswordMissingNumber):
					return reply, kerrors.BadRequest("",
						domainvo.ErrPasswordMissingNumber.Error(),
					)
				case errors.Is(err, domainvo.ErrPasswordMissingSymbol):
					return reply, kerrors.BadRequest("",
						domainvo.ErrPasswordMissingSymbol.Error(),
					)
				case errors.Is(err, domainvo.ErrPasswordContainsWhitespace):
					return reply, kerrors.BadRequest("",
						domainvo.ErrPasswordContainsWhitespace.Error(),
					)
				case errors.Is(err, domainvo.ErrEmailChangeRequestTokenInvalid):
					return reply, kerrors.BadRequest("",
						domainvo.ErrEmailChangeRequestTokenInvalid.Error(),
					)

				default:
					return reply, kerrors.InternalServer("",
						"internal server error",
					).WithCause(err)
				}
			}
		},
	)
}

func (m ErrorMapper) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
