package git

import (
	"fmt"
	"gitf/internal/cli"
	"gitf/internal/statistics"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const goroutinesCount = 10

func GetRepositoryFiles(dir string, commit string) ([]string, error) {
	lstreeCmd := exec.Command("git", "ls-tree", "--name-only", "-r", commit)
	lstreeCmd.Dir = dir
	var blobsOut strings.Builder

	lstreeCmd.Stdout = &blobsOut
	if err := lstreeCmd.Run(); err != nil {
		return []string{}, err
	}
	files := strings.FieldsFunc(blobsOut.String(), func(r rune) bool {
		return r == '\n'
	})
	return files, nil
}

func getLog(repository, commit, file string) (string, string, error) {
	logCmd := exec.Command("git", "log", "--pretty=format:%H %an", commit, "--", file)
	logCmd.Dir = repository

	var rawLogInfo strings.Builder
	logCmd.Stdout = &rawLogInfo
	if err := logCmd.Run(); err != nil {
		return "", "", err
	}

	logInfo := strings.FieldsFunc(rawLogInfo.String(), func(r rune) bool {
		return r == '\n'
	})
	logInfoArray := strings.SplitN(logInfo[0], " ", 2)
	return logInfoArray[0], logInfoArray[1], nil
}

func getFileStatistics(repository string, file string, revision string, useCommitter bool, stats map[string]*statistics.Author, mu *sync.Mutex) error {
	blameCmd := exec.Command("git", "blame", "--incremental", revision, file)
	blameCmd.Dir = repository

	var fileInside strings.Builder
	blameCmd.Stdout = &fileInside

	if err := blameCmd.Run(); err != nil {
		return err
	}
	blameInfo := strings.FieldsFunc(fileInside.String(), func(r rune) bool {
		return r == '\n'
	})

	mu.Lock()
	defer mu.Unlock()

	if len(blameInfo) == 0 {
		hash, author, err := getLog(repository, revision, file)
		if err != nil {
			return err
		}
		if stat, ok := stats[author]; !ok {
			stats[author] = &statistics.Author{
				Lines:   0,
				Commits: make(map[string]bool),
				Files:   make(map[string]bool),
			}
			stats[author].Commits[hash] = true
			stats[author].Files[file] = true
		} else {
			stat.Files[file] = true
			stat.Commits[hash] = true
		}
	} else {
		checkCommit := true

		var k int
		var authorPrefix string

		if useCommitter {
			authorPrefix = "committer "
			k = 5
		} else {
			authorPrefix = "author "
			k = 1
		}

		author := ""
		currentCommit := ""

		for i := range blameInfo {
			if checkCommit {
				commitInfo := strings.Split(blameInfo[i], " ")
				currentCommit = commitInfo[0]
				linesCount, err := strconv.Atoi(commitInfo[3])

				if err != nil {
					return err
				}

				author = strings.TrimPrefix(blameInfo[i+k], authorPrefix)

				if stat, ok := stats[author]; !ok {
					stats[author] = &statistics.Author{
						Lines:   linesCount,
						Commits: make(map[string]bool),
						Files:   make(map[string]bool),
					}
					stats[author].Commits[currentCommit] = true
					stats[author].Files[file] = true
				} else {
					stat.Lines += linesCount
					if _, ok = stat.Commits[currentCommit]; !ok {
						stat.Commits[currentCommit] = true
					}
					if _, ok = stat.Files[file]; !ok {
						stat.Files[file] = true
					}
				}
				checkCommit = false
			} else {
				if strings.HasPrefix(blameInfo[i], "filename") {
					if i < len(blameInfo)-1 {
						if strings.HasPrefix(blameInfo[i+1], currentCommit) {
							commitInfo := strings.Split(blameInfo[i+1], " ")
							linesCount, err := strconv.Atoi(commitInfo[3])

							if err != nil {
								return err
							}
							stat := stats[author]
							stat.Lines += linesCount
						} else {
							checkCommit = true
						}
					}
				}
			}
		}
	}
	return nil
}

func GetStatistics(repository, revision string, useCommitter bool, repositoryFiles []string) (map[string]*statistics.Author, []error) {
	stats := make(map[string]*statistics.Author)

	var mu sync.Mutex
	var wg sync.WaitGroup

	guard := make(chan struct{}, goroutinesCount)

	errChan := make(chan error, len(repositoryFiles))

	filesCnt := len(repositoryFiles)
	filesProcessed := 1
	var mu2 sync.Mutex

	for _, filePath := range repositoryFiles {
		wg.Add(1)
		guard <- struct{}{}

		go func(file string) {
			defer wg.Done()
			defer func() { <-guard }()
			if err := getFileStatistics(repository, file, revision, useCommitter, stats, &mu); err != nil {
				errChan <- fmt.Errorf("cannot process file %s", file)
			} else {
				mu2.Lock()
				defer mu2.Unlock()
				fmt.Printf("Processed %s. Progress: %d/%d\n", file, filesProcessed, filesCnt)
				filesProcessed++
			}
		}(filePath)
	}

	wg.Wait()
	close(errChan)

	errorList := make([]error, 0, len(repositoryFiles))

	for err := range errChan {
		if err != nil {
			errorList = append(errorList, err)
		}
	}

	cli.ClearConsole()

	return stats, errorList
}
