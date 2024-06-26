package runnreceiver // import "github.com/Arthur1/opentelemetry-collector-arthur1/receiver/runnreceiver"

import (
	"context"
	"os"
	"time"

	"github.com/Arthur1/opentelemetry-collector-arthur1/receiver/runnreceiver/internal/metadata"
	"github.com/k1LoW/runn"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
)

type runnScraper struct {
	cfg    *config
	mb     *metadata.MetricsBuilder
	opDefs []*runnOperationDef
}

type runnOperationDef struct {
	runnbookFileName   string
	runnbookFileIsTemp bool
}

var _ receiver.Metrics = new(runnScraper)

func newScraper(cfg *config, settings receiver.CreateSettings) *runnScraper {
	return &runnScraper{
		cfg: cfg,
		mb:  metadata.NewMetricsBuilder(cfg.MetricsBuilderConfig, settings),
	}
}

func (s *runnScraper) Start(ctx context.Context, _ component.Host) error {
	opDefs := make([]*runnOperationDef, 0, len(s.cfg.Runbooks))
	for _, runbook := range s.cfg.Runbooks {
		f, err := os.CreateTemp("", "runnreceiver-runbook")
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.WriteString(runbook); err != nil {
			return err
		}
		opDefs = append(opDefs, &runnOperationDef{
			runnbookFileName:   f.Name(),
			runnbookFileIsTemp: true,
		})
	}
	s.opDefs = opDefs

	return nil
}

func (s *runnScraper) Scrape(ctx context.Context) (pmetric.Metrics, error) {
	for _, opDef := range s.opDefs {
		if err := s.scrapeForRunbook(ctx, opDef); err != nil {
			return pmetric.NewMetrics(), err
		}
	}
	return s.mb.Emit(), nil
}

func (s *runnScraper) scrapeForRunbook(ctx context.Context, opDef *runnOperationDef) error {
	opts := []runn.Option{
		runn.Book(opDef.runnbookFileName),
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

	ts := pcommon.NewTimestampFromTime(time.Now())

	s.mb.RecordRunnStatusDataPoint(ts, btoi(status), result.Desc)
	for _, step := range result.StepResults {
		s.mb.RecordRunnElapsedTimeDataPoint(ts, step.Elapsed.Seconds(), result.Desc, step.Key)
	}

	return nil
}

func btoi(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func (s *runnScraper) Shutdown(ctx context.Context) error {
	for _, opDef := range s.opDefs {
		if opDef.runnbookFileIsTemp {
			os.Remove(opDef.runnbookFileName)
		}
	}
	return nil
}
