package cloudflaremetricsreceiver

import (
	"context"
	"net/http"

	"github.com/hasura/go-graphql-client"
)

type apiClient interface {
	ZoneIDByName(zoneName string) (string, error)
}

type graphQLClient interface {
	Query(ctx context.Context, q any, variables map[string]any, options ...graphql.Option) error
}

func newGraghQLClient(cfg *config) graphQLClient {
	httpCli := &http.Client{
		Transport: &bearerAuthTransport{
			BearerToken: cfg.ApiToken,
		},
	}
	return graphql.NewClient("https://api.cloudflare.com/client/v4/graphql", httpCli)
}

type bearerAuthTransport struct {
	Transport   http.RoundTripper
	BearerToken string
}

func (t *bearerAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	req.Header.Set("Authorization", "Bearer "+t.BearerToken)
	return transport.RoundTrip(req)
}
