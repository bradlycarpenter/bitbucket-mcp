package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Email     string
	Token     string
	Workspace string
}

func Load() (*Config, error) {
	email := os.Getenv("BITBUCKET_EMAIL")
	token := os.Getenv("BITBUCKET_TOKEN")
	domain := os.Getenv("BITBUCKET_DOMAIN")

	if email == "" || token == "" || domain == "" {
		return nil, fmt.Errorf("BITBUCKET_EMAIL, BITBUCKET_TOKEN, and BITBUCKET_DOMAIN must be set")
	}

	u, err := url.Parse(domain)
	if err != nil {
		return nil, fmt.Errorf("invalid BITBUCKET_DOMAIN: %w", err)
	}
	workspace := strings.Trim(u.Path, "/")
	if workspace == "" {
		return nil, fmt.Errorf("BITBUCKET_DOMAIN must include workspace path, e.g. https://bitbucket.org/myworkspace")
	}

	return &Config{Email: email, Token: token, Workspace: workspace}, nil
}
