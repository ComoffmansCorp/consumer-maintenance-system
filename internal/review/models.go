package review

import "time"

type Review struct {
	ID        int64
	RequestID int64
	ClientID  int64
	MasterID  int64
	Rating    int
	Comment   string
	Hidden    bool
	CreatedAt time.Time
}
