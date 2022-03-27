package game

import (
	"math"
)

const MaxDepth = 10
const AlphaStart = math.MaxInt
const BetaStart = math.MinInt

func (b Board) evaluate() int {
	w, r, wk, rk := b.getCounts()
	//regular pieces weight 10, kings 15
	base := (w * 10) - (r * 10) + (wk*15 - rk*15)

	//weight of each piece in the backrow
	wbr, rbr := b.getBackRowCount()
	base = base + (wbr*8 - rbr*8)

	//weight of each piece middle box position
	wmb, rmb := b.getMiddleBoxCount()
	base = base + (wmb*5 - rmb*5)

	//weight of each piece in the middle two rows
	wmr, rmr := b.getMiddleRowSideCount()
	base = base + (wmr - rmr)

	//weight of a vulnerable piece
	wvp, rvp, wvk, rvk := b.getVulnerablePieceCount()
	base = base + (wvp * -6) - (rvp * -6)
	base = base + (wvk * -1) - (rvk * -1)

	//weight of a protected piece
	wpr, rpr := b.getProtectionCount()
	base = base + (wpr*6 - rpr*6)

	return base
}

func maxOf(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func minOf(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func minimax(depth int, max bool, board Board, player bool, alpha int, beta int, m *Move) (int, *Move) {
	if depth == 0 || !board.playable() {
		return board.evaluate(), m
	}
	if max {
		maxSoFar := math.MinInt
		var bestMove *Move
		for k, v := range possibleMoves(player, board) {
			eval, _ := minimax(depth-1, false, v, player, alpha, beta, &k)
			maxSoFar = maxOf(maxSoFar, eval)
			if maxSoFar == eval {
				bestMove = &k
			}
			alpha = maxOf(alpha, eval)
			if eval >= beta {
				return maxSoFar, bestMove
			}
		}
		return maxSoFar, bestMove
	} else {
		minSoFar := math.MaxInt
		var worstMove *Move
		for k, v := range possibleMoves(player, board) {
			eval, _ := minimax(depth-1, true, v, !player, alpha, beta, &k)
			minSoFar = minOf(minSoFar, eval)
			if minSoFar == eval {
				worstMove = &k
			}
			beta = minOf(beta, eval)
			if eval <= beta {
				return minSoFar, worstMove
			}
		}
		return minSoFar, worstMove
	}
}

func unrollSkips(r, c int, player bool, board Board, prev *Move) map[Move]Board {
	m := make(map[Move]Board)
	tmp := make(Board, len(board))
	possible := tmp.getPossibleSkipsFor(r, c, player)
	for _, p := range possible {
		p.Previous = prev
		copy(tmp, board)
		followUps := tmp.applyMove(p, player)
		if followUps {
			return unrollSkips(p.ToRow, p.ToCol, player, tmp, &p)
		} else {
			m[p] = tmp
		}
	}
	return m
}

func possibleMoves(player bool, board Board) map[Move]Board {
	m := make(map[Move]Board)
	possible := board.getPossibleValidMovesForPlayer(player)
	for _, move := range possible {
		tmp := make(Board, len(board))
		copy(tmp, board)
		followUps := tmp.applyMove(move, player)
		if followUps {
			mx := unrollSkips(move.ToRow, move.ToCol, player, tmp, &move)
			if len(mx) > 0 {
				for k, v := range mx {
					m[k] = v
				}
			} else {
				m[move] = tmp
			}
		} else {
			m[move] = tmp
		}
	}
	return m
}

//getPossibleValidMovesForPlayer utility method for the ai, shuffles move for more variation
func (b Board) getPossibleValidMovesForPlayer(player bool) []Move {
	m := make([]Move, 0)
	pieces := b.allPiecesFor(player)
	for _, p := range pieces {
		moves := b.getPossibleMoves(p.Row, p.Col, player)
		m = append(m, moves...)
	}
	return filterMoves(m)
}
