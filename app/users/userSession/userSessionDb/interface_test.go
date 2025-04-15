package userSessionDb

import (
	"testing"
	"time"

	"imageresizerservice/app/users/userSession"
	"imageresizerservice/library/keyValueDb"
	"imageresizerservice/library/sqlite"
	"imageresizerservice/library/uow"
)

type Fixture struct {
	UowFactory uow.UowFactory
	SessionDb  UserSessionDb
}

func newFixture() *Fixture {
	db := sqlite.New()

	return &Fixture{
		SessionDb:  ImplKeyValueDb{Db: &keyValueDb.ImplHashMap{}},
		UowFactory: uow.UowFactory{Db: db},
	}
}

func Test_GetById(t *testing.T) {
	f := newFixture()
	uow, _ := f.UowFactory.Begin()

	// Create a session
	session := userSession.UserSession{
		Id:        "test-id",
		UserId:    "user-123",
		SessionId: "session-456",
		CreatedAt: time.Now(),
	}

	// Insert the session
	err := f.SessionDb.Upsert(uow, session)
	if err != nil {
		t.Errorf("Expected no error on insert, got %v", err)
	}

	// Get the session
	retrieved, err := f.SessionDb.GetById("test-id")
	if err != nil {
		t.Errorf("Expected no error on retrieval, got %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected to retrieve session, got nil")
	}

	if retrieved.Id != session.Id {
		t.Errorf("Expected Id to be %s, got %s", session.Id, retrieved.Id)
	}

	if retrieved.UserId != session.UserId {
		t.Errorf("Expected UserId to be %s, got %s", session.UserId, retrieved.UserId)
	}

	uow.Commit()
}

func Test_GetByIdNonExistent(t *testing.T) {
	f := newFixture()

	// Try to get a session that doesn't exist
	retrieved, err := f.SessionDb.GetById("nonexistent")

	if err != nil {
		t.Errorf("Expected no error for nonexistent session, got %v", err)
	}

	if retrieved != nil {
		t.Errorf("Expected nil for nonexistent session, got %+v", retrieved)
	}
}

func Test_UpsertNewSession(t *testing.T) {
	f := newFixture()
	uow, _ := f.UowFactory.Begin()

	// Create a session
	session := userSession.UserSession{
		Id:        "new-session",
		UserId:    "user-123",
		SessionId: "session-456",
		CreatedAt: time.Now(),
	}

	// Insert the session
	err := f.SessionDb.Upsert(uow, session)
	if err != nil {
		t.Errorf("Expected no error on insert, got %v", err)
	}

	// Verify it exists
	retrieved, err := f.SessionDb.GetById("new-session")
	if err != nil {
		t.Errorf("Expected no error on retrieval, got %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected to retrieve session, got nil")
	}

	if retrieved.Id != session.Id {
		t.Errorf("Expected Id to be %s, got %s", session.Id, retrieved.Id)
	}

	uow.Commit()
}

func Test_UpsertUpdateSession(t *testing.T) {
	f := newFixture()
	uow, _ := f.UowFactory.Begin()

	// Create initial session
	session := userSession.UserSession{
		Id:        "update-session",
		UserId:    "user-123",
		SessionId: "session-456",
		CreatedAt: time.Now(),
	}

	// Insert the session
	err := f.SessionDb.Upsert(uow, session)
	if err != nil {
		t.Errorf("Expected no error on insert, got %v", err)
	}

	// Update the session
	updatedSession := userSession.UserSession{
		Id:        "update-session",
		UserId:    "user-123",
		SessionId: "session-456",
		CreatedAt: session.CreatedAt,
		EndedAt:   time.Now(),
	}

	// Update the session
	err = f.SessionDb.Upsert(uow, updatedSession)
	if err != nil {
		t.Errorf("Expected no error on update, got %v", err)
	}

	// Verify it was updated
	retrieved, err := f.SessionDb.GetById("update-session")
	if err != nil {
		t.Errorf("Expected no error on retrieval, got %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected to retrieve session, got nil")
	}

	if retrieved.EndedAt.IsZero() {
		t.Error("Expected EndedAt to be set, but it's zero")
	}

	uow.Commit()
}
