package main

import (
	"fmt"
	"os"
	"time"

	req "pples-caravan/internal/request"
	mr "pples-caravan/mapregion"

	"github.com/jroimartin/gocui"
)

const (
	// Dogmatic offsets
	OFFSET_X = 4
	OFFSET_Y = 2

	VIEW         = "main"
	STATUS       = "status"
	CARAVAN_INFO = "caravan_info"
)

func view(g *gocui.Gui) error {
	m := mr.NewMap()
	minWidth := m.Size.Col
	minHeight := m.Size.Row

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

	// Status view
	if sv, err := g.SetView(STATUS, OFFSET_X, maxY-3, maxX-OFFSET_X, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sv.Frame = false
		sv.BgColor = gocui.ColorWhite
		sv.FgColor = gocui.ColorBlack
	}

	_ = updateStatusPos(g)

	// Caravan info view
	if civ, err := g.SetView(CARAVAN_INFO, x1+1, y0, maxX-OFFSET_X, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		civ.Title = "Caravan Info"
		civ.Wrap = true
		civ.Frame = true
		civ.Editable = false
		civ.Autoscroll = false
		civ.SetCursor(0, 0)
		civ.SetOrigin(0, 0)

		// TODO: receive url from config
		caravan := req.NewCaravanInfo("https://storage.googleapis.com/pple-media/election-2569/caravan.json")

		interval := 5
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		caravanDone = make(chan struct{})
		bgWG.Go(func() {
			defer ticker.Stop()
			for {
				select {
				case <-caravanDone:
					return
				case <-ticker.C:
					_, _, err := caravan.MakeRequest()
					if err != nil {
						fmt.Fprintf(civ, "Error fetching caravan info: %v\n", err)
						// retry on next tick
						continue
					}

					select {
					case <-caravanDone:
						return
					default:
					}

					g.Update(func(g *gocui.Gui) error {
						civ.Clear()
						fmt.Fprint(civ, caravan.String())

						return nil
					})
				}
			}
		})
	}

	g.SetCurrentView(CARAVAN_INFO)
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
		if caravanDone != nil {
			close(caravanDone)
			caravanDone = nil
		}
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

	if err := g.SetKeybinding("", 'k', gocui.ModNone, move(0, -1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone, move(0, 1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'h', gocui.ModNone, move(-1, 0)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', gocui.ModNone, move(1, 0)); err != nil {
		return err
	}

	return nil
}
