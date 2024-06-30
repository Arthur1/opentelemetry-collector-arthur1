package mackerelattributesprocessor

import (
	"github.com/Arthur1/opentelemetry-collector-arthur1/processor/mackerelattributesprocessor/internal/metadata"
	mackerelAgentConfig "github.com/mackerelio/mackerel-agent/config"
	"go.opentelemetry.io/collector/component"
)

type Config struct {
	ConfigFilePath                    string `mapstructure:"config_file_path"`
	metadata.ResourceAttributesConfig `mapstructure:"resource_attributes"`
}

var _ component.Config = (*Config)(nil)

func (c *Config) Validate() error {
	if _, err := mackerelAgentConfig.ValidateConfigFile(c.ConfigFilePath); err != nil {
		return err
	}
	return nil
}
