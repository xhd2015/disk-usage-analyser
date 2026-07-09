package usagescan

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sampleTree() TreeResult {
	// root 1000
	//   big/ 600
	//     deep.txt 400
	//     keep.bin 200
	//   small.txt 400
	root := "/tmp/fixture-root"
	return TreeResult{
		Path:      root,
		TotalSize: 1000,
		Min:       1 << 20,
		MaxDepth:  6,
		Tree: TreeNode{
			Name:  ".",
			Path:  root,
			Size:  1000,
			IsDir: true,
			Depth: 0,
			Children: []TreeNode{
				{
					Name:  "big",
					Path:  filepath.Join(root, "big"),
					Size:  600,
					IsDir: true,
					Depth: 1,
					Children: []TreeNode{
						{
							Name:  "deep.txt",
							Path:  filepath.Join(root, "big", "deep.txt"),
							Size:  400,
							IsDir: false,
							Depth: 2,
						},
						{
							Name:  "keep.bin",
							Path:  filepath.Join(root, "big", "keep.bin"),
							Size:  200,
							IsDir: false,
							Depth: 2,
						},
					},
				},
				{
					Name:  "small.txt",
					Path:  filepath.Join(root, "small.txt"),
					Size:  400,
					IsDir: false,
					Depth: 1,
				},
			},
		},
	}
}

func TestLoadTreeResult_roundTrip(t *testing.T) {
	src := sampleTree()
	raw, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	got, err := LoadTreeResult(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("LoadTreeResult: %v", err)
	}
	if got.TotalSize != 1000 || got.Path != src.Path {
		t.Fatalf("unexpected load: %+v", got)
	}
	if len(got.Tree.Children) != 2 {
		t.Fatalf("children: %d", len(got.Tree.Children))
	}
	if got.Min != src.Min {
		t.Fatalf("min: expected %d, got %d", src.Min, got.Min)
	}
}

func TestBuildView_topSkipsRoot(t *testing.T) {
	res := BuildView(sampleTree(), ViewOptions{Top: 10, TopSet: true, Min: 0, MaxDepth: 6})
	if len(res.Matches) == 0 {
		t.Fatal("expected matches")
	}
	if res.Matches[0].Path != filepath.Join("/tmp/fixture-root", "big") {
		t.Fatalf("top0: %+v", res.Matches[0])
	}
	for _, m := range res.Matches {
		if m.Depth == 0 {
			t.Fatal("root should be skipped by default")
		}
	}
}

func TestBuildView_atFocusesTree(t *testing.T) {
	root := "/tmp/fixture-root"
	res := BuildView(sampleTree(), ViewOptions{
		AtPath:   filepath.Join(root, "big"),
		Min:      0,
		MaxDepth: 1,
	})
	// Focused tree at maxDepth 1 shows deep.txt and keep.bin as children.
	if len(res.Tree.Children) != 2 {
		t.Fatalf("focused children: %d %#v", len(res.Tree.Children), res.Tree.Children)
	}
	if WantMatchSection(ViewOptions{AtPath: filepath.Join(root, "big")}) {
		t.Fatal("--at alone must not activate match section")
	}
}

func TestBuildView_findAndSuffix(t *testing.T) {
	res := BuildView(sampleTree(), ViewOptions{Find: "deep", Top: 10, TopSet: true, Min: 0, MaxDepth: 6})
	if len(res.Matches) != 1 || res.Matches[0].Name != "deep.txt" {
		t.Fatalf("find: %+v", res.Matches)
	}
	res2 := BuildView(sampleTree(), ViewOptions{Suffix: ".bin", Top: 10, TopSet: true, Min: 0, MaxDepth: 6})
	if len(res2.Matches) != 1 || res2.Matches[0].Name != "keep.bin" {
		t.Fatalf("suffix: %+v", res2.Matches)
	}
}

func TestBuildView_minSize(t *testing.T) {
	res := BuildView(sampleTree(), ViewOptions{Top: 50, TopSet: true, Min: 400, MaxDepth: 6})
	for _, m := range res.Matches {
		if m.Size < 400 {
			t.Fatalf("min size violated: %+v", m)
		}
	}
	// big=600, deep.txt=400, small.txt=400
	if len(res.Matches) != 3 {
		t.Fatalf("matches: %d %#v", len(res.Matches), res.Matches)
	}
}

func TestRunCLI_inspectTopJSON(t *testing.T) {
	raw, err := json.Marshal(sampleTree())
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "scan.json")
	if err := os.WriteFile(path, raw, 0644); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	code, err := RunCLI([]string{"--inspect", path, "--json", "--top", "2"}, CLIOptions{Stdout: &stdout})
	if err != nil || code != 0 {
		t.Fatalf("code=%d err=%v out=%s", code, err, stdout.String())
	}
	var got ViewResult
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &got); err != nil {
		t.Fatalf("json: %v\n%s", err, stdout.String())
	}
	if len(got.Matches) != 2 {
		t.Fatalf("got: %+v", got)
	}
	text := FormatViewText(got, ViewOptions{Top: 2, TopSet: true})
	if !strings.Contains(text, "PATH:") {
		t.Fatal("text format missing PATH header")
	}
	if !strings.Contains(text, "TOP 2") {
		t.Fatal("text format missing TOP 2")
	}
}
