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

// Vars are config variables
type Vars struct {
	// APIBasePath makes reference to api basepath
	APIBasePath string
	// APIPort makes reference to api port
	APIPort string
	// APIKey makes reference to api key string
	APIKey string
	// APIKeyHash makes reference to api key string in hash
	APIKeyHash string
	// DBConnString is the connection string
	DBConnString string
	// JWTSecret is the secret to generate the jwts
	JWTSecret []byte
	// ProyectName means the proyect name
	ProyectName string
	// Stage is the stage in which the app runs
	Stage string
	// ProyectPath means the absolute path of th proyect
	ProyectPath string
	// CookieSecret is the secret to encode the cookies
	CookieSecret string
	// APIVersion indicates the version of the api
	APIVersion string
	// AppName contains the name of the server including version
	AppName string
}

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
	SetConfig() (*Vars, error)
}

type config struct {
	port    string
	version string
	Vars    Vars
}

var _ Config = (*config)(nil)

// NewConfig is a constructor for config
func NewConfig(p string, v string) Config {
	return &config{
		port:    p,
		version: v,
	}
}

func (c *config) SetConfig() (*Vars, error) {
	if c.port == "" {
		return nil, fmt.Errorf("empty port parameter")
	}
	if c.version == "" {
		return nil, fmt.Errorf("empty version parameter")
	}

	if err := c.loadEnv(); err != nil {
		return nil, err
	}

	c.Vars.DBConnString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", //"<user>:<password>@tcp(127.0.0.1:3306)/<dbname>"
		os.Getenv(envUserDB), os.Getenv(envPasswordDB), os.Getenv(envHostDB), os.Getenv(envPortDB), os.Getenv(envNameDB))
	c.Vars.JWTSecret = []byte(fmt.Sprint(os.Getenv("SECRET_JWT")))
	c.Vars.Stage = strings.ToLower(os.Getenv("STAGE"))

	proyectName, err := c.loadName()
	if err != nil {
		return nil, err
	}
	c.Vars.ProyectName = proyectName
	c.Vars.APIPort = c.port
	c.Vars.APIVersion = c.version

	c.Vars.CookieSecret = os.Getenv(envCookieEncryption)
	c.Vars.APIBasePath = fmt.Sprintf("/%s/api/v%s", c.Vars.ProyectName, c.version)
	c.Vars.AppName = fmt.Sprintf("%s v%s", proyectName, c.Vars.APIVersion)

	ak := os.Getenv(envAPIKey)
	c.Vars.APIKey = ak
	sha := sha512.Sum512_256([]byte(ak))
	c.Vars.APIKeyHash = hex.EncodeToString(sha[:])

	return &c.Vars, nil
}

func (c *config) loadEnv() error {
	projectPath, err := c.getProjectPath()
	if err != nil {
		return err
	}
	c.Vars.ProyectPath = projectPath

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
