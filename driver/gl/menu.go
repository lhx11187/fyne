package gl

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"image/color"
)

type menuBarAction struct {
	Label string

	menu   *fyne.Menu
	canvas fyne.Canvas
}

// ToolbarObject gets a button to render this ToolbarAction
func (m *menuBarAction) ToolbarObject() fyne.CanvasObject {
	button := widget.NewButton(m.Label, nil)

	button.OnTapped = func() {
		pos := button.Position().Add(fyne.NewPos(0, button.Size().Height))
		showMenu(m.menu, pos, m.canvas)
	}

	return button
}

func newMenuBarAction(menu *fyne.Menu, w fyne.Window) widget.ToolbarItem {
	return &menuBarAction{menu.Label, menu, w.Canvas()}
}

func buildMenuBar(menus *fyne.MainMenu, w fyne.Window) *widget.Toolbar {
	var items []widget.ToolbarItem

	for _, menu := range menus.Items {
		items = append(items, newMenuBarAction(menu, w))
	}
	return widget.NewToolbar(items...)
}

func showMenu(menu *fyne.Menu, pos fyne.Position, c fyne.Canvas) {
	pop := widget.NewPopUpMenu(fyne.NewMenu("", menu.Items...), c)
	pop.Move(pos)
}

func (c *glCanvas) menuBar() fyne.CanvasObject {
	if c.window.mainmenu == nil {
		return nil
	}

	c.RLock()
	ret := c.menu
	c.RUnlock()
	if ret != nil {
		return ret
	}

	if hasNativeMenu() {
		setupNativeMenu(c.window.mainmenu)

		ret = canvas.NewRectangle(color.Transparent) // just a dummy value really
	} else {
		ret = buildMenuBar(c.window.mainmenu, c.window)
	}
	c.Lock()
	c.menu = ret
	c.Unlock()
	return ret
}

func (c *glCanvas) menuHeight() int {
	bar := c.menuBar()

	if bar == nil { // reserve no height for a native menus or windows with no menu
		return 0
	}

	return bar.MinSize().Height
}
