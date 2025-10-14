package notify

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

// LoadSlackBlocksFromFile reads a JSON file containing "blocks" (Block Kit template)
// and converts it to slack.Blocks so it can be used in PostMessage.
func LoadSlackBlocksFromFile(path string) (slack.Blocks, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return slack.Blocks{}, fmt.Errorf("read file error: %w", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return slack.Blocks{}, fmt.Errorf("unmarshal payload error: %w", err)
	}

	blocksRaw, ok := payload["blocks"]
	if !ok {
		return slack.Blocks{}, fmt.Errorf("missing 'blocks' field in JSON")
	}

	blocksJSON, err := json.Marshal(blocksRaw)
	if err != nil {
		return slack.Blocks{}, fmt.Errorf("marshal blocks error: %w", err)
	}

	var blocks slack.Blocks
	if err := json.Unmarshal(blocksJSON, &blocks); err != nil {
		return slack.Blocks{}, fmt.Errorf("unmarshal blocks error: %w", err)
	}

	return blocks, nil
}
