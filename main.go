package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
	mr "pples-caravan/mapregion"
)

const (
	VIEW         = "main"
	STATUS       = "status"
	CARAVAN_INFO = "caravan_info"

	// Dogmatic offsets
	OFFSET_X = 4
	OFFSET_Y = 2
)

func close(g *gocui.Gui) {
	g.Close()
	log.Println("GUI closed successfully")
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalln(err)
	}
	defer close(g)

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

func view(g *gocui.Gui) error {
	m := mr.NewMap()
	minWidth := m.Pos.Col
	minHeight := m.Pos.Row

	maxX, maxY := g.Size()

	buffer := 2

	x0 := OFFSET_X
	y0 := OFFSET_Y

	x1 := OFFSET_X + minWidth + buffer
	y1 := OFFSET_Y + minHeight

	if v, err := g.SetView(VIEW, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Fprint(os.Stdout, "SetView error:", err)
			return err
		}

		actualWidth := x1 - x0 - buffer
		actualHeight := y1 - y0 - buffer

		v.Title = fmt.Sprintf("Caravan | View (%d x %d) ", actualWidth, actualHeight)
		v.Wrap = true
		v.Frame = true
		v.Editable = false
		v.Autoscroll = true
		v.SetCursor(0, 0)

		_, _ = g.SetCurrentView(VIEW)

		for _, r := range m.Grid {
			for _, c := range r {
				if c == "" {
					fmt.Fprint(v, "    ")
				} else {
					s := c + mr.X
					fmt.Fprintf(v, "[%s]", s)
				}
			}
			fmt.Fprintln(v)
		}
	}

	if sv, err := g.SetView(STATUS, OFFSET_X, maxY-3, maxX-OFFSET_X, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sv.Frame = false
		sv.BgColor = gocui.ColorWhite
		sv.FgColor = gocui.ColorBlack
	}

	_ = updateStatusPos(g)
	return nil
}

func updateStatusPos(g *gocui.Gui) error {
	v, err := g.View(VIEW)
	if err != nil || v == nil {
		return nil
	}
	cx, cy := v.Cursor()
	sv, err := g.View(STATUS)
	if err != nil || sv == nil {
		return nil
	}
	sv.Clear()
	ox, oy := v.Origin()

	fmt.Fprintf(sv, "pos: %d,%d", cx, cy)
	fmt.Fprintf(sv, " | origin: %d,%d", ox, oy)
	fmt.Fprintf(sv, " | Press Ctrl+C to exit.")

	g.SetRune(ox, oy, ' ', gocui.Attribute(gocui.ModNone), gocui.ColorRed)

	return nil
}

func setKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}

	move := func(dx, dy int) func(*gocui.Gui, *gocui.View) error {
		return func(g *gocui.Gui, v *gocui.View) error {
			if v == nil {
				return nil
			}
			v.MoveCursor(dx, dy, false)
			return updateStatusPos(g)
		}
	}

	if err := g.SetKeybinding(VIEW, 'k', gocui.ModNone, move(0, -1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(VIEW, 'j', gocui.ModNone, move(0, 1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(VIEW, 'h', gocui.ModNone, move(-1, 0)); err != nil {
		return err
	}
	if err := g.SetKeybinding(VIEW, 'l', gocui.ModNone, move(1, 0)); err != nil {
		return err
	}

	// refresh/redraw
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return nil //TODO: call refresh function
	}); err != nil {
		return err
	}

	return nil
}
