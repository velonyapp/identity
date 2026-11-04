package entity

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

func TestUserChangeCooldowns(t *testing.T) {
	now := time.Date(2026, time.November, 4, 0, 0, 0, 0, time.UTC)
	cooldown := 2 * time.Minute
	username, _ := vo.NewUsername("original")
	token, _ := vo.NewEmailChangeRequestToken("123456")

	tests := []struct {
		name     string
		change   func(*User, int, time.Time) error
		nextTime func(*User) time.Time
		wantErr  error
	}{
		{
			name: "username",
			change: func(user *User, attempt int, now time.Time) error {
				username, _ := vo.NewUsername([]string{"first", "second"}[attempt])
				return user.ChangeUsername(username, cooldown, now)
			},
			nextTime: (*User).NextUsernameChangeTime,
			wantErr:  ErrUsernameChangeRateLimitExceeded,
		},
		{
			name: "email change request",
			change: func(user *User, attempt int, now time.Time) error {
				email, _ := vo.NewEmail([]string{"first@example.com", "second@example.com"}[attempt])
				return user.RequestEmailChange(email, token, 15*time.Minute, cooldown, now)
			},
			nextTime: (*User).NextEmailChangeRequestTime,
			wantErr:  ErrEmailChangeRequestRateLimitExceeded,
		},
		{
			name: "avatar change request",
			change: func(user *User, attempt int, now time.Time) error {
				return user.RequestAvatarChange(vo.GenerateAvatarID(), 30*time.Second, cooldown, now)
			},
			nextTime: (*User).NextAvatarChangeRequestTime,
			wantErr:  ErrAvatarChangeRequestRateLimitExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(username, now)
			if !tt.nextTime(user).Equal(now) {
				t.Fatal("new user must be eligible immediately")
			}
			if err := tt.change(user, 0, now); err != nil {
				t.Fatal(err)
			}
			if !tt.nextTime(user).Equal(now.Add(cooldown)) {
				t.Fatal("next change time must use the supplied cooldown")
			}
			if request := user.EmailChangeRequest(); request != nil && !request.ExpireTime().Equal(now.Add(15*time.Minute)) {
				t.Fatal("email request expiry must use its TTL")
			}
			if request := user.AvatarChangeRequest(); request != nil && !request.ExpireTime().Equal(now.Add(30*time.Second)) {
				t.Fatal("avatar request expiry must use its TTL")
			}

			user = ReconstituteUser(
				user.ID(),
				user.Username(),
				user.Email(),
				user.AvatarID(),
				user.CreateTime(),
				user.DeleteTime(),
				user.EmailChangeRequest(),
				user.AvatarChangeRequest(),
				user.LocalAuthStrategy(),
				user.NextUsernameChangeTime(),
				user.NextEmailChangeRequestTime(),
				user.NextAvatarChangeRequestTime(),
				user.Version(),
			)
			before := *user
			if err := tt.change(user, 1, now.Add(cooldown-time.Nanosecond)); !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
			if !reflect.DeepEqual(before, *user) {
				t.Fatal("rejected change must preserve user state")
			}
			if err := tt.change(user, 1, now.Add(cooldown)); err != nil {
				t.Fatalf("change at the cooldown boundary failed: %v", err)
			}
			if !tt.nextTime(user).Equal(now.Add(2 * cooldown)) {
				t.Fatal("successful change must advance the cooldown")
			}
		})
	}
}

func TestUserChangeCooldownsAreIndependent(t *testing.T) {
	now := time.Date(2026, time.November, 4, 0, 0, 0, 0, time.UTC)
	username, _ := vo.NewUsername("original")
	user := NewUser(username, now)
	newUsername, _ := vo.NewUsername("changed")
	email, _ := vo.NewEmail("user@example.com")
	token, _ := vo.NewEmailChangeRequestToken("123456")

	if err := user.ChangeUsername(newUsername, time.Minute, now); err != nil {
		t.Fatal(err)
	}
	if err := user.RequestEmailChange(email, token, 15*time.Minute, 2*time.Minute, now); err != nil {
		t.Fatal(err)
	}
	if err := user.RequestAvatarChange(vo.GenerateAvatarID(), 30*time.Second, 3*time.Minute, now); err != nil {
		t.Fatal(err)
	}
	if !user.NextUsernameChangeTime().Equal(now.Add(time.Minute)) ||
		!user.NextEmailChangeRequestTime().Equal(now.Add(2*time.Minute)) ||
		!user.NextAvatarChangeRequestTime().Equal(now.Add(3*time.Minute)) {
		t.Fatal("each action must maintain its own cooldown")
	}

	before := *user
	if err := user.ChangeUsername(newUsername, time.Hour, now.Add(time.Second)); !errors.Is(err, ErrUsernameChangeRateLimitExceeded) {
		t.Fatalf("expected username rate limit error for the same username, got %v", err)
	}
	if err := user.ChangeUsername(newUsername, time.Hour, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, *user) {
		t.Fatal("unchanged username must not advance the cooldown or record an event")
	}

	if err := user.ConfirmEmailChange(now); err != nil {
		t.Fatal(err)
	}
	if err := user.RequestEmailChange(email, token, 15*time.Minute, time.Minute, now.Add(time.Second)); !errors.Is(err, ErrEmailChangeRequestRateLimitExceeded) {
		t.Fatalf("expected email rate limit error for the same email, got %v", err)
	}
	if err := user.RequestEmailChange(email, token, 15*time.Minute, time.Minute, now.Add(2*time.Minute)); !errors.Is(err, ErrEmailAlreadySet) {
		t.Fatalf("expected email already set after cooldown, got %v", err)
	}

	if err := user.ConfirmAvatarChange(now); err != nil {
		t.Fatal(err)
	}
	if err := user.RequestAvatarChange(*user.AvatarID(), 30*time.Second, time.Minute, now.Add(time.Second)); !errors.Is(err, ErrAvatarChangeRequestRateLimitExceeded) {
		t.Fatalf("expected avatar rate limit error for the same avatar, got %v", err)
	}
	if err := user.RequestAvatarChange(*user.AvatarID(), 30*time.Second, time.Minute, now.Add(3*time.Minute)); !errors.Is(err, ErrAvatarAlreadySet) {
		t.Fatalf("expected avatar already set after cooldown, got %v", err)
	}
}
