package main

import (
	"errors"
	"fmt"
	"image/jpeg"
	"net/http"
	"strings"
	"time"

	age "github.com/bearbin/go-age"
	"github.com/jroimartin/gocui"
	"github.com/marc-gr/asciize"
)

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "pic" {
		_, err := g.SetCurrentView("bio")
		return err
	}
	_, err := g.SetCurrentView("pic")
	return err
}

func scrollUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		if oy <= 0 {
			return nil
		}
		v.SetOrigin(ox, oy-1)
	}
	return nil
}

func scrollDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, h := v.Size()
		// Disable infinite scrolling.
		if _, err := v.Line(h); err != nil {
			return nil
		}
		v.SetOrigin(ox, oy+1)
	}
	return nil
}

const swipe = "asdfghjkl"

func (model *RecsModel) partialSwipe(r rune) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		idx := strings.IndexRune(swipe, r)
		if idx == -1 {
			return errors.New("unexpected rune")
		}
		if model.lastSwipe.idx == -1 {
			model.lastSwipe.idx = idx
			model.lastSwipe.dir = 0
			return nil
		}
		dir := idx - model.lastSwipe.idx
		if dir == 0 {
			return nil
		}
		model.lastSwipe.idx = idx
		if model.lastSwipe.dir == 0 {
			model.lastSwipe.dir = dir
			return nil
		}
		if model.lastSwipe.dir < 0 && dir > 0 || model.lastSwipe.dir > 0 && dir < 0 {
			// Reset
			model.lastSwipe.dir = 0
			model.lastSwipe.idx = -1
			return nil
		}
		if idx > 0 && idx < len(swipe)-1 {
			return nil
		}
		model.finishSwipe(g, dir > 0)
		return nil
	}
}

func (model *RecsModel) finishSwipe(g *gocui.Gui, isRight bool) {
	if len(model.recs) > 0 {
		user := model.recs[model.userIdx]
		if isRight {
			model.client.SwipeRight(&user)
		} else {
			model.client.SwipeLeft(&user)
		}
		model.userIdx++
	}
	if model.userIdx >= len(model.recs) {
		model.fetchUsers()
	}
	model.picIdx = 0
	model.drawPhoto(g)
	model.drawBio(g)
}

func (model *RecsModel) fetchUsers() {
	recs, err := model.client.GetRecs()
	if err == nil {
		model.SetRecs(recs)
	}
}

func (model *RecsModel) nextPic() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v != nil && len(model.recs) > model.userIdx {
			if model.picIdx < len(model.recs[model.userIdx].Photos)-1 {
				model.picIdx++
				model.drawPhoto(g)
			}
		}
		return nil
	}
}

func (model *RecsModel) prevPic() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v != nil && len(model.recs) > model.userIdx {
			if model.picIdx > 0 {
				model.picIdx--
				model.drawPhoto(g)
			}
		}
		return nil
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

type RecsModel struct {
	// State variables.
	recs      []Recommendation
	userIdx   int
	picIdx    int
	lastSwipe struct {
		idx int
		dir int
	}
	client *TinderClient
}

func (model *RecsModel) SetRecs(recs []Recommendation) {
	model.recs = recs
	model.userIdx = 0
	model.picIdx = 0
}

func NewRecsModel(client *TinderClient) *RecsModel {
	model := &RecsModel{}
	model.lastSwipe.idx = -1
	model.client = client
	return model
}

func (model *RecsModel) drawPhoto(g *gocui.Gui) {
	if v, err := g.View("pic"); err == nil {
		v.Clear()
		v.SetOrigin(0, 0)
		if len(model.recs) == 0 {
			return
		}
		user := model.recs[model.userIdx]
		fmt.Fprintln(v, fmt.Sprintf("%d/%d", model.picIdx+1, len(user.Photos)))
		if len(user.Photos) == 0 {
			fmt.Fprintln(v, "User has no photos!")
			return
		}
		url := user.Photos[model.picIdx].Url
		response, err := http.Get(url)
		if err != nil {
			fmt.Fprintln(v, err.Error())
			return
		}
		defer response.Body.Close()
		img, err := jpeg.Decode(response.Body)
		if err != nil {
			fmt.Fprintln(v, err.Error())
		}
		w, _ := v.Size()
		a := asciize.NewAsciizer(asciize.Width(uint(w)))
		s, err := a.Asciize(img)
		if err != nil {
			fmt.Fprintln(v, err.Error())
		}
		fmt.Fprintln(v, s)
	}
}

func (model *RecsModel) drawBio(g *gocui.Gui) {
	if v, err := g.View("bio"); err == nil {
		v.Clear()
		v.SetOrigin(0, 0)
		if len(model.recs) == 0 {
			fmt.Fprint(v, fmt.Sprintf("No recommendations found! :(\nSwipe to try again."))
			return
		}
		user := model.recs[model.userIdx]
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
	}
}

func (model *RecsModel) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	const picHeight = 48
	if v, err := g.SetView("pic", -1, -1, maxX, picHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = false
		v.Wrap = false
		model.drawPhoto(g)
		if _, err := g.SetCurrentView("pic"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("bio", -1, picHeight, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = false
		v.Wrap = true
		model.drawBio(g)
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
	if err := g.SetKeybinding("", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'a', gocui.ModNone, model.partialSwipe('a')); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 's', gocui.ModNone, model.partialSwipe('s')); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'd', gocui.ModNone, model.partialSwipe('d')); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'f', gocui.ModNone, model.partialSwipe('f')); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'g', gocui.ModNone, model.partialSwipe('g')); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'h', gocui.ModNone, model.partialSwipe('h')); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone, model.partialSwipe('j')); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'k', gocui.ModNone, model.partialSwipe('k')); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', gocui.ModNone, model.partialSwipe('l')); err != nil {
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
