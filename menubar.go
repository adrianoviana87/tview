package tview

import (
	"github.com/gdamore/tcell/v2"
	"strings"
)

// MenuOption is one option that can be selected in a drop-down primitive.
type MenuOption struct {
	Id        string
	Text      string // The text to be displayed in the drop-down.
	Structure string // The menu structure separated by periods.
}

func NewMenuOption(id, text, structure string) *MenuOption {
	return &MenuOption{
		Id:        id,
		Text:      text,
		Structure: structure,
	}
}

// MenuBar implements a selection widget whose options become visible in a
// drop-down list when activated.
//
// See https://github.com/rivo/tview/wiki/MenuBar for an example.
type MenuBar struct {
	*Box

	// Whether or not this drop-down is disabled/read-only.
	disabled bool

	isActive bool

	topLevelMenus []*MenuItem

	focusedMenu *MenuItem

	// A callback function which is called when the user changes the drop-down's
	// selection.
	selected func(id string)
}

func getFirstLevelOptions(options []*MenuOption) []*MenuOption {
	firstLevelOptions := make([]*MenuOption, 0)
	for _, option := range options {
		if !strings.Contains(option.Structure, ".") {
			firstLevelOptions = append(firstLevelOptions, option)
		}
	}

	return firstLevelOptions
}

func getChildren(options []*MenuOption, parent *MenuOption) []*MenuOption {
	children := make([]*MenuOption, 0)
	for _, option := range options {
		// get only the first level children
		if strings.HasPrefix(option.Structure, parent.Structure) && strings.Count(option.Structure, ".") == strings.Count(parent.Structure, ".") + 1 {
			children = append(children, option)
		}
	}

	return children
}

func makeMenu(options []*MenuOption, parent *MenuOption, orientation Orientation, bar *MenuBar) *MenuItem {
	parentMenu := NewMenuItem(parent, orientation, bar)
	children := getChildren(options, parent)

	for _, child := range children {
		menu := makeMenu(options, child, OrientationHorizontal, bar)
		parentMenu.AddChild(menu)
	}

	return parentMenu
}

func NewMenuBar(options []*MenuOption) *MenuBar {
	topLevelOptions := getFirstLevelOptions(options)
	topLevelMenus := make([]*MenuItem, len(topLevelOptions))

	bar := &MenuBar{
		Box: NewBox(), 
	}

	for i, option := range topLevelOptions {
		topLevelMenus[i] = makeMenu(options, option, OrientationVertical, bar)
	}

	bar.topLevelMenus = topLevelMenus

	return bar
}

func (m *MenuBar) SetFocusedMenu(menu *MenuItem) {
	m.focusedMenu = menu
}

func (m *MenuBar) Close() {
	m.isActive = false
	for _, menu := range m.topLevelMenus {
		menu.Close()
	}

	m.SetFocusedMenu(nil)
}

// SetDisabled sets whether or not the item is disabled / read-only.
func (m *MenuBar) SetDisabled(disabled bool) *MenuBar {
	m.disabled = disabled
	return m
}

func (m *MenuBar) SetSelectedFunc(handler func(id string)) *MenuBar {
	m.selected = handler
	return m
}

func (m *MenuBar) executeMenu(menu *MenuItem) {
	if m.selected != nil {
		m.selected(menu.option.Id)
	}

	m.Close()
}

// Draw draws this primitive onto the screen.
func (m *MenuBar) Draw(screen tcell.Screen) {
	m.Box.DrawForSubclass(screen, m)

	innerRect := NewRectangleFromValues(m.GetInnerRect())

	if innerRect.IsEmpty() {
		return
	}

	// Draw the options.
	for _, menu := range m.topLevelMenus {
		width := TaggedStringWidth(menu.option.Text)
		menu.SetRect(innerRect.X, innerRect.Y, width, 1)
		menu.UpdateChildrenLocation()
		innerRect.X += width + 1
		menu.Draw(screen)
	}
}

// InputHandler returns the handler for this primitive.
func (m *MenuBar) InputHandler() func(event *tcell.EventKey, setFocus func(p Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p Primitive)) {
		if m.disabled {
			return
		}
	})
}

// Focus is called by the application when the primitive receives focus.
func (m *MenuBar) Focus(delegate func(p Primitive)) {
	m.Box.Focus(delegate)
}

// HasFocus returns whether or not this primitive has focus.
func (m *MenuBar) HasFocus() bool {
	return m.Box.HasFocus()
}

// MouseHandler returns the mouse handler for this primitive.
func (m *MenuBar) MouseHandler() func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
	return m.WrapMouseHandler(func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
		if m.disabled {
			return false, nil
		}

		for _, menu := range m.topLevelMenus {
			consumed, capture = menu.MouseHandler()(action, event, setFocus)
			if consumed {
				return consumed, capture
			}
		}

		if action == MouseLeftDown && m.isActive {
			m.Close()
			return true, m
		}

		return false, nil
	})
}
