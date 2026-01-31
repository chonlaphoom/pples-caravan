package mapregion

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	// Use standard 8-color SGR sequences
	N = "\x1b[36m" // North -> Cyan
	Y = "\x1b[33m" // Central -> Yellow
	I = "\x1b[33m" // Isan -> Yellow (approx Orange)
	E = "\x1b[32m" // East -> Green
	S = "\x1b[31m" // South -> Red
	W = "\x1b[35m" // West -> Magenta (approx Pink)
	X = "\x1b[0m"  // Reset

	SPACE_LEN = 4
	MAX_COLS  = 10
)

type MapRegion struct {
	Grid [][10]string
	Size Position
}

type Position struct {
	Row, Col int
}

type province struct {
	FullName  string
	ShortName string
	Color     string

	Pos Position
}

var provinces = []province{
	// --- NORTHERN (N)
	{ShortName: "ชร", Color: N, Pos: Position{Row: 2, Col: 1}, FullName: "เชียงราย"},
	{ShortName: "ชร", Color: N, Pos: Position{Row: 2, Col: 2}, FullName: "เชียงราย"},
	{ShortName: "มส", Color: N, Pos: Position{Row: 3, Col: 0}, FullName: "แม่ฮ่องสอน"},
	{ShortName: "ชม", Color: N, Pos: Position{Row: 3, Col: 1}, FullName: "เชียงใหม่"},
	{ShortName: "พย", Color: N, Pos: Position{Row: 3, Col: 2}, FullName: "พะเยา"},
	{ShortName: "นน", Color: N, Pos: Position{Row: 3, Col: 3}, FullName: "น่าน"},
	{ShortName: "มส", Color: N, Pos: Position{Row: 4, Col: 0}, FullName: "แม่ฮ่องสอน"},
	{ShortName: "ชม", Color: N, Pos: Position{Row: 4, Col: 1}, FullName: "เชียงใหม่"},
	{ShortName: "พย", Color: N, Pos: Position{Row: 4, Col: 2}, FullName: "พะเยา"},
	{ShortName: "นน", Color: N, Pos: Position{Row: 4, Col: 3}, FullName: "น่าน"},
	{ShortName: "ลพ", Color: N, Pos: Position{Row: 5, Col: 1}, FullName: "ลำพูน"},
	{ShortName: "ลป", Color: N, Pos: Position{Row: 5, Col: 2}, FullName: "ลำปาง"},
	{ShortName: "พร", Color: N, Pos: Position{Row: 5, Col: 3}, FullName: "แพร่"},
	{ShortName: "ลพ", Color: N, Pos: Position{Row: 6, Col: 1}, FullName: "ลำพูน"},
	{ShortName: "สท", Color: N, Pos: Position{Row: 6, Col: 2}, FullName: "สุโขทัย"},
	{ShortName: "อต", Color: N, Pos: Position{Row: 6, Col: 3}, FullName: "อุตรดิตถ์"},

	// --- UPPER ISAN & NORTHEASTERN (I)
	{ShortName: "ตก", Color: N, Pos: Position{Row: 7, Col: 1}, FullName: "ตาก"},
	{ShortName: "สข", Color: N, Pos: Position{Row: 7, Col: 2}, FullName: "สุโขทัย"},
	{ShortName: "อต", Color: N, Pos: Position{Row: 7, Col: 3}, FullName: "อุตรดิตถ์"},
	{ShortName: "บก", Color: I, Pos: Position{Row: 7, Col: 8}, FullName: "บึงกาฬ"},
	{ShortName: "นพ", Color: I, Pos: Position{Row: 7, Col: 9}, FullName: "นครพนม"},
	{ShortName: "ตก", Color: N, Pos: Position{Row: 8, Col: 1}, FullName: "ตาก"},
	{ShortName: "กพ", Color: N, Pos: Position{Row: 8, Col: 2}, FullName: "กำแพงเพชร"},
	{ShortName: "พล", Color: N, Pos: Position{Row: 8, Col: 3}, FullName: "พิษณุโลก"},
	{ShortName: "ลย", Color: I, Pos: Position{Row: 8, Col: 6}, FullName: "เลย"},
	{ShortName: "นค", Color: I, Pos: Position{Row: 8, Col: 7}, FullName: "หนองคาย"},
	{ShortName: "สน", Color: I, Pos: Position{Row: 8, Col: 8}, FullName: "สกลนคร"},
	{ShortName: "พจ", Color: Y, Pos: Position{Row: 9, Col: 2}, FullName: "พิจิตร"},
	{ShortName: "นว", Color: Y, Pos: Position{Row: 9, Col: 3}, FullName: "นครสวรรค์"},
	{ShortName: "พช", Color: I, Pos: Position{Row: 9, Col: 4}, FullName: "เพชรบูรณ์"},
	{ShortName: "หน", Color: I, Pos: Position{Row: 9, Col: 5}, FullName: "หนองบัวลำภู"},
	{ShortName: "อด", Color: I, Pos: Position{Row: 9, Col: 6}, FullName: "อุดรธานี"},
	{ShortName: "กส", Color: I, Pos: Position{Row: 9, Col: 7}, FullName: "กาฬสินธุ์"},
	{ShortName: "มฮ", Color: I, Pos: Position{Row: 9, Col: 8}, FullName: "มุกดาหาร"},

	// --- CENTRAL & WESTERN (Y & W)
	{ShortName: "กจ", Color: W, Pos: Position{Row: 10, Col: 1}, FullName: "กาญจนบุรี"},
	{ShortName: "อน", Color: Y, Pos: Position{Row: 10, Col: 2}, FullName: "อุทัยธานี"},
	{ShortName: "ชน", Color: Y, Pos: Position{Row: 10, Col: 3}, FullName: "ชัยนาท"},
	{ShortName: "สพ", Color: Y, Pos: Position{Row: 10, Col: 4}, FullName: "สุพรรณบุรี"},
	{ShortName: "ชย", Color: I, Pos: Position{Row: 10, Col: 5}, FullName: "ชัยภูมิ"},
	{ShortName: "ขก", Color: I, Pos: Position{Row: 10, Col: 6}, FullName: "ขอนแก่น"},
	{ShortName: "มค", Color: I, Pos: Position{Row: 10, Col: 7}, FullName: "มหาสารคาม"},
	{ShortName: "อำ", Color: I, Pos: Position{Row: 10, Col: 8}, FullName: "อำนาจเจริญ"},

	{ShortName: "กจ", Color: W, Pos: Position{Row: 11, Col: 1}, FullName: "กาญจนบุรี"},
	{ShortName: "อน", Color: Y, Pos: Position{Row: 11, Col: 2}, FullName: "อุทัยธานี"},
	{ShortName: "ชน", Color: Y, Pos: Position{Row: 11, Col: 3}, FullName: "ชัยนาท"},
	{ShortName: "สพ", Color: Y, Pos: Position{Row: 11, Col: 4}, FullName: "สุพรรณบุรี"},
	{ShortName: "ชย", Color: I, Pos: Position{Row: 11, Col: 5}, FullName: "ชัยภูมิ"},
	{ShortName: "ขก", Color: I, Pos: Position{Row: 11, Col: 6}, FullName: "ขอนแก่น"},
	{ShortName: "มค", Color: I, Pos: Position{Row: 11, Col: 7}, FullName: "มหาสารคาม"},
	{ShortName: "อำ", Color: I, Pos: Position{Row: 11, Col: 8}, FullName: "อำนาจเจริญ"},

	{ShortName: "กจ", Color: W, Pos: Position{Row: 12, Col: 1}, FullName: "กาญจนบุรี"},
	{ShortName: "นฐ", Color: Y, Pos: Position{Row: 12, Col: 2}, FullName: "นครปฐม"},
	{ShortName: "อย", Color: Y, Pos: Position{Row: 12, Col: 3}, FullName: "พระนครศรีอยุธยา"},
	{ShortName: "อท", Color: Y, Pos: Position{Row: 12, Col: 4}, FullName: "อ่างทอง"},
	{ShortName: "ลบ", Color: Y, Pos: Position{Row: 12, Col: 5}, FullName: "ลพบุรี"},
	{ShortName: "ขก", Color: I, Pos: Position{Row: 12, Col: 6}, FullName: "ขอนแก่น"},
	{ShortName: "รอ", Color: I, Pos: Position{Row: 12, Col: 7}, FullName: "ร้อยเอ็ด"},
	{ShortName: "ยส", Color: I, Pos: Position{Row: 12, Col: 8}, FullName: "ยโสธร"},

	// --- LOWER ISAN & CENTRAL (I & Y)
	{ShortName: "กจ", Color: W, Pos: Position{Row: 13, Col: 1}, FullName: "กาญจนบุรี"},
	{ShortName: "สพ", Color: Y, Pos: Position{Row: 13, Col: 2}, FullName: "สุพรรณบุรี"},
	{ShortName: "สบ", Color: Y, Pos: Position{Row: 13, Col: 3}, FullName: "สระบุรี"},
	{ShortName: "สบ", Color: Y, Pos: Position{Row: 13, Col: 4}, FullName: "สระบุรี"},
	{ShortName: "นม", Color: I, Pos: Position{Row: 13, Col: 5}, FullName: "นครราชสีมา"},
	{ShortName: "บร", Color: I, Pos: Position{Row: 13, Col: 6}, FullName: "บุรีรัมย์"},
	{ShortName: "สร", Color: I, Pos: Position{Row: 13, Col: 7}, FullName: "สุรินทร์"},
	{ShortName: "อบ", Color: I, Pos: Position{Row: 13, Col: 8}, FullName: "อุบลราชธานี"},

	{ShortName: "กจ", Color: W, Pos: Position{Row: 14, Col: 1}, FullName: "กาญจนบุรี"},
	{ShortName: "สพ", Color: Y, Pos: Position{Row: 14, Col: 2}, FullName: "สุพรรณบุรี"},
	{ShortName: "สบ", Color: Y, Pos: Position{Row: 14, Col: 3}, FullName: "สระบุรี"},
	{ShortName: "สบ", Color: Y, Pos: Position{Row: 14, Col: 4}, FullName: "สระบุรี"},
	{ShortName: "นม", Color: I, Pos: Position{Row: 14, Col: 5}, FullName: "นครราชสีมา"},
	{ShortName: "บร", Color: I, Pos: Position{Row: 14, Col: 6}, FullName: "บุรีรัมย์"},
	{ShortName: "สร", Color: I, Pos: Position{Row: 14, Col: 7}, FullName: "สุรินทร์"},
	{ShortName: "อบ", Color: I, Pos: Position{Row: 14, Col: 8}, FullName: "อุบลราชธานี"},

	{ShortName: "รบ", Color: W, Pos: Position{Row: 15, Col: 1}, FullName: "ราชบุรี"},
	{ShortName: "สส", Color: Y, Pos: Position{Row: 15, Col: 2}, FullName: "สมุทรสงคราม"},
	{ShortName: "กท", Color: Y, Pos: Position{Row: 15, Col: 3}, FullName: "กรุงเทพมหานคร"},
	{ShortName: "นบ", Color: Y, Pos: Position{Row: 15, Col: 4}, FullName: "นนทบุรี"},
	{ShortName: "ปท", Color: Y, Pos: Position{Row: 15, Col: 5}, FullName: "ปทุมธานี"},
	{ShortName: "นม", Color: I, Pos: Position{Row: 15, Col: 6}, FullName: "นครราชสีมา"},
	{ShortName: "ศก", Color: I, Pos: Position{Row: 15, Col: 7}, FullName: "ศรีสะเกษ"},
	{ShortName: "อบ", Color: I, Pos: Position{Row: 15, Col: 8}, FullName: "อุบลราชธานี"},

	// --- CENTRAL, EASTERN & UPPER SOUTH (W, Y, E)
	{ShortName: "พบ", Color: W, Pos: Position{Row: 16, Col: 1}, FullName: "เพชรบุรี"},
	{ShortName: "สค", Color: Y, Pos: Position{Row: 16, Col: 2}, FullName: "สมุทรสาคร"},
	{ShortName: "สป", Color: Y, Pos: Position{Row: 16, Col: 3}, FullName: "สมุทรปราการ"},
	{ShortName: "นย", Color: E, Pos: Position{Row: 16, Col: 4}, FullName: "นครนายก"},
	{ShortName: "ปจ", Color: E, Pos: Position{Row: 16, Col: 5}, FullName: "ปราจีนบุรี"},
	{ShortName: "สก", Color: E, Pos: Position{Row: 16, Col: 6}, FullName: "สระแก้ว"},

	{ShortName: "พบ", Color: W, Pos: Position{Row: 17, Col: 1}, FullName: "เพชรบุรี"},
	{ShortName: "สค", Color: Y, Pos: Position{Row: 17, Col: 2}, FullName: "สมุทรสาคร"},
	{ShortName: "สป", Color: Y, Pos: Position{Row: 17, Col: 3}, FullName: "สมุทรปราการ"},
	{ShortName: "ฉช", Color: E, Pos: Position{Row: 17, Col: 4}, FullName: "ฉะเชิงเทรา"},
	{ShortName: "ชบ", Color: E, Pos: Position{Row: 17, Col: 5}, FullName: "ชลบุรี"},
	{ShortName: "สก", Color: E, Pos: Position{Row: 17, Col: 6}, FullName: "สระแก้ว"},

	{ShortName: "ปข", Color: W, Pos: Position{Row: 18, Col: 1}, FullName: "ประจวบคีรีขันธ์"},
	{ShortName: "ฉช", Color: E, Pos: Position{Row: 18, Col: 4}, FullName: "ฉะเชิงเทรา"},
	{ShortName: "ชบ", Color: E, Pos: Position{Row: 18, Col: 5}, FullName: "ชลบุรี"},

	{ShortName: "ปข", Color: W, Pos: Position{Row: 19, Col: 1}, FullName: "ประจวบคีรีขันธ์"},
	{ShortName: "รย", Color: E, Pos: Position{Row: 19, Col: 4}, FullName: "ระยอง"},
	{ShortName: "จบ", Color: E, Pos: Position{Row: 19, Col: 5}, FullName: "จันทบุรี"},
	{ShortName: "ตร", Color: E, Pos: Position{Row: 19, Col: 6}, FullName: "ตราด"},

	// --- SOUTHERN (S)
	{ShortName: "ชพ", Color: S, Pos: Position{Row: 20, Col: 1}, FullName: "ชุมพร"},
	{ShortName: "รน", Color: S, Pos: Position{Row: 21, Col: 1}, FullName: "ระนอง"},
	{ShortName: "สฎ", Color: S, Pos: Position{Row: 21, Col: 2}, FullName: "สุราษฎร์ธานี"},
	{ShortName: "รน", Color: S, Pos: Position{Row: 22, Col: 1}, FullName: "ระนอง"},
	{ShortName: "สฎ", Color: S, Pos: Position{Row: 22, Col: 2}, FullName: "สุราษฎร์ธานี"},
	{ShortName: "พง", Color: S, Pos: Position{Row: 23, Col: 1}, FullName: "พังงา"},
	{ShortName: "กบ", Color: S, Pos: Position{Row: 23, Col: 2}, FullName: "กระบี่"},
	{ShortName: "นศ", Color: S, Pos: Position{Row: 23, Col: 3}, FullName: "นครศรีธรรมราช"},
	{ShortName: "พง", Color: S, Pos: Position{Row: 24, Col: 1}, FullName: "พังงา"},
	{ShortName: "กบ", Color: S, Pos: Position{Row: 24, Col: 2}, FullName: "กระบี่"},
	{ShortName: "นศ", Color: S, Pos: Position{Row: 24, Col: 3}, FullName: "นครศรีธรรมราช"},
	{ShortName: "ภก", Color: S, Pos: Position{Row: 25, Col: 1}, FullName: "ภูเก็ต"},
	{ShortName: "ตง", Color: S, Pos: Position{Row: 25, Col: 2}, FullName: "ตรัง"},
	{ShortName: "พท", Color: S, Pos: Position{Row: 25, Col: 3}, FullName: "พัทลุง"},
	{ShortName: "สง", Color: S, Pos: Position{Row: 25, Col: 4}, FullName: "สงขลา"},
	{ShortName: "ตง", Color: S, Pos: Position{Row: 26, Col: 2}, FullName: "ตรัง"},
	{ShortName: "พท", Color: S, Pos: Position{Row: 26, Col: 3}, FullName: "พัทลุง"},
	{ShortName: "สง", Color: S, Pos: Position{Row: 26, Col: 4}, FullName: "สงขลา"},
	{ShortName: "สต", Color: S, Pos: Position{Row: 27, Col: 2}, FullName: "สตูล"},
	{ShortName: "สง", Color: S, Pos: Position{Row: 27, Col: 3}, FullName: "สงขลา"},
	{ShortName: "ปน", Color: S, Pos: Position{Row: 27, Col: 4}, FullName: "ปัตตานี"},
	{ShortName: "ยล", Color: S, Pos: Position{Row: 28, Col: 3}, FullName: "ยะลา"},
	{ShortName: "ปน", Color: S, Pos: Position{Row: 28, Col: 4}, FullName: "ปัตตานี"},
	{ShortName: "ยล", Color: S, Pos: Position{Row: 29, Col: 3}, FullName: "ยะลา"},
	{ShortName: "นธ", Color: S, Pos: Position{Row: 29, Col: 4}, FullName: "นราธิวาส"},
	{ShortName: "นธ", Color: S, Pos: Position{Row: 30, Col: 4}, FullName: "นราธิวาส"},
}

var (
	provinceByFullname = map[string]*province{}
	provinceByCoord    = map[int]*province{}
	gridCache          [][10]string
	initGridOnce       sync.Once
)

func initCaches() {
	maxRow := 0
	for i := range provinces {
		p := &provinces[i]
		if p.Pos.Row > maxRow {
			maxRow = p.Pos.Row
		}
		if p.FullName != "" {
			provinceByFullname[p.FullName] = p
		}
		key := p.Pos.Row*MAX_COLS + p.Pos.Col
		provinceByCoord[key] = p
	}
	rows := maxRow + 1
	grid := make([][10]string, rows)
	for _, p := range provinces {
		r := p.Pos.Row
		c := p.Pos.Col
		if r < 0 || r >= rows || c < 0 || c >= MAX_COLS {
			continue
		}
		grid[r][c] = p.Color + p.ShortName
	}
	gridCache = grid
}

func GetProvinceByFullname(fullname string) *province {
	initGridOnce.Do(initCaches)
	return provinceByFullname[fullname]
}

func GetProvinceAt(row, col int) *province {
	initGridOnce.Do(initCaches)
	key := row*MAX_COLS + col
	return provinceByCoord[key]
}

func NewMap() *MapRegion {
	initGridOnce.Do(initCaches)
	return &MapRegion{
		Grid: gridCache,
		Size: Position{
			Row: len(gridCache),
			Col: MAX_COLS * SPACE_LEN,
		},
	}
}

func (m *MapRegion) Debug() {
	for _, r := range m.Grid {
		for _, p := range r {
			if p == "" {
				fmt.Print(strings.Repeat(" ", SPACE_LEN))
			} else {
				s := fmt.Sprintf("[%s%s]", p, X)
				fmt.Print(s)
			}
		}
		fmt.Println()
	}
	os.Exit(0)
}
