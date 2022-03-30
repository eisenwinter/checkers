# Checkers

 ![Checkers](/doc/animation.gif)

>> IN THE MIDDLE OF REFACTORING - NOT WORKING RN

Simple international checkers, playing arround with minimax heuristics.

It can be launched in full ai vs ai mode with the -ai flag.
The main goal here wasnt the gameitself but rather exploring minimax and heuristics.

I started with adapted heurstics from https://github.com/kevingregor/Checkers/blob/master/Final%20Project%20Report.pdf but its ever evolving and changing right now.

Due to the protection heuristic some AI vs AI games will result in some wall hugging.

If you might wonder why each ai run yields different results despite no shuffelinig, 
its because the go map is intentionally random when iterated with range,
thus if there are many equal moves it will pick a random one and there is no need to shuffle.


## Game Rules

Source: https://en.wikipedia.org/wiki/International_draughts

- The game is played on a board with 10×10 squares, alternatingly dark and light. The lower-leftmost square should be dark.
- Each player has 20 pieces. In the starting position  the pieces are placed on the first four rows closest to the players. This leaves two central rows empty.



- The player with the light pieces moves first. Then turns alternate.
- Ordinary pieces move one square diagonally forward to an unoccupied square.
Enemy pieces can and must be captured by jumping over the enemy piece, two squares forward or backward to an unoccupied square immediately beyond. If a jump is possible it must be done, even if doing so incurs a disadvantage.
- - Multiple successive jumps forward or backward in a single turn can and must be made if after each jump there is an unoccupied square immediately beyond the enemy piece. It is compulsory to jump over as many pieces as possible. One must play with the piece that can make the maximum number of captures.
- - A jumped piece is removed from the board at the end of the turn. (So for a multi-jump move, jumped pieces are not removed during the move, they are removed only after the entire multi-jump move is complete.)
The same piece may not be jumped more than once.
- A piece is crowned if it stops on the far edge of the board at the end of its turn (that is, not if it reaches the edge but must then jump another piece backward). Another piece is placed on top of it to mark it. Crowned pieces, sometimes called kings, can move freely multiple steps in any direction and may jump over and hence capture an opponent piece some distance away and choose where to stop afterwards, but must still capture the maximum number of pieces possible.

- A player with no valid move remaining loses. This occurs if the player has no pieces left, or if all the player's pieces are obstructed from moving by opponent pieces.
- A game is a draw if neither opponent has the possibility to win the game.
- The game is considered a draw when the same position repeats itself for the third time (not necessarily consecutive), with the same player having the move each time. _this one is not implemented yet I may or may not decide to do this_
- A king-versus-king endgame is automatically declared a draw, as is any other position proven to be a draw

## Building

```
go mod download
go build main.go
```


## Used Packages

https://github.com/faiface/pixel  - used to draw the Board

## License

BSD-2-Clause