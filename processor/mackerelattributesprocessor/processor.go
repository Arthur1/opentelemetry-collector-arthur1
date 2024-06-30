package mackerelattributesprocessor

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/Arthur1/opentelemetry-collector-arthur1/processor/mackerelattributesprocessor/internal/metadata"
	mackerelAgentConfig "github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-client-go"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type mackerelProcessor struct {
	cfg               *Config
	rb                *metadata.ResourceBuilder
	hostID            string
	orgName           string
	mackerelHostName  string
	hostIDStorage     *mackerelAgentConfig.FileSystemHostIDStorage
	configFileModTime time.Time
	hostIDFileModTime time.Time
	cancel            context.CancelFunc
}

func newMackerelProcessor(cfg *Config) *mackerelProcessor {
	rb := metadata.NewResourceBuilder(cfg.ResourceAttributesConfig)
	return &mackerelProcessor{
		cfg: cfg,
		rb:  rb,
	}
}

func (p *mackerelProcessor) Start(ctx context.Context, _ component.Host) error {
	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	if err := p.sync(ctx); err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := p.sync(ctx); err != nil {
					fmt.Printf("error: %v\n", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (p *mackerelProcessor) Shutdown(_ context.Context) error {
	if p.cancel != nil {
		p.cancel()
	}
	return nil
}

func (p *mackerelProcessor) sync(_ context.Context) error {
	configFileInfo, err := os.Stat(p.cfg.ConfigFilePath)
	if err != nil {
		return err
	}
	if currentModTime := configFileInfo.ModTime(); currentModTime.After(p.configFileModTime) {
		cfg, err := mackerelAgentConfig.LoadConfig(p.cfg.ConfigFilePath)
		if err != nil {
			return err
		}
		if cfg.Apibase == mackerelAgentConfig.DefaultConfig.Apibase {
			// Usually the domain of the api and the domain of the web console are different.
			p.mackerelHostName = "https://mackerel.io"
		} else {
			p.mackerelHostName = cfg.Apibase
		}

		p.hostIDStorage = &mackerelAgentConfig.FileSystemHostIDStorage{
			Root: cfg.Root,
		}

		cli, err := mackerel.NewClientWithOptions(cfg.Apikey, cfg.Apibase, false)
		if err != nil {
			return err
		}

		org, err := cli.GetOrg()
		if err != nil {
			return err
		}
		p.orgName = org.Name

		p.configFileModTime = currentModTime
	}

	if p.hostIDStorage == nil {
		return fmt.Errorf("hostIDStorage is not prepared")
	}

	hostIDFileInfo, err := os.Stat(p.hostIDStorage.HostIDFile())
	if err != nil {
		return err
	}
	if currentModTime := hostIDFileInfo.ModTime(); currentModTime.After(p.hostIDFileModTime) {
		hostID, err := p.hostIDStorage.LoadHostID()
		if err != nil {
			return err
		}
		p.hostID = hostID

		p.hostIDFileModTime = currentModTime
	}

	return nil
}

func (p *mackerelProcessor) hostURL() (string, error) {
	return url.JoinPath(p.mackerelHostName, "orgs", p.orgName, "hosts", p.hostID)
}

func (p *mackerelProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		resource := rss.At(i).Resource()
		p.processResource(ctx, resource)
	}
	return td, nil
}

func (p *mackerelProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	rl := ld.ResourceLogs()
	for i := 0; i < rl.Len(); i++ {
		resource := rl.At(i).Resource()
		p.processResource(ctx, resource)
	}
	return ld, nil
}

func (p *mackerelProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	rm := md.ResourceMetrics()
	for i := 0; i < rm.Len(); i++ {
		resource := rm.At(i).Resource()
		p.processResource(ctx, resource)
	}
	return md, nil
}

func (p *mackerelProcessor) processResource(_ context.Context, resource pcommon.Resource) {
	if _, ok := resource.Attributes().Get("mackerelio.org.name"); !ok && p.orgName != "" {
		resource.Attributes().PutStr("mackerelio.org.name", p.orgName)
	}

	if _, ok := resource.Attributes().Get("mackerelio.host.id"); !ok && p.hostID != "" {
		resource.Attributes().PutStr("mackerelio.host.id", p.hostID)
	}

	hostURL, err := p.hostURL()
	if _, ok := resource.Attributes().Get("mackerelio.host.url"); !ok && err == nil {
		resource.Attributes().PutStr("mackerelio.host.url", hostURL)
	}
}
