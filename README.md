# TODO Command Line App

[![Go Report](https://goreportcard.com/badge/github.com/cklukas/todo)](https://goreportcard.com/report/github.com/cklukas/todo)

A simple Kanban board for your terminal.

* Stores data in a simple JSON document in `$HOME/.todo/todo.json`
* Makes a daily backup of the data in `$HOME/.todo/backup/` (at first start on a particular day)
* Contains a function to archive an todo item in `$HOME/.todo/archive`
* If a non-default mode is used (see below), the files and folders for that mode (`todo.json`, `backup`, `archive`) are saved under `$home/.todo/mode/[mode]`

* Allows input of topic and second description line
* Provides function to view/edit a longer note for each item in vim
* All changes are immediately saved (no save command)
* The application can be started multiple times, modifications performed in one instance are detected in other instances (through monitoring changes to the active `todo.json` file)
* Use [red], [blue] etc. to colorize your item text
* New: Hotkeys F1..F10 are shown in status bar

## Modes

The app may work with multiple todo lists. By default the mode "main" is activated. By launching the program with a single parameter (e.g. 'private' or 'work'), a new todo list is created and used for the particular execution of the program. If no argument is provided, the default list is used, indicated in the status line as 'main' (after the F10 Exit command). From version 1.0.11 on, you can also press 'm' to show the mode selection dialog. This dialog is also shown if you click on the mode name in the status bar. The mode selection dialog allows selection of all existing modes (which do not start with a dot), or by clicking 'Add' the creation of a new mode.

<img src="https://user-images.githubusercontent.com/11664020/207910707-c72c1b17-5550-4806-9d63-85d835427e61.png" width="75%" height="75%"/>

## Compatibility

* Linux (release `todo` executable), requires installed `vim` editor for editing longer todo item note text (hotkey 'n')
* Windows (release `todo.exe`), editing notes (hotkey 'n') is performed in Notepad
* macOS (arm64) (also uses `vim` as the note text editor, vim is installed by default)

## Screenshots

Remark: The current version shows available commands and active mode (see above), in a status line at the bottom of the screen.

## Help

![image](https://user-images.githubusercontent.com/11664020/173088701-9043227a-9e86-4319-b04d-f33103c82c72.png)

## Archive item

![image](https://user-images.githubusercontent.com/11664020/173088646-1ac573d3-c34d-44ad-9b9b-1f963602e206.png)

## Add item

![image](https://user-images.githubusercontent.com/11664020/173089014-685a21c1-6eb8-4a40-ad00-29f2abb817e0.png)
