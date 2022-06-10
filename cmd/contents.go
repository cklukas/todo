package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

type Item struct {
	Title     string
	Secondary string
	Note      string
}

type Content struct {
	Titles        []string
	Items         [][]Item
	fname         string `json:"-"`
	archiveFolder string `json:"-"`
	backupFolder  string `json:"-"`
}

func NewContentIo(r io.Reader) *Content {
	c := &Content{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(c); err != nil {
		return nil
	}
	return c
}

func NewContentDefault() *Content {
	ret := &Content{}
	ret.Titles = []string{"To Do", "Doing", "Done"}
	ret.Items = make([][]Item, 3)
	return ret
}

func (c *Content) GetNumLanes() int {
	return len(c.Titles)
}

func (c *Content) GetLaneTitle(idx int) string {
	return fmt.Sprintf(" %v (%v) ", c.Titles[idx], len(c.Items[idx]))
}

func (c *Content) GetLaneItems(idx int) []Item {
	return c.Items[idx]
}

func (c *Content) MoveItem(fromlane, fromidx, tolane, toidx int) {
	item := c.Items[fromlane][fromidx]
	// https://github.com/golang/go/wiki/SliceTricks
	c.Items[fromlane] = append(c.Items[fromlane][:fromidx], c.Items[fromlane][fromidx+1:]...)
	c.Items[tolane] = append(c.Items[tolane][:toidx], append([]Item{item}, c.Items[tolane][toidx:]...)...)
}

func (c *Content) DelItem(lane, idx int) {
	c.Items[lane] = append(c.Items[lane][:idx], c.Items[lane][idx+1:]...)
}

func (c *Content) ArchiveItem(lane, idx int) error {
	now := time.Now()
	archiveItemFileName := fmt.Sprintf("%v.%v.json", now.Format("2006-01-02 15_04_05.000"), c.Titles[lane])

	cnt, _ := json.MarshalIndent(c.Items[lane][idx], "", " ")
	err := ioutil.WriteFile(path.Join(c.archiveFolder, archiveItemFileName), cnt, 0644)
	if err != nil {
		return err
	}
	c.Items[lane] = append(c.Items[lane][:idx], c.Items[lane][idx+1:]...)
	return nil
}

func (c *Content) AddItem(lane, idx int, title string, secondary string) {
	c.Items[lane] = append(c.Items[lane][:idx], append([]Item{{title, secondary, ""}}, c.Items[lane][idx:]...)...)
}

func (c *Content) Read() {
	f, err := os.OpenFile(c.fname, os.O_RDONLY, os.FileMode(int(0600)))
	if err == nil {
		decoder := json.NewDecoder(f)
		if err := decoder.Decode(c); err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
}

func (c *Content) SetFileName(fname, archiveFolder, backupFolder string) {
	c.fname = fname
	c.archiveFolder = archiveFolder
	c.backupFolder = backupFolder
}

func (c *Content) Save() error {
	cnt, _ := json.MarshalIndent(c, "", " ")

	now := time.Now()
	dayFileName := path.Join(c.backupFolder, fmt.Sprintf("%v.json", now.Format("2006-01-02")))
	if _, err := os.Stat(dayFileName); errors.Is(err, os.ErrNotExist) {
		err = ioutil.WriteFile(dayFileName, cnt, 0644)
		if err != nil {
			return err
		}
	}

	err := ioutil.WriteFile(c.fname, cnt, 0644)
	if err != nil {
		return err
	}

	return nil
}
