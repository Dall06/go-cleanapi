package utils_test

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/utils"
	"testing"
)

func TestCreateUserJWT(test *testing.T) {
	conf := config.NewConfig("8080")
	err := conf.SetConfig()
	if err != nil {
		test.Fatalf("failed to jwt config: %v", err)
	}

	jwt := utils.NewJWT()
	if jwt == nil {
		test.Fatalf("failed to jwt jwt repo")
	}
	test.Run("it should crreate a token", func(t *testing.T) {
		token, err := jwt.CreateUserJWT("myUid")
		if err != nil {
			t.Fatalf("failed to jwt %s", err)
		}

		t.Log("jwt: ", token)
	})

	test.Run("it should fail token creation", func(t *testing.T) {
		_, err := jwt.CreateUserJWT("")
		if err == nil {
			t.Fatalf("failed to fail")
		}

		t.Log("failed successfully: ", err)
	})
}

func TestCheckUserJWT(test *testing.T) {
	jwt := utils.NewJWT()
	if jwt == nil {
		test.Fatalf("failed to jwt jwt repo")
	}
	test.Run("it should succesfully do a token check", func(t *testing.T) {
		token, err := jwt.CreateUserJWT("myUid")
		if err != nil {
			t.Fatalf("failed to jwt: %s", err)
		}

		t.Log("jwt: ", token)

		ok, err := jwt.CheckUserJwt(token)
		if err != nil {
			t.Fatalf("failed to Check: %s", err)
		}

		if !ok {
			t.Fatalf("expected to get true, invalid token")
		}

		t.Log("check ok: ", ok)
	})

	test.Run("it should fail token check, is empty", func(t *testing.T) {
		ok, err := jwt.CheckUserJwt("")
		if err == nil {
			t.Fatalf("failed to fail")
		}

		if ok {
			t.Fatalf("expected to get false")
		}

		t.Log("failed successfully", err)
	})

	test.Run("it should fail token check, is not correct formatted", func(t *testing.T) {
		ok, err := jwt.CheckUserJwt("dawedawedawedawed")
		if err == nil {
			t.Fatalf("failed to fail")
		}

		if ok {
			t.Fatalf("expected to get false")
		}

		t.Log("failed successfully", err)
	})
}

func TestCreateApiJWT(test *testing.T) {
	jwt := utils.NewJWT()
	if jwt == nil {
		test.Fatalf("failed to jwt jwt repo")
	}

	test.Run("it should create a token", func(t *testing.T) {
		t.Parallel()
		conf := config.NewConfig("8080")
		err := conf.SetConfig()
		if err != nil {
			test.Fatalf("failed to jwt config: %v", err)
		}

		token, err := jwt.CreateApiJWT()
		if err != nil {
			t.Fatalf("failed to jwt %s", err)
		}

		t.Log("jwt: ", token)
	})

	test.Run("it should fail token creation", func(t *testing.T) {
		t.Parallel()
		_, err := jwt.CreateApiJWT()
		if err == nil {
			t.Fatalf("failed to fail")
		}

		t.Log("failed successfully: ", err)
	})
}

func TestCheckApiJWT(test *testing.T) {
	jwt := utils.NewJWT()
	if jwt == nil {
		test.Fatalf("failed to CreateUserJWT jwt repo")
	}

	test.Run("it should succesfully do a token check", func(t *testing.T) {
		t.Parallel()
		conf := config.NewConfig("8080")
		err := conf.SetConfig()
		if err != nil {
			test.Fatalf("failed to jwt config: %v", err)
		}

		token, err := jwt.CreateApiJWT()
		if err != nil {
			t.Fatalf("failed to create api jwt: %s", err)
		}

		t.Log("jwt: ", token)

		ok, err := jwt.CheckApiJwt(token)
		if err != nil {
			t.Fatalf("failed to Check: %s", err)
		}

		if !ok {
			t.Fatalf("expected to get true, invalid token")
		}

		t.Log("check ok: ", ok)
	})

	test.Run("it should fail token check, is empty", func(t *testing.T) {
		t.Parallel()
		conf := config.NewConfig("8080")
		err := conf.SetConfig()
		if err != nil {
			test.Fatalf("failed to jwt config: %v", err)
		}

		ok, err := jwt.CheckUserJwt("")
		if err == nil {
			t.Fatalf("failed to fail")
		}

		if ok {
			t.Fatalf("expected to get false")
		}

		t.Log("failed successfully", err)
	})

	test.Run("it should fail token check, is config", func(t *testing.T) {
		t.Parallel()
		ok, err := jwt.CheckUserJwt("jfrjfrnrf")
		if err == nil {
			t.Fatalf("failed to fail")
		}

		if ok {
			t.Fatalf("expected to get false")
		}

		t.Log("failed successfully", err)
	})

	test.Run("it should fail token check, is not correct formatted", func(t *testing.T) {
		t.Parallel()
		conf := config.NewConfig("8080")
		err := conf.SetConfig()
		if err != nil {
			test.Fatalf("failed to jwt config: %v", err)
		}
		ok, err := jwt.CheckUserJwt("dawedawedawedawed")
		if err == nil {
			t.Fatalf("failed to fail")
		}

		if ok {
			t.Fatalf("expected to get false")
		}

		t.Log("failed successfully", err)
	})
}
