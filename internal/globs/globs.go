package globs

import (
	"path/filepath"
)

func Exclude(patterns, files []string) []string {
	if len(patterns) == 0 {
		return files
	}

	files_ := make([]string, 0, len(files))

	match := func(file string) bool {
		for _, pattern := range patterns {
			if match, _ := filepath.Match(pattern, file); match {
				return false
			}
		}
		return true
	}

	for _, file := range files {
		if ok := match(file); ok {
			files_ = append(files_, file)
		}
	}
	return files_
}

func RestrictTo(patterns, files []string) []string {
	if len(patterns) == 0 {
		return files
	}

	files_ := make([]string, 0, len(files))

	restrict := func(file string) bool {
		for _, pattern := range patterns {
			if match, _ := filepath.Match(pattern, file); match {
				return true
			}
		}
		return false
	}

	for _, file := range files {
		if match := restrict(file); match {
			files_ = append(files_, file)
		}
	}
	return files_
}
