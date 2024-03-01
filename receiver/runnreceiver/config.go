package runnreceiver

import (
	"github.com/Arthur1/opentelemetry-collector-arthur1/receiver/runnreceiver/internal/metadata"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

type config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
	metadata.MetricsBuilderConfig           `mapstructure:",squash"`
	Runbooks                                []*runbookConfig `mapstructure:"runbooks"`
}

type runbookConfig struct {
	Body string `mapstructure:"body"`
}
