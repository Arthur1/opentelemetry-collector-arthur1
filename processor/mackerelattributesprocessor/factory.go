package mackerelattributesprocessor

import (
	"context"
	"fmt"

	"github.com/Arthur1/opentelemetry-collector-arthur1/processor/mackerelattributesprocessor/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

var consumerCapabilities = consumer.Capabilities{MutatesData: true}

func NewFactory() processor.Factory {
	return processor.NewFactory(
		metadata.Type,
		createDefaultConfig,
		processor.WithTraces(createTracesProcessor, metadata.TracesStability),
		processor.WithLogs(createLogsProcessor, metadata.LogsStability),
		processor.WithMetrics(createMetricsProcessor, metadata.MetricsStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ConfigFilePath: "/etc/mackerel-agent/mackerel-agent.conf",
	}
}

func createTracesProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	next consumer.Traces,
) (processor.Traces, error) {
	c, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid config")
	}
	mp := newMackerelProcessor(c)
	return processorhelper.NewTracesProcessor(
		ctx, set, cfg, next, mp.processTraces,
		processorhelper.WithCapabilities(consumerCapabilities),
		processorhelper.WithStart(mp.Start),
		processorhelper.WithShutdown(mp.Shutdown),
	)
}

func createLogsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	next consumer.Logs,
) (processor.Logs, error) {
	c, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid config")
	}
	mp := newMackerelProcessor(c)
	return processorhelper.NewLogsProcessor(
		ctx, set, cfg, next, mp.processLogs,
		processorhelper.WithCapabilities(consumerCapabilities),
		processorhelper.WithStart(mp.Start),
		processorhelper.WithShutdown(mp.Shutdown),
	)
}

func createMetricsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	next consumer.Metrics,
) (processor.Metrics, error) {
	c, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid config")
	}
	mp := newMackerelProcessor(c)
	return processorhelper.NewMetricsProcessor(
		ctx, set, cfg, next, mp.processMetrics,
		processorhelper.WithCapabilities(consumerCapabilities),
		processorhelper.WithStart(mp.Start),
		processorhelper.WithShutdown(mp.Shutdown),
	)
}
