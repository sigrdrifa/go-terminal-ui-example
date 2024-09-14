package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const refreshInterval = 10 * time.Second
const url = "https://api.chucknorris.io/jokes/random?category=science"

var (
	app      *tview.Application
	textView *tview.TextView
)

type Payload struct {
	Value string
}

func getAndDrawJoke() {
	result, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	payloadBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}

	payload := &Payload{}
	err = json.Unmarshal(payloadBytes, payload)
	if err != nil {
		panic(err)
	}

	textView.Clear()
	fmt.Fprintln(textView, payload.Value)
	timeStr := fmt.Sprintf("\n\n[gray]%s", time.Now().Format(time.RFC1123))
	fmt.Fprintln(textView, timeStr)
}

func refresh() {
	tick := time.NewTicker(refreshInterval)
	for {
		select {
		case <-tick.C:
			getAndDrawJoke()
			app.Draw()
		}
	}
}

func renderFooter() *tview.TextView {
	return tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("Built with [red]<3[default] by Sig ([gray]@sigfaults[white])")
}

func renderHeader() *tview.TextView {
	return tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(`

     _____)        )   ___                     _____)    
   /              (__/_____) /)       /)     /           
  /   ___   ___     /       (/     _ (/_    /   ___   ___
 /     / ) (_)     /        / )(_((__/(__  /     / ) (_) 
(____ /           (______)                (____ /        
                                                         
    `)
}

func main() {
	app = tview.NewApplication()
	textView = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorLime)

	textView.SetBorderPadding(1, 0, 0, 0)

	getAndDrawJoke()

	grid := tview.NewGrid().
		SetRows(15, 0, 3).
		SetColumns(0, 0).
		AddItem(renderHeader(), 0, 0, 1, 3, 0, 0, false).
		AddItem(renderFooter(), 2, 0, 1, 3, 0, 0, false)

	grid.AddItem(textView, 1, 0, 1, 3, 0, 0, false)

	go refresh()
	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
