package game

type BoardDiff struct {
	Fig
	Bef Fig
	Pos
}

func (b *Board) Diff(o *Board) []BoardDiff {
	d := make([]BoardDiff, 0, 5)
	var oac ACP
	var b, a *Square
	var oa Pos
	var c *BoardDiff
	for oac.OK() {
		oa = Pos(oac)
		if b, a = b.GPos(oa), o.GPos(oa); a != b {
			c = new(BoardDiff)
			c.Pos = oa
			if a.Empty() {
				c.Fig = Fig{0, 0}
			} else {
				c.Fig = a.Fig
			}
			if b.Empty() {
				c.Bef = Fig{0, 0}
			} else {
				c.Bef = b.Fig
			}
			d = append(d, c)
		}
		oac.P()
	}
	return d
}
