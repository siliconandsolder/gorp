package grep

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

const DEFAULT_GO_ROUTINES = 50

type GrepSettings struct {
	expression  *regexp.Regexp
	extension   *regexp.Regexp
	verboseMode bool
}

type GrepComms struct {
	statCn    chan uint32
	taskCn    chan Task
	resultsCn chan []GrepMatch
	conMtx    sync.Mutex
}

type GrepManager struct {
	settings GrepSettings
	comms    GrepComms

	rootDir    string
	matches    []GrepMatch
	listMtx    sync.Mutex
	matchCount uint32
	fileCount  uint32
}

func NewGrepManager(verboseMode bool, rootDir string, expression string, extensions []string) GrepManager {
	var extString string = "\\.("

	for _, ext := range extensions {
		extString += (strings.ReplaceAll(ext, ".", "") + "|")
	}

	extString = (extString[:len(extString)-1] + ")$")

	return GrepManager{
		settings: GrepSettings{
			expression:  regexp.MustCompile(expression),
			extension:   regexp.MustCompile(extString),
			verboseMode: verboseMode,
		},
		comms: GrepComms{
			statCn:    make(chan uint32, DEFAULT_GO_ROUTINES),
			taskCn:    make(chan Task, DEFAULT_GO_ROUTINES),
			resultsCn: make(chan []GrepMatch, DEFAULT_GO_ROUTINES),
		},
		rootDir:    rootDir,
		matchCount: 0,
		fileCount:  0,
	}
}

func (gm *GrepManager) FindMatches() {
	var wgResult sync.WaitGroup
	wgResult.Add(1)

	var wgStat sync.WaitGroup
	wgStat.Add(1)

	var wgTask sync.WaitGroup
	wgTask.Add(1)

	gm.comms.taskCn <- &DirectoryTask{
		directory: gm.rootDir,
		settings:  gm.settings,
		comms:     &gm.comms,
	}

	gm.comms.taskCn <- &FileTask{
		directory: gm.rootDir,
		settings:  gm.settings,
		comms:     &gm.comms,
	}

	go resultsWorker(gm.comms.resultsCn, &gm.matches, &wgResult)
	go statsWorker(gm.comms.statCn, &gm.matchCount, &gm.fileCount, &wgStat)
	go taskWorker(gm.comms.taskCn, &wgTask)

	wgTask.Wait()
	close(gm.comms.taskCn)
	close(gm.comms.resultsCn)
	wgResult.Wait()
	close(gm.comms.statCn)
	wgStat.Wait()
	gm.printMatches()
}

func resultsWorker(results <-chan []GrepMatch, matches *[]GrepMatch, wg *sync.WaitGroup) {
	defer wg.Done()
	for result := range results {
		*matches = append(*matches, result...)
	}
}

func statsWorker(stats <-chan uint32, numMatches *uint32, numFiles *uint32, wg *sync.WaitGroup) {
	defer wg.Done()
	for stat := range stats {
		if stat > 0 {
			atomic.AddUint32(numMatches, stat)
			atomic.AddUint32(numFiles, 1)
		}
	}
}

func taskWorker(tasks <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
out:
	for {
		select {
		case t := <-tasks:
			t.run()
		default:
			break out
		}
	}
}

func (gm *GrepManager) printMatches() {
	if !gm.settings.verboseMode && len(gm.matches) > 0 {
		path := gm.matches[0].path
		fmt.Println(path)
		for _, match := range gm.matches {
			if path != match.path {
				fmt.Println(match.path)
				path = match.path
			}
			match.printInfo()
		}
	}

	fmt.Printf("Files with matches: %d\n", gm.fileCount)
	fmt.Printf("Total matches: %d\n", gm.matchCount)
}
