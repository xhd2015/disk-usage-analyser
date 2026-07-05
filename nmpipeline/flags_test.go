package nmpipeline

import (
	"bytes"
	"testing"
)

func TestParsePipelineFlagsRequiresInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	_, _, code, err := ParsePipelineFlags([]string{"--dry-run"}, help, usageLine, &stdout, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
}

func TestParsePipelineFlagsArgPath(t *testing.T) {
	var stdout, stderr bytes.Buffer
	_, input, code, err := ParsePipelineFlags([]string{"/tmp/foo/node_modules"}, help, usageLine, &stdout, &stderr)
	if err != nil || code != -1 {
		t.Fatalf("code=%d err=%v", code, err)
	}
	if len(input.Paths) != 1 || input.Paths[0] != "/tmp/foo/node_modules" {
		t.Fatalf("input=%#v", input)
	}
}

func TestParsePipelineFlagsRecords(t *testing.T) {
	var stdout, stderr bytes.Buffer
	_, input, code, err := ParsePipelineFlags([]string{"--node-modules-records", "inventory.json"}, help, usageLine, &stdout, &stderr)
	if err != nil || code != -1 {
		t.Fatalf("code=%d err=%v", code, err)
	}
	if input.RecordsFile != "inventory.json" || len(input.Paths) != 0 {
		t.Fatalf("input=%#v", input)
	}
}