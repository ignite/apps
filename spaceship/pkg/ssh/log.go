package ssh

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

const (
	logExtension = ".log"
)

// log represents a log file with its name and modification time.
type log struct {
	name string
	time time.Time
}

// logs implements sort.Interface based on the Time field.
type logs []log

func (a logs) Len() int           { return len(a) }
func (a logs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a logs) Less(i, j int) bool { return a[i].time.Before(a[j].time) }

func (s *SSH) LatestLog() ([]byte, error) {
	logFiles, err := s.getLogFiles()
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
	return s.readFileToBytes(latestLogFile.name)
}

// readFileToBytes reads the contents of a file and returns them as a byte slice.
func (s *SSH) readFileToBytes(filePath string) ([]byte, error) {
	file, err := s.sftpClient.OpenFile(filePath, os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// getLogFiles fetches all log files from the specified directory.
func (s *SSH) getLogFiles() (logs, error) {
	dir := s.Log()

	files, err := s.sftpClient.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	logFiles := make([]log, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Assuming log files have a .log extension
		if filepath.Ext(file.Name()) == logExtension {
			logFiles = append(logFiles, log{
				name: filepath.Join(dir, file.Name()),
				time: file.ModTime(),
			})
		}
	}
	return logFiles, nil
}
