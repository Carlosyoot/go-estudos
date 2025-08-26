package infra

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

//===================Global vars ====================//

var leitor sync.RWMutex
var MapValue = make(map[string]struct{})

//-=================================================-//

func Contagem() int {
	leitor.RLock()
	defer leitor.RUnlock()
	return len(MapValue)
}

func Observar(ctx context.Context, dir string) error {
	if err := os.MkdirAll(dir, 00755); err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := watcher.Add(dir); err != nil {
		_ = watcher.Close()
		return err
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case ev := <-watcher.Events:

				if ev.Has(fsnotify.Create) {
					if strings.EqualFold(filepath.Ext(ev.Name), ".rem") {
						leitor.RLock()
						MapValue[ev.Name] = struct{}{}
						leitor.Unlock()
					}
				}

			case <-watcher.Errors:

			case <-ctx.Done():
				return
			}

		}
	}()
	return nil
}
