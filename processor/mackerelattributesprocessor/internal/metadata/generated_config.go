// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/confmap"
)

// ResourceAttributeConfig provides common config for a particular resource attribute.
type ResourceAttributeConfig struct {
	Enabled bool `mapstructure:"enabled"`

	enabledSetByUser bool
}

func (rac *ResourceAttributeConfig) Unmarshal(parser *confmap.Conf) error {
	if parser == nil {
		return nil
	}
	err := parser.Unmarshal(rac)
	if err != nil {
		return err
	}
	rac.enabledSetByUser = parser.IsSet("enabled")
	return nil
}

// ResourceAttributesConfig provides config for mackerelattributes resource attributes.
type ResourceAttributesConfig struct {
	MackerelioHostID  ResourceAttributeConfig `mapstructure:"mackerelio.host.id"`
	MackerelioHostURL ResourceAttributeConfig `mapstructure:"mackerelio.host.url"`
	MackerelioOrgName ResourceAttributeConfig `mapstructure:"mackerelio.org.name"`
}

func DefaultResourceAttributesConfig() ResourceAttributesConfig {
	return ResourceAttributesConfig{
		MackerelioHostID: ResourceAttributeConfig{
			Enabled: true,
		},
		MackerelioHostURL: ResourceAttributeConfig{
			Enabled: true,
		},
		MackerelioOrgName: ResourceAttributeConfig{
			Enabled: true,
		},
	}
}
