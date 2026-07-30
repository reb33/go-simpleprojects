package api

import (
	"3-struct/config"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Api struct {
	config config.Config
	client *http.Client
}

func NewApi(config config.Config) *Api {
	return &Api{
		config: config,
		client: &http.Client{},
	}
}

type Metadata struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Private   bool      `json:"private"`
}

type ResponseWithMetadata struct {
	Metadata Metadata `json:"metadata"`
}

func (api *Api) CreateBin(data []byte) (*Metadata, error) {
	req, err := http.NewRequest("POST", api.config.ApiUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", api.config.Key)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("bad status code")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response ResponseWithMetadata
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response.Metadata, nil
}

func (api *Api) buildURL(id string) string {
	u, _ := url.Parse(api.config.ApiUrl)
	u.Path = path.Join(u.Path, id)
	return u.String()
}

func (api *Api) GetBin(id string) ([]byte, error) {
	req, err := http.NewRequest("GET", api.buildURL(id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Master-Key", api.config.Key)
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (api *Api) DeleteBin(id string) error {
	req, err := http.NewRequest("DELETE", api.buildURL(id), nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Master-Key", api.config.Key)
	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

func (api *Api) UpdateBin(id string, data []byte) error {
	req, err := http.NewRequest("PUT", api.buildURL(id), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", api.config.Key)

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}
