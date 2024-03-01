package runnreceiver // import "github.com/Arthur1/opentelemetry-collector-arthur1/receiver/runnreceiver"

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/Arthur1/opentelemetry-collector-arthur1/receiver/runnreceiver/internal/metadata"
	"github.com/k1LoW/runn"
	"github.com/k1LoW/stopw"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
)

type runnScraper struct {
	cfg *config
	mb  *metadata.MetricsBuilder
}

var _ receiver.Metrics = new(runnScraper)

func newScraper(cfg *config, settings receiver.CreateSettings) *runnScraper {
	return &runnScraper{
		cfg: cfg,
		mb:  metadata.NewMetricsBuilder(cfg.MetricsBuilderConfig, settings),
	}
}

func (s *runnScraper) Start(ctx context.Context, _ component.Host) error {
	return nil
}

func (s *runnScraper) Scrape(ctx context.Context) (pmetric.Metrics, error) {
	for _, runbook := range s.cfg.Runbooks {
		if err := s.scrapeForRunbook(ctx, runbook); err != nil {
			return pmetric.NewMetrics(), err
		}
	}
	return s.mb.Emit(), nil
}

func (s *runnScraper) scrapeForRunbook(ctx context.Context, runbook *runbookConfig) error {
	f, err := os.CreateTemp("", "runnreceiver-runbook")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	f.WriteString(runbook.Body)

	opts := []runn.Option{
		runn.Book(f.Name()),
		runn.Profile(true),
		runn.Scopes("read:parent"),
	}
	o, err := runn.New(opts...)
	if err != nil {
		return err
	}

	status := true
	if err := o.Run(ctx); err != nil {
		status = false
	}
	result := o.Result()

	var (
		buf     bytes.Buffer
		profile stopw.Span
	)
	if err := o.DumpProfile(&buf); err != nil {
		return err
	}
	if err := json.NewDecoder(&buf).Decode(&profile); err != nil {
		return err
	}
	ts := pcommon.NewTimestampFromTime(profile.StartedAt)

	s.mb.RecordRunnStatusDataPoint(ts, btoi(status), result.Desc)

	return nil
}

func btoi(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func (s *runnScraper) Shutdown(ctx context.Context) error {
	return nil
}
