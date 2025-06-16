package testapi

import (
	"os"
	"testing"
	"time"

	"github.com/dongnguyen248/simple_bank/api"
	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	"github.com/dongnguyen248/simple_bank/util"
	"github.com/gin-gonic/gin"
)

func NewTestServer(t *testing.T, store db.Store) *api.Server {
	config := util.Config{
		TokenSymmetricKey: util.RandomString(32),
		TokenDuration:     time.Minute,
	}
	server, err := api.NewServer(store, config)
	if err != nil {
		t.Fatalf("cannot create server: %v", err)
	}
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}
