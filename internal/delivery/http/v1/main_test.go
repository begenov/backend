package v1

import (
	"os"
	"testing"

	"github.com/begenov/backend/pkg/hash"
	"github.com/gin-gonic/gin"
)

var h hash.PasswordHasher

func TestMain(m *testing.M) {
	h = hash.NewHash()
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
