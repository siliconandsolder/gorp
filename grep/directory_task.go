package grep

import (
	"fmt"
	"io/ioutil"
)

type DirectoryTask struct {
	directory string
	settings  GrepSettings
	comms     *GrepComms
}

func (dt *DirectoryTask) run() {
	files, err := ioutil.ReadDir(dt.directory)
	if err != nil {
		dt.comms.conMtx.Lock()
		fmt.Printf("Could not parse directory: %s\n", dt.directory)
		dt.comms.conMtx.Unlock()
	}

	for _, file := range files {
		if file.IsDir() {

			dt.comms.taskCn <- &FileTask{
				directory: dt.directory,
				settings:  dt.settings,
				comms:     dt.comms,
			}

			dt.comms.taskCn <- &DirectoryTask{
				directory: dt.directory,
				settings:  dt.settings,
				comms:     dt.comms,
			}
		}
	}
}
