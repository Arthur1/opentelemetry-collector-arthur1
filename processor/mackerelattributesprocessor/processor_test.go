package mackerelattributesprocessor

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/stretchr/testify/assert"
)

var configTmpl, _ = template.New("mackerelAgentConf").Parse(`
apikey = "{{ .ApiKey }}"
root = "{{ .Root }}"
apibase = "{{ .ApiBase }}"
`)

func TestMackerelProcessorSync(t *testing.T) {
	d := t.TempDir()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		org := mackerel.Org{Name: "test-org"}
		b, err := json.Marshal(org)
		assert.NoError(t, err)
		w.Write(b)
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	configFilePath := path.Join(d, "mackerel-agent.conf")
	configFile, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE, 0600)
	assert.NoError(t, err)
	err = configTmpl.Execute(configFile, struct{ ApiKey, Root, ApiBase string }{
		ApiKey:  "dummy",
		Root:    d,
		ApiBase: ts.URL,
	})
	assert.NoError(t, err)

	hostIDFilePath := path.Join(d, "id")
	err = os.WriteFile(hostIDFilePath, []byte("12345678"), 0600)
	assert.NoError(t, err)

	mp := newMackerelProcessor(&Config{
		ConfigFilePath: configFilePath,
	})
	err = mp.sync(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "test-org", mp.orgName)
	assert.Equal(t, "12345678", mp.hostID)

	hostURL, err := mp.hostURL()
	assert.NoError(t, err)
	assert.Equal(t, ts.URL+"/orgs/test-org/hosts/12345678", hostURL)
}
