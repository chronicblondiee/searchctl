package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	APIVersion     string    `yaml:"apiVersion"`
	Kind           string    `yaml:"kind"`
	CurrentContext string    `yaml:"current-context"`
	Contexts       []Context `yaml:"contexts"`
	Clusters       []Cluster `yaml:"clusters"`
	Users          []User    `yaml:"users"`
}

type Context struct {
	Name    string        `yaml:"name"`
	Context ContextConfig `yaml:"context"`
}

type ContextConfig struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type Cluster struct {
	Name    string        `yaml:"name"`
	Cluster ClusterConfig `yaml:"cluster"`
}

type ClusterConfig struct {
	Server                string `yaml:"server"`
	CertificateAuthority  string `yaml:"certificate-authority,omitempty"`
	InsecureSkipTLSVerify bool   `yaml:"insecure-skip-tls-verify,omitempty"`
}

type User struct {
	Name string     `yaml:"name"`
	User UserConfig `yaml:"user"`
}

type UserConfig struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	APIKey   string `yaml:"api-key,omitempty"`
}

var config *Config

func InitConfig(cfgFile string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		searchctlDir := filepath.Join(home, ".searchctl")
		if err := os.MkdirAll(searchctlDir, 0755); err != nil {
			return err
		}

		viper.AddConfigPath(searchctlDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create default
			return createDefaultConfig()
		}
		return err
	}

	// Use direct viper access for hyphenated keys since viper may not unmarshal them correctly
	config = &Config{
		APIVersion:     viper.GetString("apiVersion"),
		Kind:           viper.GetString("kind"),
		CurrentContext: viper.GetString("current-context"),
	}

	// Unmarshal complex structures
	if err := viper.UnmarshalKey("contexts", &config.Contexts); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("clusters", &config.Clusters); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("users", &config.Users); err != nil {
		return err
	}

	return nil
}

func createDefaultConfig() error {
	defaultConfig := &Config{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentContext: "default",
		Contexts: []Context{
			{
				Name: "default",
				Context: ContextConfig{
					Cluster: "default",
					User:    "default",
				},
			},
		},
		Clusters: []Cluster{
			{
				Name: "default",
				Cluster: ClusterConfig{
					Server:                "http://localhost:9200",
					InsecureSkipTLSVerify: true,
				},
			},
		},
		Users: []User{
			{
				Name: "default",
				User: UserConfig{},
			},
		},
	}

	config = defaultConfig
	return nil
}

func GetConfig() *Config {
	return config
}

func GetCurrentContext() (*Context, error) {
	if config == nil {
		return nil, fmt.Errorf("config not initialized")
	}

	// Check if context was explicitly set via flag, otherwise use config's current context
	currentContextName := config.CurrentContext
	if viper.IsSet("context") && viper.GetString("context") != "" {
		currentContextName = viper.GetString("context")
	}

	if currentContextName == "" {
		return nil, fmt.Errorf("no current context available")
	}

	for _, ctx := range config.Contexts {
		if ctx.Name == currentContextName {
			return &ctx, nil
		}
	}

	return nil, fmt.Errorf("context %q not found", currentContextName)
}

func GetCluster(name string) (*Cluster, error) {
	if config == nil {
		return nil, fmt.Errorf("config not initialized")
	}

	for _, cluster := range config.Clusters {
		if cluster.Name == name {
			return &cluster, nil
		}
	}

	return nil, fmt.Errorf("cluster %q not found", name)
}

func GetUser(name string) (*User, error) {
	if config == nil {
		return nil, fmt.Errorf("config not initialized")
	}

	for _, user := range config.Users {
		if user.Name == name {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user %q not found", name)
}
