package grep

import (
	"container/list"
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

	matches    list.List
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
		matches:    *list.New(),
		matchCount: 0,
		fileCount:  0,
	}
}

func (gm *GrepManager) FindMatches() {
	var runningTasks uint32 = 0
	var wg sync.WaitGroup
	for i := 0; i < DEFAULT_GO_ROUTINES; i++ {

		go resultsWorker(gm.comms.resultsCn, &gm.matches, &gm.listMtx)

		go statsWorker(gm.comms.statCn, &gm.matchCount, &gm.fileCount)

		wg.Add(1)
		go taskWorker(gm.comms.taskCn, &runningTasks, &wg)
	}
	wg.Wait()
	close(gm.comms.resultsCn)
	close(gm.comms.statCn)
	close(gm.comms.taskCn)
}

func resultsWorker(results <-chan []GrepMatch, matches *list.List, mtx *sync.Mutex) {
	for result := range results {
		mtx.Lock()
		matches.PushBack(result)
		mtx.Unlock()
	}
}

func statsWorker(stats <-chan uint32, numMatches *uint32, numFiles *uint32) {
	for stat := range stats {
		if *numMatches > 0 {
			atomic.AddUint32(numMatches, stat)
			atomic.AddUint32(numFiles, 1)
		}
	}
}

func taskWorker(tasks <-chan Task, numRunningTasks *uint32, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tasks {
		atomic.AddUint32(numRunningTasks, 1)
		t.run()
		atomic.AddUint32(numRunningTasks, ^uint32(0))

		if *numRunningTasks == 0 && len(tasks) == 0 {
			break
		}
	}
}

func (gm GrepManager) printMatches() {
	if !gm.settings.verboseMode && gm.matches.Len() > 0 {
		for item := gm.matches.Front(); item != nil; item.Next() {
			item.Value.(*GrepMatch).printInfo()
		}
	}

	fmt.Printf("Files with matches: %d\n", gm.fileCount)
	fmt.Printf("Total matches: %d\n", gm.matchCount)
}
