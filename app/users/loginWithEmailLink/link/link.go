package link

import (
	"time"

	"imageresizerservice/library/id"
)

type Link struct {
	Id           string
	EmailAddress string
	CreatedAt    time.Time
	UsedAt       time.Time
}

func New(email string) Link {
	return Link{
		Id:           id.Gen(),
		EmailAddress: email,
		CreatedAt:    time.Now(),
	}
}

func MarkAsUsed(l Link) Link {
	return Link{
		Id:           l.Id,
		EmailAddress: l.EmailAddress,
		CreatedAt:    l.CreatedAt,
		UsedAt:       time.Now(),
	}
}

func WasUsed(l *Link) bool {
	return !l.UsedAt.IsZero()
}
