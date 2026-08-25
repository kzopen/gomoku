package common

type Board [BoardSize][BoardSize]int

// Point 棋盘坐标（FindWinLine 返回值）
type Point struct {
	X int
	Y int
}

// FindWinLine 返回 (x,y) 落 color 后形成的连线坐标（含落子点，≥WinCount 个）。
// 未成连返回 nil。判胜与制胜连线下发的唯一核心实现。
func FindWinLine(b *Board, x, y int, color int) []Point {
	var dir = [4][2]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}}
	for _, d := range dir {
		line := []Point{{x, y}}
		// 正向
		for nx, ny := x+d[0], y+d[1]; inBoard(nx, ny) && b[nx][ny] == color; nx, ny = nx+d[0], ny+d[1] {
			line = append(line, Point{nx, ny})
		}
		// 反向（前插，保持连线一端到另一端的顺序）
		for nx, ny := x-d[0], y-d[1]; inBoard(nx, ny) && b[nx][ny] == color; nx, ny = nx-d[0], ny-d[1] {
			line = append([]Point{{nx, ny}}, line...)
		}
		if len(line) >= WinCount {
			return line
		}
	}
	return nil
}

// CheckWin 是否五连：FindWinLine 的薄封装（调用方只需 bool 时用）
func CheckWin(b *Board, x, y int, color int) bool {
	return FindWinLine(b, x, y, color) != nil
}

func inBoard(x int, y int) bool {
	return x >= 0 && x < BoardSize && y >= 0 && y < BoardSize
}
