package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func RecentDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`%s\AppData\Roaming\Microsoft\Windows\Recent`, home), nil
}

func LoadRecent() ([]string, error) {
	recentDir, err := RecentDir()
	if err != nil {
		return nil, err
	}
	dir, err := os.Open(recentDir)
	if err != nil {
		return nil, err
	}
	fileNames, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	wshell, err := NewWscriptShell()
	if err != nil {
		log.Fatal(err)
	}
	defer wshell.Close()

	paths := make([]string, 0, len(fileNames))
	for _, f := range fileNames {
		tpath, _, err := wshell.ShortcutInfo(filepath.Join(recentDir, f))
		if err != nil || tpath == "" {
			continue
		}
		paths = append(paths, tpath)
	}
	return paths, nil
}
