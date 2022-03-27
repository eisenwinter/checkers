package game

//getBackRowCount gets the number of pieces in the backrow heustric
func (b Board) getBackRowCount() (white int, red int) {
	white = 0
	red = 0
	//red backrow
	for i := 0; i < width; i++ {
		if !has(b[i], Empty) {
			//only if backrow is not a king
			if !has(b[i], Player) && !has(b[i], King) {
				red++
			}
		}
	}
	//white backrow
	for i := len(b) - width - 1; i < len(b); i++ {
		if !has(b[i], Empty) {
			if has(b[i], Player) && !has(b[i], King) {
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
	fromRow := middleRow - 1
	for i := fromRow; i <= (fromRow + 1); i++ {
		//width without left and right side  -> enemy can only pass on the side
		for j := 3; j <= (width - 4); j++ {
			idx := IndexOf(i, j)
			if !has(b[idx], Empty) {
				if has(b[i], Player) {
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
	fromRow := middleRow - 1
	for i := fromRow; i <= (fromRow + 1); i++ {
		for j := 0; j < 3; j++ {
			idx := IndexOf(i, j)
			if !has(b[idx], Empty) {
				if has(b[i], Player) {
					white++
				} else {
					red++
				}
			}
		}
		for j := width - 3; j < width; j++ {
			idx := IndexOf(i, j)
			if !has(b[idx], Empty) {
				if has(b[i], Player) {
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
					!has(b[IndexOf(r-1, c-1)], Empty) &&
					has(b[IndexOf(r-1, c-1)], Player) &&
					has(b[IndexOf(r+1, c+1)], Empty) {
					red++
					if has(b[IndexOf(r, c)], King) {
						rking++
					}
				} else if b.canDrawTo(r-1, c+1) &&
					b.canDrawTo(r+1, c-1) &&
					!has(b[IndexOf(r-1, c+1)], Empty) &&
					has(b[IndexOf(r-1, c+1)], Player) &&
					has(b[IndexOf(r+1, c-1)], Empty) {
					red++
					if has(b[IndexOf(r, c)], King) {
						rking++
					}
				} else if b.canDrawTo(r+1, c+1) &&
					b.canDrawTo(r-1, c-1) &&
					!has(b[IndexOf(r+1, c+1)], Empty) &&
					has(b[IndexOf(r+1, c+1)], Player) &&
					has(b[IndexOf(r-1, c-1)], Empty) &&
					has(b[IndexOf(r+1, c+1)], King) {
					red++
					if has(b[IndexOf(r, c)], King) {
						rking++
					}
				} else if b.canDrawTo(r+1, c-1) &&
					b.canDrawTo(r-1, c+1) &&
					!has(b[IndexOf(r+1, c-1)], Empty) &&
					has(b[IndexOf(r+1, c-1)], Player) &&
					has(b[IndexOf(r-1, c+1)], Empty) &&
					has(b[IndexOf(r+1, c-1)], King) {
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
