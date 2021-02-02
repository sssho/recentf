package lib

import (
	"github.com/fsnotify/fsnotify"
)

func Observe(targetDir string, addedPath chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	err = watcher.Add(targetDir)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				addedPath <- event.Name
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			panic(err)
		}
	}
}
