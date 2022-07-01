package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"time"
)

import "golang.design/x/clipboard"

type App struct {
	ctx context.Context
}

type Entry struct {
	Content    string    `json:"content"`
	LastCopied int64     `json:"lastCopied"`
	UUID       uuid.UUID `json:"id"`
}

type Entries struct {
	Entries []Entry `json:"entries"`
}

func NewApp() *App {
	return &App{}
}

var ENTRIES = map[string]*Entry{}

func pollClipboardManual(clipboardChan chan<- string, duration time.Duration) {
	// i'd rather have used the packaged clipboard.Watch function, but its ticker duration is 1 second which is too slow
	lastclipItem := string(clipboard.Read(clipboard.FmtText))

	for {
		clipItem := string(clipboard.Read(clipboard.FmtText))
		if clipItem != lastclipItem {
			clipboardChan <- clipItem
			lastclipItem = clipItem
		}
		time.Sleep(duration)
	}
}

func handleChangedClipboard(ctx context.Context, entryPub <-chan string) {
	for clipBoardItem := range entryPub {
		val, ok := ENTRIES[clipBoardItem]
		var msg string
		if ok {
			newTime := time.Now().Unix()
			msg = fmt.Sprintf("updated from old time %v to new time %v", val.LastCopied, newTime)
			ENTRIES[clipBoardItem].LastCopied = newTime
		} else {
			msg = "is a new clipboard item"
			ENTRIES[clipBoardItem] = &Entry{
				Content:    clipBoardItem,
				LastCopied: time.Now().Unix(),
				UUID:       uuid.New(),
			}
		}

		entries := &Entries{
			Entries: []Entry{*ENTRIES[clipBoardItem]},
		}
		fmt.Printf("%v %v . sending to frontend\n", entries, msg)
		runtime.EventsEmit(ctx, "newClips", entries)
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	clipboardChan := make(chan string)
	clipboardPollInterval := time.Millisecond * 100
	go pollClipboardManual(clipboardChan, clipboardPollInterval)
	go handleChangedClipboard(ctx, clipboardChan)
}
