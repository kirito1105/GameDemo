package world

import (
	"bufio"
	"image"
	"image/color"
	"image/png"
	rand "myGameDemo/myRand"
	"os"
	"sync"
	"time"
)

const (
	Size           = 40
	GRID_PER_BLOCK = 10
	PIONT_PER_GRID = 100
)

type Point struct {
	BlockX int
	BlockY int
	GridX  int
	GridY  int
	X      int
	Y      int
}

func NewPoint(x int, y int) *Point {
	p := Point{}
	p.BlockX = x / (PIONT_PER_GRID * GRID_PER_BLOCK)
	p.BlockY = y / (PIONT_PER_GRID * GRID_PER_BLOCK)
	x = x % (PIONT_PER_GRID * GRID_PER_BLOCK)
	y = y % (PIONT_PER_GRID * GRID_PER_BLOCK)
	p.GridX = x / PIONT_PER_GRID
	p.GridY = y / PIONT_PER_GRID
	p.X = x % PIONT_PER_GRID
	p.Y = y % PIONT_PER_GRID
	return &p
}

func (p Point) ToUnity() (int, int) {
	x := p.X + p.GridX*PIONT_PER_GRID + p.BlockX*PIONT_PER_GRID*GRID_PER_BLOCK
	y := p.Y + p.GridY*PIONT_PER_GRID + p.BlockY*PIONT_PER_GRID*GRID_PER_BLOCK
	return x, y
}

type BlockCreate struct {
	world [Size][Size]bool
	count [Size][Size]int
}

type World struct {
	blocks [Size][Size]Block
	spawn  *Point
	num    int
}
type Block struct {
	typeOfBlock int
	Objs        []Obj
}
type Obj struct {
	x       int
	y       int
	objType string
}

func CreateWorld(blocks [Size][Size]bool) *World {
	w := World{}
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			if blocks[i][j] {
				w.blocks[i][j].typeOfBlock = 0
			} else {
				w.blocks[i][j].typeOfBlock = 1
			}
		}
	}
	return &w
}

func NewWorld() *World {
	w := BlockCreate{}
	w.Init()
	for i := 2; i < 6; i++ {
		w.Loop(3, 4)
	}
	for i := 6; i < 10; i++ {
		w.Loop(0, 4)
	}
	world1 := CreateWorld(w.GetWorld())
	world1.Init()
	return world1
}

func (this *World) Init() {
	var once sync.Once
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			if this.blocks[i][j].typeOfBlock == 0 {
				continue
			}
			for x := 0; x < GRID_PER_BLOCK; x++ {
				for y := 0; y < GRID_PER_BLOCK; y++ {
					r := rand.Intn(10000)
					num := 0
					for _, rate := range *GetList() {
						if r < num+rate.rate {
							this.blocks[i][j].Objs = append(this.blocks[i][j].Objs, Obj{
								objType: rate.objType,
								x:       PIONT_PER_GRID/2 + x*GRID_PER_BLOCK,
								y:       PIONT_PER_GRID/2 + y*GRID_PER_BLOCK,
							})
							this.num++
							break
						}
						num = num + rate.rate
					}
					once.Do(func() {
						this.spawn = &Point{
							BlockX: i,
							BlockY: j,
							GridX:  x,
							GridY:  y,
							X:      50,
							Y:      50,
						}
					})
				}
			}
		}
	}
}

func (this *World) GetBlock(x int, y int) Block {
	return this.blocks[x][y]
}

func (this *World) GetSpawn() Point {
	return *this.spawn
}

func (this *BlockCreate) Init() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			num := rand.Intn(100)
			if num < 40 {
				this.world[i][j] = true
			}
		}
	}
	for i := 0; i < Size; i++ {
		this.world[0][i] = true
		this.world[Size-1][i] = true
		this.world[i][Size-1] = true
		this.world[i][0] = true
	}
}

func (this *BlockCreate) CountPoint(i int, j int) int {
	num := 0
	if this.world[i-1][j-1] {
		num++
	}
	if this.world[i-1][j] {
		num++
	}
	if this.world[i-1][j+1] {
		num++
	}
	if this.world[i][j+1] {
		num++
	}
	if this.world[i+1][j+1] {
		num++
	}
	if this.world[i+1][j] {
		num++
	}
	if this.world[i+1][j-1] {
		num++
	}
	if this.world[i][j-1] {
		num++
	}
	return num
}

func (this *BlockCreate) ToCount() {
	for i := 1; i < Size-1; i++ {
		for j := 1; j < Size-1; j++ {

			this.count[i][j] = this.CountPoint(i, j)
		}
	}
}

func (this *BlockCreate) Loop(min int, max int) {
	this.ToCount()
	for i := 1; i < Size-1; i++ {
		for j := 1; j < Size-1; j++ {
			if this.count[i][j] > max || this.count[i][j] < min {
				this.world[i][j] = true
			} else if this.count[i][j] == 5 {
			} else {
				this.world[i][j] = false
			}
		}
	}
}
func (this *BlockCreate) LoopWithOutCount(min int, max int) {
	for i := 1; i < Size-1; i++ {
		for j := 1; j < Size-1; j++ {
			if this.CountPoint(i, j) > max || this.CountPoint(i, j) < min {
				this.world[i][j] = true
			} else if this.count[i][j] == 5 {
			} else {
				this.world[i][j] = false
			}
		}
	}
}

func (this *BlockCreate) GetWorld() [Size][Size]bool {
	return this.world
}

func (this *BlockCreate) ToImage(str string) {
	img := image.NewGray(image.Rect(0, 0, Size, Size))
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			if !this.GetWorld()[x][y] {
				img.Set(x, y, color.White)
			}
		}
	}
	f, err := os.Create(str)
	if err != nil {
		return
	}
	b := bufio.NewWriter(f)
	err = png.Encode(b, img)
	if err != nil {
		return
	}
	defer f.Close()
	b.Flush()
}
