//go:build !go1.11
// +build !go1.11

package tester

import (
	"testing"

	"github.com/forthxu/gin"
)

func testOptionSameSitego(t *testing.T, r *gin.Engine) {
	// not supported
}
