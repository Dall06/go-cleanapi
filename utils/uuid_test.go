package utils_test

import (
	"dall06/go-cleanapi/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUID(test *testing.T){
	uuidGen := utils.NewUUIDGenerator()

	test.Run("it should return a string", func(t *testing.T) {
		s := uuidGen.NewString()

		assert.NotEmpty(t, s)
	})

	test.Run("it should return a UUID", func(t *testing.T) {
		uid := uuidGen.NewUUID()

		assert.NotEmpty(t, uid)
	})
}