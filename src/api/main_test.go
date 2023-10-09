package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// os.exit() will return the status code of the tests
	// to the OS (e.x. if the tests fail, exit code 1 will be returned)
	os.Exit(m.Run())
}
