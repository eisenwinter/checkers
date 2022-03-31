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

func (b Board) getLeftSideCount() (white int, red int) {
	white = 0
	red = 0
	for i := 0; i <= (height - 1); i++ {
		for j := 0; j < 2; j++ {
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

func (b Board) checkForPattern(pattern []Coordinate, player bool) bool {
	for _, v := range pattern {
		ok, field := b.at(v)
		if !ok {
			return false
		}
		if field.isEmpty() {
			return false
		}
		if field.isWhitePiece() != player {
			return false
		}
	}
	return true
}

func (b Board) getFullSquares() (white int, red int) {
	white = 0
	red = 0
	for i, v := range b {
		if !v.isEmpty() {
			r, c := reverseIndexOf(i)
			coord := Coordinate{r, c}
			shape := []Coordinate{
				coord,
				coord.Shift(-1, +1),
				coord.Shift(+1, +1),
				coord.Shift(-2, +2),
				coord.Shift(+2, +2),
				coord.Shift(0, +2),
				coord.Shift(0, +4),
				coord.Shift(-1, +3),
				coord.Shift(1, +3)}
			if b.checkForPattern(shape, v.isWhitePiece()) {
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

func (b Board) getHalfSquare() (white int, red int) {
	white = 0
	red = 0
	for i, v := range b {
		if !v.isEmpty() {
			r, c := reverseIndexOf(i)
			coord := Coordinate{r, c}
			shape := []Coordinate{
				coord,
				coord.Shift(+2, 0),
				coord.Shift(+1, +1),
				coord.Shift(0, +2),
				coord.Shift(+2, +2),
			}
			if b.checkForPattern(shape, v.isWhitePiece()) {
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

func (b Board) getFullGates() (white int, red int) {
	white = 0
	red = 0
	for i, v := range b {
		if !v.isEmpty() {
			r, c := reverseIndexOf(i)
			coord := Coordinate{r, c}
			shape := []Coordinate{
				coord,
				coord.Shift(+2, 0),
				coord.Shift(0, +2),
				coord.Shift(+2, +2),
			}
			if b.checkForPattern(shape, v.isWhitePiece()) {
				if v.isWhitePiece() {
					white++
				} else {
					red++
				}
			}

			shape = []Coordinate{
				coord,
				coord.Shift(+2, 0),
				coord.Shift(+1, +1),
				coord.Shift(+2, +2),
			}
			if b.checkForPattern(shape, v.isWhitePiece()) {
				if v.isWhitePiece() {
					white++
				} else {
					red++
				}
			}
			shape = []Coordinate{
				coord,
				coord.Shift(+2, 0),
				coord.Shift(+1, +1),
				coord.Shift(0, +2),
			}
			if b.checkForPattern(shape, v.isWhitePiece()) {
				if v.isWhitePiece() {
					white++
				} else {
					red++
				}
			}

			shape = []Coordinate{
				coord,
				coord.Shift(-2, 0),
				coord.Shift(-1, +1),
				coord.Shift(0, +2),
			}
			if b.checkForPattern(shape, v.isWhitePiece()) {
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

func (b Board) getHalfGates() (white int, red int) {
	white = 0
	red = 0
	for i, v := range b {
		if !v.isEmpty() {
			r, c := reverseIndexOf(i)
			coord := Coordinate{r, c}
			shape := []Coordinate{
				coord,
				coord.Shift(+1, +1),
				coord.Shift(+2, +2),
			}
			if b.checkForPattern(shape, v.isWhitePiece()) {
				if v.isWhitePiece() {
					white++
				} else {
					red++
				}
			}

			shape = []Coordinate{
				coord,
				coord.Shift(-1, -1),
				coord.Shift(-2, -2),
			}
			if b.checkForPattern(shape, v.isWhitePiece()) {
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

func (b Board) getPincers() (white int, red int) {
	white = 0
	red = 0
	for i, v := range b {
		if !v.isEmpty() {
			r, c := reverseIndexOf(i)
			if v.isRedPiece() {
				coord := Coordinate{r, c}
				shape := []Coordinate{
					coord,
					coord.Shift(+1, +1),
					coord.Shift(+1, +3),
					coord.Shift(0, +4),
				}
				if b.checkForPattern(shape, v.isWhitePiece()) {
					red++
				}
			} else {
				coord := Coordinate{r, c}
				shape := []Coordinate{
					coord,
					coord.Shift(-1, 1),
					coord.Shift(-1, +3),
					coord.Shift(0, +4),
				}
				if b.checkForPattern(shape, v.isWhitePiece()) {
					white++
				}
			}

		}
	}
	return
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
