package core

type point struct {
	x, y int
}

func (p point) equals(q point) bool {
	return p.x == q.x && p.y == q.y
}
