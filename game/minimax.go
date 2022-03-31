package game

import (
	"fmt"
	"log"
	"math"
	"strings"
)

const MaxDepth = 4
const AlphaStart = math.MinInt
const BetaStart = math.MaxInt

func (b Board) evaluate() int {
	w, r, wk, rk := b.getCounts()
	//regular pieces weight 10, kings 15
	base := (w * 14) - (r * 14)
	base = base + (wk*40 - rk*40)

	//weight of each piece in the backrow
	// wbr, rbr := b.getBackRowCount()
	// base = base + (wbr*14 - rbr*14)

	wbr, rbr := b.getGoldenStoneCount()
	base = base + (wbr*14 - rbr*14)

	//weight of each piece middle box position
	wmb, rmb := b.getMiddleBoxCount()
	base = base + (wmb*7 - rmb*7)

	wms, rms := b.getMiddleCount()
	base = base + (wms*3 - rms*3)

	wls, rls := b.getLeftSideCount()
	base = base + (wls*1 - rls*1)

	wrs, rrs := b.getRightSideCount()
	base = base + (wrs*1 - rrs*1)

	//weight of a protected piece
	wpr, rpr := b.getProtectionCount()
	base = base + (wpr*4 - rpr*4)

	wst, rst, wstk, rstk := b.getStuckPiecesCount()
	base = base + (wst * -1) - (rst * -1) + (wstk*-2 - rstk*-2)

	wfs, rfs := b.getFullSquares()
	base = base + (wfs*14 - rfs*14)

	whs, rhs := b.getHalfSquare()
	base = base + (whs*12 - rhs*12)

	wgs, rgs := b.getFullGates()
	base = base + (wgs*9 - rgs*9)

	whgs, rhgs := b.getHalfGates()
	base = base + (whgs*7 - rhgs*7)

	wps, rps := b.getPincers()
	base = base + (wps*5 - rps*5)

	return base
}

type HeustricStat struct {
	name  string
	white int
	red   int
}

func (b Board) LogBoardHeurstics() {
	stats := make([]HeustricStat, 0)
	wpr, rpr := b.getProtectionCount()
	stats = append(stats, HeustricStat{"protection count", wpr, rpr})
	wst, rst, _, _ := b.getStuckPiecesCount()
	stats = append(stats, HeustricStat{"stuck pieces", wst, rst})
	wfs, rfs := b.getFullSquares()
	stats = append(stats, HeustricStat{"full squares", wfs, rfs})
	whs, rhs := b.getHalfSquare()
	stats = append(stats, HeustricStat{"half squares", whs, rhs})
	wgs, rgs := b.getFullGates()
	stats = append(stats, HeustricStat{"full gates", wgs, rgs})
	whgs, rhgs := b.getHalfGates()
	stats = append(stats, HeustricStat{"half gates", whgs, rhgs})
	wps, rps := b.getPincers()
	stats = append(stats, HeustricStat{"pincers", wps, rps})
	var whiteStats strings.Builder
	var redStats strings.Builder
	for _, v := range stats {
		fmt.Fprintf(&whiteStats, " %s: %02d |", v.name, v.white)
		fmt.Fprintf(&redStats, " %s: %02d |", v.name, v.red)
	}
	log.Printf("White| %s", whiteStats.String())
	log.Printf("Red  | %s", redStats.String())
	// log.Printf("White: Backrow %d | Box: %d | M.Side: %d | Vulnerable: %d (K: %d) | Protected: %d | Stuck: %d (K: %d) | Fortified: %d | Diamonds: %d | Runaway: %d", wbr, wmb, wmr, wvp, wvk, wpr, wst, wstk, wfr, wds, wrun)
	// log.Printf("Red:   Backrow %d | Box: %d | M.Side: %d | Vulnerable: %d (K: %d) | Protected: %d | Stuck: %d (K: %d) | Fortified: %d | Diamonds: %d | Runaway: %d", rbr, rmb, rmr, rvp, rvk, rpr, rst, rstk, rfr, rds, rrun)
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

func minimax(depth int, board Board, player bool, alpha int, beta int, m *Move) (int, *Move) {
	if depth == 0 || !board.playable() {
		return board.evaluate(), m
	}
	if player {
		value := math.MinInt
		var move *Move
		for k, v := range possibleMoves(player, board) {
			eval, _ := minimax(depth-1, v, !player, alpha, beta, &k)
			value = maxOf(value, eval)
			if value == eval {
				move = &k
			}
			alpha = maxOf(alpha, eval)
			if eval >= beta {
				return value, move
			}
		}
		return value, move
	} else {
		value := math.MaxInt
		var move *Move
		for k, v := range possibleMoves(player, board) {
			eval, _ := minimax(depth-1, v, !player, alpha, beta, &k)
			value = minOf(value, eval)
			if value == eval {
				move = &k
			}
			beta = minOf(beta, eval)
			if eval <= alpha {
				return value, move
			}
		}
		return value, move
	}
}

func unrollMove(b *Board, move Move, player bool, maxDepth int) {
	if move.Previous != nil {
		unrollMove(b, *move.Previous, player, maxDepth)
	}
	k := b.applyMove(move, player)
	if move.Depth == maxDepth && k {
		b.promoteToKing(move.To)
	}
}

func possibleMoves(player bool, board Board) map[Move]Board {
	m := make(map[Move]Board)
	possible := board.getPossibleValidMovesForPlayer(player)
	for _, move := range possible {
		tmp := board.copy()
		unrollMove(&tmp, move, player, move.Depth)
		m[move] = tmp
	}
	return m
}

//getPossibleValidMovesForPlayer utility method for the ai
func (b Board) getPossibleValidMovesForPlayer(player bool) []Move {
	m := make([]Move, 0)
	pieces := b.allPiecesFor(player)
	for _, p := range pieces {
		moves := b.getPossibleMoves(p, player)
		m = append(m, moves...)
	}
	return filterMoves(m)
}
