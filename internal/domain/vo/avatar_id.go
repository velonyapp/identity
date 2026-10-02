package vo

type AvatarID struct {
	value string
}

func NewAvatarID(value string) AvatarID {
	return AvatarID{value: value}
}

func (id AvatarID) Value() string {
	return id.value
}

func (id AvatarID) String() string {
	return id.value
}
