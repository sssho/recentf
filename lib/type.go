package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/google/uuid"
)

type Filetype string

const (
	Word    Filetype = "word"
	Excel            = "excel"
	PPT              = "powerpoint"
	Visio            = "visio"
	Text             = "text"
	Image            = "image"
	Folder           = "folder"
	Unknown          = "unknown"
)

type MyFileInfo struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Dir          string    `json:"dir"`
	IsDir        bool      `json:"isDir"`
	Type         Filetype  `json:"fileType"`
	LastAccessed time.Time `json:"lastAccessed"`
}

// In-Memory fileinfo data
type IM struct {
	Data []MyFileInfo
}

func NewIM() *IM {
	im := &IM{}
	var data []MyFileInfo
	data, err := LoadAll()
	if err != nil {
		data = nil
	}
	im.Data = data
	im.Sort()
	im.Clean()
	return im
}

// Sort by LastAccessed
func (im *IM) Sort() {
	sort.Slice(im.Data, func(i, j int) bool {
		return im.Data[i].LastAccessed.After(im.Data[j].LastAccessed)
	})
}

// Exclude non existense files
func (im *IM) Clean() {
	cleaned := make([]MyFileInfo, 0, len(im.Data))
	for _, item := range im.Data {
		_, err := os.Stat(filepath.Join(item.Dir, item.Name))
		if err == os.ErrNotExist {
			continue
		} else if err != nil {
			continue
		}
		cleaned = append(cleaned, item)
	}
	im.Data = cleaned
}

// Append
func (im *IM) Append(path string) {
	info := NewMyFileInfo(path)
	if info == nil {
		return
	}
	im.Data = append(im.Data, *info)
}

// Update lastAccessed by path
func (im *IM) Update(path string, lastAccessed time.Time) {
	index, err := im.Index(path)
	if err != nil {
		return
	}
	im.Data[index].LastAccessed = lastAccessed
}

// Return item index
func (im *IM) Index(path string) (int, error) {
	for i, item := range im.Data {
		if path == filepath.Join(item.Dir, item.Name) {
			return i, nil
		}
	}
	return 0, errors.New("item not found")
}

// Debug print
func (im *IM) P() {
	for _, i := range im.Data {
		fmt.Println(i)
	}
}

func GetFileType(path string) Filetype {
	ext := filepath.Ext(path)
	switch ext {
	case ".doc", ".docx":
		return Word
	case ".xlsx", ".xlsm":
		return Excel
	case ".ppt":
		return PPT
	case ".vsd", ".vsdx":
		return Visio
	case ".txt":
		return Text
	case ".png", ".jpg", ".jpeg", ".bmp":
		return Image
	}
	return Unknown
}

func NewMyFileInfo(path string) *MyFileInfo {
	info, err := os.Stat(path)
	if err == os.ErrNotExist {
		return nil
	} else if err != nil {
		return nil
	}

	var t Filetype
	if info.IsDir() {
		t = Folder
	} else {
		t = GetFileType(path)
	}
	id := uuid.New()
	return &MyFileInfo{id.String(), info.Name(), filepath.Dir(path), info.IsDir(), t, info.ModTime()}
}

func LoadAll() ([]MyFileInfo, error) {
	// history may be nil
	historyItems, _ := LoadHistory()
	recentPaths, err := LoadRecent()
	if err != nil {
		return nil, err
	}

	// pick up only in recentPaths
	onlyInRecent := make([]string, 0, len(recentPaths))
	for _, p := range recentPaths {
		found := false
		for _, item := range historyItems {
			if p == filepath.Join(item.Dir, item.Name) {
				found = true
				break
			}
		}
		if !found {
			onlyInRecent = append(onlyInRecent, p)
		}
	}
	recentItems := make([]MyFileInfo, 0, len(onlyInRecent))
	for _, p := range onlyInRecent {
		info := NewMyFileInfo(p)
		if info == nil {
			continue
		}
		recentItems = append(recentItems, *info)
	}

	return append(historyItems, recentItems...), nil
}
