package link

import (
	"imageresizerservice/app/ctx/sessionID"
	"imageresizerservice/app/users/login/link/linkID"
	"imageresizerservice/library/email/emailAddress"
	"time"
)

type Link struct {
	ID           linkID.LinkID
	EmailAddress emailAddress.EmailAddress
	CreatedAt    time.Time
	UsedAt       time.Time
	SessionID    sessionID.SessionID
}

func New(emailAddress emailAddress.EmailAddress, sessionID sessionID.SessionID) Link {
	return Link{
		ID:           linkID.Gen(),
		EmailAddress: emailAddress,
		CreatedAt:    time.Now(),
		SessionID:    sessionID,
	}
}

func MarkAsUsed(l Link) Link {
	return Link{
		ID:           l.ID,
		EmailAddress: l.EmailAddress,
		CreatedAt:    l.CreatedAt,
		UsedAt:       time.Now(),
		SessionID:    l.SessionID,
	}
}

func WasUsed(l *Link) bool {
	return !l.UsedAt.IsZero()
}
