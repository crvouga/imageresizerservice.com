package appCtx

import (
	"database/sql"
	"imageresizerservice/app/users/loginWithEmailLink/link/linkDb"
	"imageresizerservice/library/email/emailOutbox"
	"imageresizerservice/library/email/sendEmail"
	"imageresizerservice/library/keyValueDb"
	"imageresizerservice/library/uow"
)

type AppCtx struct {
	SendEmail   sendEmail.SendEmail
	LinkDb      linkDb.LinkDb
	UowFactory  uow.UowFactory
	EmailOutbox emailOutbox.EmailOutbox
	KeyValueDb  keyValueDb.KeyValueDb
}

func New(db *sql.DB) AppCtx {

	keyValueDbHashMap := keyValueDb.ImplHashMap{}

	return AppCtx{
		SendEmail:   &sendEmail.ImplFake{},
		LinkDb:      &linkDb.ImplKeyValueDb{Db: &keyValueDbHashMap},
		UowFactory:  uow.UowFactory{Db: db},
		KeyValueDb:  &keyValueDbHashMap,
		EmailOutbox: &emailOutbox.ImplKeyValueDb{Db: &keyValueDbHashMap},
	}

}
