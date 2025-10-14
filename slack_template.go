package notify

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"

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

	// cache to avoid reading the same template multiple times
var templateCache = struct {
	sync.RWMutex
	items map[string]string
}{
	items: make(map[string]string),
}

// LoadSlackTemplate reads a JSON template file and supports:
// - Replacing placeholders {{key}} in the content
// - Caching the parsed template to avoid repeated reads
// - Loading from embed.FS (if fsys != nil)
func LoadSlackTemplate(path string, data map[string]string, fsys fs.FS) (slack.Blocks, error) {
	var jsonStr string
	var err error

	// 1️⃣ — Check cache
	templateCache.RLock()
	cached, ok := templateCache.items[path]
	templateCache.RUnlock()

	if ok {
		jsonStr = cached
	} else {
		// 2️⃣ — Read file (from embed.FS or OS)
		var content []byte
		if fsys != nil {
			content, err = fs.ReadFile(fsys, path)
		} else {
			content, err = os.ReadFile(path)
		}
		if err != nil {
			return slack.Blocks{}, fmt.Errorf("read file error: %w", err)
		}

		jsonStr = string(content)

		// Save cache
		templateCache.Lock()
		templateCache.items[path] = jsonStr
		templateCache.Unlock()
	}

	// 3️⃣ — Replace placeholder {{key}} with value
	for k, v := range data {
		placeholder := fmt.Sprintf("{{%s}}", k)
		jsonStr = strings.ReplaceAll(jsonStr, placeholder, v)
	}

	// 4️⃣ — Parse to slack.Blocks
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &payload); err != nil {
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