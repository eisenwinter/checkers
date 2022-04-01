package game

import (
	"fmt"
	"log"
	"time"
)

//Game represents a game
type Game struct {
	turn       int
	player     bool // true = white, false = red
	board      Board
	running    bool
	started    time.Time
	redCount   int
	whiteCount int
	redKings   int
	whiteKings int
	state      GameState
	boardQueue []Board
}

//Start starts the game
func (g *Game) Start() {
	g.running = true
	g.started = time.Now()
	g.boardQueue = make([]Board, 0)
	g.board.LogBoardHeurstics()
	log.Printf("Starting evulation: %d", g.CurrentEvaulation())
}

//IsRunning indicates if the game is still running
func (g *Game) IsRunning() bool {
	return g.running
}

//SetupGame creates a new game instance ready to start
func SetupGame() *Game {
	b := boardSetup(make(Board, height*width))
	w, r, wk, rk := b.getCounts()
	return &Game{
		turn:       0,
		redCount:   r,
		whiteCount: w,
		redKings:   rk,
		whiteKings: wk,
		player:     true, //white goes first
		board:      b,
		state:      GameStateRunning,
	}
}

//refreshCount updates the current board counts
func (g *Game) refreshCount() {
	w, r, wk, rk := g.board.getCounts()
	g.whiteCount = w
	g.redCount = r
	g.redKings = rk
	g.whiteKings = wk
}

//CurrentBoard returns the current game board
func (g *Game) CurrentBoard() Board {
	return g.board
}

func (g *Game) CurrentEvaulation() int {
	return g.board.evaluate()
}

//GameState is the state the game is currently in
type GameState string

var GameStateRunning GameState = "Running"
var GameStateRedWins GameState = "RedWin"
var GameStateWhiteWins GameState = "WhiteWin"
var GameStateDraw GameState = "Draw"

//GameState returns the current gamestate
func (g *Game) GameState() GameState {
	if g.state != GameStateRunning {
		return g.state
	}
	if g.whiteCount == 0 {
		log.Printf("White has no more pieces left, red wins")
		g.state = GameStateRedWins
		return GameStateRedWins
	}
	if g.redCount == 0 {
		log.Printf("Red has no more pieces left, white wins")
		g.state = GameStateWhiteWins
		return GameStateWhiteWins
	}

	//get pieces of other player first
	otherPieces := g.board.allPiecesFor(!g.player)
	om := make([]Move, 0)
	for _, p := range otherPieces {
		o := filterMoves(g.board.getPossibleMoves(p, !g.player))
		om = append(om, o...)
	}
	if len(om) == 0 {
		log.Printf("Its currently %+v player ", g.player)
		if !g.player {
			log.Printf("White has no more moves left, red wins")
			g.state = GameStateRedWins
			return GameStateRedWins
		} else {
			log.Printf("Red has no more moves left, white wins")
			g.state = GameStateWhiteWins
			return GameStateWhiteWins
		}
	}
	return GameStateRunning
}

//StatusDisplay is used to return the current game stats
func (g *Game) StatusDisplay() string {
	if g.state == GameStateRedWins {
		return "Red Wins"
	}
	if g.state == GameStateWhiteWins {
		return "White Wins"
	}
	if g.state == GameStateDraw {
		return "Draw"
	}
	return fmt.Sprintf("White: %d (%d Kings) | Red: %d (%d Kings) | Turn: %d | Time (s): %d", g.whiteCount, g.whiteKings, g.redCount, g.redKings, g.Turn(), g.Time())
}

//checkForcedMove checks if thats a MUST play move
func (g *Game) checkForcedMove(r, c int, player bool) bool {
	im := g.board.getPossibleMoves(Coordinate{r, c}, g.player)
	for _, v := range im {
		if v.Depth > 0 {
			return true
		}
	}
	return false
}

//checkForcedMoves returns all possible forced moves (must play moves)
func (g *Game) checkForcedMoves() []Move {
	m := make([]Move, 0)
	p := g.board.allPiecesFor(g.player)
	for _, v := range p {
		im := g.board.getPossibleMoves(v, g.player)
		for _, v := range im {
			if v.Takes != nil || v.Depth > 0 {
				m = append(m, v)
			}
		}
	}
	return m
}

func pathFromMove(m Move, p *Path) {
	if m.Previous != nil {
		pathFromMove(*m.Previous, p)
	} else {
		p.Coordinates = append(p.Coordinates, m.From)
	}
	p.Coordinates = append(p.Coordinates, m.To)
}
func pathsFromMoves(m []Move) []Path {
	p := make([]Path, 0)
	for _, v := range m {
		path := Path{
			Coordinates: make([]Coordinate, 0),
		}
		pathFromMove(v, &path)
		p = append(p, path)
	}
	return p
}
func mapPossibleMove(m []Move) []PossibleMove {
	pos := make([]PossibleMove, 0)
	for _, v := range m {
		path := Path{
			Coordinates: make([]Coordinate, 0),
		}
		pathFromMove(v, &path)
		pos = append(pos, PossibleMove{v, path})
	}

	return pos
}

type PossibleMove struct {
	Move Move
	Path Path
}

//GetPossibleMoves returns the possible moves for that given field
func (g *Game) GetPossibleMoves() []PossibleMove {

	m := make([]Move, 0)
	p := g.board.allPiecesFor(g.player)
	for _, v := range p {
		im := g.board.getPossibleMoves(v, g.player)
		m = append(m, im...)
	}
	m = filterMoves(m)
	return mapPossibleMove(m)
}

//Turn returns the current turn
func (g *Game) Turn() int {
	return g.turn
}

//Time returns the current time in seconds
func (g *Game) Time() int {
	return int(time.Since(g.started).Seconds())
}

func (g *Game) unrollMove(m *Move, maxDepth int) {
	if m.Previous != nil {
		g.unrollMove(m.Previous, maxDepth)
	}
	k := g.makeMove(*m)
	if m.Depth == maxDepth && k {
		g.board.promoteToKing(m.To)
	}
	g.boardQueue = append(g.boardQueue, g.board.copy())
	if maxDepth == m.Depth {
		g.refreshCount()
		g.endTurn()
	}
}

//MakeAIMove triggers a computer move
func (g *Game) MakeAIMove() {
	if g.GameState() == GameStateRunning {
		_, m := minimax(MaxDepth, g.board, g.player, AlphaStart, BetaStart, nil)
		if m != nil {
			g.unrollMove(m, m.Depth)
		} else {
			panic("well well well this should not happend - no solution found")
		}
	}
}

func (g *Game) MakeMove(m Move) {
	g.unrollMove(&m, m.Depth)
}

//MakeMove applies the given move to the board
//if true is returned the corresponding player has to make another move
//the next move has to be a forced move
func (g *Game) makeMove(m Move) bool {
	promoted := g.board.applyMove(m, g.player)
	return promoted

}

func (g *Game) endTurn() {
	if g.GameState() == GameStateRunning {
		g.turn++
		g.player = !g.player
		log.Printf("Turn: %d | Current board eval: %d | Whites turn: %v", g.turn, g.board.evaluate(), g.player)
		g.board.LogBoardHeurstics()
	} else {
		g.running = false
		log.Printf("Final board eval: %d | Whites turn: %v", g.board.evaluate(), g.player)
	}
	log.Print(g.StatusDisplay())
}

//Player indiciates wich players turn it is (True = White, False = Red)
func (g *Game) Player() bool {
	return g.player
}

func (g *Game) HasBoardInQueue() bool {
	return len(g.boardQueue) > 0
}

func (g *Game) DequeueBoard() Board {
	b := g.boardQueue[0]
	g.boardQueue = g.boardQueue[1:]
	return b
}
