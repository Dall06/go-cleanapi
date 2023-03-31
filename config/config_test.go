package config_test

import (
	"dall06/go-cleanapi/config"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(test *testing.T) {
	test.Run("it should load successfully conf test", func(t *testing.T) {
		conf := config.NewConfig()
		err := conf.SetConfig()

		if err != nil {
			t.Fatalf("LoadStrings failed: %v", err)
		}
		
		assert.NotEqual(t, "", config.DBConnString)
		assert.NotEqual(t, []byte(""), config.JwtSecret)
		assert.NotEqual(t, "", config.ProyectName)
		assert.NotEqual(t, "", config.Stage)

		fmt.Println("db " + config.DBConnString)
		fmt.Println("jwt " + string(config.JwtSecret))
		fmt.Println("proyect " + config.ProyectName)
		fmt.Println("stage " + config.Stage)
	})
}