package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/chronicblondiee/searchctl/pkg/config"
)

type Factory struct {
	httpClient *http.Client
	baseURL    string
	username   string
	password   string
	apiKey     string
}

func NewFactory() (*Factory, error) {
	ctx, err := config.GetCurrentContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get current context: %w", err)
	}

	cluster, err := config.GetCluster(ctx.Context.Cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster config: %w", err)
	}

	user, err := config.GetUser(ctx.Context.User)
	if err != nil {
		return nil, fmt.Errorf("failed to get user config: %w", err)
	}

	httpClient := &http.Client{}

	if cluster.Cluster.InsecureSkipTLSVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	return &Factory{
		httpClient: httpClient,
		baseURL:    strings.TrimSuffix(cluster.Cluster.Server, "/"),
		username:   user.User.Username,
		password:   user.User.Password,
		apiKey:     user.User.APIKey,
	}, nil
}

func (f *Factory) HTTPClient() *http.Client {
	return f.httpClient
}

func (f *Factory) BaseURL() string {
	return f.baseURL
}

func (f *Factory) Username() string {
	return f.username
}

func (f *Factory) Password() string {
	return f.password
}

func (f *Factory) APIKey() string {
	return f.apiKey
}