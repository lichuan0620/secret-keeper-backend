package queueclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lichuan0620/secret-keeper-backend/pkg/models"
	"github.com/pkg/errors"
)

type Interface interface {
	Sync(ctx context.Context, box *models.SyncRequest) error
	Dequeue(ctx context.Context) (*models.DequeueResponse, error)
}

type Type struct {
	endpoint string
	client   *http.Client
}

func New(url string) Interface {
	return &Type{
		endpoint: url,
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   time.Second,
		},
	}
}

func (t *Type) Sync(ctx context.Context, reqBody *models.SyncRequest) error {
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(reqBody); err != nil {
		return errors.Wrap(err, "encode JSON object")
	}
	req, err := http.NewRequest(http.MethodPost, t.buildURL("Sync", t.endpoint), body)
	if err != nil {
		return errors.Wrap(err, "build request")
	}
	resp, err := t.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("non-OK response status code %d", resp.StatusCode)
	}
	return nil
}

func (t *Type) Dequeue(ctx context.Context) (*models.DequeueResponse, error) {
	req, err := http.NewRequest(http.MethodPost, t.buildURL("Dequeue", t.endpoint), nil)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}
	resp, err := t.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("non-OK response status code %d", resp.StatusCode)
	}
	var respBody models.DequeueResponse
	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, errors.Wrap(err, "decode JSON object")
	}
	return &respBody, nil
}

func (t *Type) buildURL(action, version string) string {
	return fmt.Sprintf("%s?action=%s&version=%s", t.endpoint, action, version)
}
