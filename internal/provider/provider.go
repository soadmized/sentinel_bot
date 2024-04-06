//nolint:stylecheck
package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/soadmized/sentinel/pkg/dataset"

	"sentinel_bot/internal/config"
)

const (
	lastValues = "/last_values"
	sensorIDs  = "/sensor_ids"
)

// Provider provides data from Sentinel server.
type Provider struct {
	SentinelUrl  string
	SentinelUser string
	SentinelPass string
}

func New(conf config.Config) *Provider {
	return &Provider{
		SentinelUrl:  conf.SentinelServerUrl,
		SentinelUser: conf.SentinelUser,
		SentinelPass: conf.SentinelPass,
	}
}

// LastValues returns last values received from sensor.
func (p *Provider) LastValues(ctx context.Context, sensorID string) (*dataset.Dataset, error) {
	body, err := lastValuesReqBody(sensorID)
	if err != nil {
		return nil, err
	}

	req, err := p.makeReq(ctx, lastValues, body)
	if err != nil {
		return nil, err
	}

	var response dataset.Dataset

	err = p.sendReq(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *Provider) SensorIDs(ctx context.Context) ([]string, error) {
	req, err := p.makeReq(ctx, sensorIDs, nil)
	if err != nil {
		return nil, err
	}

	var response []string

	err = p.sendReq(req, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func lastValuesReqBody(id string) (io.Reader, error) {
	body := struct {
		ID string `json:"id"`
	}{
		ID: id,
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, errors.Wrap(err, "encode body")
	}

	return &buf, nil
}

func (p *Provider) makeReq(ctx context.Context, endpoint string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", p.SentinelUrl, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "create request")
	}

	req.SetBasicAuth(p.SentinelUser, p.SentinelPass)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (p *Provider) sendReq(req *http.Request, response any) error {
	client := http.Client{} //nolint:exhaustruct

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return errors.Wrap(err, "decode response")
	}

	return nil
}
