package common

type Board [BoardSize][BoardSize]int

func CheckWin(b *Board, x int, y int, color int) bool {
	var dir = [4][2]int{
		{1, 0},
		{0, 1},
		{1, 1},
		{1, -1},
	}
	for _, d := range dir {
		count := 1
		count += countDir(b, x, y, color, d[0], d[1])
		count += countDir(b, x, y, color, -d[0], -d[1]) // 反方向
		if count >= WinCount {
			return true
		}
	}
	return false
}

func countDir(b *Board, x, y int, color, dx, dy int) int {
	n := 0
	for nx, ny := x+dx, y+dy; inBoard(nx, ny) && b[nx][ny] == color; nx, ny = nx+dx, ny+dy {
		n++
	}
	return n
}

func inBoard(x int, y int) bool {
	return x >= 0 && x < BoardSize && y >= 0 && y < BoardSize
}
