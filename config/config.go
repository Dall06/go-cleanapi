package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var (
	DBConnString = ""
	JwtSecret    = []byte("")
	ProyectName  = ""
	Stage        = ""
	ProyectPath  = ""
)

const (
    file = ".env"

    envUserDB     = "USER_DB"
    envPasswordDB = "PASSWORD_DB"
    envHostDB     = "HOST_DB"
    envPortDB     = "PORT_DB"
    envNameDB     = "NAME_DB"
    envSecretJWT  = "SECRET_JWT"
    envStage      = "STAGE"
)

type ConfigRepository interface {
	SetConfig() error
}

type Config struct{}

var _ ConfigRepository = (*Config)(nil)

func NewConfig() ConfigRepository {
	return &Config{}
}

func (c *Config) SetConfig() error {
	if err := c.loadEnv(); err != nil {
		return err
	}

	DBConnString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", //"<user>:<password>@tcp(127.0.0.1:3306)/<dbname>"
		os.Getenv(envUserDB), os.Getenv(envPasswordDB), os.Getenv(envHostDB), os.Getenv(envPortDB), os.Getenv(envNameDB))
	JwtSecret = []byte(fmt.Sprint(os.Getenv("SECRET_JWT")))
	Stage = strings.ToLower(os.Getenv("STAGE"))

	name, err := c.loadName()
	if err != nil {
		return err
	}
	ProyectName = name

	return nil
}

func (c *Config) loadEnv() error {
	projectPath, err := c.getProjectPath()
	if err != nil {
		return err
	}
	ProyectPath = projectPath

	filePath := filepath.Join(projectPath, file)

	err = godotenv.Load(filePath)
	if err != nil {
		return err
	}

	return nil
}


func (c *Config) getProjectPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	for dir := cwd; dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to check directory: %w", err)
		}
	}

	return "", fmt.Errorf("failed to find project root directory")
}

func (c *Config) loadName() (string, error) {
	cmd := exec.Command("go", "list", "-m")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// The output of the `go list -m` command will contain the module name
	// and version, separated by a space. We only want the module name.
	module := strings.Split(string(out), " ")[0]
	module = strings.Split(string(module), "/")[1]
	module = strings.ReplaceAll(module, "\n", "")
	return module, nil
}