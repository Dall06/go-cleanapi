package config_test

import (
	"dall06/go-cleanapi/config"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(test *testing.T) {
	test.Run("it should load successfully conf test", func(t *testing.T) {
		conf := config.NewConfig("8080")
		err := conf.SetConfig()
		if err != nil {
			t.Fatalf("LoadStrings failed: %v", err)
		}
		
		assert.NotEqual(t, "", config.DBConnString)
		assert.NotEqual(t, []byte(""), config.JwtSecret)
		assert.NotEqual(t, "", config.ProyectName)
		assert.NotEqual(t, "", config.Stage)
		assert.NotEqual(t, "", config.ProyectPath)
		assert.NotEqual(t, "", config.CookieSecret)
		assert.NotEqual(t, "", config.ApiBasePath)

		fmt.Println("db " + config.DBConnString)
		fmt.Println("jwt " + string(config.JwtSecret))
		fmt.Println("proyect " + config.ProyectName)
		fmt.Println("stage " + config.Stage)
	})

	test.Run("it should not load successfully conf test, set config missing", func(t *testing.T) {
		config.NewConfig("8080")

		assert.Empty(t, config.DBConnString)
		assert.Empty(t, config.JwtSecret)
		assert.Empty(t, config.ProyectName)
		assert.Empty(t, config.Stage)
		assert.Empty(t, config.ProyectPath)
		assert.Empty(t, config.CookieSecret)
		assert.Empty(t, config.ApiBasePath)
	})

	test.Run("it should not load successfully conf test, new config missing", func(t *testing.T) {
		assert.Empty(t, config.DBConnString)
		assert.Empty(t, config.JwtSecret)
		assert.Empty(t, config.ProyectName)
		assert.Empty(t, config.Stage)
		assert.Empty(t, config.ProyectPath)
		assert.Empty(t, config.CookieSecret)
		assert.Empty(t, config.ApiBasePath)
	})
}