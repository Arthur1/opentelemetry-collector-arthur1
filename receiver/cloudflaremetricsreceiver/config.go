package cloudflaremetricsreceiver

import (
	"time"

	"github.com/Arthur1/opentelemetry-collector-arthur1/receiver/cloudflaremetricsreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

type config struct {
	ApiToken    string   `mapstructure:"api_token"`
	ZoneDomains []string `mapstructure:"zone_domains"`

	scraperhelper.ControllerConfig `mapstructure:",squash"`
	metadata.MetricsBuilderConfig  `mapstructure:",squash"`
}

func defaultConfig() component.Config {
	cc := scraperhelper.NewDefaultControllerConfig()
	cc.CollectionInterval = time.Hour
	mbc := metadata.DefaultMetricsBuilderConfig()
	return &config{
		ControllerConfig:     cc,
		MetricsBuilderConfig: mbc,
	}
}
