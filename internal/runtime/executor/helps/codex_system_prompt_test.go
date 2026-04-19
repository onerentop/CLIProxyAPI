package helps

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/config"
	"github.com/tidwall/gjson"
)

func writePromptFile(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "prompt.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write prompt file: %v", err)
	}
	return path
}

func TestInjectCodexSystemPrompt_DisabledIsNoop(t *testing.T) {
	cfg := &config.Config{}
	body := []byte(`{"input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"hi"}]}]}`)
	got := InjectCodexSystemPrompt(body, cfg, "gpt-5.4")
	if string(got) != string(body) {
		t.Fatalf("expected body unchanged when disabled; got %s", string(got))
	}
}

func TestInjectCodexSystemPrompt_PrependsToArrayInput(t *testing.T) {
	dir := t.TempDir()
	path := writePromptFile(t, dir, "DEV RULES")
	cfg := &config.Config{}
	cfg.CodexSystemPrompt = config.CodexSystemPrompt{Enabled: true, File: path}

	body := []byte(`{"input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"hi"}]}]}`)
	got := InjectCodexSystemPrompt(body, cfg, "gpt-5.4")

	arr := gjson.GetBytes(got, "input").Array()
	if len(arr) != 2 {
		t.Fatalf("input len = %d, want 2; body=%s", len(arr), string(got))
	}
	if role := arr[0].Get("role").String(); role != "developer" {
		t.Fatalf("first role = %q, want developer", role)
	}
	if txt := arr[0].Get("content.0.text").String(); txt != "DEV RULES" {
		t.Fatalf("first text = %q, want DEV RULES", txt)
	}
	if role := arr[1].Get("role").String(); role != "user" {
		t.Fatalf("second role = %q, want user", role)
	}
}

func TestInjectCodexSystemPrompt_WrapsStringInput(t *testing.T) {
	dir := t.TempDir()
	path := writePromptFile(t, dir, "DEV RULES")
	cfg := &config.Config{}
	cfg.CodexSystemPrompt = config.CodexSystemPrompt{Enabled: true, File: path}

	body := []byte(`{"input":"hello"}`)
	got := InjectCodexSystemPrompt(body, cfg, "gpt-5.4")

	arr := gjson.GetBytes(got, "input").Array()
	if len(arr) != 2 {
		t.Fatalf("input len = %d, want 2; body=%s", len(arr), string(got))
	}
	if txt := arr[1].Get("content.0.text").String(); txt != "hello" {
		t.Fatalf("wrapped user text = %q, want hello", txt)
	}
}

func TestInjectCodexSystemPrompt_ModelPrefixFilterSkips(t *testing.T) {
	dir := t.TempDir()
	path := writePromptFile(t, dir, "DEV RULES")
	cfg := &config.Config{}
	cfg.CodexSystemPrompt = config.CodexSystemPrompt{
		Enabled: true,
		File:    path,
		Models:  []string{"codex"},
	}

	body := []byte(`{"input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"hi"}]}]}`)
	got := InjectCodexSystemPrompt(body, cfg, "gpt-4o")
	if string(got) != string(body) {
		t.Fatalf("expected no injection when model prefix does not match; got %s", string(got))
	}

	got2 := InjectCodexSystemPrompt(body, cfg, "codex-latest")
	arr := gjson.GetBytes(got2, "input").Array()
	if len(arr) != 2 {
		t.Fatalf("expected injection for matching model; got %s", string(got2))
	}
}

func TestInjectCodexSystemPrompt_DoesNotTouchInstructions(t *testing.T) {
	dir := t.TempDir()
	path := writePromptFile(t, dir, "DEV RULES")
	cfg := &config.Config{}
	cfg.CodexSystemPrompt = config.CodexSystemPrompt{Enabled: true, File: path}

	body := []byte(`{"instructions":"ORIGINAL","input":[]}`)
	got := InjectCodexSystemPrompt(body, cfg, "gpt-5.4")
	if inst := gjson.GetBytes(got, "instructions").String(); inst != "ORIGINAL" {
		t.Fatalf("instructions mutated: got %q", inst)
	}
}

func TestInjectCodexSystemPrompt_CustomRole(t *testing.T) {
	dir := t.TempDir()
	path := writePromptFile(t, dir, "RULES")
	cfg := &config.Config{}
	cfg.CodexSystemPrompt = config.CodexSystemPrompt{Enabled: true, File: path, Role: "user"}

	body := []byte(`{"input":[]}`)
	got := InjectCodexSystemPrompt(body, cfg, "gpt-5.4")
	if role := gjson.GetBytes(got, "input.0.role").String(); role != "user" {
		t.Fatalf("role = %q, want user", role)
	}
}
