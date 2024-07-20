package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

const (
	logExtension = ".log"
)

// Log represents a log file with its name and modification time.
type Log struct {
	Name string
	Time time.Time
}

// Logs implements sort.Interface based on the Time field.
type Logs []Log

func (a Logs) Len() int           { return len(a) }
func (a Logs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Logs) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }

func (s *SSH) LatestLog() ([]byte, error) {
	logFiles, err := getLogFiles(s.Log())
	if err != nil {
		return nil, errors.Wrap(err, "error fetching log files")
	}
	if len(logFiles) == 0 {
		return nil, errors.Wrap(err, "no log files found")
	}
	// Sort log files by modification time.
	sort.Sort(logFiles)

	// Get the latest log file
	latestLogFile := logFiles[len(logFiles)-1]
	fmt.Printf("Latest log file: %s\n", latestLogFile.Name)

	// Read the file and return its contents as bytes.
	data, err := os.ReadFile(latestLogFile.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading the latest log file %s", latestLogFile.Name)
	}

	return data, nil
}

// getLogFiles fetches all log files from the specified directory.
func getLogFiles(dir string) (Logs, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	logFiles := make([]Log, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Assuming log files have a .log extension
		if filepath.Ext(file.Name()) == logExtension {
			info, err := file.Info()
			if err != nil {
				return nil, err
			}
			logFiles = append(logFiles, Log{
				Name: filepath.Join(dir, file.Name()),
				Time: info.ModTime(),
			})
		}
	}
	return logFiles, nil
}
