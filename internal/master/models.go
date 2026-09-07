package master

import "time"

type Profile struct {
	UserID      int64
	City        string
	Bio         string
	AvatarURL   *string
	RatingAvg   float64
	RatingCount int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
