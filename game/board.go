package game

var height = 10
var width = 10

//MoveType is the kind of move
type MoveType uint8

//InvalidMove is not valid
const InvalidMove = 0

//Jump move is a normal movevment
const JumpMove = 1

//SkipMove is a move where a piece gets taken
const SkipMove = 2

func IndexOf(r, c int) int {
	return r*width + c
}

func reverseIndexOf(i int) (int, int) {
	return i / width, i % width
}

//Field is a borad field
type Field uint8

//basic bitflags
const (
	//Empty indicates a empty a field
	Empty Field = 1 << iota
	//Player indicates wheter its a red or white field (true=white)
	Player
	//King inidicates that piece is a king
	King
	//Marked
	Marked
)

//set sets the bitflag supplied
func set(b, flag Field) Field { return b | flag }

//clear removes the bitflag
func clear(b, flag Field) Field { return b &^ flag }

//has checks the bitflag supplied
func has(b, flag Field) bool { return b&flag != 0 }

//IsEmptyField checks if a field is a empty field
func IsEmptyField(f Field) bool {
	return has(f, Empty)
}

//IsPlayer checks if a field is a player field
func IsPlayer(f Field) bool {
	return has(f, Player)
}

func IsKing(f Field) bool {
	return has(f, King)
}

func (f Field) isRedPiece() bool {
	return !has(f, Empty) && !has(f, Player)
}

func (f Field) isEmpty() bool {
	return has(f, Empty)
}

func (f Field) isWhitePiece() bool {
	return !has(f, Empty) && has(f, Player)
}

func (f Field) isMarked() bool {
	return !has(f, Empty) && has(f, Marked)
}

func (f Field) mark() Field {
	return set(f, Marked)
}

func (f Field) isKing() bool {
	return !has(f, Empty) && has(f, King)
}

//Coordinate in the form Row Column
type Coordinate struct {
	Row int
	Col int
}

func (c Coordinate) Shift(rows, cols int) Coordinate {
	return Coordinate{c.Row + rows, c.Col + cols}
}

func (c Coordinate) neighbourhood() []Coordinate {
	n := make([]Coordinate, 0)
	if c.Row > 0 && c.Col > 0 {
		n = append(n, Coordinate{c.Row - 1, c.Col - 1})
	}
	if c.Row > 0 && c.Col < (width-1) {
		n = append(n, Coordinate{c.Row - 1, c.Col + 1})
	}
	if c.Row < height-1 && c.Col < (width-1) {
		n = append(n, Coordinate{c.Row + 1, c.Col + 1})
	}
	if c.Row < height-1 && c.Col > 0 {
		n = append(n, Coordinate{c.Row + 1, c.Col - 1})
	}
	return n
}

type Path struct {
	Coordinates []Coordinate
}

func (c Coordinate) ToIndex() int {
	return IndexOf(c.Row, c.Col)
}

func (c Coordinate) clone() Coordinate {
	return Coordinate{c.Row, c.Col}
}

func (c Coordinate) upwards(pos Coordinate) bool {
	return c.Row > pos.Row
}

func (c Coordinate) downwards(pos Coordinate) bool {
	return c.Row < pos.Row
}

func (c Coordinate) leftwards(pos Coordinate) bool {
	return c.Col > pos.Col
}

func (c Coordinate) rightwards(pos Coordinate) bool {
	return c.Col < pos.Col
}

func northeast(c Coordinate) (bool, Coordinate) {
	return c.northEastOf()
}

func (c Coordinate) northEastOf() (bool, Coordinate) {
	if c.Row > 0 && c.Col < width-1 {
		return true, c.Shift(-1, +1)
	}
	return false, Coordinate{}
}

func northwest(c Coordinate) (bool, Coordinate) {
	return c.northWestOf()
}

func (c Coordinate) northWestOf() (bool, Coordinate) {
	if c.Row > 0 && c.Col > 0 {
		return true, c.Shift(-1, -1)
	}
	return false, Coordinate{}
}

func southeast(c Coordinate) (bool, Coordinate) {
	return c.southEastOf()
}

func (c Coordinate) southEastOf() (bool, Coordinate) {
	if c.Row < (height-1) && c.Col < (width-1) {
		return true, c.Shift(+1, +1)
	}
	return false, Coordinate{}
}

func southwest(c Coordinate) (bool, Coordinate) {
	return c.southWestOf()
}

func (c Coordinate) southWestOf() (bool, Coordinate) {
	if c.Row < (height-1) && c.Col > 0 {
		return true, c.Shift(+1, -1)
	}
	return false, Coordinate{}
}

func coordinateFromIndex(i int) Coordinate {
	r, c := reverseIndexOf(i)
	return Coordinate{r, c}
}

//Move is a game move from a coordinate to another and information if its
//as MUST play move and if a checker is taken
type Move struct {
	From     Coordinate
	To       Coordinate
	Takes    *Coordinate
	Previous *Move
	Depth    int
}

//Board is the basic game board structure
type Board []Field

//canDrawTo check if the move is even `physically` possible
func (b Board) canDrawTo(r int, c int) bool {
	return r >= 0 && c >= 0 && r < height && c < width
}

//allPiecesFor gets all remaining pieces of a the player
func (b Board) allPiecesFor(player bool) []Coordinate {
	p := make([]Coordinate, 0)
	for i, v := range b {
		if !has(v, Empty) && has(v, Player) == player {
			r, c := reverseIndexOf(i)
			p = append(p, Coordinate{r, c})
		}
	}
	return p
}

func (b Board) copy() Board {
	nextBoard := make(Board, len(b))
	copy(nextBoard, b)
	return nextBoard
}

//playable indicates the board is still playable
func (b Board) playable() bool {
	w, r, wk, rk := b.getCounts()
	if w == 0 {
		return false
	}
	if r == 0 {
		return false
	}
	//draw
	if w == 1 && r == 1 && wk == 1 && rk == 1 {
		return false
	}

	return true
}

func (b Board) at(pos Coordinate) (bool, Field) {
	if pos.Col >= 0 && pos.Col < width && pos.Row >= 0 && pos.Row < height {
		return true, b[IndexOf(pos.Row, pos.Col)]
	}
	return false, 0
}

func (b Board) must(pos Coordinate) Field {
	return b[IndexOf(pos.Row, pos.Col)]
}

func (b Board) topLeftOf(pos Coordinate) (bool, Field, Coordinate) {
	p := pos.Shift(-1, -1)
	ok, field := b.at(p)
	return ok, field, p
}

func (b Board) bottomLeftOf(pos Coordinate) (bool, Field, Coordinate) {
	p := pos.Shift(+1, -1)
	ok, field := b.at(p)
	return ok, field, p
}

func (b Board) topRightOf(pos Coordinate) (bool, Field, Coordinate) {
	p := pos.Shift(-1, +1)
	ok, field := b.at(p)
	return ok, field, p
}

func (b Board) bottomRightOf(pos Coordinate) (bool, Field, Coordinate) {
	p := pos.Shift(+1, +1)
	ok, field := b.at(p)
	return ok, field, p
}

func (b Board) isEmptyField(pos Coordinate) bool {
	if ok, f := b.at(pos); ok && f.isEmpty() {
		return true
	}
	return false
}

//removePiece removes the piece at the given position
func (b Board) removePiece(pos Coordinate) {
	b[IndexOf(pos.Row, pos.Col)] = set(b[IndexOf(pos.Row, pos.Col)], Empty)
	b[IndexOf(pos.Row, pos.Col)] = clear(b[IndexOf(pos.Row, pos.Col)], Player)
	b[IndexOf(pos.Row, pos.Col)] = clear(b[IndexOf(pos.Row, pos.Col)], King)
}

//promoteToKing promotes the field to king
func (b Board) promoteToKing(pos Coordinate) {
	b[IndexOf(pos.Row, pos.Col)] = set(b[IndexOf(pos.Row, pos.Col)], King)
}

//isBoardEnd checks if its the board end
func (b Board) isBoardEnd(r int, player bool) bool {
	if player && r == 0 {
		return true
	}
	if !player && r == (height-1) {
		return true
	}
	return false
}

//movePiece moves a piece on the board
func (b Board) movePiece(from, to Coordinate, player bool) {
	isKing := b.must(from).isKing()
	b[IndexOf(from.Row, from.Col)] = set(b[IndexOf(from.Row, from.Col)], Empty)
	b[IndexOf(from.Row, from.Col)] = clear(b[IndexOf(from.Row, from.Col)], King)
	b[IndexOf(from.Row, from.Col)] = clear(b[IndexOf(from.Row, from.Col)], Player)

	b[IndexOf(to.Row, to.Col)] = clear(b[IndexOf(to.Row, to.Col)], Empty)
	if player {
		b[IndexOf(to.Row, to.Col)] = set(b[IndexOf(to.Row, to.Col)], Player)
	} else {
		b[IndexOf(to.Row, to.Col)] = clear(b[IndexOf(to.Row, to.Col)], Player)
	}
	if isKing {
		b[IndexOf(to.Row, to.Col)] = set(b[IndexOf(to.Row, to.Col)], King)
	}
}

func (b Board) applyMove(m Move, player bool) bool {
	kingPromoted := false
	if b.canDrawTo(m.To.Row, m.To.Col) {
		if m.Takes != nil {
			b.removePiece(*m.Takes)
		}
		b.movePiece(m.From, m.To, player)
		if b.isBoardEnd(m.To.Row, player) && !b.must(m.To).isKing() {
			kingPromoted = true
		}
	}
	return kingPromoted
}

//getMoveType returns if the move would be valid in terms of gameplay and returns the move type and row and column
func (b Board) getMoveType(from, to Coordinate, player bool) (MoveType, Coordinate) {
	ok, f := b.at(from)
	if !ok {
		return InvalidMove, Coordinate{}
	}
	ok, t := b.at(to)
	if !ok {
		return InvalidMove, Coordinate{}
	}
	if (t.isEmpty() || t.isMarked()) && player && (from.upwards(to) || f.isKing()) {
		return JumpMove, to
	}
	if (t.isEmpty() || t.isMarked()) && !player && (from.downwards(to) || f.isKing()) {
		return JumpMove, to
	}

	if from.leftwards(to) && from.upwards(to) {
		ok, tl, ctl := b.topLeftOf(to)
		//the not marked on the field is for the rule: are not removed during the move, they are removed only after the entire multi-jump move is complete
		if ok && !t.isEmpty() && t.isWhitePiece() != player && tl.isEmpty() && !tl.isMarked() && !t.isMarked() {
			return SkipMove, ctl
		}

	}
	if from.leftwards(to) && from.downwards(to) {
		ok, bl, cbl := b.bottomLeftOf(to)
		if ok && !t.isEmpty() && t.isWhitePiece() != player && bl.isEmpty() && !bl.isMarked() && !t.isMarked() {
			return SkipMove, cbl
		}
	}
	if from.rightwards(to) && from.upwards(to) {
		ok, tr, ctr := b.topRightOf(to)
		if ok && !t.isEmpty() && t.isWhitePiece() != player && tr.isEmpty() && !tr.isMarked() && !t.isMarked() {
			return SkipMove, ctr
		}
	}
	if from.rightwards(to) && from.downwards(to) {
		ok, br, cbr := b.bottomRightOf(to)
		if ok && !t.isEmpty() && t.isWhitePiece() != player && br.isEmpty() && !br.isMarked() && !t.isMarked() {
			return SkipMove, cbr
		}
	}
	return InvalidMove, Coordinate{}
}

func (b Board) lineOfSightSkip(dir func(Coordinate) (bool, Coordinate), pos Coordinate, player bool, prev *Move) []Move {
	m := make([]Move, 0)
	for ok, current := dir(pos); ok; ok, current = dir(current) {
		if mt, cord := b.getMoveType(pos, current, player); mt != InvalidMove {
			move := Move{pos, cord, nil, nil, 0}
			if prev != nil {
				move.Previous = prev
				move.Depth = prev.Depth + 1
			}
			if mt == SkipMove {
				tmp := current.clone()
				move.Takes = &tmp
				m = append(m, move)

				nextBoard := boardForNextSkip(b, pos, cord, *move.Takes, player)
				nextMoves := nextBoard.getPossibleSkipsFor(cord, player, &move)
				for _, v := range nextMoves {
					m = append(m, v)
				}
				return m
			}
		}
	}
	return m
}

//getPossibleSkipsFor returns all possible skips (take moves)
func (b Board) getPossibleSkipsFor(pos Coordinate, player bool, prev *Move) []Move {
	//Todo this needs board copies so it doesnt jump the pieces all ofer again
	m := make([]Move, 0)
	ok, f := b.at(pos)
	if !ok {
		return m
	}
	if !f.isEmpty() {
		//piece does not belong to player
		if f.isWhitePiece() != player {
			return m
		}
		if f.isKing() {
			nw := b.lineOfSightSkip(northwest, pos, player, prev)
			m = append(m, nw...)

			ne := b.lineOfSightSkip(northeast, pos, player, prev)
			m = append(m, ne...)

			se := b.lineOfSightSkip(southeast, pos, player, prev)
			m = append(m, se...)

			sw := b.lineOfSightSkip(southwest, pos, player, prev)
			m = append(m, sw...)
		} else {
			nbs := pos.neighbourhood()
			for _, c := range nbs {
				if mt, cord := b.getMoveType(pos, c, player); mt == SkipMove {
					depth := 0
					if prev != nil {
						depth = prev.Depth + 1
					}
					tmp := c.clone()
					move := Move{pos, cord, &tmp, prev, depth}
					nextBoard := boardForNextSkip(b, pos, cord, *move.Takes, player)
					nextMoves := nextBoard.getPossibleSkipsFor(cord, player, &move)
					for _, v := range nextMoves {
						m = append(m, v)
					}
				}
			}
		}
	}
	return m
}

func (b Board) lineOfSightMoves(dir func(Coordinate) (bool, Coordinate), pos Coordinate, player bool) []Move {
	m := make([]Move, 0)
	for ok, current := dir(pos); ok; ok, current = dir(current) {
		if mt, cord := b.getMoveType(pos, current, player); mt != InvalidMove {
			move := Move{pos, cord, nil, nil, 0}
			if mt == SkipMove {
				tmp := current.clone()
				move.Takes = &tmp
				m = append(m, move)

				nextBoard := boardForNextSkip(b, pos, cord, *move.Takes, player)
				nextMoves := nextBoard.getPossibleSkipsFor(cord, player, &move)
				for _, v := range nextMoves {
					m = append(m, v)
				}
				return m
			}
			m = append(m, move)
		} else {
			return m
		}
	}
	return m
}

func boardForNextSkip(b Board, from, to, taken Coordinate, player bool) Board {
	nextBoard := b.copy()
	nextBoard[IndexOf(taken.Row, taken.Col)] = nextBoard.must(taken).mark()
	nextBoard.movePiece(from, to, player)
	return nextBoard

}

//getPossibleMoves returns any possible moves for that field
func (b Board) getPossibleMoves(pos Coordinate, player bool) []Move {
	m := make([]Move, 0)
	ok, f := b.at(pos)
	if !ok {
		return m
	}
	if !f.isEmpty() {
		//piece does not belong to player
		if f.isWhitePiece() != player {
			return m
		}
		if f.isKing() {
			nw := b.lineOfSightMoves(northwest, pos, player)
			m = append(m, nw...)

			ne := b.lineOfSightMoves(northeast, pos, player)
			m = append(m, ne...)

			se := b.lineOfSightMoves(southeast, pos, player)
			m = append(m, se...)

			sw := b.lineOfSightMoves(southwest, pos, player)
			m = append(m, sw...)
		} else {
			nbs := pos.neighbourhood()
			for _, c := range nbs {
				if mt, cord := b.getMoveType(pos, c, player); mt != InvalidMove {
					move := Move{pos, cord, nil, nil, 0}
					if mt == SkipMove {
						tmp := c.clone()
						move.Takes = &tmp
					}
					m = append(m, move)
					if mt == SkipMove {
						nextBoard := boardForNextSkip(b, pos, cord, *move.Takes, player)
						nextMoves := nextBoard.getPossibleSkipsFor(cord, player, &move)
						for _, v := range nextMoves {
							m = append(m, v)
						}
					}
				}
			}
		}
	}
	return m
}

//filterMoves prunes any non must moves when must moves are in the list
func filterMoves(move []Move) []Move {
	highestDepth := 0
	take := false
	for _, v := range move {
		highestDepth = maxOf(v.Depth, highestDepth)
		if v.Takes != nil {
			take = true
		}
	}
	filtered := make([]Move, 0)
	for _, v := range move {
		if highestDepth == v.Depth {
			if !take || (take && v.Takes != nil) {
				filtered = append(filtered, v)
			}
		}
	}
	return filtered
}

// getCounts retruns white piece count, red piece count, white king count and red king count
func (b Board) getCounts() (white int, red int, wking int, rking int) {
	white = 0
	red = 0
	wking = 0
	rking = 0
	for _, v := range b {
		if !has(v, Empty) {
			if has(v, Player) {
				white++
				if has(v, King) {
					wking++
				}
			} else {
				red++
				if has(v, King) {
					rking++
				}
			}

		}
	}
	return white, red, wking, rking
}

// boardSetup creates a starting board
func boardSetup(board Board) Board {
	for i := range board {
		board[i] = set(0, Empty)
	}

	for i := 0; i < 20; i++ {
		r, _ := reverseIndexOf(i * 2)
		if (r+1)%2 == 0 {
			board[i*2] = clear(board[i*2], Empty)
		} else {
			board[i*2+1] = clear(board[i*2+1], Empty)
		}
	}

	for i := IndexOf(6, 0); i < (IndexOf(6, 0) + 20*2); i = i + 2 {
		r, _ := reverseIndexOf(i)
		if (r+1)%2 == 0 {
			board[i] = set(clear(board[i], Empty), Player)
		} else {
			board[i+1] = set(clear(board[i+1], Empty), Player)
		}

	}

	return board
}
