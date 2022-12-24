# TODO Command Line App

[![Go Report](https://goreportcard.com/badge/github.com/cklukas/todo)](https://goreportcard.com/report/github.com/cklukas/todo)

A simple Kanban board for your terminal.

* Stores data in a simple JSON document in `$HOME/.todo/todo.json`
* Makes a daily backup of the data in `$HOME/.todo/backup/` (at first start on a particular day)
* Contains a function to archive an todo item in `$HOME/.todo/archive`
* If a non-default mode is used (see below), the files and folders for that mode (`todo.json`, `backup`, `archive`) are saved under `$home/.todo/mode/[mode]`

* Allows input of topic and second description line
* Provides function to view/edit a longer note for each item in vim (or other editor, as defined by the `EDITOR` environment variable)
* All changes are immediately saved (no save command)
* The application can be started multiple times, modifications performed in one instance are detected in other instances (through monitoring changes to the active `todo.json` file)
* Use [red], [blue] etc. to colorize your item text
* New: Hotkeys F1..F10 are shown in status bar

# Documentation

## Installation / start of the program

After download, you may launch the application from the command line. Depending on your operating systems, additional steps may be needed:

## Linux

On Linux, you need to mark the downloaded program as executable (`chmod +x todo`). To simplify the start, move the program to a directory, which is contained in the search path (e.g., `sudo mv todo /usr/bin/todo`).

## First start on Mac OS (Ventura)

On Mac OS, for security reasons, the direct start of downloaded applications is blocked. 

1. First, ensure that you are able to start apps which are not downloaded from App Store. Open the `Systems Settings`, select `Privacy & Security`. On the right scroll down to `Security`and select `Allow applications downloaded from: App Store and identified developers`.
2. Move the downloaded executable (`todo_mac_arm64`) to the desired target folder and renamed it to `todo`:
    - `mkdir -p ~/bin; mv ~/Downloads/todo_mac_arm64 ~/bin/todo`)
3. Mark the program as executable: `chmod +x ~/bin/todo`
4. Open the '~/bin' folder in the finder (GUI), control-click the app icon, then choose Open from the shortcut menu. A dialog will be shown, informing you that the app is not signed, choose `Open` to start the program, you may then close the app by clicking the X window button, or pressing 'q' and then 'Enter'.

Once you completed the above steps, the program can be directly opened from the terminal, without additional steps.

## Modes

The app may work with multiple todo lists. By default the mode "main" is activated. By launching the program with a single parameter (e.g. 'private' or 'work'), a new todo list is created and used for the particular execution of the program. If no argument is provided, the default list is used, indicated in the status line as 'main' (after the F10 Exit command). From version 1.0.11 on, you can also press 'm' to show the mode selection dialog. This dialog is also shown if you click on the mode name in the status bar. The mode selection dialog allows selection of all existing modes (which do not start with a dot), or by clicking 'Add' the creation of a new mode.

<img src="https://user-images.githubusercontent.com/11664020/207910707-c72c1b17-5550-4806-9d63-85d835427e61.png" width="75%" height="75%"/>

## Compatibility

* Linux (release `todo` executable), requires installed `vim` editor for editing longer todo item note text (hotkey 'n')
* Windows (release `todo.exe`), editing notes (hotkey 'n') is performed in Notepad
* macOS (arm64) (also uses `vim` as the note text editor, vim is installed by default)

You may set the environment variable `EDITOR`, to use a different editor in Linux, Mac or Windows.
Alternatively, you may set the environment variable `VISUAL` to set a graphical editor. Once invoked, the user interface of the ToDo Appp will be blocked and a info message will be shown. Once the editor is closed, the ToDo app will read the temporary editor file and proceed operation.

Example call on Mac OS, to use the default text editor GUI to edit notes:

```bash
$ VISUAL="/usr/bin/open -e -W" todo
```

## Screenshots

Remark: The current version shows available commands and active mode (see above), in a status line at the bottom of the screen.

## Help

![image](https://user-images.githubusercontent.com/11664020/173088701-9043227a-9e86-4319-b04d-f33103c82c72.png)

## Archive item

![image](https://user-images.githubusercontent.com/11664020/173088646-1ac573d3-c34d-44ad-9b9b-1f963602e206.png)

## Add item

![image](https://user-images.githubusercontent.com/11664020/173089014-685a21c1-6eb8-4a40-ad00-29f2abb817e0.png)
