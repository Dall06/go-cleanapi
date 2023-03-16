package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type EnvRepository interface {
	LoadStrings() error
}

type Env struct {
	stage string
}

var _ EnvRepository = (*Env)(nil)

func NewEnv(s string) *Env {
	return &Env{
		stage: s,
	}
}


func (e *Env) LoadStrings() error {
	projectName := regexp.MustCompile(`^(.*` + ProjectDirName + `)`)
    currentWorkDirectory, _ := os.Getwd()
    rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env.` + e.stage)

	if err != nil {
		return err
	}

	DBDNS = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", //"<user>:<password>@tcp(127.0.0.1:3306)/<dbname>"
		os.Getenv("USER_DB"),
		os.Getenv("PASSWORD_DB"),
		os.Getenv("HOST_DB"),
		os.Getenv("PORT_DB"),
		os.Getenv("NAME_DB"))

	JWTSecret = []byte(fmt.Sprint(os.Getenv("SECRET_JWT")))

	return nil
}
