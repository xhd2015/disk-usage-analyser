package usagescan

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

var scanSem = make(chan struct{}, 20)

func ScanFlat(path string, opts Options) (Result, error) {
	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return Result{}, fmt.Errorf("invalid path %s: %w", path, err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Result{}, fmt.Errorf("path does not exist: %s", absPath)
		}
		return Result{}, fmt.Errorf("cannot access %s: %w", absPath, err)
	}
	if !info.IsDir() {
		return Result{}, fmt.Errorf("not a directory: %s", absPath)
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return Result{}, fmt.Errorf("cannot read directory %s: %w", absPath, err)
	}

	var subDirs []fs.DirEntry
	var files []fs.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			subDirs = append(subDirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	items := make([]Item, 0, len(entries))
	emit := func(item Item) {
		if opts.OnItem != nil {
			opts.OnItem(item)
		}
	}

	for _, entry := range files {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		item := Item{
			Name:  entry.Name(),
			Size:  info.Size(),
			IsDir: false,
		}
		items = append(items, item)
		emit(item)
	}

	for _, entry := range subDirs {
		pending := Item{
			Name:  entry.Name(),
			Size:  0,
			IsDir: true,
		}
		emit(pending)
	}

	type dirResult struct {
		name string
		size int64
	}

	resultChan := make(chan dirResult, len(subDirs))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)

	for _, dir := range subDirs {
		wg.Add(1)
		go func(d fs.DirEntry) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic scanning %s: %v", d.Name(), r)
				}
			}()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()

			fullPath := filepath.Join(absPath, d.Name())
			onProgress := func(currentSize int64) {
				emit(Item{
					Name:  d.Name(),
					Size:  currentSize,
					IsDir: true,
				})
			}

			size := getDirSizeWithCache(ctx, fullPath, onProgress)
			select {
			case resultChan <- dirResult{name: d.Name(), size: size}:
			case <-ctx.Done():
			}
		}(dir)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for res := range resultChan {
		item := Item{
			Name:  res.name,
			Size:  res.size,
			IsDir: true,
		}
		items = append(items, item)
		emit(item)
	}

	sortItems(items)

	var totalSize int64
	for _, item := range items {
		totalSize += item.Size
	}

	return Result{
		Path:      absPath,
		TotalSize: totalSize,
		Items:     items,
	}, nil
}

func sortItems(items []Item) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Size != items[j].Size {
			return items[i].Size > items[j].Size
		}
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir
		}
		return items[i].Name < items[j].Name
	})
}

func getDirSizeWithCache(ctx context.Context, path string, onProgress func(int64)) int64 {
	entry, exists := GlobalCache.GetOrCreateEntry(path)

	if !exists {
		go scanDirRecursive(ctx, path, entry)
	}

	unsubscribe := entry.Subscribe(func(s int64) {
		if onProgress != nil {
			onProgress(s)
		}
	})
	defer unsubscribe()

	select {
	case <-entry.doneCh:
		return entry.Size
	case <-ctx.Done():
		return entry.Size
	}
}

func scanDirRecursive(ctx context.Context, dirPath string, entry *CacheEntry) {
	defer entry.MarkDone()

	select {
	case scanSem <- struct{}{}:
	case <-ctx.Done():
		return
	}

	entries, err := os.ReadDir(dirPath)
	<-scanSem

	if err != nil {
		log.Printf("Error reading %s: %v", dirPath, err)
		return
	}

	var (
		mu          sync.Mutex
		filesSize   int64
		subDirSizes = make(map[string]int64)
		dirty       bool
		wg          sync.WaitGroup
	)

	ticker := time.NewTicker(200 * time.Millisecond)
	doneCh := make(chan struct{})

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-doneCh:
				return
			case <-ticker.C:
				mu.Lock()
				if dirty {
					total := filesSize
					for _, s := range subDirSizes {
						total += s
					}
					entry.UpdateSize(total)
					dirty = false
				}
				mu.Unlock()
			}
		}
	}()

	updateLocal := func(name string, size int64) {
		mu.Lock()
		subDirSizes[name] = size
		dirty = true
		mu.Unlock()
	}

	for _, e := range entries {
		if ctx.Err() != nil {
			break
		}

		if !e.IsDir() {
			info, err := e.Info()
			if err == nil {
				mu.Lock()
				filesSize += info.Size()
				dirty = true
				mu.Unlock()
			}
		} else {
			subPath := filepath.Join(dirPath, e.Name())
			subName := e.Name()

			wg.Add(1)

			subEntry, exists := GlobalCache.GetOrCreateEntry(subPath)

			if !exists {
				go scanDirRecursive(ctx, subPath, subEntry)
			}

			unsub := subEntry.Subscribe(func(size int64) {
				updateLocal(subName, size)
			})

			go func() {
				defer wg.Done()
				defer unsub()
				subEntry.Wait()
			}()
		}
	}

	wg.Wait()
	close(doneCh)

	mu.Lock()
	total := filesSize
	for _, s := range subDirSizes {
		total += s
	}
	entry.UpdateSize(total)
	mu.Unlock()
}