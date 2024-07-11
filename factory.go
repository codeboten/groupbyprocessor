package groupbyprocessor

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"

	"github.com/codeboten/groupbyprocessor/internal/metadata"
)

var (
	consumerCapabilities = consumer.Capabilities{MutatesData: true}
)

// createDefaultConfig creates the default configuration for the processor.
func createDefaultConfig() component.Config {
	return &Config{MaxLogsBuffered: 64}
}

func createGroupByProcessor(set processor.Settings, cfg *Config) (*groupByProcessor, error) {
	cache, _ := lru.New[string, plog.LogRecord](int(cfg.MaxLogsBuffered))

	telemetryBuilder, err := metadata.NewTelemetryBuilder(set.TelemetrySettings)
	if err != nil {
		return nil, err
	}
	return &groupByProcessor{telemetryBuilder: telemetryBuilder, cache: cache}, nil
}

// NewFactory returns a new factory for the Filter processor.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		metadata.Type,
		createDefaultConfig,
		processor.WithLogs(createLogsProcessor, metadata.LogsStability))
}

// createLogsProcessor creates a logs processor based on this config.
func createLogsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Logs) (processor.Logs, error) {
	oCfg := cfg.(*Config)
	gap, err := createGroupByProcessor(set, oCfg)
	if err != nil {
		return nil, err
	}

	return processorhelper.NewLogsProcessor(
		ctx,
		set,
		cfg,
		nextConsumer,
		gap.processLogs,
		processorhelper.WithCapabilities(consumerCapabilities))
}
