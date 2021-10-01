package grep

import (
	"fmt"
	"os"
	"path/filepath"
)

type DirectoryTask struct {
	directory string
	settings  GrepSettings
	comms     *GrepComms
}

func (dt *DirectoryTask) run() {
	files, err := os.ReadDir(dt.directory)
	if err != nil {
		dt.comms.conMtx.Lock()
		fmt.Printf("Could not parse directory: %s\n", dt.directory)
		dt.comms.conMtx.Unlock()
	}

	for _, file := range files {
		if file.IsDir() {

			dt.comms.taskCn <- &FileTask{
				directory: filepath.Join(dt.directory, file.Name()),
				settings:  dt.settings,
				comms:     dt.comms,
			}

			dt.comms.taskCn <- &DirectoryTask{
				directory: filepath.Join(dt.directory, file.Name()),
				settings:  dt.settings,
				comms:     dt.comms,
			}
		}
	}
}
