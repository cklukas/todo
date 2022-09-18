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
	mode := "main"

	if len(args) == 1 {
		todoDir = path.Join(todoDir, "mode", args[0])
		mode = args[0]
	}

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
	lanes := NewLanes(content, app, mode)

	for idx := range lanes.lanes {
		if lanes.active == idx {
			lanes.lanes[idx].SetSelectedBackgroundColor(tcell.ColorLightBlue)
			lanes.lanes[idx].SetSelectedTextColor(tcell.ColorBlack)
		} else {
			lanes.lanes[idx].SetSelectedStyle(tcell.StyleDefault)
		}
	}

	bAbout := tview.NewButton("[brown::-]F1 [black::-]About")
	bAbout.SetBackgroundColor(tcell.ColorLightGray)
	bAbout.SetSelectedFunc(lanes.CmdAbout)

	bAddToDo := tview.NewButton("[red::-]F2 [black::-]Add Task")
	bAddToDo.SetBackgroundColor(tcell.ColorLightGray)
	bAddToDo.SetSelectedFunc(lanes.CmdAddTask)

	bEditToDo := tview.NewButton("[red::-]F3 [black::-]Edit")
	bEditToDo.SetBackgroundColor(tcell.ColorLightGray)
	bEditToDo.SetSelectedFunc(lanes.CmdEditTask)

	bNoteToDo := tview.NewButton("[red::-]F4 [black::-]Note")
	bNoteToDo.SetBackgroundColor(tcell.ColorLightGray)
	bNoteToDo.SetSelectedFunc(lanes.CmdEditNote)

	bArchiveToDo := tview.NewButton("[red::-]F5 [black::-]Archive")
	bArchiveToDo.SetBackgroundColor(tcell.ColorLightGray)
	bArchiveToDo.SetSelectedFunc(lanes.CmdArchiveNote)

	bSelectToDo := tview.NewButton("[red::-]F6 [black::-]Select")
	bSelectToDo.SetBackgroundColor(tcell.ColorLightGray)
	bSelectToDo.SetSelectedFunc(lanes.CmdSelectNote)

	// bAddColumn := tview.NewButton("[darkblue::-]F8 [black::-]Add Column")
	// bAddColumn.SetBackgroundColor(tcell.ColorLightGray)

	// bDeleteColumn := tview.NewButton("[darkblue::-]F8 [black::-]Delete")
	// bDeleteColumn.SetBackgroundColor(tcell.ColorLightGray)

	// bRenameColumn := tview.NewButton("[darkblue::-]F9 [black::-]Rename")
	// bRenameColumn.SetBackgroundColor(tcell.ColorLightGray)

	bExit := tview.NewButton("[brown::-]F10 [black::-]Exit")
	bExit.SetBackgroundColor(tcell.ColorLightGray)
	bExit.SetSelectedFunc(lanes.CmdExit)

	bMode := tview.NewButton("[blue::-]" + mode) // [blue::-]F12
	bMode.SetBackgroundColor(tcell.ColorLightGray)

	bMoveHelp := tview.NewButton("")
	bMoveHelp.SetBackgroundColor(tcell.ColorLightGray)
	lanes.bMoveHelp = bMoveHelp

	info := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(bAbout, 10, 1, false).
		AddItem(bAddToDo, 13, 1, false).
		AddItem(bEditToDo, 9, 1, false).
		AddItem(bNoteToDo, 9, 1, false).
		AddItem(bArchiveToDo, 13, 1, false).
		AddItem(bSelectToDo, 10, 1, false).
		// AddItem(bAddColumn, 15, 1, false).
		// AddItem(bDeleteColumn, 10, 1, false).
		// AddItem(bRenameColumn, 10, 1, false).
		AddItem(bExit, 10, 1, false).
		AddItem(bMode, 2+len(mode), 1, false).
		AddItem(bMoveHelp, 38, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(lanes.GetUi(), 0, 1, true).
		AddItem(info, 1, 1, false)
	app.SetRoot(layout, true).EnableMouse(true)

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
	Short:        "TODO App",
	Long:         "ToDo Main View - optional program argument: mode (e.g. 'private' or 'work')",
	Version:      AppVersion,
	SilenceUsage: true,
	RunE:         main,
	Args:         cobra.RangeArgs(0, 1),
}

func init() {

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
