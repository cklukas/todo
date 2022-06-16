package cmd

import (
	"errors"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func CreateDir(path string) (string, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return path, err
		}
	}

	return path, nil
}

// AppVersion contains the version information and is set from build.sh
var AppVersion string = ""

func main(cmd *cobra.Command, args []string) error {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	todoDir := ".todo"
	archiveFolder := path.Join(usr.HomeDir, todoDir, "archive")
	archiveDir, err := CreateDir(archiveFolder)
	if err != nil {
		log.Fatal(err)
	}

	backupFolder := path.Join(usr.HomeDir, todoDir, "backup")
	backupDir, err := CreateDir(backupFolder)
	if err != nil {
		log.Fatal(err)
	}

	fname := path.Join(usr.HomeDir, todoDir, "todo.json")
	var content *Content
	f, err := os.OpenFile(fname, os.O_RDONLY, os.FileMode(int(0600)))
	if err == nil {
		content = NewContentIo(f)
		f.Close()
	}

	if content == nil {
		content = NewContentDefault()
	}

	content.SetFileName(fname, archiveDir, backupDir)
	content.Save()

	app := tview.NewApplication()
	// app.EnableMouse(true) // modal dialogs are not handled well
	lanes := NewLanes(content, app)
	for idx := range lanes.lanes {
		if lanes.active == idx {
			lanes.lanes[idx].SetSelectedBackgroundColor(tcell.ColorWhite)
		} else {
			lanes.lanes[idx].SetSelectedBackgroundColor(tcell.ColorGray)
		}
	}

	app.SetRoot(lanes.GetUi(), true)

	// watch todo.json for changes
	if _, err := os.Stat(fname); err == nil {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					// log.Println("event:", event)
					if event.Op&fsnotify.Write == fsnotify.Write {
						content.Read()
						lanes.RedrawAllLanes()
						app.ForceDraw()
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}()

		err = watcher.Add(fname)
		if err != nil {
			log.Fatal(err)
		}
	}
	// end of watch todo.json

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v\n", err)
	}

	content.Save()
	return nil
}

var rootCmd = &cobra.Command{
	Use:          "todo",
	Short:        "Kanban board",
	Long:         "kanban todo list",
	Version:      AppVersion,
	SilenceUsage: true,
	RunE:         main,
}

func init() {
	// empty
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
