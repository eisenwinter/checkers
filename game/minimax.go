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

const pieceBaseVaue = 20

func (b Board) evaluate() int {
	w, r, wk, rk := b.getCounts()

	if w == 0 {
		return math.MinInt32
	}
	if r == 0 {
		return math.MaxInt32
	}

	base := (w * pieceBaseVaue) - (r * pieceBaseVaue)
	base = base + (wk*pieceBaseVaue*2 - rk*pieceBaseVaue*2)

	wbr, rbr := b.getGoldenStoneCount()
	base = base + (wbr - rbr)

	wmb, rmb := b.getMiddleBoxCount()
	base = base + (wmb*3 - rmb*3)
	wms, rms := b.getMiddleCount()
	base = base + (wms*2 - rms*2)
	wls, rls := b.getLeftSideCount()
	base = base + (wls*1 - rls*1)
	wrs, rrs := b.getRightSideCount()
	base = base + (wrs*1 - rrs*1)
	wpr, rpr := b.getProtectionCount()
	base = base + (wpr*4 - rpr*4)

	wst, rst, wstk, rstk := b.getStuckPiecesCount()
	base = base + (wst * -1 * (pieceBaseVaue / 2)) - (rst * -1 * (pieceBaseVaue / 2)) + (wstk*-2*(pieceBaseVaue/2) - rstk*-2*(pieceBaseVaue/2))
	if wst == w {
		return math.MinInt32
	}
	if rst == r {
		return math.MaxInt32
	}

	wlgr, rlgr := b.getLeggardAndGrapeCount()
	base = base + (wlgr * -1 * pieceBaseVaue) - (rlgr * -1 * pieceBaseVaue)

	wfs, rfs := b.getFullSquares()
	base = base + (wfs*9*pieceBaseVaue - rfs*9*pieceBaseVaue)

	whs, rhs := b.getHalfSquare()
	base = base + (whs*5*pieceBaseVaue - rhs*5*pieceBaseVaue)

	wgs, rgs := b.getFullGates()
	base = base + (wgs*4*pieceBaseVaue - rgs*4*pieceBaseVaue)

	whgs, rhgs := b.getHalfGates()
	base = base + (whgs*3*pieceBaseVaue - rhgs*3*pieceBaseVaue)

	wps, rps := b.getPincers()
	base = base + (wps*4*pieceBaseVaue - rps*4*pieceBaseVaue)

	llw, lrw := b.getLargestConnectedField()
	base = base + (llw*2*pieceBaseVaue - lrw*2*pieceBaseVaue)

	wvp, rvp := b.getVulnerablePiecesCount()
	base = base + (wvp * -50) - (rvp * -50)

	wsc, rsc := b.getSuicidalPiecesCount()
	base = base + (wsc * -20) - (rsc * -20)

	//white := b.getPossibleValidMovesForPlayer(true)
	//red := b.getPossibleValidMovesForPlayer(false)

	//rethinking this
	//so we basically check all moves once and apply a possible evaulation
	//of the resutls
	//a move is good IF
	//the move saves a check from beeing taken
	//the move protects a checker
	//the move lead to a long jump

	//a move is bad IF
	//the move leads to the checker beeing taken (anti suicide measure)
	//the move makes a checker vulnerable (ends protection of checker)

	return base
}

type HeuristicStat struct {
	name  string
	white int
	red   int
}

func (b Board) LogBoardHeurstics() {
	stats := make([]HeuristicStat, 0)

	wbr, rbr := b.getGoldenStoneCount()
	stats = append(stats, HeuristicStat{"g.stones", wbr, rbr})

	wlgr, rlgr := b.getLeggardAndGrapeCount()
	stats = append(stats, HeuristicStat{"l&g", wlgr, rlgr})

	wmb, rmb := b.getMiddleBoxCount()
	stats = append(stats, HeuristicStat{"m.box", wmb, rmb})

	wms, rms := b.getMiddleCount()
	stats = append(stats, HeuristicStat{"m.", wms, rms})

	wls, rls := b.getLeftSideCount()
	stats = append(stats, HeuristicStat{"l.", wls, rls})

	wrs, rrs := b.getRightSideCount()
	stats = append(stats, HeuristicStat{"r.", wrs, rrs})

	wpr, rpr := b.getProtectionCount()
	stats = append(stats, HeuristicStat{"protection count", wpr, rpr})
	wst, rst, _, _ := b.getStuckPiecesCount()
	stats = append(stats, HeuristicStat{"stuck pieces", wst, rst})
	wfs, rfs := b.getFullSquares()
	stats = append(stats, HeuristicStat{"full squares", wfs, rfs})
	whs, rhs := b.getHalfSquare()
	stats = append(stats, HeuristicStat{"half squares", whs, rhs})
	wgs, rgs := b.getFullGates()
	stats = append(stats, HeuristicStat{"full gates", wgs, rgs})
	whgs, rhgs := b.getHalfGates()
	stats = append(stats, HeuristicStat{"half gates", whgs, rhgs})
	wps, rps := b.getPincers()
	stats = append(stats, HeuristicStat{"pincers", wps, rps})

	llw, lrw := b.getLargestConnectedField()
	stats = append(stats, HeuristicStat{"largest field", llw, lrw})

	wvp, rvp := b.getVulnerablePiecesCount()
	stats = append(stats, HeuristicStat{"vulnerable pieces", wvp, rvp})

	wsc, rsc := b.getSuicidalPiecesCount()
	stats = append(stats, HeuristicStat{"suicidal pieces", wsc, rsc})

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
