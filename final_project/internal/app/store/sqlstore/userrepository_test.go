package sqlstore_test

import (
	"context"
	"final_project/internal/app/model"
	"final_project/internal/app/store/sqlstore"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	ctx := context.Background()
	dburl := os.Getenv("DBCONN")
	db, td := sqlstore.TestDB(t, ctx, dburl) //databaseURL
	defer td("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(ctx, model.TestUser(t)))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	ctx := context.Background()
	dburl := os.Getenv("DBCONN")
	db, td := sqlstore.TestDB(t, ctx, dburl) //databaseURL
	defer td("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(ctx, u)
	u2, err := s.User().FindByEmail(ctx, u.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
func TestUserRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	dburl := os.Getenv("DBCONN")
	db, td := sqlstore.TestDB(t, ctx, dburl) //databaseURL
	defer td("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(ctx, u)
	u2, err := s.User().FindByID(ctx, u.Id.String())
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
