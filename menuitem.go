package tview

import (
	"strings"
	"fmt"
	"github.com/gdamore/tcell/v2"
)

type MenuItem struct {
	*Box
	option              *MenuOption
	children            []*MenuItem
	parent              *MenuItem
	bar                 *MenuBar
	childrenOrientation Orientation
	activeStyle         tcell.Style
	disabledStyle       tcell.Style
	defaultStyle        tcell.Style
	disabled            bool
	isOpen              bool
}

func NewMenuItem(option *MenuOption, orientation Orientation, bar *MenuBar) *MenuItem {
	box := NewBox()
	box.SetRect(0, 0, TaggedStringWidth(option.Text), 1)
	return &MenuItem{
		Box:                 box,
		option:              option,
		disabled:            false,
		isOpen:              false,
		childrenOrientation: orientation,
		defaultStyle: tcell.StyleDefault.
			Foreground(tcell.ColorWhite),
		activeStyle: tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorRed),
		disabledStyle: tcell.StyleDefault.
			Foreground(tcell.ColorGray).
			Background(tcell.ColorBlack),
		bar: bar,
	}
}

func (m *MenuItem) Close() {
	m.isOpen = false
	if len(m.children) > 0 {
		for _, child := range m.children {
			child.Close()
		}
	}
}

func (m *MenuItem) SetX(x int) {
	_, y, width, height := m.GetRect()
	m.Box.SetRect(x, y, width, height)
}

func (m *MenuItem) SetTopLeft(x, y int) {
	_, _, width, height := m.GetRect()
	m.Box.SetRect(x, y, width, height)
}

func (m *MenuItem) AddChild(child *MenuItem) {
	m.children = append(m.children, child)
	child.parent = m
}

func (m *MenuItem) GetWidestLength() int {
	if len(m.children) == 0 {
		return 0
	}

	widest := 0

	for _, child := range m.children {
		width := TaggedStringWidth(child.Text())
		if width > widest {
			widest = width
		}
	}

	return widest
}

func (m *MenuItem) UpdateChildrenLocation() {
	if len(m.children) == 0 || !m.isOpen {
		return
	}

	rect := NewRectangleFromValues(m.GetRect())
	width := m.GetWidestLength() + 1

	switch m.childrenOrientation {
	case OrientationVertical:
		for i, child := range m.children {
			child.SetRect(rect.X, rect.Y+i+1, width, 1)
			child.UpdateChildrenLocation()
		}
	case OrientationHorizontal:
		x := rect.Right()
		for i, child := range m.children {
			child.SetRect(x, rect.Y+i, width, 1)
			child.UpdateChildrenLocation()
		}
	}
}

func (m *MenuItem) IsFocused() bool {
	if m.bar.focusedMenu == m {
		return true
	}

	for _, child := range m.children {
		if child.IsFocused() {
			return true
		}
	}

	return false
}

func (m *MenuItem) Draw(screen tcell.Screen) {
	m.Box.DrawForSubclass(screen, m)
	innerRect := NewRectangleFromValues(m.GetRect())

	if innerRect.IsEmpty() {
		panic(fmt.Sprintf("MenuItem.Draw: %s at %v", m.option.Text, innerRect))
	}

	style := func() tcell.Style {
		if m.disabled {
			return m.disabledStyle
		}

		if m.IsFocused() {
			return m.activeStyle
		}

		return m.defaultStyle
	}()

	// fill text with spaces
	text := func() string {
		if m.parent == nil {
			return m.Text()
		}

		return " " + m.Text()
	}()

	textWidth := TaggedStringWidth(text)
	if textWidth < innerRect.Width {
		text += strings.Repeat(" ", innerRect.Width-textWidth)
	}
	printWithStyle(screen, text, innerRect.X, innerRect.Y, 0, innerRect.Width, AlignLeft, style, false)
	if m.parent != nil {
		shadowStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGray)
		screen.SetContent(innerRect.Right(), innerRect.Y, '░', nil, shadowStyle)
		printWithStyle(screen, strings.Repeat("░", innerRect.Width), innerRect.X + 1, innerRect.Y + 1, 0, innerRect.Width, AlignLeft, shadowStyle, false)
	}
	if m.isOpen {
		for _, child := range m.children {
			// panic("draw child")
			child.Draw(screen)
		}
	}
}

func (m *MenuItem) Text() string {
	if len(m.children) > 0 {
		return m.option.Text + " ⯈"
	}

	return m.option.Text
}

func (m *MenuItem) Open() {
	if m.parent == nil {
		m.bar.Close()
	}

	m.bar.isActive = true
	m.isOpen = true
}

// MouseHandler returns the mouse handler for this primitive.
func (m *MenuItem) MouseHandler() func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
	return m.WrapMouseHandler(func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
		if m.disabled {
			return false, nil
		}

		if m.isOpen && len(m.children) > 0 {
			for _, item := range m.children {
				cons, capt := item.MouseHandler()(action, event, setFocus)
				if cons {
					return cons, capt
				}
			}
		}

		// if use this is a top level menu item
		if m.InRect(event.Position()) {
			m.bar.SetFocusedMenu(m)
			switch action {
			case MouseLeftDown:
				if m.parent != nil && len(m.children) > 0 {
					m.Open()
				}

				return true, m
			case MouseLeftUp:
				if len(m.children) > 0 {
					m.Open()

					return true, m
				}

				m.bar.executeMenu(m)

				return true, m
			default:
				if len(m.children) > 0 && m.bar.isActive {
					m.Open()
				}

				return true, m
			}
		} else if m.isOpen {
			m.Close()

			return true, m
		}

		return false, nil
	})
}
