// Package config contains the method to load config variables that will be used in the app
package config

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var (
	// APIBasePath makes reference to api basepath
	APIBasePath = ""
	// APIPort makes reference to api port
	APIPort = ""
	// APIKey makes reference to api key string
	APIKey = ""
	// APIKeyHash makes reference to api key string in hash
	APIKeyHash = ""
	// DBConnString is the connection string
	DBConnString = ""
	// JWTSecret is the secret to generate the jwts
	JWTSecret = []byte("")
	// ProyectName means the proyect name
	ProyectName = ""
	// Stage is the stage in which the app runs
	Stage = ""
	// ProyectPath means the absolute path of th proyect
	ProyectPath = ""
	// CookieSecret is the secret to encode the cookies
	CookieSecret = ""
)

const (
	file = ".env"

	envUserDB           = "USER_DB"
	envPasswordDB       = "PASSWORD_DB"
	envHostDB           = "HOST_DB"
	envPortDB           = "PORT_DB"
	envNameDB           = "NAME_DB"
	envSecretJWT        = "SECRET_JWT"
	envStage            = "STAGE"
	envCookieEncryption = "COOKIE_ENCRYPTION"
	envAPIKey           = "API_KEY"
)

// Config is an interface that extends config
type Config interface {
	SetConfig() error
}

type config struct {
	port string
}

var _ Config = (*config)(nil)

// NewConfig is a constructor for config
func NewConfig(p string) Config {
	return &config{
		port: p,
	}
}

func (c *config) SetConfig() error {
	if err := c.loadEnv(); err != nil {
		return err
	}

	DBConnString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", //"<user>:<password>@tcp(127.0.0.1:3306)/<dbname>"
		os.Getenv(envUserDB), os.Getenv(envPasswordDB), os.Getenv(envHostDB), os.Getenv(envPortDB), os.Getenv(envNameDB))
	JWTSecret = []byte(fmt.Sprint(os.Getenv("SECRET_JWT")))
	Stage = strings.ToLower(os.Getenv("STAGE"))

	proyectName, err := c.loadName()
	if err != nil {
		return err
	}
	ProyectName = proyectName
	APIPort = c.port

	CookieSecret = os.Getenv(envCookieEncryption)
	APIBasePath = fmt.Sprintf("/%s/api", ProyectName)

	ak := os.Getenv(envAPIKey)
	APIKey = ak
	sha := sha512.Sum512_256([]byte(ak))
	APIKeyHash = hex.EncodeToString(sha[:])

	return nil
}

func (c *config) loadEnv() error {
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

func (c *config) getProjectPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	for dir := cwd; dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		_, err := os.Stat(filepath.Join(dir, "go.mod"))
		if err == nil {
			return dir, nil
		}

		if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to check directory: %w", err)
		}
	}

	return "", fmt.Errorf("failed to find project root directory")
}

func (c *config) loadName() (string, error) {
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
