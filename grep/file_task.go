package grep

import (
	"fmt"
	"io/ioutil"
)

type FileTask struct {
	directory string
	settings  GrepSettings
	comms     *GrepComms
}

func (ft *FileTask) run() {
	files, err := ioutil.ReadDir(ft.directory)
	if err != nil {
		ft.comms.conMtx.Lock()
		fmt.Printf("Could not parse directory: %s\n", ft.directory)
		ft.comms.conMtx.Unlock()
	}

	for _, file := range files {
		if !file.IsDir() {
			extMatches := ft.settings.extension.FindString(file.Name())
			if len(extMatches) > 0 {
				ft.comms.taskCn <- &GrepTask{
					results:  make([]GrepMatch, 0),
					settings: ft.settings,
					comms:    ft.comms,
					fileName: file.Name(),
				}
			}
		}
	}
}