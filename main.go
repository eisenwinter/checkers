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
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/eisenwinter/checkers/game"
)

var fullAIMode = false
var showEvalMode = false

const moveSeconds = 0.3

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
	highlight := []game.Path{}
	last := time.Now()
	moving := 0.0
	var currentBoard game.Board

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	overlayText := text.New(pixel.V(10, 10), atlas)
	overlayText.Color = colornames.Magenta
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
										if v.To.Col == int(col) && v.To.Row == int(row) {
											g.MakeMove(v)
											done = true
											moves = []game.Move{}
											highlight = []game.Path{}
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
				moving = moveSeconds
			} else {
				DrawBoard(grid, g.CurrentBoard(), moves, highlight)
			}
		}

		grid.Draw(win)
		if fullAIMode || showEvalMode {
			overlayText.Clear()
			fmt.Fprintf(overlayText, "%d", g.CurrentEvaulation())
			overlayText.Draw(win, pixel.IM)
		}
		win.Update()
	}
}

func DrawBoard(imd *imdraw.IMDraw, board game.Board, moves []game.Move, hl []game.Path) {
	cellSize := 60
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if (j+i%2)%2 == 0 {
				imd.Color = colornames.White
			} else {
				imd.Color = colornames.Black

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

	if len(hl) > 0 {
		for _, p := range hl {

			for i, v := range p.Coordinates {
				if i == 0 {
					imd.Color = colornames.Lightsalmon
					imd.Push(
						pixel.V(float64(v.Row*cellSize+3), float64(v.Col*cellSize+3)),
						pixel.V(float64(v.Row*cellSize+cellSize-3), float64(v.Col*cellSize+cellSize-3)),
					)
					imd.Rectangle(4)
				} else {
					imd.Color = colornames.Lightblue
					imd.Push(
						pixel.V(float64(v.Row*cellSize+cellSize/2), float64(v.Col*cellSize+cellSize/2)),
					)
					imd.Circle(float64(cellSize/3), 0)
				}

			}
		}
	}

	if len(moves) > 0 {
		for _, v := range moves {
			imd.Color = colornames.Lightgreen
			imd.Push(
				pixel.V(float64(v.To.Row*cellSize+cellSize/2), float64(v.To.Col*cellSize+cellSize/2)),
			)
			imd.Circle(float64(cellSize/3), 0)
		}

	}

}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logger))
	aiMode := flag.Bool("ai", false, "full auto ai flag")
	scoreMode := flag.Bool("s", false, "show score")
	flag.Parse()
	fullAIMode = *aiMode
	showEvalMode = *scoreMode
	log.Printf("AI only mode: %v", fullAIMode)
	pixelgl.Run(run)
}

type logger struct {
}

func (writer logger) Write(bytes []byte) (int, error) {
	return fmt.Printf("[DBG] @ %s | %s", time.Now().UTC().Format("15:04:05"), string(bytes))
}
