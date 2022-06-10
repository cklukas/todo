# TODO Command Line App

A simple Kanban board for your terminal.

* Stores data in a simple JSON document in $HOME/.todo/todo.json
* Makes a daily backup of the data in $HOME/.todo/backup/
* Contains a function to archive an todo item in $HOME/.todo/archive

* Allows input of topic and second description line
* Provides function to view/edit a longer note for each item in vim
* Use [red], [blue] etc. to colorize your item text

# Compatibility

* Linux (release `todo` executable), requires installed `vim` editor for editing longer todo item note text (hotkey 'n')
* Windows (release `todo.exe`), editing notes is currently not working (tries to start `vim`, which is likely not available)

# Screenshots

## Help

![image](https://user-images.githubusercontent.com/11664020/173088701-9043227a-9e86-4319-b04d-f33103c82c72.png)

## Archive item

![image](https://user-images.githubusercontent.com/11664020/173088646-1ac573d3-c34d-44ad-9b9b-1f963602e206.png)

## Add item

![image](https://user-images.githubusercontent.com/11664020/173089014-685a21c1-6eb8-4a40-ad00-29f2abb817e0.png)
