package api

import (
	"encoding/json"
	"fmt"
)

func (c *Client) ListBuilds(appSlug string, next string) ([]Build, Pagination, error) {
	path := fmt.Sprintf("/apps/%s/builds?sort_by=created_at&limit=50", appSlug)
	if next != "" {
		path = fmt.Sprintf("%s&next=%s", path, next)
	}

	data, err := c.get(path)
	if err != nil {
		return nil, Pagination{}, err
	}

	builds, err := decode[[]Build](data, "data")
	if err != nil {
		return nil, Pagination{}, err
	}

	paging, _ := decodePagination(data)
	return builds, paging, nil
}

func (c *Client) GetBuild(appSlug, buildSlug string) (Build, error) {
	path := fmt.Sprintf("/apps/%s/builds/%s", appSlug, buildSlug)
	data, err := c.get(path)
	if err != nil {
		return Build{}, err
	}
	return decode[Build](data, "data")
}

func (c *Client) TriggerBuild(appSlug string, params TriggerBuildParams) (Build, error) {
	buildParams := map[string]string{
		"branch":     params.Branch,
		"workflow_id": params.WorkflowID,
	}
	if params.CommitMessage != "" {
		buildParams["commit_message"] = params.CommitMessage
	}

	payload := map[string]any{
		"hook_info": map[string]string{
			"type": "bitrise",
		},
		"build_params": buildParams,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return Build{}, err
	}

	resp, err := c.post(fmt.Sprintf("/apps/%s/builds", appSlug), string(data))
	if err != nil {
		return Build{}, err
	}

	return decode[Build](resp, "build_data")
}

func (c *Client) AbortBuild(appSlug, buildSlug string, params AbortBuildParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.post(fmt.Sprintf("/apps/%s/builds/%s/abort", appSlug, buildSlug), string(data))
	return err
}
