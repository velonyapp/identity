package vo

type GoogleAuthStrategy struct {
	Sub GoogleSub
}

func NewGoogleAuthStrategy(
	sub GoogleSub,
) GoogleAuthStrategy {
	return GoogleAuthStrategy{
		Sub: sub,
	}
}
