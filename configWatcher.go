package main

import (
	"log"

	"github.com/howeyc/fsnotify"
)

func (a *application) WatchForConfigChanges(configPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Problem watching filesystem for changes, mockingjay wont bother", err)
	} else {

		go func() {
			for {
				select {
				case <-watcher.Event:
					a.updateServer()
				case err := <-watcher.Error:
					log.Println("error:", err)
				}
			}
		}()

		err = watcher.Watch(configPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
