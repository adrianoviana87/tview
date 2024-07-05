package tview

type Point struct {
  X, Y int
}

func NewPoint(x, y int) Point {
  return Point{
    x,
    y,
  }
}

type Rectangle struct {
  X, Y, Width, Height int
}

func NewRectangle(x, y, width, height int) Rectangle {
  return Rectangle{x, y, width, height}
}

func NewRectangleFromValues(x, y, width, height int) Rectangle {
  return Rectangle{x, y, width, height}
}

func (r Rectangle) Right() int {
  return r.X + r.Width
}

func (r Rectangle) Bottom() int {
  return r.Y + r.Height
}

func (r Rectangle) Contains(p Point) bool {
  return p.X >= r.X && p.X < r.Right() && p.Y >= r.Y && p.Y < r.Bottom()
}

func (r Rectangle) IsEmpty() bool {
  return r.Width <= 0 || r.Height <= 0
}

func (r Rectangle) Overlaps(other Rectangle) bool {
  return r.X < other.Right() && r.Right() > other.X && r.Y < other.Bottom() && r.Bottom() > other.Y
}
