package config_test

import (
	"dall06/go-cleanapi/config"
	"testing"
)

func TestLoadString(test *testing.T) {
	test.Run("it should load successfully env dev", func(t *testing.T) {
		env := config.NewEnv("dev")
		err := env.LoadStrings()

		if err != nil {
			t.Fatalf("LoadStrings failed: %v", err)
		}
		
		t.Log("loaded succesfully")
	})

	test.Run("it should load successfully env test", func(t *testing.T) {
		env := config.NewEnv("test")
		err := env.LoadStrings()

		if err != nil {
			t.Fatalf("LoadStrings failed: %v", err)
		}
		
		t.Log("loaded succesfully")
	})

	test.Run("it should load successfully env staging", func(t *testing.T) {
		env := config.NewEnv("staging")
		err := env.LoadStrings()

		if err != nil {
			t.Fatalf("LoadStrings failed: %v", err)
		}
		
		t.Log("loaded succesfully")
	})

	test.Run("it should load successfully env prod", func(t *testing.T) {
		env := config.NewEnv("prod")
		err := env.LoadStrings()

		if err != nil {
			t.Fatalf("LoadStrings failed: %v", err)
		}
		
		t.Log("loaded succesfully")
	})
}