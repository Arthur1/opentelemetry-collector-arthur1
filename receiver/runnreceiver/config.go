package runnreceiver

import (
	"github.com/Arthur1/opentelemetry-collector-arthur1/receiver/runnreceiver/internal/metadata"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

type config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`
	metadata.MetricsBuilderConfig  `mapstructure:",squash"`
	Runbooks                       []string `mapstructure:"runbooks"`
}
