package api

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func (c *Client) GetBuildLog(appSlug, buildSlug string) (string, error) {
	path := fmt.Sprintf("/apps/%s/builds/%s/log", appSlug, buildSlug)
	data, err := c.get(path)
	if err != nil {
		return "", err
	}

	var log BuildLog
	if err := json.Unmarshal(data, &log); err != nil {
		return "", fmt.Errorf("parsing log response: %w", err)
	}

	// If there's a raw log URL, fetch it directly
	if log.ExpiringRawLogURL != "" {
		raw, err := c.getRaw(log.ExpiringRawLogURL)
		if err != nil {
			return "", fmt.Errorf("fetching raw log: %w", err)
		}
		return string(raw), nil
	}

	// Otherwise assemble from chunks
	if len(log.LogChunks) == 0 {
		return "(no log output yet)", nil
	}

	sort.Slice(log.LogChunks, func(i, j int) bool {
		return log.LogChunks[i].Position < log.LogChunks[j].Position
	})

	var sb strings.Builder
	for _, chunk := range log.LogChunks {
		sb.WriteString(chunk.Chunk)
	}
	return sb.String(), nil
}
