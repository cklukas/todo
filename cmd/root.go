package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/flytam/filenamify"
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
	todoDir := ".todo"
	mode := "main"

	todoDirModes := path.Join(todoDir, "mode")

	if len(args) == 1 && args[0] != "main" {
		saveName, err := filenamify.FilenamifyV2(args[0], func(options *filenamify.Options) {
			options.Replacement = "_"
		})
		if err != nil {
			log.Fatal(err)
		}

		todoDir = path.Join(todoDirModes, saveName)
		mode = saveName
	}

	saveName := mode

	var err error

	nextModeLaneFocus := 0
	for {
		var nextMode string
		nextMode, nextModeLaneFocus, err = launchGui(todoDir, todoDirModes, saveName, nextModeLaneFocus)
		if len(nextMode) == 0 {
			break
		}
		if nextMode == "main" {
			todoDir = ".todo"
			saveName = "main"
			continue
		}
		saveName, err = filenamify.FilenamifyV2(nextMode, func(options *filenamify.Options) {
			options.Replacement = "_"
		})
		if err != nil {
			log.Fatal(err)
		}
		todoDir = path.Join(todoDirModes, saveName)
	}

	return err
}

func getStatusBar(lanes *Lanes, mode string) *tview.Flex {
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

	bLanesCommands := tview.NewButton("[red::-]F7 [black::-]Lane")
	bLanesCommands.SetBackgroundColor(tcell.ColorLightGray)
	bLanesCommands.SetSelectedFunc(lanes.CmdLanesCmds)

	bExit := tview.NewButton("[brown::-]F10 [black::-]Exit")
	bExit.SetBackgroundColor(tcell.ColorLightGray)
	bExit.SetSelectedFunc(lanes.CmdExit)

	bMode := tview.NewButton("[blue::-]" + mode) // [blue::-]F12
	bMode.SetBackgroundColor(tcell.ColorLightGray)
	bMode.SetSelectedFunc(lanes.CmdSelectModeDialog)

	bMoveHelp := tview.NewButton("")
	bMoveHelp.SetBackgroundColor(tcell.ColorLightGray)
	lanes.bMoveHelp = bMoveHelp

	defaultStatusBarMenuItems := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(bAbout, 10, 1, false).
		AddItem(bAddToDo, 13, 1, false).
		AddItem(bEditToDo, 9, 1, false).
		AddItem(bNoteToDo, 9, 1, false).
		AddItem(bArchiveToDo, 13, 1, false).
		AddItem(bSelectToDo, 10, 1, false).
		AddItem(bLanesCommands, 9, 1, false).
		AddItem(bExit, 10, 1, false).
		AddItem(bMode, 2+len(mode), 1, false).
		AddItem(bMoveHelp, 38, 1, false)

	return defaultStatusBarMenuItems
}

func JsonWatcher(watcher *fsnotify.Watcher, content *ToDoContent, lanes *Lanes, app *tview.Application) {
	appLocked := false
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if filepath.Base(event.Name) == "todo.json" && (event.Has(fsnotify.Remove)) {
				if !appLocked {
					app.Lock()
					appLocked = true
				}
			}

			if filepath.Base(event.Name) == "todo.json" && (event.Has(fsnotify.Write) || event.Has(fsnotify.Create)) {
				if appLocked {
					app.Unlock()
					appLocked = false
				}

				err := content.Read()
				if err != nil {
					app.Stop()
					log.Fatal(err)
				}

				lanes.RedrawAllLanes()
				app.ForceDraw()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			app.Stop()
			log.Fatalf("JsonWatcher, error monitoring settings file and folders: %v", err)
		}
	}
}

func launchGui(todoDir, todoDirModes, mode string, nextModeLaneFocus int) (string, int, error) {
	usr, errU := user.Current()
	if errU != nil {
		log.Fatal(errU)
	}

	archiveDir, err := CreateDir(path.Join(usr.HomeDir, todoDir, "archive"))
	if err != nil {
		log.Fatal(err)
	}

	backupDir, err := CreateDir(path.Join(usr.HomeDir, todoDir, "backup"))
	if err != nil {
		log.Fatal(err)
	}

	fname := path.Join(usr.HomeDir, todoDir, "todo.json")

	content := new(ToDoContent)
	err = content.ReadFromFile(fname)
	if err != nil {
		content.InitializeNew()
	}

	content.SetFileName(fname, archiveDir, backupDir)
	err = content.Save()
	if err != nil {
		log.Fatal(fmt.Errorf("could not save todos in '%v': %w", fname, err))
	}

	app := tview.NewApplication()
	lanes := NewLanes(content, app, mode, path.Join(usr.HomeDir, todoDirModes))

	// lanes.active = nextModeLaneFocus
	// lanes.lastActive = nextModeLaneFocus

	for idx := range lanes.lanes {
		if lanes.active == idx {
			lanes.lanes[idx].SetSelectedBackgroundColor(tcell.ColorLightBlue)
			lanes.lanes[idx].SetSelectedTextColor(tcell.ColorBlack)
		} else {
			lanes.lanes[idx].SetSelectedStyle(tcell.StyleDefault)
		}
	}

	defaultStatusBarMenuItems := getStatusBar(lanes, mode)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(lanes.GetUi(), 0, 1, true).
		AddItem(defaultStatusBarMenuItems, 1, 1, false)
	app.SetRoot(layout, true).EnableMouse(true)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// monitor changes to todo.json in background
	defer watcher.Close()
	go JsonWatcher(watcher, content, lanes, app)

	// watch directory
	err = watcher.Add(filepath.Dir(fname))
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v\n", err)
	}

	return lanes.nextMode, lanes.nextLaneFocus, content.Save()
}

var rootCmd = &cobra.Command{
	Use:          "todo",
	Short:        "ToDo App",
	Long:         "ToDo Main View - optional program argument: mode (e.g. 'private' or 'work')",
	Version:      AppVersion,
	SilenceUsage: true,
	RunE:         main,
	Args:         cobra.RangeArgs(0, 1),
}

func init() {
	// empty
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
