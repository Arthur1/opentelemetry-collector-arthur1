package cloudflaremetricsreceiver

import (
	"context"
	"errors"
	"time"

	"github.com/Arthur1/opentelemetry-collector-arthur1/receiver/cloudflaremetricsreceiver/internal/metadata"
	"github.com/cloudflare/cloudflare-go"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
)

type cloudflareScraper struct {
	cfg    *config
	mb     *metadata.MetricsBuilder
	apiCli apiClient
	gqlCli graphQLClient
}

func newScraper(cfg *config, settings receiver.CreateSettings) *cloudflareScraper {
	apiCli, _ := cloudflare.NewWithAPIToken(cfg.ApiToken)
	return &cloudflareScraper{
		cfg:    cfg,
		mb:     metadata.NewMetricsBuilder(cfg.MetricsBuilderConfig, settings),
		apiCli: apiCli,
		gqlCli: newGraghQLClient(cfg),
	}
}

var query struct {
	Viewer struct {
		Zones []struct {
			HTTPRequests1HGroups []struct {
				Sum struct {
					Bytes          int64
					CachedBytes    int64
					Requests       int64
					CachedRequests int64
				}
				Dimensions struct {
					Datetime string
				}
				Uniq struct {
					Uniques int64
				}
			} `graphql:"httpRequests1hGroups(filter: {date: $date}, limit: 1, orderBy: [datetime_DESC])"`
		} `graphql:"zones(filter: {zoneTag: $zoneId})"`
	}
}

func (s *cloudflareScraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	now := time.Now().UTC()
	var errs error
	for _, zoneDomain := range s.cfg.ZoneDomains {
		zoneID, err := s.apiCli.ZoneIDByName(zoneDomain)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		variables := map[string]any{
			"zoneId": zoneID,
			// TODO: specify datetime
			"date": now.Format(time.DateOnly),
		}

		if err := s.gqlCli.Query(ctx, &query, variables); err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		if len(query.Viewer.Zones) < 1 {
			continue
		}
		if len(query.Viewer.Zones[0].HTTPRequests1HGroups) > 0 {
			t, err := time.Parse(time.RFC3339, query.Viewer.Zones[0].HTTPRequests1HGroups[0].Dimensions.Datetime)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			ts := pcommon.NewTimestampFromTime(t)

			bytes := query.Viewer.Zones[0].HTTPRequests1HGroups[0].Sum.Bytes
			cachedBytes := query.Viewer.Zones[0].HTTPRequests1HGroups[0].Sum.CachedBytes
			uncachedBytes := bytes - cachedBytes
			s.mb.RecordCloudflareHTTPBytesDataPoint(ts, cachedBytes, true)
			s.mb.RecordCloudflareHTTPBytesDataPoint(ts, uncachedBytes, false)

			requests := query.Viewer.Zones[0].HTTPRequests1HGroups[0].Sum.Requests
			cachedRequests := query.Viewer.Zones[0].HTTPRequests1HGroups[0].Sum.CachedRequests
			uncachedRequests := requests - cachedRequests
			s.mb.RecordCloudflareHTTPRequestsDataPoint(ts, cachedRequests, true)
			s.mb.RecordCloudflareHTTPRequestsDataPoint(ts, uncachedRequests, false)

			uniqueIPs := query.Viewer.Zones[0].HTTPRequests1HGroups[0].Uniq.Uniques
			s.mb.RecordCloudflareHTTPUniqueIpsDataPoint(ts, uniqueIPs)
		}

		rb := s.mb.NewResourceBuilder()
		rb.SetCloudflareZoneDomain(zoneDomain)
		s.mb.EmitForResource(metadata.WithResource(rb.Emit()))
	}
	return s.mb.Emit(), errs
}
