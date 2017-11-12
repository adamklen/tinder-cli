package main

import (
    "fmt"
	"time"

    "github.com/jroimartin/gocui"
    age "github.com/bearbin/go-age"
)

func scrollUp(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        ox, oy := v.Origin()
        if oy <= 0 {
            return nil;
        }
        v.SetOrigin(ox, oy-1)
    }
    return nil
}

func scrollDown(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        ox, oy := v.Origin()
        _, h := v.Size();
        // Disable infinite scrolling.
        if _, err := v.Line(h); err != nil {
            return nil
        }
        v.SetOrigin(ox, oy+1)
    }
    return nil
}

func (model *RecsModel) nextPic() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			if model.picIdx < len(model.recs[model.userIdx].Photos) - 1 {
				model.picIdx++
			}
			model.drawPhoto(g)
		}
		return nil
	}
}

func (model *RecsModel) prevPic() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			if model.picIdx > 0 {
				model.picIdx--
			}
			model.drawPhoto(g)
		}
		return nil
	}
}

func getLine(g *gocui.Gui, v *gocui.View) error {
    var l string
    var err error

    _, cy := v.Cursor()
    if l, err = v.Line(cy); err != nil {
        l = ""
    }

    maxX, maxY := g.Size()
    if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        fmt.Fprintln(v, l)
        if _, err := g.SetCurrentView("msg"); err != nil {
            return err
        }
    }
    return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
    if err := g.DeleteView("msg"); err != nil {
        return err
    }
    if _, err := g.SetCurrentView("side"); err != nil {
        return err
    }
    return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}

type RecsModel struct {
	// State variables.
	recs []Recommendation
	userIdx int
	picIdx int
}

func (model *RecsModel) SetRecs(recs []Recommendation) {
	model.recs = recs
	model.userIdx = 0
	model.picIdx = 0
}

func NewRecsModel() *RecsModel {
	return &RecsModel{}
}

func (model *RecsModel) drawPhoto(g *gocui.Gui) {
	if v, err := g.View("pic"); err == nil {
		v.Clear()
		// TODO print ascii image
		user := model.recs[model.userIdx]
		fmt.Fprintln(v, fmt.Sprintf("%d/%d", model.picIdx + 1, len(user.Photos)))
		fmt.Fprintln(v, user.Photos[model.picIdx].Url)
	}
}

func (model *RecsModel) Layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
	user := model.recs[model.userIdx]
    if _, err := g.SetView("pic", -1, -1, maxX, 32); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
		model.drawPhoto(g)
    }
    if v, err := g.SetView("bio", -1, 32, maxX, maxY); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Editable = false
        v.Wrap = true
		birthDate, _ := time.Parse(time.RFC3339, user.BirthDate)
		age := age.Age(birthDate)
		fmt.Fprint(v, fmt.Sprintf("Id: %s\n", user.Id))
		fmt.Fprint(v, fmt.Sprint("Name: ", model.recs[model.userIdx].Name))
		fmt.Fprint(v, fmt.Sprint(" [", age))
		if user.Gender == 0 {
			fmt.Fprint(v, "M")
		} else if user.Gender == 1 {
			fmt.Fprint(v, "F")
		}
		fmt.Fprint(v, "]\n")
		fmt.Fprint(v, fmt.Sprintf("%d miles away\n\n", user.DistanceMi))
		fmt.Fprintln(v, user.Bio)
        if _, err := g.SetCurrentView("bio"); err != nil {
            return err
        }
    }
    return nil
}

func keybindings(g *gocui.Gui, model *RecsModel) error {
    if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, scrollDown); err != nil {
        return err
    }
    if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, scrollUp); err != nil {
        return err
    }
    if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, model.nextPic()); err != nil {
        return err
    }
    if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, model.prevPic()); err != nil {
        return err
    }
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        return err
    }
    if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
        return err
    }
    return nil
}


func Run(recsModel *RecsModel) {
	g, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        panic(err)
    }
    defer g.Close()

    g.Cursor = false
    g.SetManager(recsModel)

	if err := keybindings(g, recsModel); err != nil {
        panic(err)
    }

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
        panic(err)
    }
}

