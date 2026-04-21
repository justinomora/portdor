package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const DefaultAddr = "http://localhost:4242"

type Client struct {
	base string
	http *http.Client
}

func New(addr string) *Client {
	return &Client{base: addr, http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) do(method, path string, body any) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, err := http.NewRequest(method, c.base+path, &buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.http.Do(req)
}

func (c *Client) IsReachable() bool {
	resp, err := c.http.Get(c.base + "/api/status")
	return err == nil && resp.StatusCode == 200
}

func decode(resp *http.Response, v any) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		var e map[string]string
		json.NewDecoder(resp.Body).Decode(&e)
		return fmt.Errorf("server error %d: %s", resp.StatusCode, e["error"])
	}
	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

func (c *Client) ListServices() ([]map[string]any, error) {
	resp, err := c.do("GET", "/api/services", nil)
	if err != nil {
		return nil, err
	}
	var result struct{ Services []map[string]any }
	return result.Services, decode(resp, &result)
}

func (c *Client) RegisterService(name, command, cwd string, port int, project string) error {
	resp, err := c.do("POST", "/api/services", map[string]any{
		"name": name, "command": command, "cwd": cwd, "port": port, "project": project,
	})
	if err != nil {
		return err
	}
	return decode(resp, nil)
}

func (c *Client) UnregisterService(name string) error {
	resp, err := c.do("DELETE", "/api/services/"+name, nil)
	if err != nil {
		return err
	}
	return decode(resp, nil)
}

func (c *Client) UpdateService(name string, fields map[string]any) error {
	resp, err := c.do("PUT", "/api/services/"+name, fields)
	if err != nil {
		return err
	}
	return decode(resp, nil)
}

func (c *Client) StopService(name string) error {
	resp, err := c.do("POST", "/api/services/"+name+"/stop", nil)
	if err != nil {
		return err
	}
	return decode(resp, nil)
}

func (c *Client) KillService(name string) error {
	resp, err := c.do("POST", "/api/services/"+name+"/kill", nil)
	if err != nil {
		return err
	}
	return decode(resp, nil)
}

func (c *Client) RestartService(name string) error {
	resp, err := c.do("POST", "/api/services/"+name+"/restart", nil)
	if err != nil {
		return err
	}
	return decode(resp, nil)
}
