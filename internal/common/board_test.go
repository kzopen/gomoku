package common

import "testing"

func makeBoard(pieces [][3]int) *Board {
	b := &Board{}
	for _, p := range pieces {
		b[p[0]][p[1]] = p[2]
	}
	return b
}

func TestCheckWin(t *testing.T) {
	tests := []struct {
		name   string
		pieces [][3]int
		x, y   int
		color  int
		want   bool
	}{
		{"水平五连", [][3]int{{3, 7, 1}, {4, 7, 1}, {5, 7, 1}, {6, 7, 1}, {7, 7, 1}}, 5, 7, 1, true},
		{"垂直五连", [][3]int{{7, 3, 2}, {7, 4, 2}, {7, 5, 2}, {7, 6, 2}, {7, 7, 2}}, 7, 5, 2, true},
		{"主对角五连", [][3]int{{3, 3, 1}, {4, 4, 1}, {5, 5, 1}, {6, 6, 1}, {7, 7, 1}}, 5, 5, 1, true},
		{"副对角五连", [][3]int{{3, 7, 2}, {4, 6, 2}, {5, 5, 2}, {6, 4, 2}, {7, 3, 2}}, 5, 5, 2, true},
		{"水平四连非胜", [][3]int{{3, 7, 1}, {4, 7, 1}, {5, 7, 1}, {6, 7, 1}}, 5, 7, 1, false},
		{"六连长连仍算胜", [][3]int{{2, 7, 1}, {3, 7, 1}, {4, 7, 1}, {5, 7, 1}, {6, 7, 1}, {7, 7, 1}}, 5, 7, 1, true},
		{"贴左上角", [][3]int{{0, 0, 1}, {1, 0, 1}, {2, 0, 1}, {3, 0, 1}, {4, 0, 1}}, 0, 0, 1, true},
		{"贴右下角", [][3]int{{10, 14, 2}, {11, 14, 2}, {12, 14, 2}, {13, 14, 2}, {14, 14, 2}}, 14, 14, 2, true},
		{"颜色参数与落子不符", [][3]int{{3, 7, 1}, {4, 7, 1}, {5, 7, 1}, {6, 7, 1}, {7, 7, 1}}, 5, 7, 2, false},
		{"空棋盘落子", nil, 7, 7, 1, false},
		{"交叉各4连均不足5", [][3]int{{5, 4, 1}, {6, 4, 1}, {7, 4, 1}, {4, 5, 1}, {4, 6, 1}, {4, 7, 1}}, 4, 4, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := makeBoard(tt.pieces)
			if got := CheckWin(b, tt.x, tt.y, tt.color); got != tt.want {
				t.Errorf("CheckWin(x=%d, y=%d, color=%d) = %v, want %v", tt.x, tt.y, tt.color, got, tt.want)
			}
		})
	}
}
