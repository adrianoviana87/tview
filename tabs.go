package tview

import (
	"github.com/gdamore/tcell/v2"
)

type Tabs struct {
	flex              *Flex
	buttonFlex        *Flex
	tabContentButtons map[string]*Button
	tabCloseButtons   map[string]*Button
	contents          map[string]Primitive
	contentContainer  *Flex
}

func NewTabs() *Tabs {
	contentContainer := NewFlex().
		SetDirection(FlexColumnCSS)
	buttonFlex := NewFlex().
		SetDirection(FlexRowCSS)
	flex := NewFlex().
		SetDirection(FlexColumnCSS).
		AddItem(contentContainer, 0, 1, false).
		AddItem(buttonFlex, 1, 0, false)

	t := &Tabs{
		flex:              flex,
		buttonFlex:        buttonFlex,
		contents:          make(map[string]Primitive),
		tabContentButtons: make(map[string]*Button),
		tabCloseButtons:   make(map[string]*Button),
		contentContainer:  contentContainer,
	}
	return t
}

func (t *Tabs) AddTab(name string, p Primitive) *Tabs {
	tabButton := NewButton(name).SetSelectedFunc(func() {
		t.ShowTab(name)
	})
	tabButton.SetStyle(tabButton.style.Background(tcell.ColorNone))

	tabCloseButton := NewButton("X").SetSelectedFunc(func() {
		t.RemoveTab(name)
	})
	tabCloseButton.SetStyle(tabCloseButton.
		style.Background(tcell.ColorNone).
		Foreground(tcell.ColorDarkRed))

	t.tabContentButtons[name] = tabButton
	t.tabCloseButtons[name] = tabCloseButton
	t.contents[name] = p
	width := TaggedStringWidth(name)
	t.buttonFlex.AddItem(NewFlex().SetDirection(FlexRowCSS).
		AddItem(tabButton, TaggedStringWidth(name), 0, false).
		AddItem(tabCloseButton, 1, 0, false),
		width+2, 0, false)

	return t
}

func (t *Tabs) getCurrentTabContent() Primitive {
  if t.contentContainer.GetItemCount() == 0 {
    return nil
  }

  return t.contentContainer.GetItem(0)
}

func (t *Tabs) RemoveTab(name string) *Tabs {
	itemToRemove := func() Primitive {
		for _, item := range t.buttonFlex.items {
			for _, subItem := range item.Item.(*Flex).items {
				if subItem.Item.(*Button).GetLabel() == name {
					return item.Item
				}
			}
		}

		return nil
	}()
	t.tabCloseButtons[name] = nil
	t.tabContentButtons[name] = nil
	t.buttonFlex.RemoveItem(itemToRemove)
	if currentItem := t.getCurrentTabContent(); currentItem == t.contents[name] {
    for key := range t.contents {
      t.ShowTab(key)
      break
    }
	}
	t.contents[name] = nil
	if len(t.contents) == 0 {
	  t.contentContainer.Clear()
	}
	return t
}

func (t *Tabs) ShowTab(tabName string) {
	if t.contentContainer.GetItemCount() > 0 {
		t.contentContainer.GetItem(0).Blur()
	}

	t.contentContainer.Clear()
	content := t.FindTab(tabName)

	t.contentContainer.AddItem(content, 0, 1, true)
	content.Focus(nil)
}

func (t *Tabs) Exists(tabName string) bool {
	return t.contents[tabName] != nil
}

func (t *Tabs) FindTab(tabName string) Primitive {
	return t.contents[tabName]
}

func (t *Tabs) Draw(screen tcell.Screen) {
	t.flex.Draw(screen)
}

func (t *Tabs) GetRect() (int, int, int, int) {
	return t.flex.GetRect()
}

func (t *Tabs) SetRect(x, y, width, height int) {
	t.flex.SetRect(x, y, width, height)
}

func (t *Tabs) InputHandler() func(event *tcell.EventKey, setFocus func(p Primitive)) {
	return t.flex.InputHandler()
}

func (t *Tabs) Focus(delegate func(p Primitive)) {
	t.flex.Focus(delegate)
}

func (t *Tabs) HasFocus() bool {
	return t.flex.HasFocus()
}

func (t *Tabs) Blur() {
	t.flex.Blur()
}

func (t *Tabs) MouseHandler() func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
	return t.flex.MouseHandler()
}

func (t *Tabs) PasteHandler() func(text string, setFocus func(p Primitive)) {
	return t.flex.PasteHandler()
}
