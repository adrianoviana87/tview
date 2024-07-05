package tview

type Direction int
const (
  DirectionNone Direction = iota
  DirectionUp
  DirectionDown
  DirectionLeft
  DirectionRight
)

type Orientation int
const (
  OrientationHorizontal Orientation = iota
  OrientationVertical
)
