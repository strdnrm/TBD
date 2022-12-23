package sqlstore_test

import (
	"context"
	"final_project/internal/app/model"
	"final_project/internal/app/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	ctx := context.Background()
	db, td := sqlstore.TestDB(t, ctx, databaseURL)
	defer td("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(ctx, model.TestUser(t)))
	assert.NotNil(t, u)
}
