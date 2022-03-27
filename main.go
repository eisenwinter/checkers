package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/eisenwinter/checkers/game"
)

var fullAIMode = false

const aiMoveSeconds = 0.3

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Checkers",
		Bounds: pixel.R(0, 0, 600, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Skyblue)
	win.SetSmooth(true)
	grid := imdraw.New(nil)
	g := game.SetupGame()
	g.Start()
	moves := []game.Move{}
	highlight := []game.Coordinate{}
	last := time.Now()
	moving := 0.0
	var currentBoard game.Board
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		if moving > 0 {
			moving = moving - dt
			if moving < 0 {
				moving = 0
			}
		}
		last = time.Now()
		if g.GameState() == game.GameStateRunning {
			if !g.HasBoardInQueue() {
				if !g.Player() {
					//ai
					g.MakeAIMove()
				} else {
					if fullAIMode {
						g.MakeAIMove()
					} else {
						if win.JustPressed(pixelgl.MouseButtonLeft) {
							vec := win.MousePosition()
							col := math.Floor(vec.X / 60)
							row := math.Floor((win.Bounds().H() - vec.Y) / 60)
							if row < 10 && col < 10 {
								done := false
								if len(moves) > 0 {
									for _, v := range moves {
										if v.ToCol == int(col) && v.ToRow == int(row) {
											followUp := g.MakeMove(v)
											if !followUp {
												done = true
												moves = []game.Move{}
												highlight = []game.Coordinate{}
											}
										}
									}
								}
								if !done {
									moves, highlight = g.GetPossibleMoves(int(row), int(col))
								}
							}
						}
					}
				}
			}
		}
		grid.Clear()
		mat := pixel.IM
		mat = mat.Rotated(win.Bounds().Center(), -math.Pi/2)
		grid.SetMatrix(mat)
		if moving > 0 {
			DrawBoard(grid, currentBoard, moves, highlight)
		} else {
			if g.HasBoardInQueue() {
				currentBoard = g.DequeueBoard()
				DrawBoard(grid, currentBoard, moves, highlight)
				moving = aiMoveSeconds
			} else {
				DrawBoard(grid, g.CurrentBoard(), moves, highlight)
			}
		}

		grid.Draw(win)
		win.Update()
	}
}

func DrawBoard(imd *imdraw.IMDraw, board game.Board, moves []game.Move, hl []game.Coordinate) {
	cellSize := 60
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if (j+i%2)%2 == 0 {
				imd.Color = colornames.Black
			} else {
				imd.Color = colornames.White
			}
			imd.Push(
				pixel.V(float64(i*cellSize), float64(j*cellSize)),
				pixel.V(float64(i*cellSize+cellSize), float64(j*cellSize+cellSize)),
			)
			imd.Rectangle(0)

		}
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			idx := game.IndexOf(i, j)
			if !game.IsEmptyField(board[idx]) {
				if game.IsPlayer(board[idx]) {
					imd.Color = colornames.Darkgray
					imd.Push(
						pixel.V(float64(i*cellSize+cellSize/2), float64(j*cellSize+cellSize/2)),
					)
					if game.IsKing(board[idx]) {
						imd.Circle(float64(cellSize/3), 8)
					} else {
						imd.Circle(float64(cellSize/3), 0)
					}

				} else {
					imd.Color = colornames.Red
					imd.Push(
						pixel.V(float64(i*cellSize+cellSize/2), float64(j*cellSize+cellSize/2)),
					)
					if game.IsKing(board[idx]) {
						imd.Circle(float64(cellSize/3), 8)
					} else {
						imd.Circle(float64(cellSize/3), 0)
					}
				}
			}
		}
	}

	if len(moves) > 0 {
		imd.Color = colornames.Lightgreen
		for _, v := range moves {
			imd.Push(
				pixel.V(float64(v.ToRow*cellSize+cellSize/2), float64(v.ToCol*cellSize+cellSize/2)),
			)
			imd.Circle(float64(cellSize/3), 0)
		}
	}

	if len(hl) > 0 {
		imd.Color = colornames.Lightsalmon
		for _, v := range hl {
			imd.Push(
				pixel.V(float64(v.Row*cellSize+3), float64(v.Col*cellSize+3)),
				pixel.V(float64(v.Row*cellSize+cellSize-3), float64(v.Col*cellSize+cellSize-3)),
			)
			imd.Rectangle(4)
		}
	}

}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logger))
	aiMode := flag.Bool("ai", false, "full auto ai flag")
	flag.Parse()
	fullAIMode = *aiMode
	log.Printf("AI only mode: %v", fullAIMode)
	pixelgl.Run(run)
}

type logger struct {
}

func (writer logger) Write(bytes []byte) (int, error) {
	return fmt.Printf("[DBG] @ %s | %s", time.Now().UTC().Format("15:04:05"), string(bytes))
}
