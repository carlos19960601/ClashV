package hub

import (
	"github.com/carlos19960601/ClashV/config"
	"github.com/carlos19960601/ClashV/hub/executor"
	"github.com/carlos19960601/ClashV/hub/route"
	"github.com/carlos19960601/ClashV/log"
)

type Option func(*config.Config)

func WithExternalUI(externalUI string) Option {
	return func(cfg *config.Config) {
		cfg.General.ExternalUI = externalUI
	}
}

func WithExternalController(externalController string) Option {
	return func(cfg *config.Config) {
		cfg.General.ExternalController = externalController
	}
}

func Parse(options ...Option) error {
	cfg, err := executor.Parse()
	if err != nil {
		return err
	}

	for _, option := range options {
		option(cfg)
	}

	if cfg.General.ExternalController != "" {
		go route.Start(cfg.General.ExternalController, cfg.General.LogLevel == log.DEBUG)
	}

	executor.ApplyConfig(cfg, true)

	return nil
}
