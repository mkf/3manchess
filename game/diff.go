package game

//© Copyright 2015-2016 Michał Krzysztof Feiler & Paweł Zacharek

//BoardDiff describes a single difference between boards
type BoardDiff struct {
	Fig
	Bef Fig
	Pos
}

//Diff returns differences between Boards
func (b *Board) Diff(o *Board) []BoardDiff {
	d := make([]BoardDiff, 0, 5)
	var p, a *Square
	var c *BoardDiff
	for _, oa := range ALLPOS {
		if p, a = b.GPos(oa), o.GPos(oa); a != p {
			c = new(BoardDiff)
			c.Pos = oa
			if a.Empty() {
				c.Fig = Fig{0, 0, false}
			} else {
				c.Fig = a.Fig
			}
			if p.Empty() {
				c.Bef = Fig{0, 0, false}
			} else {
				c.Bef = p.Fig
			}
			d = append(d, *c)
		}
	}
	return d
}
