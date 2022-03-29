package game

//getBackRowCount gets the number of pieces in the backrow heustric
func (b Board) getBackRowCount() (white int, red int) {
	white = 0
	red = 0
	//red backrow
	for i := 0; i < width; i++ {
		idx := IndexOf(0, i)
		if !has(b[idx], Empty) && !has(b[idx], King) {
			if !has(b[idx], Player) {
				red++
			}
		}
	}
	//white backrow
	for i := 0; i < width; i++ {
		idx := IndexOf((height - 1), i)
		if !has(b[idx], Empty) {
			if has(b[idx], Player) && !has(b[idx], King) {
				white++
			}
		}
	}
	return
}

//getMiddleBoxCount gets the number of pieces in the middle box
func (b Board) getMiddleBoxCount() (white int, red int) {
	white = 0
	red = 0
	middleRow := (height / 2) - 1
	for i := middleRow; i <= (middleRow + 1); i++ {
		//width without left and right side  -> enemy can only pass on the side
		for j := 2; j <= (width - 3); j++ {
			idx := IndexOf(i, j)
			if !has(b[idx], Empty) {
				if has(b[idx], Player) {
					white++
				} else {
					red++
				}
			}
		}
	}
	return
}

//getMiddleRowSideCount gets the number of pieces in the of the middle rows without the box
func (b Board) getMiddleRowSideCount() (white int, red int) {
	white = 0
	red = 0
	middleRow := (height / 2) - 1
	for i := middleRow; i <= (middleRow + 1); i++ {
		for j := 0; j < 2; j++ {
			idx := IndexOf(i, j)
			if !has(b[idx], Empty) {
				if has(b[idx], Player) {
					white++
				} else {
					red++
				}
			}
		}
		for j := width - 2; j < width; j++ {
			idx := IndexOf(i, j)
			if !has(b[idx], Empty) {
				if has(b[idx], Player) {
					white++
				} else {
					red++
				}
			}
		}
	}
	return
}

// getProtectedPieceCount returns the count of protected pieces for the heuristic
func (b Board) getProtectionCount() (white int, red int) {
	white = 0
	red = 0
	for i, v := range b {
		if !has(v, Empty) {
			r, c := reverseIndexOf(i)
			if has(v, Player) {
				//white
				//on edge
				if c == 0 || c == (width-1) {
					white++
				} else if r == 0 || r == (height-1) {
					white++
				} else if b.canDrawTo(r-1, c-1) &&
					!has(b[IndexOf(r-1, c-1)], Empty) &&
					has(b[IndexOf(r-1, c-1)], Player) {
					white++
				} else if b.canDrawTo(r-1, c+1) &&
					!has(b[IndexOf(r-1, c+1)], Empty) &&
					has(b[IndexOf(r-1, c+1)], Player) {
					white++
				} else if b.canDrawTo(r+1, c+1) &&
					!has(b[IndexOf(r+1, c+1)], Empty) &&
					has(b[IndexOf(r+1, c+1)], Player) {
					white++
				} else if b.canDrawTo(r+1, c-1) &&
					!has(b[IndexOf(r+1, c-1)], Empty) &&
					has(b[IndexOf(r+1, c-1)], Player) {
					white++
				}

			} else {
				//red
				if c == 0 || c == (width-1) {
					red++
				} else if r == 0 || r == (height-1) {
					red++
				} else if b.canDrawTo(r-1, c-1) &&
					!has(b[IndexOf(r-1, c-1)], Empty) &&
					!has(b[IndexOf(r-1, c-1)], Player) {
					red++
				} else if b.canDrawTo(r-1, c+1) &&
					!has(b[IndexOf(r-1, c+1)], Empty) &&
					!has(b[IndexOf(r-1, c+1)], Player) {
					red++
				} else if b.canDrawTo(r+1, c+1) &&
					!has(b[IndexOf(r+1, c+1)], Empty) &&
					!has(b[IndexOf(r+1, c+1)], Player) {
					red++
				} else if b.canDrawTo(r+1, c-1) &&
					!has(b[IndexOf(r+1, c-1)], Empty) &&
					!has(b[IndexOf(r+1, c-1)], Player) {
					red++
				}
			}
		}
	}
	return
}

// getVulnerablePieceCount returns the count of vulnerbale pieces for the heuristic
func (b Board) getVulnerablePieceCount() (white int, red int, wking int, rking int) {
	white = 0
	red = 0
	wking = 0
	rking = 0
	for i, v := range b {
		if !has(v, Empty) {
			r, c := reverseIndexOf(i)
			if has(v, Player) {
				//white
				//can be taken from diagonially right
				if b.canDrawTo(r-1, c-1) &&
					b.canDrawTo(r+1, c+1) &&
					!has(b[IndexOf(r-1, c-1)], Empty) &&
					!has(b[IndexOf(r-1, c-1)], Player) &&
					has(b[IndexOf(r+1, c+1)], Empty) {
					white++
					if has(b[IndexOf(r, c)], King) {
						wking++
					}
				} else if b.canDrawTo(r-1, c+1) && //can be taken diagnionally left
					b.canDrawTo(r+1, c-1) &&
					!has(b[IndexOf(r-1, c+1)], Empty) &&
					!has(b[IndexOf(r-1, c+1)], Player) &&
					has(b[IndexOf(r+1, c-1)], Empty) {
					white++
					if has(b[IndexOf(r, c)], King) {
						wking++
					}
				} else if b.canDrawTo(r+1, c+1) && //can be taken diagionally left reverse (king)
					b.canDrawTo(r-1, c-1) &&
					!has(b[IndexOf(r+1, c+1)], Empty) &&
					!has(b[IndexOf(r+1, c+1)], Player) &&
					has(b[IndexOf(r-1, c-1)], Empty) &&
					has(b[IndexOf(r+1, c+1)], King) {
					white++
					if has(b[IndexOf(r, c)], King) {
						wking++
					}
				} else if b.canDrawTo(r+1, c-1) &&
					b.canDrawTo(r-1, c+1) &&
					!has(b[IndexOf(r+1, c-1)], Empty) &&
					!has(b[IndexOf(r+1, c-1)], Player) &&
					has(b[IndexOf(r-1, c+1)], Empty) &&
					has(b[IndexOf(r+1, c-1)], King) {
					white++
					if has(b[IndexOf(r, c)], King) {
						wking++
					}
				}
			} else {
				//red
				if b.canDrawTo(r-1, c-1) &&
					b.canDrawTo(r+1, c+1) &&
					!has(b[IndexOf(r+1, c-1)], Empty) &&
					has(b[IndexOf(r+1, c-1)], Player) &&
					has(b[IndexOf(r-1, c+1)], Empty) {
					red++
					if has(b[IndexOf(r, c)], King) {
						rking++
					}
				} else if b.canDrawTo(r-1, c+1) &&
					b.canDrawTo(r+1, c-1) &&
					!has(b[IndexOf(r+1, c+1)], Empty) &&
					has(b[IndexOf(r+1, c+1)], Player) &&
					has(b[IndexOf(r-1, c-1)], Empty) {
					red++
					if has(b[IndexOf(r, c)], King) {
						rking++
					}
				} else if b.canDrawTo(r+1, c+1) &&
					b.canDrawTo(r-1, c-1) &&
					!has(b[IndexOf(r-1, c+1)], Empty) &&
					has(b[IndexOf(r-1, c+1)], Player) &&
					has(b[IndexOf(r+1, c-1)], Empty) &&
					has(b[IndexOf(r-1, c+1)], King) {
					red++
					if has(b[IndexOf(r, c)], King) {
						rking++
					}
				} else if b.canDrawTo(r+1, c-1) &&
					b.canDrawTo(r-1, c+1) &&
					!has(b[IndexOf(r-1, c-1)], Empty) &&
					has(b[IndexOf(r-1, c-1)], Player) &&
					has(b[IndexOf(r+1, c+1)], Empty) &&
					has(b[IndexOf(r-1, c-1)], King) {
					red++
					if has(b[IndexOf(r, c)], King) {
						rking++
					}
				}
			}
		}
	}
	return
}

//getFortressCount returns absolute fortied pieces
func (b Board) getFortressCount() (white int, red int) {
	white = 0
	red = 0
	isSafe := func(r, c int, player bool) bool {
		if player {
			return !b.canDrawTo(r, c) || (!has(b[IndexOf(r, c)], Empty) && has(b[IndexOf(r, c)], Player))
		}
		return !b.canDrawTo(r, c) || (!has(b[IndexOf(r, c)], Empty) && !has(b[IndexOf(r, c)], Player))
	}
	for i, v := range b {
		if !has(v, Empty) {
			r, c := reverseIndexOf(i)
			if has(v, Player) {
				//white
				if isSafe(r-1, c-1, true) && isSafe(r-1, c+1, true) && isSafe(r+1, c+1, true) && isSafe(r+1, c-1, true) {
					white++
				}
			} else {
				if isSafe(r-1, c-1, false) && isSafe(r-1, c+1, false) && isSafe(r+1, c+1, false) && isSafe(r+1, c-1, false) {
					red++
				}

			}
		}
	}
	return
}

//getDiamondShapes returns piece count participating in diamond shapes
func (b Board) getDiamondShapes() (white int, red int) {
	white = 0
	red = 0
	isSafe := func(r, c int, player bool) bool {
		if player {
			return !b.canDrawTo(r, c) || (!has(b[IndexOf(r, c)], Empty) && has(b[IndexOf(r, c)], Player))
		}
		return !b.canDrawTo(r, c) || (!has(b[IndexOf(r, c)], Empty) && !has(b[IndexOf(r, c)], Player))
	}
	for i, v := range b {
		if !has(v, Empty) {
			r, c := reverseIndexOf(i)
			if has(v, Player) {
				//white
				if isSafe(r-1, c-1, true) && isSafe(r+1, c-1, true) && isSafe(r, c-2, true) {
					white++
				}
				if isSafe(r-1, c+1, true) && isSafe(r+1, c+1, true) && isSafe(r, c+2, true) {
					white++
				}
				if isSafe(r-1, c-1, true) && isSafe(r-1, c+1, true) && isSafe(r-2, c, true) {
					white++
				}
				if isSafe(r+1, c-1, true) && isSafe(r+1, c+1, true) && isSafe(r+2, c, true) {
					white++
				}
			} else {
				if isSafe(r-1, c-1, false) && isSafe(r+1, c-1, false) && isSafe(r, c-2, false) {
					red++
				}
				if isSafe(r-1, c+1, false) && isSafe(r+1, c+1, false) && isSafe(r, c+2, false) {
					red++
				}
				if isSafe(r-1, c-1, false) && isSafe(r-1, c+1, false) && isSafe(r-2, c, false) {
					red++
				}
				if isSafe(r+1, c-1, false) && isSafe(r+1, c+1, false) && isSafe(r+2, c, false) {
					red++
				}

			}
		}
	}
	return
}

func findRunAway(b Board, r, c int, up bool) bool {
	if !b.canDrawTo(r, c) {
		return false
	}
	i := IndexOf(r, c)
	if !has(b[i], Empty) {
		return false
	}
	if up && r == 0 && has(b[i], Empty) {
		return true
	}
	if !up && r == height-1 && has(b[i], Empty) {
		return true
	}
	if up {
		return findRunAway(b, r-1, c-1, up) || findRunAway(b, r-1, c+1, up)
	} else {
		return findRunAway(b, r+1, c-1, up) || findRunAway(b, r+1, c+1, up)
	}
}

//getRunAwayCount  unimpeded path to be kinged
func (b Board) getRunAwayCount() (white int, red int) {
	white = 0
	red = 0
	for i, v := range b {
		if !has(v, Empty) && !has(v, King) {
			r, c := reverseIndexOf(i)
			if has(v, Player) {
				if findRunAway(b, r-1, c+1, true) || findRunAway(b, r-1, c-1, true) {
					white++
				}
			} else {
				if findRunAway(b, r+1, c+1, false) || findRunAway(b, r+1, c-1, false) {
					red++
				}
			}
		}
	}
	return
}

//getStuckPiecesCount gets the number of pieces in the stuck piecies heustric
func (b Board) getStuckPiecesCount() (white int, red int, wking int, rking int) {
	white = 0
	red = 0
	wking = 0
	rking = 0
	for i, v := range b {
		if !has(v, Empty) {
			r, c := reverseIndexOf(i)
			if !has(v, King) && has(v, Player) &&
				(!b.canDrawTo(r-1, c-1) || !has(b[IndexOf(r-1, c-1)], Empty)) &&
				(!b.canDrawTo(r-1, c+1) || !has(b[IndexOf(r-1, c+1)], Empty)) {
				white++

			}
			if !has(v, King) && !has(v, Player) &&
				(!b.canDrawTo(r+1, c-1) || !has(b[IndexOf(r+1, c-1)], Empty)) &&
				(!b.canDrawTo(r+1, c+1) || !has(b[IndexOf(r+1, c+1)], Empty)) {
				red++

			}
			if has(v, King) &&
				(!b.canDrawTo(r-1, c-1) || !has(b[IndexOf(r-1, c-1)], Empty)) &&
				(!b.canDrawTo(r-1, c+1) || !has(b[IndexOf(r-1, c+1)], Empty)) &&
				(!b.canDrawTo(r+1, c-1) || !has(b[IndexOf(r+1, c-1)], Empty)) &&
				(!b.canDrawTo(r+1, c+1) || !has(b[IndexOf(r+1, c+1)], Empty)) {
				if has(v, Player) {
					wking++
				} else {
					rking++
				}
			}

		}
	}
	return white, red, wking, rking
}
