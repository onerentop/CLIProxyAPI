package helps

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// codexPromptFileCache caches the contents of the Codex system-prompt file,
// keyed by (path, mtime) so edits are picked up without a process restart.
type codexPromptFileCache struct {
	mu      sync.RWMutex
	path    string
	mtime   time.Time
	content string
}

var codexPromptCache codexPromptFileCache

func readCodexPromptFile(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return ""
	}
	stat, err := os.Stat(trimmed)
	if err != nil {
		log.Debugf("codex-system-prompt: stat %q failed: %v", trimmed, err)
		return ""
	}

	codexPromptCache.mu.RLock()
	if codexPromptCache.path == trimmed && codexPromptCache.mtime.Equal(stat.ModTime()) {
		cached := codexPromptCache.content
		codexPromptCache.mu.RUnlock()
		return cached
	}
	codexPromptCache.mu.RUnlock()

	data, err := os.ReadFile(trimmed)
	if err != nil {
		log.Debugf("codex-system-prompt: read %q failed: %v", trimmed, err)
		return ""
	}
	text := string(data)

	codexPromptCache.mu.Lock()
	codexPromptCache.path = trimmed
	codexPromptCache.mtime = stat.ModTime()
	codexPromptCache.content = text
	codexPromptCache.mu.Unlock()
	return text
}

func codexPromptModelMatches(prefixes []string, model string) bool {
	if len(prefixes) == 0 {
		return true
	}
	lower := strings.ToLower(strings.TrimSpace(model))
	if lower == "" {
		return false
	}
	for _, p := range prefixes {
		if prefix := strings.ToLower(strings.TrimSpace(p)); prefix != "" && strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

// InjectCodexSystemPrompt prepends a developer-role message to the Responses
// API input[] array when configured. It never touches the `instructions`
// field (server-side validated on the ChatGPT-Codex backend).
func InjectCodexSystemPrompt(body []byte, cfg *config.Config, model string) []byte {
	if cfg == nil || !cfg.CodexSystemPrompt.Enabled {
		return body
	}
	if !codexPromptModelMatches(cfg.CodexSystemPrompt.Models, model) {
		return body
	}
	text := strings.TrimSpace(readCodexPromptFile(cfg.CodexSystemPrompt.File))
	if text == "" {
		return body
	}

	role := strings.TrimSpace(cfg.CodexSystemPrompt.Role)
	if role == "" {
		role = "developer"
	}

	item := map[string]any{
		"type": "message",
		"role": role,
		"content": []map[string]any{
			{"type": "input_text", "text": text},
		},
	}
	itemJSON, err := json.Marshal(item)
	if err != nil {
		log.Debugf("codex-system-prompt: marshal item failed: %v", err)
		return body
	}

	input := gjson.GetBytes(body, "input")
	switch {
	case !input.Exists(), input.Type == gjson.Null:
		arrRaw := "[" + string(itemJSON) + "]"
		if out, errSet := sjson.SetRawBytes(body, "input", []byte(arrRaw)); errSet == nil {
			return out
		}
	case input.Type == gjson.String:
		// Wrap the original string input as a user message and prepend ours.
		userItem := map[string]any{
			"type": "message",
			"role": "user",
			"content": []map[string]any{
				{"type": "input_text", "text": input.String()},
			},
		}
		userJSON, errMarshal := json.Marshal(userItem)
		if errMarshal != nil {
			log.Debugf("codex-system-prompt: marshal user item failed: %v", errMarshal)
			return body
		}
		arrRaw := "[" + string(itemJSON) + "," + string(userJSON) + "]"
		if out, errSet := sjson.SetRawBytes(body, "input", []byte(arrRaw)); errSet == nil {
			return out
		}
	case input.IsArray():
		inner := strings.TrimSpace(input.Raw)
		inner = strings.TrimPrefix(inner, "[")
		inner = strings.TrimSuffix(inner, "]")
		inner = strings.TrimSpace(inner)
		var arrRaw string
		if inner == "" {
			arrRaw = "[" + string(itemJSON) + "]"
		} else {
			arrRaw = "[" + string(itemJSON) + "," + inner + "]"
		}
		if out, errSet := sjson.SetRawBytes(body, "input", []byte(arrRaw)); errSet == nil {
			return out
		}
	}
	return body
}
