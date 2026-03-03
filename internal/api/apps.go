package api

import "fmt"

func (c *Client) ListApps(next string) ([]App, Pagination, error) {
	path := "/apps?sort_by=last_build_at&limit=50"
	if next != "" {
		path = fmt.Sprintf("%s&next=%s", path, next)
	}

	data, err := c.get(path)
	if err != nil {
		return nil, Pagination{}, err
	}

	apps, err := decode[[]App](data, "data")
	if err != nil {
		return nil, Pagination{}, err
	}

	paging, _ := decodePagination(data)
	return apps, paging, nil
}
