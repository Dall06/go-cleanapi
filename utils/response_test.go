package utils_test

import (
	"dall06/go-cleanapi/utils"
	"fmt"
	"testing"
)

func TestResponse(test *testing.T) {
	response := utils.NewResponsesUtils()
	if response == nil {
		test.Fatalf("failed to create response repo")
	}

	test.Run("it should create response", func(t *testing.T) {
		r, err := response.Ok("ok", "happy hello")
		if err != nil {
			t.Errorf("failed to create ok response: %s", err)
		}
		fmt.Println(r)
	})

	test.Run("it should create fail response", func(t *testing.T) {
		r, err := response.Fail("ok", "sad hello")
		if err != nil {
			t.Errorf("failed to generate fail response: %s", err)
		}

		fmt.Println(r)
	})

	test.Run("it should fail response-ok 1", func(t *testing.T) {
		r, err := response.Ok("", "")
		if err == nil {
			t.Errorf("failed to not generate ok response: %s", err)
		}
		fmt.Println(r)
	})

	test.Run("it should fail response-ok 2", func(t *testing.T) {
		r, err := response.Ok("ok", nil)
		if err == nil {
			t.Errorf("failed to not generate ok response: %s", err)
		}
		fmt.Println(r)
	})

	test.Run("it should fail response-fail 1", func(t *testing.T) {
		r, err := response.Fail("", "sad hello")
		if err == nil {
			t.Errorf("failed to not generate fail response: %s", err)
		}

		fmt.Println(r)
	})

	test.Run("it should fail response-fail 2", func(t *testing.T) {
		r, err := response.Fail("fail", "")
		if err == nil {
			t.Errorf("failed to not generate fail response: %s", err)
		}

		fmt.Println(r)
	})
}