// Copyright 2021, Justen Walker and the goodwill contributors
// SPDX-License-Identifier: Apache-2.0

package mage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"time"
)

type ConcordEnv struct {
	Endpoint string
	APIKey   string
	AgentKey string
	Org      string
	Project  string
}

func NewConcordProcess(ctx context.Context, env ConcordEnv, concordYAML io.Reader, files []ZipFile) (string, error) {
	var buf bytes.Buffer
	mpw := multipart.NewWriter(&buf)
	err := mpw.WriteField("org", env.Org)
	if err != nil {
		return "", err
	}
	err = mpw.WriteField("project", env.Project)
	if err != nil {
		return "", err
	}
	payload, err := mpw.CreateFormFile("archive", "payload.zip")
	if err != nil {
		return "", err
	}
	files = append(files, ZipFile{
		SourceReader: concordYAML,
		Dest:         "concord.yml",
	})
	if err := WriteZip(payload, files); err != nil {
		return "", fmt.Errorf("payload file: %w", err)
	}
	if err := mpw.Close(); err != nil {
		return "", fmt.Errorf("close multipart file: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, env.Endpoint+"/api/v1/process", &buf)
	if err != nil {
		return "", fmt.Errorf("could not create http request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+env.APIKey)
	req.Header.Set("Content-Type", mime.FormatMediaType("multipart/form-data", map[string]string{"boundary": mpw.Boundary()}))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not send http request: %w", err)
	}
	defer closeLogError(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	var response ConcordProcess
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("could not parse concord response: %w\nresponse:\n%s", err, string(body))
	}
	return response.InstanceID, nil
}

type ConcordProcess struct {
	InstanceID       string   `json:"instanceId"`
	ParentInstanceID string   `json:"parentInstanceId"`
	ProjectName      string   `json:"projectName"`
	CreatedAt        string   `json:"createdAt"`
	Initiator        string   `json:"initiator"`
	LastUpdatedAt    string   `json:"lastUpdatedAt"`
	Status           string   `json:"status"`
	ChildrenIds      []string `json:"childrenIds"`
}

func GetConcordProcess(ctx context.Context, env ConcordEnv, processID string) (*ConcordProcess, error) {
	req, err := http.NewRequest(http.MethodGet, env.Endpoint+"/api/v1/process/"+processID, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create http request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+env.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send http request: %w", err)
	}
	defer closeLogError(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
	}

	var response ConcordProcess
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("could not parse concord response: %w\nresponse:\n%s", err, string(body))
	}
	return &response, nil
}

func ConcordPing(ctx context.Context, env ConcordEnv) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, env.Endpoint+"/api/v1/server/ping", nil)
	if err != nil {
		return false, fmt.Errorf("could not create http request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+env.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("could not send http request: %w", err)
	}
	defer closeLogError(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("%s: %s", resp.Status, string(body))
	}

	var response struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return false, fmt.Errorf("could not parse concord response: %w\nresponse:\n%s", err, string(body))
	}
	return response.OK, nil
}

func WaitConcordRunning(ctx context.Context, env ConcordEnv) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			ok, _ := ConcordPing(ctx, env)
			if ok {
				return nil
			}
		}
	}
}

func WaitConcordProcess(ctx context.Context, env ConcordEnv, processID string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			//continue
		}
		proc, err := GetConcordProcess(ctx, env, processID)
		if err != nil {
			return err
		}
		switch proc.Status {
		case "FINISHED":
			return nil
		case "FAILED":
			return errors.New("process failed")
		case "CANCELLED":
			return errors.New("process cancelled")
		case "TIMED_OUT":
			return errors.New("process timed out")
		}
	}
}

func closeLogError(c io.Closer) {
	if err := c.Close(); err != nil {
		debug.Println("close error:", err)
	}
}
