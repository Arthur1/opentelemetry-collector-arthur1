package runnreceiver

import (
	"context"
	"errors"
	"time"

	"github.com/Arthur1/opentelemetry-collector-arthur1/receiver/runnreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

var errConfigNotRunnReceiver = errors.New("config was not a runn receiver config")

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, metadata.MetricsStability),
	)
}

func createDefaultConfig() component.Config {
	cc := scraperhelper.NewDefaultControllerConfig()
	cc.CollectionInterval = 60 * time.Second
	mbc := metadata.DefaultMetricsBuilderConfig()
	return &config{
		ControllerConfig:     cc,
		MetricsBuilderConfig: mbc,
		Runbooks:             []string{},
	}
}

func createMetricsReceiver(
	_ context.Context,
	settings receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	c, ok := cfg.(*config)
	if !ok {
		return nil, errConfigNotRunnReceiver
	}
	runnScraper := newScraper(c, settings)
	scraper, err := scraperhelper.NewScraper(
		metadata.Type.String(), runnScraper.Scrape, scraperhelper.WithStart(runnScraper.Start),
	)
	if err != nil {
		return nil, err
	}
	return scraperhelper.NewScraperControllerReceiver(
		&c.ControllerConfig, settings, consumer, scraperhelper.AddScraper(scraper),
	)
}
