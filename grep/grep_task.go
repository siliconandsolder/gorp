package grep

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type GrepMatch struct {
	path    string
	lineNum uint32
	matches uint32
	content string
}

func (match GrepMatch) printInfo() {
	fmt.Printf("\tLine [%d] %s\n\tMatches found: %d\n", match.lineNum, match.content, match.matches)
}

func (match GrepMatch) printVerbose() {
	fmt.Printf("Path: %s\n\tLine [%d] %s\n", match.path, match.lineNum, match.content)
}

type GrepTask struct {
	results  []GrepMatch
	settings GrepSettings
	comms    *GrepComms
	fileName string
}

func (gt *GrepTask) run() {

	if gt.settings.verboseMode {
		fmt.Printf("Scanning %s\n", gt.fileName)
	}

	file, err := os.Open(gt.fileName)

	if err != nil {
		fmt.Printf("Could not open %s\n", gt.fileName)
		return
	}

	defer file.Close()

	var lineNum uint32 = 1
	var totalMatches uint32 = 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		matches := gt.settings.expression.FindAllString(lineText, -1)
		var lineMatches uint32 = uint32(len(matches))
		totalMatches += lineMatches

		if lineMatches > 0 {
			if gt.settings.verboseMode {
				GrepMatch{gt.fileName, lineNum, lineMatches, lineText}.printVerbose()
			} else {
				gt.results = append(gt.results, GrepMatch{gt.fileName, lineNum, lineMatches, lineText})
			}
		}

		lineNum += 1
	}

	if len(gt.results) > 0 {
		sort.Slice(gt.results, func(i, j int) bool {
			if gt.results[i].path == gt.results[j].path {
				return gt.results[i].lineNum < gt.results[j].lineNum
			}
			return gt.results[i].path < gt.results[j].path
		})

		gt.comms.resultsCn <- gt.results
	}

	gt.comms.statCn <- totalMatches

	//time.Sleep(5 * time.Second)
}
