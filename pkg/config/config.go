package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TaskDefinition string `envconfig:"task_definition" required:"true"`
	Container      string `envconfig:"container" required:"true"`
	Cluster        string `envconfig:"cluster" required:"true"`
	Service        string `envconfig:"service"`
	Image          string `envconfig:"image" required:"true"`
	WaitSleep      int    `envconfig:"wait_sleep" default:"3"`
	WaitCount      int    `envconfig:"wait_count" default:"30"`
	SkipVerify     bool   `envconfig:"skip_verify" default:"false"`
}

var conf Config

func LoadConf() error {
	// load .env file
	searchPath := []string{".env", "/config/.env"}
	for _, p := range searchPath {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err != nil {
				return err
			}
		}
	}

	if err := envconfig.Process("", &conf); err != nil {
		return err
	}

	return nil
}

func GetConf() *Config {
	return &conf
}

func ParseFlags() {
	var (
		td        = flag.String("task-definition", "", "Task Definition Name")
		container = flag.String("container", "", "Container Name")
		cluster   = flag.String("cluster", "", "Cluster Name")
		service   = flag.String("service", "", "Service Name")
		image     = flag.String("image", "", "Image URL")
		ws        = flag.Int("wait-sleep", 3, "Wait Sleep")
		wc        = flag.Int("wait-count", 30, "Wait Count")
		sv        = flag.Bool("skip-verify", false, "Skip Service Verify")
	)
	flag.Parse()

	if td != nil {
		_ = os.Setenv("TASK_DEFINITION", *td)
	}
	if container != nil {
		_ = os.Setenv("CONTAINER", *container)
	}
	if cluster != nil {
		_ = os.Setenv("CLUSTER", *cluster)
	}
	if service != nil {
		_ = os.Setenv("SERVICE", *service)
	}
	if image != nil {
		_ = os.Setenv("IMAGE", *image)
	}
	if ws != nil {
		_ = os.Setenv("WAIT_SLEEP", strconv.Itoa(*ws))
	}
	if wc != nil {
		_ = os.Setenv("WAIT_COUNT", strconv.Itoa(*wc))
	}
	if sv != nil {
		_ = os.Setenv("SKIP_VERIFY", strconv.FormatBool(*sv))
	}
}
