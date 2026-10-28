package common

import "time"

type UserResult struct {
	ID         string
	Username   string
	Email      *string
	AvatarKey  *string
	CreateTime time.Time
}
