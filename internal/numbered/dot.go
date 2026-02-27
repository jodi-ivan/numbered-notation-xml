package numbered

type DotPosition struct {
	beforeXpos int
	afterXPos  int
	Address    []*int
}

func (dt *DotPosition) Reset(startPosition int) {
	dt.beforeXpos = startPosition
	dt.Address = []*int{}
}

func (dt *DotPosition) Render(endPosition int) {
	if len(dt.Address) == 0 {
		return
	}
	dt.afterXPos = endPosition
	space := (dt.afterXPos - dt.beforeXpos) / (len(dt.Address) + 1)
	for i, d := range dt.Address {
		*d = (dt.beforeXpos + (space * (i + 1)))
	}

	// reset here
	dt.beforeXpos = endPosition
	dt.Address = []*int{}
}
