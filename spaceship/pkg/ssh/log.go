package ssh

import (
	"bufio"
	"context"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

// LatestLog returns the last n lines from the latest log file.
func (s *SSH) LatestLog(n int) (string, error) {
	logFiles, err := s.getLogFiles()
	if err != nil {
		return "", errors.Wrap(err, "error fetching log files")
	}
	if len(logFiles) == 0 {
		return "", errors.Wrap(err, "no log files found")
	}
	// Sort log files by modification time.
	sort.Sort(logFiles)

	// Get the latest log file
	latestLogFile := logFiles[len(logFiles)-1]
	lines, err := s.readLastNLines(latestLogFile.name, n)
	if err != nil {
		return "", errors.Wrap(err, "error reading last log line")
	}
	return strings.Join(lines, "\n"), nil
}

// readLastNLines reads the last n lines from the specified file.
func (s *SSH) readLastNLines(filePath string, n int) ([]string, error) {
	file, err := s.sftpClient.OpenFile(filePath, os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Use a buffer to read the file
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)

	// Read lines into a buffer
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) == n {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// FollowLog follows the latest log file and sends new lines to the provided channel in real-time.
func (s *SSH) FollowLog(ctx context.Context, ch chan<- string) error {
	logFiles, err := s.getLogFiles()
	if err != nil {
		return errors.Wrap(err, "error fetching log files")
	}

	if len(logFiles) == 0 {
		return errors.Wrap(err, "no log files found")
	}
	// Sort log files by modification time.
	sort.Sort(logFiles)

	// Get the latest log file
	latestLogFile := logFiles[len(logFiles)-1]
	file, err := s.sftpClient.OpenFile(latestLogFile.name, os.O_RDONLY)
	if err != nil {
		return err
	}
	defer file.Close()

	// Seek to the end of the file
	file.Seek(0, io.SeekEnd)

	reader := bufio.NewReader(file)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(1 * time.Second)
					continue
				}
				return err
			}
			ch <- line
		}
	}
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
