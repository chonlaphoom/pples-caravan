package main

import (
	"log"
	"sync"

	"github.com/jroimartin/gocui"
)

var caravanDone chan struct{}
var bgWG sync.WaitGroup

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalln(err)
	}
	defer closeGUI(g)

	g.Cursor = true
	g.SetManagerFunc(view)

	if err := setKeybindings(g); err != nil {
		log.Println("setKeybindings:", err)
		return
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Println("main loop error:", err)
	}
}

func closeGUI(g *gocui.Gui) {
	bgWG.Wait()
	log.Println("Background tasks completed")
	g.Close()
	log.Println("GUI closed successfully")
}
