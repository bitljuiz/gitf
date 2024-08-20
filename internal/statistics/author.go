package statistics

import (
	"fmt"
	"sort"
	"strings"
)

type Author struct {
	Lines   int
	Commits map[string]bool
	Files   map[string]bool
}

var (
	sortBy = "lines"
)

type AuthorPair struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

type SortedStatistics []AuthorPair

func Sort(orderBy string, stats map[string]*Author) SortedStatistics {
	ss := make(SortedStatistics, len(stats))

	if sortBy != orderBy {
		sortBy = orderBy
	}

	i := 0
	for k, v := range stats {
		ss[i] = AuthorPair{
			Name:    k,
			Lines:   v.Lines,
			Commits: len(v.Commits),
			Files:   len(v.Files),
		}
		i++
	}
	sort.Sort(ss)
	return ss
}

func (ss SortedStatistics) Len() int {
	return len(ss)
}
func (ss SortedStatistics) Less(i, j int) bool {
	linesEq := ss[i].Lines == ss[j].Lines
	linesCmp := ss[i].Lines > ss[j].Lines
	commitsEq := ss[i].Commits == ss[j].Commits
	commitsCmp := ss[i].Commits > ss[j].Commits
	filesCmp := ss[i].Files > ss[j].Files
	filesEq := ss[i].Files == ss[j].Files
	stringsCmp := strings.Compare(ss[i].Name, ss[j].Name)
	switch sortBy {
	case "lines":
		return linesCmp || (linesEq && (commitsCmp || (commitsEq && (filesCmp || (filesEq && (stringsCmp == -1))))))
	case "commits":
		return commitsCmp || (commitsEq && (linesCmp || (linesEq && (filesCmp || (filesEq && (stringsCmp == -1))))))
	case "files":
		return filesCmp || (filesEq && (linesCmp || (linesEq && (commitsCmp || (commitsEq && (stringsCmp == -1))))))
	default:
		panic(fmt.Sprintf("cannot sort by %s", sortBy))
	}
}
func (ss SortedStatistics) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}
