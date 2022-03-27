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

//Coordinate in the form Row Column
type Coordinate struct {
	Row int
	Col int
}

//Move is a game move from a coordinate to another and information if its
//as MUST play move and if a checker is taken
type Move struct {
	FromRow int
	FromCol int

	ToRow    int
	ToCol    int
	Must     bool
	Takes    *Coordinate
	Previous *Move
}

//Board is the basic game board structure
type Board []Field

//canDrawTo check if the move is even `physically` possible
func (b Board) canDrawTo(r int, c int) bool {
	if r >= 0 && c >= 0 && r < height && c < width {
		idx := IndexOf(r, c)
		return idx > 0 && idx < len(b)
	}
	return false
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

//playable indicates the board is still playable
func (b Board) playable() bool {
	w, r, _, _ := b.getCounts()
	if w == 0 {
		return false
	}
	if r == 0 {
		return false
	}
	return true
}

//movePiece moves a piece on the board
func (b Board) movePiece(fromRow, fromCol, toRow, toCol int, player bool) {
	isKing := has(b[IndexOf(fromRow, fromCol)], King)
	b[IndexOf(fromRow, fromCol)] = set(b[IndexOf(fromRow, fromCol)], Empty)
	b[IndexOf(fromRow, fromCol)] = clear(b[IndexOf(fromRow, fromCol)], King)
	b[IndexOf(toRow, toCol)] = clear(b[IndexOf(fromRow, fromCol)], Empty)
	if player {
		b[IndexOf(fromRow, fromCol)] = clear(b[IndexOf(fromRow, fromCol)], Player)
	} else {
		b[IndexOf(fromRow, fromCol)] = clear(b[IndexOf(fromRow, fromCol)], Player)
	}
	if isKing {
		b[IndexOf(toRow, toCol)] = set(b[IndexOf(toRow, toCol)], King)
	}
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

func (b Board) applyMove(m Move, player bool) bool {
	pieceTaken := false
	if b.canDrawTo(m.ToRow, m.ToCol) {
		if m.Must && m.Takes != nil {
			b.removePiece(m.Takes.Row, m.Takes.Col)
			pieceTaken = true
		}
		b.movePiece(m.FromRow, m.FromCol, m.ToRow, m.ToCol, player)
		if b.isBoardEnd(m.ToRow, player) && !b.isKing(m.ToRow, m.ToCol) {
			b.promoteToKing(m.ToRow, m.ToCol)
		}
	}
	return pieceTaken
}

//removePiece removes the piece at the given position
func (b Board) removePiece(r, c int) {
	b[IndexOf(r, c)] = set(b[IndexOf(r, c)], Empty)
	b[IndexOf(r, c)] = clear(b[IndexOf(r, c)], Player)
	b[IndexOf(r, c)] = clear(b[IndexOf(r, c)], King)
}

//promoteToKing promotes the field to king
func (b Board) promoteToKing(r, c int) {
	b[IndexOf(r, c)] = set(b[IndexOf(r, c)], King)
}

//isKing checks if field contains a king
func (b Board) isKing(r, c int) bool {
	return has(b[IndexOf(r, c)], King)
}

//validateMove returns if the move would be valid in terms of gameplay and returns the move type and row and column
func (b Board) validateMove(fromRow, fromCol, toRow, toCol int, player bool) (MoveType, int, int) {
	i := IndexOf(toRow, toCol)
	f := b[i]
	if has(f, Empty) {
		return JumpMove, toRow, toCol
	}
	if !has(f, Empty) && has(f, Player) != player {
		//check if a skip is possible
		//need to evulate direction the jump is going
		if fromRow > toRow {
			//down
			if toCol < fromCol {
				//left
				if b.canDrawTo(toRow-1, toCol-1) && has(b[IndexOf(toRow-1, toCol-1)], Empty) {
					return SkipMove, toRow - 1, toCol - 1
				}
			} else {
				//right
				if b.canDrawTo(toRow-1, toCol+1) && has(b[IndexOf(toRow-1, toCol+1)], Empty) {
					return SkipMove, toRow - 1, toCol + 1
				}
			}
		} else {
			//up
			if toCol < fromCol {
				//left
				if b.canDrawTo(toRow+1, toCol-1) && has(b[IndexOf(toRow+1, toCol-1)], Empty) {
					return SkipMove, toRow + 1, toCol - 1
				}
			} else {
				//right
				if b.canDrawTo(toRow+1, toCol+1) && has(b[IndexOf(toRow+1, toCol+1)], Empty) {
					return SkipMove, toRow + 1, toCol + 1
				}
			}
		}
	}
	return InvalidMove, -1, -1
}

//getPossibleSkipsFor returns all possible skips (take moves)
func (b Board) getPossibleSkipsFor(r int, c int, player bool) []Move {
	m := make([]Move, 0)
	i := IndexOf(r, c)
	f := b[i]
	if !has(f, Empty) {
		//piece does not belong to player
		if has(f, Player) != player {
			return m
		}
		if has(f, Player) {
			//moving upwards
			if b.canDrawTo(r-1, c-1) {
				//valid movement left upwards
				if t, rf, cf := b.validateMove(r, c, r-1, c-1, player); t == SkipMove {
					m = append(m, Move{r, c, rf, cf, true, &Coordinate{r - 1, c - 1}, nil})
				}
			}
			if b.canDrawTo(r-1, c+1) {
				//valid movement right upwards
				if t, rf, cf := b.validateMove(r, c, r-1, c+1, player); t == SkipMove {
					m = append(m, Move{r, c, rf, cf, true, &Coordinate{r - 1, c + 1}, nil})
				}
			}
			if has(f, King) {
				//Kings can oposite direction as well
				//moving downwards
				if b.canDrawTo(r+1, c-1) {
					//valid movement left downards
					if t, rf, cf := b.validateMove(r, c, r+1, c-1, player); t == SkipMove {
						m = append(m, Move{r, c, rf, cf, true, &Coordinate{r + 1, c - 1}, nil})
					}
				}
				if b.canDrawTo(r+1, c+1) {
					//valid movement right downwards
					if t, rf, cf := b.validateMove(r, c, r+1, c+1, player); t == SkipMove {
						m = append(m, Move{r, c, rf, cf, true, &Coordinate{r + 1, c + 1}, nil})
					}
				}
			}
		} else {
			//moving downwards
			if b.canDrawTo(r+1, c-1) {
				//valid movement left downards
				if t, rf, cf := b.validateMove(r, c, r+1, c-1, player); t == SkipMove {
					m = append(m, Move{r, c, rf, cf, true, &Coordinate{r + 1, c - 1}, nil})
				}
			}
			if b.canDrawTo(r+1, c+1) {
				//valid movement right downwards
				if t, rf, cf := b.validateMove(r, c, r+1, c+1, player); t == SkipMove {
					m = append(m, Move{r, c, rf, cf, true, &Coordinate{r + 1, c + 1}, nil})
				}
			}
			if has(f, King) {
				//Kings can oposite direction as well
				//moving upwards
				if b.canDrawTo(r-1, c-1) {
					//valid movement left upwards
					if t, rf, cf := b.validateMove(r, c, r-1, c-1, player); t == SkipMove {
						m = append(m, Move{r, c, rf, cf, true, &Coordinate{r - 1, c - 1}, nil})
					}
				}
				if b.canDrawTo(r-1, c+1) {
					//valid movement right upwards
					if t, rf, cf := b.validateMove(r, c, r-1, c+1, player); t == SkipMove {
						m = append(m, Move{r, c, rf, cf, true, &Coordinate{r - 1, c + 1}, nil})
					}
				}
			}
		}
	}
	return m
}

//getPossibleMoves returns any possible moves for that field
func (b Board) getPossibleMoves(r int, c int, player bool) []Move {
	m := make([]Move, 0)
	i := IndexOf(r, c)
	f := b[i]
	if !has(f, Empty) {
		//piece does not belong to player
		if has(f, Player) != player {
			return m
		}
		if has(f, Player) {
			//moving upwards
			if b.canDrawTo(r-1, c-1) {
				//valid movement left upwards
				if t, rf, cf := b.validateMove(r, c, r-1, c-1, player); t != InvalidMove {
					if t == SkipMove {
						m = append(m, Move{r, c, rf, cf, true, &Coordinate{r - 1, c - 1}, nil})
					} else {
						m = append(m, Move{r, c, rf, cf, false, nil, nil})
					}

				}
			}
			if b.canDrawTo(r-1, c+1) {
				//valid movement right upwards
				if t, rf, cf := b.validateMove(r, c, r-1, c+1, player); t != InvalidMove {
					if t == SkipMove {
						m = append(m, Move{r, c, rf, cf, true, &Coordinate{r - 1, c + 1}, nil})
					} else {
						m = append(m, Move{r, c, rf, cf, false, nil, nil})
					}

				}
			}
			if has(f, King) {
				//Kings can oposite direction as well
				//moving downwards
				if b.canDrawTo(r+1, c-1) {
					//valid movement left downards
					if t, rf, cf := b.validateMove(r, c, r+1, c-1, player); t != InvalidMove {
						if t == SkipMove {
							m = append(m, Move{r, c, rf, cf, true, &Coordinate{r + 1, c - 1}, nil})
						} else {
							m = append(m, Move{r, c, rf, cf, false, nil, nil})
						}
					}
				}
				if b.canDrawTo(r+1, c+1) {
					//valid movement right downwards
					if t, rf, cf := b.validateMove(r, c, r+1, c+1, player); t != InvalidMove {
						if t == SkipMove {
							m = append(m, Move{r, c, rf, cf, true, &Coordinate{r + 1, c + 1}, nil})
						} else {
							m = append(m, Move{r, c, rf, cf, false, nil, nil})
						}
					}
				}
			}
		} else {
			//moving downwards
			if b.canDrawTo(r+1, c-1) {
				//valid movement left downards
				if t, rf, cf := b.validateMove(r, c, r+1, c-1, player); t != InvalidMove {
					if t == SkipMove {
						m = append(m, Move{r, c, rf, cf, true, &Coordinate{r + 1, c - 1}, nil})
					} else {
						m = append(m, Move{r, c, rf, cf, false, nil, nil})
					}
				}
			}
			if b.canDrawTo(r+1, c+1) {
				//valid movement right downwards
				if t, rf, cf := b.validateMove(r, c, r+1, c+1, player); t != InvalidMove {
					if t == SkipMove {
						m = append(m, Move{r, c, rf, cf, true, &Coordinate{r + 1, c + 1}, nil})
					} else {
						m = append(m, Move{r, c, rf, cf, false, nil, nil})
					}
				}
			}
			if has(f, King) {
				//Kings can oposite direction as well
				//moving upwards
				if b.canDrawTo(r-1, c-1) {
					//valid movement left upwards
					if t, rf, cf := b.validateMove(r, c, r-1, c-1, player); t != InvalidMove {
						if t == SkipMove {
							m = append(m, Move{r, c, rf, cf, true, &Coordinate{r - 1, c - 1}, nil})
						} else {
							m = append(m, Move{r, c, rf, cf, false, nil, nil})
						}

					}
				}
				if b.canDrawTo(r-1, c+1) {
					//valid movement right upwards
					if t, rf, cf := b.validateMove(r, c, r-1, c+1, player); t != InvalidMove {
						if t == SkipMove {
							m = append(m, Move{r, c, rf, cf, true, &Coordinate{r - 1, c + 1}, nil})
						} else {
							m = append(m, Move{r, c, rf, cf, false, nil, nil})
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
	must := false
	for _, v := range move {
		if v.Must {
			must = true
			break
		}
	}
	if must {
		filtered := make([]Move, 0)
		for _, v := range move {
			if v.Must {
				filtered = append(filtered, v)
			}
		}
		return filtered
	}
	return move
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
