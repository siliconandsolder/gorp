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
	fmt.Printf("\tLine [%d] %s\n\tMatches found: %d", match.lineNum, match.content, match.matches)
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
		gt.comms.conMtx.Lock()
		fmt.Printf("Scanning %s\n", gt.fileName)
		gt.comms.conMtx.Unlock()
	}

	file, err := os.Open(gt.fileName)

	if err != nil {
		gt.comms.conMtx.Lock()
		fmt.Printf("Could not open %s\n", gt.fileName)
		gt.comms.conMtx.Unlock()
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
				gt.comms.conMtx.Lock()
				GrepMatch{gt.fileName, lineNum, lineMatches, lineText}.printVerbose()
				gt.comms.conMtx.Unlock()
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
}
