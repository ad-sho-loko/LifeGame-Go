package main

import (
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

const(
	BackGround = termbox.ColorBlack
	Alive = termbox.ColorWhite
	Interval = 100 * time.Millisecond
)

type Game struct{
	width int
	height int
	field [][]bool
	queue chan termbox.Event
}

func newField(h, w int) [][]bool{
	row := make([][]bool, h)
	for i:=0; i<h; i++{
		row[i] = make([]bool, w)
	}
	return row
}

func pollEvent() chan termbox.Event{
	q := make(chan termbox.Event)
	go func(){
		for{
			q <- termbox.PollEvent()
		}
	}()
	return q
}

func NewGame(h, w int) *Game{
	return &Game{
		height:h,
		width:w,
		field:newField(h, w),
		queue:pollEvent(),
	}
}

func NewGameRand(h, w int) *Game{
	rand.Seed(time.Now().UnixNano())
	g := NewGame(h, w)
	for i:=0; i<h; i++{
		for j:=0; j<w; j++{
			if rand.Intn(10) == 0{
				g.field[i][j] = true
			}else{
				g.field[i][j] = false
			}
		}
	}
	return g
}

func (g *Game) count(r int, c int) int{
	if r < 0 || c < 0 || r >= g.height || c >= g.width{
		return 0
	}

	if g.field[r][c]{
		return 1
	}

	return 0
}

func (g *Game) countNear(r int, c int) int{
	var sum = 0
	for i:=-1; i<2; i++{
		for j:=-1; j<2; j++{
			if i == 0 && j == 0{
				continue
			}
			sum += g.count(r+i, c+j)
		}
	}
	return sum
}

func (g *Game) UpdateField() {
	var field = newField(g.height, g.width)
	var count = 0
	for i:=0; i<g.height; i++{
		for j:=0; j<g.width; j++{
			// 近傍をチェックする
			count = g.countNear(i, j)

			if !g.field[i][j]{
				// 誕生
				field[i][j] = count == 3
			}else{
				// 生存, 過疎, 過密
				field[i][j] = count == 2 || count == 3
			}
		}
	}
	g.field = field
}

func (g *Game) RenderField(){
	termbox.Clear(BackGround, BackGround)
	for i:=0; i<g.height; i++{
		for j:=0; j<g.width; j++{
			if g.field[i][j] {
				termbox.SetCell(i, j, '█', Alive, BackGround)
			}
		}
	}
	termbox.Flush()
}


func main(){
	if err := termbox.Init(); err != nil{
		panic(err)
	}
	defer termbox.Close()
	w, h := termbox.Size()

	g := NewGameRand(h, w)

	for{
		select{

		case ev:= <-g.queue:
			if ev.Key == termbox.KeyEsc{
				return
			}

		default:
			time.Sleep(Interval)
			g.RenderField()
			g.UpdateField()
		}
	}
}
