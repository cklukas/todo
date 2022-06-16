module github.com/cklukas/todo

go 1.18

replace github.com/rivo/tview v0.0.0-20220307222120-9994674d60a8 => github.com/cklukas/tview v0.0.0-20220523201936-989a18252a05

require github.com/spf13/cobra v1.4.0

require (
	github.com/fsnotify/fsnotify v1.5.4
	github.com/gdamore/tcell/v2 v2.5.1
	github.com/rivo/tview v0.0.0-20220307222120-9994674d60a8
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
)
