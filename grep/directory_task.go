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
		fmt.Printf("Could not parse directory: %s\n", dt.directory)
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
