package notify

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/fsnotify/fsnotify"
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

var embeddedTemplates embed.FS

type Config struct {
	DevMode     bool   // true -> auto reload template on change
	TemplateDir string // local path to watch (for dev)
}

var (
	config        = Config{DevMode: false, TemplateDir: "templates"}
	cacheMu       sync.RWMutex
	templateCache = make(map[string]*template.Template)
	funcMap       = template.FuncMap{}

	watcherOnce sync.Once
)

// ---------------------- PUBLIC API ----------------------

func SetConfig(c Config) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	config = c
	if config.DevMode {
		startWatcher()
	}
}

func AddTemplateFunc(name string, fn interface{}) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	funcMap[name] = fn
}

func LoadTemplate(name string, data interface{}) ([]slack.Block, error) {
	tmpl, err := getTemplate(name)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		return nil, fmt.Errorf("unmarshal rendered JSON: %w", err)
	}

	blocksRaw, ok := payload["blocks"]
	if !ok {
		return nil, fmt.Errorf("missing 'blocks' field")
	}

	blocksJSON, err := json.Marshal(blocksRaw)
	if err != nil {
		return nil, err
	}

	var blocks slack.Blocks
	if err := json.Unmarshal(blocksJSON, &blocks); err != nil {
		return nil, err
	}
	return blocks.BlockSet, nil
}

// ---------------------- INTERNAL ----------------------

func getTemplate(name string) (*template.Template, error) {
	cacheMu.RLock()
	tmpl, ok := templateCache[name]
	cacheMu.RUnlock()

	if ok && !config.DevMode {
		return tmpl, nil
	}

	var content []byte
	var err error

	// Try embedded first
	content, err = embeddedTemplates.ReadFile(filepath.Join("templates", name))
	if err != nil {
		// Fallback to local file
		content, err = os.ReadFile(filepath.Join(config.TemplateDir, name))
		if err != nil {
			return nil, fmt.Errorf("template not found: %s", name)
		}
	}

	tmpl, err = template.New(name).Funcs(funcMap).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	cacheMu.Lock()
	templateCache[name] = tmpl
	cacheMu.Unlock()

	return tmpl, nil
}

// ---------------------- WATCHER ----------------------

func startWatcher() {
	watcherOnce.Do(func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Println("⚠️  fsnotify watcher error:", err)
			return
		}

		go func() {
			defer watcher.Close()
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
						file := filepath.Base(event.Name)
						cacheMu.Lock()
						delete(templateCache, file)
						cacheMu.Unlock()
						fmt.Printf("🔄 Reloaded Slack template: %s (%s)\n", file, event.Op)
					}
				case err, ok := <-watcher.Errors:
					if ok {
						fmt.Println("watch error:", err)
					}
				}
			}
		}()

		err = watcher.Add(config.TemplateDir)
		if err != nil {
			fmt.Println("⚠️  Cannot watch directory:", err)
		} else {
			fmt.Println("👀 Watching template directory:", config.TemplateDir)
		}
	})
}
