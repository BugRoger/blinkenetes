package main

func (blink *blinkenetes) showPacMan() {
	var grid [8][8][3]int

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			grid[x][y] = [3]int{0, 0, 0}
		}
	}

	grid[0][0] = [3]int{0, 0, 0}
	grid[1][0] = [3]int{0, 0, 0}
	grid[2][0] = [3]int{0, 0, 16}
	grid[3][0] = [3]int{4, 4, 63}
	grid[4][0] = [3]int{4, 4, 63}
	grid[5][0] = [3]int{4, 4, 63}
	grid[6][0] = [3]int{0, 0, 0}
	grid[7][0] = [3]int{0, 0, 0}

	grid[0][1] = [3]int{0, 0, 0}
	grid[1][1] = [3]int{0, 0, 16}
	grid[2][1] = [3]int{4, 4, 63}
	grid[3][1] = [3]int{4, 4, 63}
	grid[4][1] = [3]int{4, 4, 63}
	grid[5][1] = [3]int{4, 4, 63}
	grid[6][1] = [3]int{4, 4, 63}
	grid[7][1] = [3]int{0, 0, 0}

	grid[0][2] = [3]int{0, 0, 16}
	grid[1][2] = [3]int{32, 32, 32}
	grid[2][2] = [3]int{0, 0, 0}
	grid[3][2] = [3]int{4, 4, 63}
	grid[4][2] = [3]int{63, 63, 63}
	grid[5][2] = [3]int{0, 0, 0}
	grid[6][2] = [3]int{4, 4, 63}
	grid[7][2] = [3]int{4, 4, 63}

	grid[0][3] = [3]int{0, 0, 16}
	grid[1][3] = [3]int{32, 32, 32}
	grid[2][3] = [3]int{32, 32, 32}
	grid[3][3] = [3]int{4, 4, 63}
	grid[4][3] = [3]int{63, 63, 63}
	grid[5][3] = [3]int{63, 63, 63}
	grid[6][3] = [3]int{4, 4, 63}
	grid[7][3] = [3]int{4, 4, 63}

	grid[0][4] = [3]int{0, 0, 16}
	grid[1][4] = [3]int{0, 0, 16}
	grid[2][4] = [3]int{4, 4, 63}
	grid[3][4] = [3]int{4, 4, 63}
	grid[4][4] = [3]int{4, 4, 63}
	grid[5][4] = [3]int{4, 4, 63}
	grid[6][4] = [3]int{4, 4, 63}
	grid[7][4] = [3]int{4, 4, 63}

	grid[0][5] = [3]int{0, 0, 16}
	grid[1][5] = [3]int{0, 0, 16}
	grid[2][5] = [3]int{4, 4, 63}
	grid[3][5] = [3]int{4, 4, 63}
	grid[4][5] = [3]int{4, 4, 63}
	grid[5][5] = [3]int{4, 4, 63}
	grid[6][5] = [3]int{4, 4, 63}
	grid[7][5] = [3]int{4, 4, 63}

	grid[0][6] = [3]int{0, 0, 16}
	grid[1][6] = [3]int{0, 0, 16}
	grid[2][6] = [3]int{4, 4, 63}
	grid[3][6] = [3]int{4, 4, 63}
	grid[4][6] = [3]int{4, 4, 63}
	grid[5][6] = [3]int{4, 4, 63}
	grid[6][6] = [3]int{4, 4, 63}
	grid[7][6] = [3]int{4, 4, 63}

	grid[0][7] = [3]int{0, 0, 16}
	grid[1][7] = [3]int{0, 0, 16}
	grid[2][7] = [3]int{4, 4, 63}
	grid[3][7] = [3]int{0, 0, 0}
	grid[4][7] = [3]int{4, 4, 63}
	grid[5][7] = [3]int{4, 4, 63}
	grid[6][7] = [3]int{0, 0, 0}
	grid[7][7] = [3]int{4, 4, 63}

	for x, col := range grid {
		for y, cell := range col {
			blink.pad.Light(x+1, y+1, cell[0], cell[1], cell[2])
		}
	}
}
