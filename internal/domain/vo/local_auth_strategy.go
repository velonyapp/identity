package vo

type LocalAuthStrategy struct {
	passwordHash PasswordHash
}

func NewLocalAuthStrategy(
	passwordHash PasswordHash,
) LocalAuthStrategy {
	return LocalAuthStrategy{
		passwordHash: passwordHash,
	}
}

func (s LocalAuthStrategy) PasswordHash() PasswordHash {
	return s.passwordHash
}

func (s LocalAuthStrategy) Equal(other LocalAuthStrategy) bool {
	return s.passwordHash.Equal(other.passwordHash)
}
