package game

func (b Board) getGoldenStoneCount() (white int, red int) {
	white = 0
	red = 0
	if b.must(Coordinate{0, 5}).isRedPiece() {
		red = 1
	}
	if b.must(Coordinate{height - 1, 4}).isWhitePiece() {
		white = 1
	}
	return
}

//leggards and grapes
func (b Board) getLeggardAndGrapeCount() (white int, red int) {
	white = 0
	red = 0
	for i, v := range b {
		if !v.isEmpty() {
			r, c := reverseIndexOf(i)
			coord := Coordinate{r, c}

			nwok, nwc := coord.northWestOf()
			neok, nec := coord.northEastOf()
			seok, sec := coord.southEastOf()
			swok, swc := coord.southWestOf()

			northWest := (!nwok || (nwok && b.must(nwc).isEmpty()))
			northEast := (!neok || (neok && b.must(nec).isEmpty()))
			southEast := (!seok || (seok && b.must(sec).isEmpty()))
			southWest := (!swok || (swok && b.must(swc).isEmpty()))

			emptySquare := 0
			blockedSquares := 0
			if northEast {
				emptySquare++
				if !neok {
					blockedSquares++
				}
			}
			if northWest {
				emptySquare++
				if !nwok {
					blockedSquares++
				}
			}
			if southEast {
				emptySquare++
				if seok {
					blockedSquares++
				}
			}
			if southWest {
				emptySquare++
				if swok {
					blockedSquares++
				}
			}

			if emptySquare == 3 && blockedSquares == 0 {
				if v.isWhitePiece() {
					white++
				} else {
					red++
				}
			} else if emptySquare == 4 && blockedSquares >= 1 {
				if v.isWhitePiece() {
					white++
				} else {
					red++
				}
			}
		}
	}
	return
}

func (b Board) getLeftSideCount() (white int, red int) {
	white = 0
	red = 0
	for i := 0; i <= (height - 1); i++ {
		for j := 0; j < 3; j++ {
			idx := IndexOf(i, j)
			if !has(b[idx], Empty) {
				if has(b[idx], Player) {
					white++
				} else {
					red++
				}
			}
		}
	}
	return
}

func (b Board) getMiddleCount() (white int, red int) {
	white = 0
	red = 0
	for i := 0; i <= (height - 1); i++ {
		for j := 3; j < width-3; j++ {
			idx := IndexOf(i, j)
			if !has(b[idx], Empty) {
				if has(b[idx], Player) {
					white++
				} else {
					red++
				}
			}
		}
	}
	return
}

func (b Board) getRightSideCount() (white int, red int) {
	white = 0
	red = 0
	for i := 0; i <= (height - 1); i++ {
		for j := width - 3; j < width; j++ {
			idx := IndexOf(i, j)
			if !has(b[idx], Empty) {
				if has(b[idx], Player) {
					white++
				} else {
					red++
				}
			}
		}
	}
	return
}

func (b Board) getFullSquares() (white int, red int) {
	white = 0
	red = 0
	mask := Mask{
		{E, E, P, E, E},
		{E, P, E, P, E},
		{P, E, P, E, P},
		{E, P, E, P, E},
		{E, E, P, E, E},
	}

	white = findMask(b, mask, true, true)
	red = findMask(b, mask, false, true)
	return
}

func (b Board) getHalfSquare() (white int, red int) {
	white = 0
	red = 0
	mask := Mask{
		{E, E, E, E, E},
		{E, P, E, P, E},
		{E, E, P, E, E},
		{E, P, E, P, E},
		{E, E, E, E, E},
	}

	white = findMask(b, mask, true, false)
	red = findMask(b, mask, false, false)
	return
}

func (b Board) getFullGates() (white int, red int) {
	white = 0
	red = 0
	lmask := Mask{
		{P, E, P},
		{E, P, E},
		{P, E, E},
	}
	rmask := Mask{
		{P, E, P},
		{E, P, E},
		{E, E, P},
	}

	white = findMask(b, lmask, true, false) + findMask(b, rmask, true, false)
	red = findMask(b, lmask.FlipHorizontally(), false, false) + findMask(b, rmask.FlipHorizontally(), false, false)
	return
}

func (b Board) getHalfGates() (white int, red int) {
	white = 0
	red = 0
	mask := Mask{
		{E, E, P},
		{E, P, E},
		{P, E, E},
	}
	maskFlipped := mask.FlipHorizontally()
	white = findMask(b, mask, true, false) + findMask(b, maskFlipped, true, false)
	red = findMask(b, mask, false, false) + findMask(b, maskFlipped, false, false)
	return
}

func (b Board) getPincers() (white int, red int) {
	white = 0
	red = 0
	mask := Mask{
		{E, P, E, P, E},
		{P, E, E, E, P},
	}
	white = findMask(b, mask, true, false)
	red = findMask(b, mask.FlipHorizontally(), false, false)
	return
}

type MaskElement bool

const E MaskElement = false
const P MaskElement = true

type Mask [][]MaskElement

func (m Mask) FlipHorizontally() Mask {
	l := len(m)
	new := make(Mask, l)
	for j, v := range m {
		new[(l-1)-j] = v

	}
	return new
}

func (m Mask) XOR(other Mask) bool {
	match := false
	for i, v := range m {
		for j, vc := range v {
			match = (match || (other[i][j] != vc))
		}
	}
	return match
}

func (m Mask) AlphaMatch(other Mask) bool {
	for i, v := range m {
		for j, vc := range v {
			if bool(vc) && !bool(other[i][j]) {
				return false
			}

		}
	}
	return true
}

func findMask(board Board, mask Mask, player bool, alpha bool) int {
	getMask := func(i, j, w, h int) Mask {
		new := make(Mask, h)
		ix := 0
		for im := i; im < (i + h); im++ {
			new[ix] = make([]MaskElement, w)
			jx := 0
			for jm := j; jm < (j + w); jm++ {
				new[ix][jx] = MaskElement(!board[IndexOf(im, jm)].isEmpty() && board[IndexOf(im, jm)].isWhitePiece() == player)
				jx++
			}
			ix++
		}
		return new
	}
	count := 0
	for i := 0; i <= height-len(mask); i++ {
		cols := len(mask[0])
		for j := 0; j <= width-cols; j++ {
			if alpha {
				if mask.AlphaMatch(getMask(i, j, cols, len(mask))) {
					count++
				}
			} else {
				if !getMask(i, j, cols, len(mask)).XOR(mask) {
					count++
				}
			}

		}
	}
	return count
}

//getMiddleBoxCount gets the number of pieces in the middle box
func (b Board) getMiddleBoxCount() (white int, red int) {
	white = 0
	red = 0
	middleRow := (height / 2) - 1
	for i := middleRow; i <= (middleRow + 1); i++ {
		//width without left and right side  -> enemy can only pass on the side
		for j := 2; j <= (width - 3); j++ {
			idx := IndexOf(i, j)
			if !has(b[idx], Empty) {
				if has(b[idx], Player) {
					white++
				} else {
					red++
				}
			}
		}
	}
	return
}

func (b Board) hasKingInLineOfSight(pos Coordinate, player bool) (bool, Direction) {
	next, c := pos.northWestOf()
	for next {
		field := b.must(c)
		if !field.isEmpty() && field.isWhitePiece() == player && field.isKing() {
			return true, DirectionNortWest
		}
		next, c = c.northWestOf()
	}
	next, c = pos.northEastOf()
	for next {
		field := b.must(c)
		if !field.isEmpty() && field.isWhitePiece() == player && field.isKing() {
			return true, DirectionNortEast
		}
		next, c = c.northEastOf()
	}

	next, c = pos.southEastOf()
	for next {
		field := b.must(c)
		if !field.isEmpty() && field.isWhitePiece() == player && field.isKing() {
			return true, DirectionSouthEast
		}
		next, c = c.southEastOf()
	}

	next, c = pos.southWestOf()
	for next {
		field := b.must(c)
		if !field.isEmpty() && field.isWhitePiece() == player && field.isKing() {
			return true, DirectionSouthWest
		}
		next, c = c.southWestOf()
	}
	return false, DirectionNortEast
}

func canBeTakenByPawn(b Board, p Coordinate, player bool) bool {
	sek, se := p.southEastOf()
	swk, sw := p.southWestOf()
	nek, ne := p.northEastOf()
	nwk, nw := p.northWestOf()

	if sek && nwk {
		if b.must(se).isEmpty() && !b.must(nw).isEmpty() && b.must(nw).isWhitePiece() != player {
			return true
		}
		if b.must(nw).isEmpty() && !b.must(se).isEmpty() && b.must(se).isWhitePiece() != player {
			return true
		}
	}

	if swk && nek {
		if b.must(sw).isEmpty() && !b.must(ne).isEmpty() && b.must(ne).isWhitePiece() != player {
			return true
		}
		if b.must(ne).isEmpty() && !b.must(sw).isEmpty() && b.must(sw).isWhitePiece() != player {
			return true
		}
	}
	return false
}

func pieceCanBeTaken(b Board, p Coordinate, player bool) bool {
	ok, _ := b.at(p)
	if !ok {
		return false
	}
	k, dir := b.hasKingInLineOfSight(p, !player)
	if k {
		switch dir {
		case DirectionNortEast:
			if k, a := b.at(p.Shift(+1, -1)); k {
				if a.isEmpty() {
					return true
				}
			}
		case DirectionNortWest:
			if k, a := b.at(p.Shift(+1, +1)); k {
				if a.isEmpty() {
					return true
				}
			}
		case DirectionSouthEast:
			if k, a := b.at(p.Shift(-1, -1)); k {
				if a.isEmpty() {
					return true
				}
			}
		case DirectionSouthWest:
			if k, a := b.at(p.Shift(-1, +1)); k {
				if a.isEmpty() {
					return true
				}
			}
		}
	}

	return canBeTakenByPawn(b, p, player)
}

func (b Board) getVulnerablePiecesCount() (white int, red int) {
	white = 0
	red = 0

	wskipps := b.getAllPossibleSkips(true)
	for _, s := range wskipps {
		red += (1 + s.Depth)
	}
	rskipps := b.getAllPossibleSkips(false)
	for _, s := range rskipps {
		white += white + (1 + s.Depth)
	}
	return
}

func (b Board) getSuicidalPiecesCount() (white int, red int) {
	white = 0
	red = 0

	whiteMoves := b.getPossibleValidMovesForPlayer(true)
	for _, s := range whiteMoves {
		tmp := b.copy()
		unrollMove(&tmp, s, true, s.Depth)
		enemy := b.getAllPossibleSkips(false)
		for _, s := range enemy {
			white += (1 + s.Depth)
		}
	}
	rskipps := b.getPossibleValidMovesForPlayer(false)
	for _, s := range rskipps {
		tmp := b.copy()
		unrollMove(&tmp, s, true, s.Depth)
		enemy := b.getAllPossibleSkips(true)
		for _, s := range enemy {
			red += (1 + s.Depth)
		}
	}
	return
}

// getProtectedPieceCount returns the count of protected pieces for the heuristic
func (b Board) getProtectionCount() (white int, red int) {
	white = 0
	red = 0
	for i, v := range b {
		if !has(v, Empty) {
			r, c := reverseIndexOf(i)
			if has(v, Player) {
				//white
				if b.canDrawTo(r-1, c-1) &&
					!has(b[IndexOf(r-1, c-1)], Empty) &&
					has(b[IndexOf(r-1, c-1)], Player) {
					white++
				} else if b.canDrawTo(r-1, c+1) &&
					!has(b[IndexOf(r-1, c+1)], Empty) &&
					has(b[IndexOf(r-1, c+1)], Player) {
					white++
				} else if b.canDrawTo(r+1, c+1) &&
					!has(b[IndexOf(r+1, c+1)], Empty) &&
					has(b[IndexOf(r+1, c+1)], Player) {
					white++
				} else if b.canDrawTo(r+1, c-1) &&
					!has(b[IndexOf(r+1, c-1)], Empty) &&
					has(b[IndexOf(r+1, c-1)], Player) {
					white++
				} else if k, _ := b.hasKingInLineOfSight(Coordinate{r, c}, true); k {
					white++
				}

			} else {
				//red
				if b.canDrawTo(r-1, c-1) &&
					!has(b[IndexOf(r-1, c-1)], Empty) &&
					!has(b[IndexOf(r-1, c-1)], Player) {
					red++
				} else if b.canDrawTo(r-1, c+1) &&
					!has(b[IndexOf(r-1, c+1)], Empty) &&
					!has(b[IndexOf(r-1, c+1)], Player) {
					red++
				} else if b.canDrawTo(r+1, c+1) &&
					!has(b[IndexOf(r+1, c+1)], Empty) &&
					!has(b[IndexOf(r+1, c+1)], Player) {
					red++
				} else if b.canDrawTo(r+1, c-1) &&
					!has(b[IndexOf(r+1, c-1)], Empty) &&
					!has(b[IndexOf(r+1, c-1)], Player) {
					red++
				} else if k, _ := b.hasKingInLineOfSight(Coordinate{r, c}, false); k {
					red++
				}
			}
		}
	}
	return
}

//getStuckPiecesCount gets the number of pieces in the stuck piecies heustric
func (b Board) getStuckPiecesCount() (white int, red int, wking int, rking int) {
	white = 0
	red = 0
	wking = 0
	rking = 0
	for i, v := range b {
		if !has(v, Empty) {
			r, c := reverseIndexOf(i)
			if !has(v, King) && has(v, Player) &&
				(!b.canDrawTo(r-1, c-1) || !has(b[IndexOf(r-1, c-1)], Empty)) &&
				(!b.canDrawTo(r-1, c+1) || !has(b[IndexOf(r-1, c+1)], Empty)) {
				white++

			} else if !has(v, King) && !has(v, Player) &&
				(!b.canDrawTo(r+1, c-1) || !has(b[IndexOf(r+1, c-1)], Empty)) &&
				(!b.canDrawTo(r+1, c+1) || !has(b[IndexOf(r+1, c+1)], Empty)) {
				red++

			}
			if has(v, King) &&
				(!b.canDrawTo(r-1, c-1) || !has(b[IndexOf(r-1, c-1)], Empty)) &&
				(!b.canDrawTo(r-1, c+1) || !has(b[IndexOf(r-1, c+1)], Empty)) &&
				(!b.canDrawTo(r+1, c-1) || !has(b[IndexOf(r+1, c-1)], Empty)) &&
				(!b.canDrawTo(r+1, c+1) || !has(b[IndexOf(r+1, c+1)], Empty)) {
				if has(v, Player) {
					wking++
				} else {
					rking++
				}
			}

		}
	}
	return white, red, wking, rking
}

func neighbourhoodCount(b Board, c Coordinate, player bool) int {
	ok, f := b.at(c)
	if !ok {
		return 0
	}
	if !f.isEmpty() && !f.isMarked() && f.isWhitePiece() == player {
		b[IndexOf(c.Row, c.Col)] = b[IndexOf(c.Row, c.Col)].mark()
		return 1 + neighbourhoodCount(b, c.Shift(1, 1), player) + neighbourhoodCount(b, c.Shift(1, -1), player) + neighbourhoodCount(b, c.Shift(-1, 1), player) + neighbourhoodCount(b, c.Shift(-1, -1), player)
	}
	return 0

}

func (b Board) getLargestConnectedField() (white int, red int) {
	white = 0
	red = 0
	tmp := b.copy()
	for i, v := range b {
		if !has(v, Empty) {
			r, c := reverseIndexOf(i)
			if has(v, Player) {
				w := neighbourhoodCount(tmp, Coordinate{r, c}, true)
				white = maxOf(white, w)
			} else {
				r := neighbourhoodCount(tmp, Coordinate{r, c}, false)
				red = maxOf(red, r)
			}
		}
	}
	return
}
