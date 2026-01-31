package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

const (
	VIEW     = "view"
	OFFSET_X = 4
	OFFSET_Y = 4
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalln(err)
	}
	defer close(g)

	g.SetManagerFunc(view)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		log.Fatalln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}

func close(g *gocui.Gui) {
	g.Close()
	log.Println("GUI closed successfully")
}

func view(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(VIEW, OFFSET_X, OFFSET_Y, maxX-OFFSET_X, maxY-OFFSET_Y); err != nil {
		if err != gocui.ErrUnknownView {
			log.Fatalln(err)
			return err
		}
		v.Title = "..."
		v.Wrap = false
		v.Frame = true
		coordinate := fmt.Sprintln("maxX:", maxX, "maxY:", maxY)
		v.Autoscroll = true

		v.Write([]byte("Press Ctrl+C to exit.\n"))
		v.Write([]byte(coordinate))
		v.Cursor()
	}

	return nil
}
