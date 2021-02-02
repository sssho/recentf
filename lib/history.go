package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func HistoryFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`%s\.cache\todo\history.json`, home), nil
}

func LoadHistory() ([]MyFileInfo, error) {
	var cachedItems []MyFileInfo
	history, _ := HistoryFile()
	file, err := os.Open(history)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("history file read error: %s", err)
	}

	err = json.Unmarshal(b, &cachedItems)
	if err != nil {
		return nil, fmt.Errorf("history file json unmarshal error: %s", err)
	}

	return cachedItems, nil
}

func DumpHistory(data []MyFileInfo, outfile string) error {
	jsonMarshalled, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	file, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, bytes.NewReader(jsonMarshalled))
	if err != nil {
		return err
	}
	return nil
}
