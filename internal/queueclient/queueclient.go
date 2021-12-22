package queueclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/lichuan0620/secret-keeper-backend/pkg/models"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service/model"
	"github.com/lichuan0620/secret-keeper-backend/pkg/telemetry/log"
	"github.com/pkg/errors"
)

type Interface interface {
	Sync(ctx context.Context, box *models.SyncRequest) error
	Dequeue(ctx context.Context) (*models.DequeueResponse, error)
}

type Type struct {
	endpoint string
	client   *http.Client
	logger   logr.Logger
}

func New(url string) Interface {
	return &Type{
		endpoint: url,
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   3 * time.Second,
		},
		logger: log.New().WithName("queueclient"),
	}
}

func (t *Type) Sync(ctx context.Context, reqBody *models.SyncRequest) error {
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(reqBody); err != nil {
		return errors.Wrap(err, "encode JSON object")
	}
	url := t.buildURL("Sync", models.Version)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return errors.Wrapf(err, "build request %s", url)
	}
	req.Header.Set(model.HeaderContentType, model.ContentTypeJSON)
	resp, err := t.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("request %s returns non-OK response status code %d", url, resp.StatusCode)
	}
	return nil
}

func (t *Type) Dequeue(ctx context.Context) (*models.DequeueResponse, error) {
	url := t.buildURL("Dequeue", models.Version)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "build request %s", url)
	}
	req.Header.Set(model.HeaderContentType, model.ContentTypeJSON)
	resp, err := t.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("request %s returns non-OK response status code %d", url, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}
	respBody := struct {
		Result models.DequeueResponse `json:"Result"`
	}{}
	if err = json.Unmarshal(body, &respBody); err != nil {
		return nil, errors.Wrap(err, "decode JSON object")
	}
	return &respBody.Result, nil
}

func (t *Type) buildURL(action, version string) string {
	return fmt.Sprintf("%s?Action=%s&Version=%s", t.endpoint, action, version)
}
