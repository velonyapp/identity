package common

import "time"

const (
	UsernameChangeCooldown      = 7 * 24 * time.Hour
	EmailChangeRequestCooldown  = time.Minute
	AvatarChangeRequestCooldown = 10 * time.Second
	EmailChangeRequestTTL       = 15 * time.Minute
	AvatarChangeRequestTTL      = 30 * time.Second
)
