package elasticsearch

import "encoding/json"

const (
	StatusGreen  = "green"
	StatusYellow = "yellow"
	StatusRed    = "red"
)

func (c *Client) Health() (string, error) {
	res, err := c.get("_cluster/health", nil)
	if err != nil {
		return StatusRed, err
	}
	health := map[string]interface{}{}
	if err := json.Unmarshal(res, &health); err != nil {
		return StatusRed, err
	}
	if status, ok := health["status"].(string); ok && (status == StatusGreen || status == StatusYellow) {
		return status, nil
	}
	return StatusRed, nil
}
