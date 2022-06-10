# TODO Command Line App

A simple Kanban board for your terminal.

* Stores data in a simple JSON document in $HOME/.todo/todo.json
* Makes a daily backup of the data in $HOME/.todo/backup/
* Contains a function to archive an todo item in $HOME/.todo/archive

* Allows input of topic and second description line
* Provides function to view/edit a longer note for each item in vim
* Use [red], [blue] etc. to colorize your item text

Keys:

* a: Archive item
* +/Ins: Add item
* d: Delete item
* e: Edit item
* n: View/edit notes for item
* h/?: Help
* q: Quit
* Enter: Select / deselect item (selected items can be moved)
* Arrows: Move around
* Tab: Move in forms

