// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceBuilder(t *testing.T) {
	for _, test := range []string{"default", "all_set", "none_set"} {
		t.Run(test, func(t *testing.T) {
			cfg := loadResourceAttributesConfig(t, test)
			rb := NewResourceBuilder(cfg)
			rb.SetMackerelioHostID("mackerelio.host.id-val")
			rb.SetMackerelioHostURL("mackerelio.host.url-val")
			rb.SetMackerelioOrgName("mackerelio.org.name-val")

			res := rb.Emit()
			assert.Equal(t, 0, rb.Emit().Attributes().Len()) // Second call should return empty Resource

			switch test {
			case "default":
				assert.Equal(t, 3, res.Attributes().Len())
			case "all_set":
				assert.Equal(t, 3, res.Attributes().Len())
			case "none_set":
				assert.Equal(t, 0, res.Attributes().Len())
				return
			default:
				assert.Failf(t, "unexpected test case: %s", test)
			}

			val, ok := res.Attributes().Get("mackerelio.host.id")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "mackerelio.host.id-val", val.Str())
			}
			val, ok = res.Attributes().Get("mackerelio.host.url")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "mackerelio.host.url-val", val.Str())
			}
			val, ok = res.Attributes().Get("mackerelio.org.name")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "mackerelio.org.name-val", val.Str())
			}
		})
	}
}