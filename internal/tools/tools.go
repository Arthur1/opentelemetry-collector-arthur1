//go:build tools

package tools

import (
	_ "github.com/mfridman/tparse"
	_ "go.opentelemetry.io/collector/cmd/builder"
	_ "go.opentelemetry.io/collector/cmd/mdatagen"
)
