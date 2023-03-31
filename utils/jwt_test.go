package utils_test

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/utils"
	"testing"
)

func TestJWTCreate(test *testing.T) {
	conf := config.NewConfig()
	err := conf.SetConfig()
	if err != nil {
		test.Fatalf("failed to create config: %v", err)
	}
	
	jwt := utils.NewJWT()
	if jwt == nil {
		test.Fatalf("failed to create jwt repo")
	}
	test.Run("it should crreate a token", func(t *testing.T) {
		token, err := jwt.Create("myUid")
		if err != nil {
			t.Fatalf("failed to create %s", err)
		}

		t.Log("created: ", token)
	})

	test.Run("it should fail token creation", func(t *testing.T) {
		_, err := jwt.Create("")
		if err == nil {
			t.Fatalf("failed to fail")
		}

		t.Log("failed successfully: ", err)
	})
}

func TestJWTCheckToken(test *testing.T) {
	jwt := utils.NewJWT()
	if jwt == nil {
		test.Fatalf("failed to create jwt repo")
	}
	test.Run("it should succesfully do a token check", func(t *testing.T) {
		token, err := jwt.Create("myUid")
		if err != nil {
			t.Fatalf("failed to create: %s", err)
		}

		t.Log("created: ", token)

		checkToken, err := jwt.Check(token)
		if err != nil {
			t.Fatalf("failed to Check: %s", err)
		}

		t.Log("check ok: ", checkToken)
	})

	test.Run("it should fail token check", func(t *testing.T) {
		_, err := jwt.Check("")
		if err == nil {
			t.Fatalf("failed to fail")
		}

		t.Log("failed successfully", err)
	})
}
