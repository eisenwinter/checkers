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

	if w == 0 {
		return math.MinInt32
	}
	if r == 0 {
		return math.MaxInt32
	}

	base := (w * 20) - (r * 20)
	base = base + (wk*25 - rk*25)

	wbr, rbr := b.getGoldenStoneCount()
	base = base + (wbr*7 - rbr*7)

	wmb, rmb := b.getMiddleBoxCount()
	base = base + (wmb*7 - rmb*7)

	wms, rms := b.getMiddleCount()
	base = base + (wms*2 - rms*2)

	wls, rls := b.getLeftSideCount()
	base = base + (wls*1 - rls*1)

	wrs, rrs := b.getRightSideCount()
	base = base + (wrs*1 - rrs*1)

	wpr, rpr := b.getProtectionCount()
	base = base + (wpr*4 - rpr*4)

	wst, rst, wstk, rstk := b.getStuckPiecesCount()
	base = base + (wst * -1) - (rst * -1) + (wstk*-2 - rstk*-2)

	if wst == w {
		return math.MinInt32
	}

	if rst == r {
		return math.MaxInt32
	}

	wlgr, rlgr := b.getLeggardAndGrapeCount()
	base = base + (wlgr * -1) - (rlgr * -1)

	wfs, rfs := b.getFullSquares()
	base = base + (wfs*7 - rfs*7)

	whs, rhs := b.getHalfSquare()
	base = base + (whs*4 - rhs*4)

	wgs, rgs := b.getFullGates()
	base = base + (wgs*3 - rgs*3)

	whgs, rhgs := b.getHalfGates()
	base = base + (whgs*2 - rhgs*2)

	wps, rps := b.getPincers()
	base = base + (wps*3 - rps*2)

	llw, lrw := b.getLargestConnectedField()
	base = base + (llw - lrw)

	wvp, rvp := b.getVulnerablePiecesCount()
	base = base + (wvp * -50) - (rvp * -50)

	wsc, rsc := b.getSuicidalPiecesCount()
	base = base + (wsc * -20) - (rsc * -20)

	return base
}

type HeustricStat struct {
	name  string
	white int
	red   int
}

func (b Board) LogBoardHeurstics() {
	stats := make([]HeustricStat, 0)

	wbr, rbr := b.getGoldenStoneCount()
	stats = append(stats, HeustricStat{"g.stones", wbr, rbr})

	wlgr, rlgr := b.getLeggardAndGrapeCount()
	stats = append(stats, HeustricStat{"l&g", wlgr, rlgr})

	wmb, rmb := b.getMiddleBoxCount()
	stats = append(stats, HeustricStat{"m.box", wmb, rmb})

	wms, rms := b.getMiddleCount()
	stats = append(stats, HeustricStat{"m.", wms, rms})

	wls, rls := b.getLeftSideCount()
	stats = append(stats, HeustricStat{"l.", wls, rls})

	wrs, rrs := b.getRightSideCount()
	stats = append(stats, HeustricStat{"r.", wrs, rrs})

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

	llw, lrw := b.getLargestConnectedField()
	stats = append(stats, HeustricStat{"largest field", llw, lrw})

	wvp, rvp := b.getVulnerablePiecesCount()
	stats = append(stats, HeustricStat{"vulnerable pieces", wvp, rvp})

	wsc, rsc := b.getSuicidalPiecesCount()
	stats = append(stats, HeustricStat{"suicidal pieces", wsc, rsc})

	var whiteStats strings.Builder
	var redStats strings.Builder
	for _, v := range stats {
		fmt.Fprintf(&whiteStats, " %s: %02d |", v.name, v.white)
		fmt.Fprintf(&redStats, " %s: %02d |", v.name, v.red)
	}
	log.Printf("White| %s", whiteStats.String())
	log.Printf("Red  | %s", redStats.String())
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
	terminal := !board.playable()
	if depth == 0 || terminal {
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
