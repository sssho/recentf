package main

import (
	"time"

	"github.com/sssho/recentf/lib"
)

func main() {
	im := lib.NewIM()
	go lib.StartServer(im)

	recent, _ := lib.RecentDir()
	addedFile := make(chan string)
	go lib.Observe(recent, addedFile)

	for {
		select {
		case result := <-addedFile:
			path, _, err := lib.ResolveShortcut(result)
			if err != nil {
				continue
			}
			_, err = im.Index(path)
			if err != nil {
				im.Append(path)
			}
			im.Update(path, time.Now())
			im.Sort()
		}
	}
}
