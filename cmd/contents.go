package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/flytam/filenamify"
)

type Item struct {
	Title     string
	Secondary string
	Note      string
}

type ToDoContent struct {
	Titles         []string
	Items          [][]Item
	fname          string     `json:"-"`
	archiveFolder  string     `json:"-"`
	backupFolder   string     `json:"-"`
	readWriteMutex sync.Mutex `json:"-"`
}

func (c *ToDoContent) InitializeNew() {
	c.Titles = []string{"To Do", "Doing", "Done"}
	c.Items = make([][]Item, 3)
}

func (c *ToDoContent) ReadFromFile(fname string) error {
	f, err := os.OpenFile(fname, os.O_RDONLY, os.FileMode(int(0600)))
	if err != nil {
		return err
	}

	err = json.NewDecoder(f).Decode(c)

	f.Close()

	return err
}

func (c *ToDoContent) GetNumLanes() int {
	return len(c.Titles)
}

func (c *ToDoContent) GetLaneTitle(idx int) string {
	return fmt.Sprintf(" %v (%v) ", c.Titles[idx], len(c.Items[idx]))
}

func (c *ToDoContent) SetLaneTitle(idx int, title string) {
	c.Titles[idx] = title
}

func (c *ToDoContent) GetLaneItems(idx int) []Item {
	return c.Items[idx]
}

func (c *ToDoContent) RemoveLane(lane int) {
	c.Titles = append(c.Titles[:lane], c.Titles[lane+1:]...)
	c.Items = append(c.Items[:lane], c.Items[lane+1:]...)
}

func (c *ToDoContent) InsertNewLane(addToLeft bool, laneTitle string, relativeToLaneIdx int) int {
	i := relativeToLaneIdx
	if !addToLeft {
		i++
	}

	newItemList := make([][]Item, 0)
	newItemList = append(newItemList, []Item{})
	c.Items = append(c.Items[:i], append(newItemList, c.Items[i:]...)...)
	c.Titles = append(c.Titles[:i], append([]string{laneTitle}, c.Titles[i:]...)...)

	return i
}

func (c *ToDoContent) MoveItem(fromlane, fromidx, tolane, toidx int) {
	item := c.Items[fromlane][fromidx]
	// https://github.com/golang/go/wiki/SliceTricks
	c.Items[fromlane] = append(c.Items[fromlane][:fromidx], c.Items[fromlane][fromidx+1:]...)
	c.Items[tolane] = append(c.Items[tolane][:toidx], append([]Item{item}, c.Items[tolane][toidx:]...)...)
}

func (c *ToDoContent) DelItem(lane, idx int) {
	c.Items[lane] = append(c.Items[lane][:idx], c.Items[lane][idx+1:]...)
}

func (c *ToDoContent) ArchiveItem(lane, idx int) error {
	now := time.Now()
	saveName, err := filenamify.FilenamifyV2(c.Titles[lane], func(options *filenamify.Options) {
		options.Replacement = "_"
	})
	if err != nil {
		return err
	}
	archiveItemFileName := fmt.Sprintf("%v.%v.json", now.Format("2006-01-02 15_04_05.000"), saveName)

	cnt, _ := json.MarshalIndent(c.Items[lane][idx], "", " ")
	err = os.WriteFile(path.Join(c.archiveFolder, archiveItemFileName), cnt, 0644)
	if err != nil {
		return err
	}
	c.Items[lane] = append(c.Items[lane][:idx], c.Items[lane][idx+1:]...)
	return nil
}

func (c *ToDoContent) AddItem(lane, idx int, title string, secondary string) {
	c.Items[lane] = append(c.Items[lane][:idx], append([]Item{{title, secondary, ""}}, c.Items[lane][idx:]...)...)
}

func (c *ToDoContent) Read() error {
	c.readWriteMutex.Lock()
	defer c.readWriteMutex.Unlock()
	f, err := os.OpenFile(c.fname, os.O_RDONLY, os.FileMode(int(0600)))
	if err == nil {
		decoder := json.NewDecoder(f)
		if err := decoder.Decode(c); err != nil {
			log.Fatal(err)
		}
		f.Close()
	} else {
		return err
	}

	return nil
}

func (c *ToDoContent) SetFileName(fname, archiveFolder, backupFolder string) {
	c.fname = fname
	c.archiveFolder = archiveFolder
	c.backupFolder = backupFolder
}

func (c *ToDoContent) Save() error {
	c.readWriteMutex.Lock()
	defer c.readWriteMutex.Unlock()

	cnt, _ := json.MarshalIndent(c, "", " ")

	now := time.Now()
	dayFileName := path.Join(c.backupFolder, fmt.Sprintf("%v.json", now.Format("2006-01-02")))
	if _, err := os.Stat(dayFileName); errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(dayFileName, cnt, 0644)
		if err != nil {
			return err
		}
	}

	err := os.WriteFile(c.fname, cnt, 0644)
	if err != nil {
		return err
	}

	return nil
}
