package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"syscall"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
)

type Config struct {
	file        *os.File
	ClientID    string   `json:"clientId"`
	Scopes      []string `json:"scopes"`
	RedirectURI string   `json:"redirectURI"`
	Token       *any     `json:"token"`
}

func OpenConfig(name string) (*Config, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		file.Close()
		return nil, err
	}

	cache := &Config{file: file}
	data, err := io.ReadAll(file)
	if err != nil {
		file.Close()
		return nil, err
	}
	if err := json.Unmarshal(data, cache); err != nil {
		file.Close()
		return nil, err
	}

	return cache, nil
}

func (config *Config) Close() error {
	return config.file.Close()
}

func (config *Config) Replace(
	ctx context.Context,
	cache cache.Unmarshaler,
	hints cache.ReplaceHints,
) error {
	if config.Token == nil {
		return cache.Unmarshal([]byte{})
	}
	data, err := json.Marshal(config.Token)
	if err != nil {
		return err
	}
	return cache.Unmarshal(data)
}

func (config *Config) Export(
	ctx context.Context,
	cache cache.Marshaler,
	hints cache.ExportHints,
) error {
	data, err := cache.Marshal()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &config.Token); err != nil {
		return err
	}
	dump, err := json.Marshal(config)
	if err != nil {
		return err
	}
	length, err := config.file.WriteAt(dump, 0)
	if err != nil {
		return err
	}
	return config.file.Truncate(int64(length))
}
