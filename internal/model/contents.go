package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"sync"
	"time"

	"github.com/cklukas/todo/internal/util"
	"github.com/flytam/filenamify"
	"github.com/google/uuid"
)

type Item struct {
	Title         string
	Secondary     string
	Note          string
	Guid          string
	Priority      int
	IsArchived    bool
	Created       string
	LastUpdate    string
	Due           string
	Color         string
	UserName      string
	UpdatedByName string
	Mode          string
}

type ToDoContent struct {
	Titles         []string
	Items          [][]Item
	SortModes      []string
	fname          string     `json:"-"`
	archiveFolder  string     `json:"-"`
	backupFolder   string     `json:"-"`
	readWriteMutex sync.Mutex `json:"-"`
}

func (c *ToDoContent) Lock() {
	c.readWriteMutex.Lock()
}

func (c *ToDoContent) Unlock() {
	c.readWriteMutex.Unlock()
}

func (c *ToDoContent) InitializeNew() {
	c.Titles = []string{"To Do", "Doing", "Done"}
	c.Items = make([][]Item, 3)
	c.SortModes = make([]string, 3)
}

func (c *ToDoContent) ReadFromFile(fname string) error {
	f, err := os.OpenFile(fname, os.O_RDONLY, os.FileMode(int(0600)))
	if err != nil {
		return err
	}

	err = json.NewDecoder(f).Decode(c)
	f.Close()
	if err != nil {
		return err
	}
	c.normalize()
	return nil
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
	if len(c.SortModes) > lane {
		c.SortModes = append(c.SortModes[:lane], c.SortModes[lane+1:]...)
	}
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
	c.SortModes = append(c.SortModes[:i], append([]string{""}, c.SortModes[i:]...)...)

	return i
}

func (c *ToDoContent) MoveItem(fromlane, fromidx, tolane, toidx int) {
	item := c.Items[fromlane][fromidx]
	// https://github.com/golang/go/wiki/SliceTricks
	c.Items[fromlane] = append(c.Items[fromlane][:fromidx], c.Items[fromlane][fromidx+1:]...)
	c.Items[tolane] = append(c.Items[tolane][:toidx], append([]Item{item}, c.Items[tolane][toidx:]...)...)
}

func (c *ToDoContent) SetLaneSort(idx int, mode string) {
	if idx >= 0 && idx < len(c.SortModes) {
		c.SortModes[idx] = mode
	}
}

func (c *ToDoContent) SortLane(idx int) {
	if idx >= 0 && idx < len(c.Items) {
		sortItems(c.Items[idx], c.SortModes[idx])
	}
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

	item := c.Items[lane][idx]
	item.IsArchived = true
	item.LastUpdate = now.UTC().Format(time.RFC3339)
	if usr, errU := user.Current(); errU == nil {
		item.UpdatedByName = usr.Username
	}

	cnt, _ := json.MarshalIndent(item, "", " ")
	err = os.WriteFile(path.Join(c.archiveFolder, archiveItemFileName), cnt, 0644)
	if err != nil {
		return err
	}
	c.Items[lane] = append(c.Items[lane][:idx], c.Items[lane][idx+1:]...)
	return nil
}

func (c *ToDoContent) AddItem(lane, idx int, title string, secondary string, priority int, due, color string) {
	now := time.Now().UTC().Format(time.RFC3339)
	usr, err := user.Current()
	userName := ""
	if err == nil {
		userName = usr.Username
	}

	newItem := Item{
		Title:         title,
		Secondary:     secondary,
		Note:          "",
		Guid:          uuid.NewString(),
		Priority:      priority,
		IsArchived:    false,
		Created:       now,
		LastUpdate:    now,
		Due:           due,
		Color:         color,
		UserName:      userName,
		UpdatedByName: userName,
		Mode:          "",
	}

	c.Items[lane] = append(c.Items[lane][:idx], append([]Item{newItem}, c.Items[lane][idx:]...)...)
}

func (c *ToDoContent) normalize() {
	now := time.Now().UTC().Format(time.RFC3339)
	usr, err := user.Current()
	userName := ""
	if err == nil {
		userName = usr.Username
	}

	if len(c.SortModes) != len(c.Titles) {
		c.SortModes = make([]string, len(c.Titles))
	}

	for li := range c.Items {
		for ii := range c.Items[li] {
			item := &c.Items[li][ii]
			if item.Guid == "" {
				item.Guid = uuid.NewString()
			}
			if item.Created == "" {
				item.Created = now
			}
			if item.LastUpdate == "" {
				if item.Created != "" {
					item.LastUpdate = item.Created
				} else {
					item.LastUpdate = now
				}
			}
			if item.Color == "" {
				var col string
				col, item.Title = util.ParsePrefix(item.Title)
				item.Color = col
			}
			if item.UserName == "" {
				item.UserName = userName
			}
			if item.UpdatedByName == "" {
				if item.UserName != "" {
					item.UpdatedByName = item.UserName
				} else {
					item.UpdatedByName = userName
				}
			}
		}
	}
}

func (c *ToDoContent) Read() error {
	c.readWriteMutex.Lock()
	defer c.readWriteMutex.Unlock()
	f, err := os.OpenFile(c.fname, os.O_RDONLY, os.FileMode(int(0600)))
	if err == nil {
		decoder := json.NewDecoder(f)
		if err := decoder.Decode(c); err != nil {
			f.Close()
			return err
		}
		f.Close()
		c.normalize()
		return nil
	}
	return err
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
