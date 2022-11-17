package store

import (
	"buy_list/bot/models"
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func TestAddUsertg(t *testing.T) {
	type args struct {
		ctx context.Context
		u   *models.Usertg
	}
	tests := []struct {
		name    string
		s       *Store
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddUsertg(tt.args.ctx, tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("Store.AddUsertg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
