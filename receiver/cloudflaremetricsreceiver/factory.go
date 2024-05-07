package cloudflaremetricsreceiver

import (
	"context"
	"fmt"

	"github.com/Arthur1/opentelemetry-collector-arthur1/receiver/cloudflaremetricsreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		defaultConfig,
		receiver.WithMetrics(
			createMetricsReceiver,
			metadata.MetricsStability,
		),
	)
}

func createMetricsReceiver(
	_ context.Context, settings receiver.CreateSettings,
	cfg component.Config, consumer consumer.Metrics,
) (receiver.Metrics, error) {
	c, ok := cfg.(*config)
	if !ok {
		return nil, fmt.Errorf("invalid config")
	}
	s := newScraper(c, settings)
	scraper, err := scraperhelper.NewScraper(metadata.Type.String(), s.scrape)
	if err != nil {
		return nil, err
	}
	return scraperhelper.NewScraperControllerReceiver(
		&c.ControllerConfig, settings, consumer, scraperhelper.AddScraper(scraper),
	)
}
