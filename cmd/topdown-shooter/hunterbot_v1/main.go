package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/BotBattleArena/ArenaFramework/pkg/arena"
)

type BotPlayer struct {
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	HP      int     `json:"hp"`
	Kills   int     `json:"kills"`
	Deaths  int     `json:"deaths"`
	ShootCD int     `json:"shoot_cd"`
	DashCD  int     `json:"dash_cd"`
	AimX    float64 `json:"aim_x"`
	AimY    float64 `json:"aim_y"`
}

type MapObject struct {
	X      float64      `json:"x"`
	Y      float64      `json:"y"`
	W      float64      `json:"w"`
	H      float64      `json:"h"`
	R      float64      `json:"r"`
	DX     float64      `json:"dx"`
	DY     float64      `json:"dy"`
	Points [][2]float64 `json:"points"`
}

type StaticScene struct {
	Rect   []MapObject `json:"rect"`
	Circle []MapObject `json:"circle"`
	Poly   []MapObject `json:"poly"`
}

type DynamicScene struct {
	Rect    []MapObject `json:"rect"`
	Circle  []MapObject `json:"circle"`
	Poly    []MapObject `json:"poly"`
	Bullets []MapObject `json:"bullets"`
}

type BotFrame struct {
	Type    string                `json:"type"`
	Tick    int                   `json:"tick"`
	Players map[string]*BotPlayer `json:"players"`
	Scene   struct {
		Dynamic DynamicScene `json:"dynamic"`
	} `json:"scene"`
}

func getMyID() string {
	name := filepath.Base(os.Args[0])
	name = strings.TrimSuffix(name, filepath.Ext(name))
	return name
}

func main() {
	log.SetOutput(os.Stderr)
	myID := getMyID()
	log.Printf("hunterbot: started as %s", myID)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*256), 1024*256)

	var staticMap StaticScene
	mapW, mapH := 2000.0, 2000.0
	_ = staticMap

	for scanner.Scan() {
		var msg arena.ServerMessage
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			log.Printf("hunterbot: parse error: %v", err)
			continue
		}

		switch msg.Type {
		case "start":
			log.Println("hunterbot: game started")

		case "state":
			var base map[string]interface{}
			if err := json.Unmarshal(msg.State, &base); err != nil {
				continue
			}

			frameType, _ := base["type"].(string)

			if frameType == "setup" {
				log.Println("hunterbot: received setup")
				var setup struct {
					Scene struct {
						MapW   float64     `json:"mapw"`
						MapH   float64     `json:"maph"`
						Static StaticScene `json:"static"`
					} `json:"scene"`
				}
				json.Unmarshal(msg.State, &setup)
				staticMap = setup.Scene.Static
				if setup.Scene.MapW > 0 {
					mapW = setup.Scene.MapW
				}
				if setup.Scene.MapH > 0 {
					mapH = setup.Scene.MapH
				}
				continue
			}

			if frameType != "tick" {
				continue
			}

			var frame BotFrame
			if err := json.Unmarshal(msg.State, &frame); err != nil {
				log.Printf("hunterbot: state error: %v", err)
				continue
			}

			var me *BotPlayer
			if player, exists := frame.Players[myID]; exists {
				me = player
			}

			if me == nil || me.HP <= 0 {
				resp := arena.InputMessage{Axes: map[string]float32{
					"move_x": 0, "move_y": 0,
					"aim_x": 0, "aim_y": 1,
					"shoot": 0, "dash": 0,
				}}
				out, _ := json.Marshal(resp)
				fmt.Println(string(out))
				continue
			}

			// AGGRESSIVE HUNTING LOGIC
			closest := math.MaxFloat64
			var target *BotPlayer
			for id, p := range frame.Players {
				if id == myID || p.HP <= 0 {
					continue
				}
				dx, dy := p.X-me.X, p.Y-me.Y
				d := math.Sqrt(dx*dx + dy*dy)
				if d < closest {
					closest = d
					target = p
				}
			}

			var mx, my, ax, ay float32
			shoot := float32(0)
			dash := float32(0)

			if target != nil {
				dx, dy := target.X-me.X, target.Y-me.Y
				dist := math.Sqrt(dx*dx + dy*dy)

				// Move towards target
				if dist > 0.01 {
					mx = float32(dx / dist)
					my = float32(dy / dist)
					ax, ay = mx, my
				}

				// Shoot if in range
				if dist < 600 && me.ShootCD == 0 {
					shoot = 1
				}

				// Dash to close the gap if far away
				if dist > 300 && me.DashCD == 0 && rand.Float64() < 0.1 {
					dash = 1
				}
			} else {
				// No target? Move towards center
				cx, cy := mapW/2-me.X, mapH/2-me.Y
				d := math.Sqrt(cx*cx + cy*cy)
				if d > 10 {
					mx, my = float32(cx/d), float32(cy/d)
				}
				ax, ay = 0, 1
			}

			// EVASION LOGIC (Dodge incoming bullets)
			for _, b := range frame.Scene.Dynamic.Bullets {
				dx, dy := b.X-me.X, b.Y-me.Y
				d := math.Sqrt(dx*dx + dy*dy)
				if d < 150 {
					dot := dx*b.DX + dy*b.DY
					// Bullet moving towards us?
					if dot < 0 {
						evadeX, evadeY := -b.DY, b.DX
						if rand.Float64() < 0.5 {
							evadeX, evadeY = -evadeX, -evadeY
						}
						// Strong override to dodge
						mx, my = float32(evadeX), float32(evadeY)
						if d < 70 && me.DashCD == 0 {
							dash = 1
						}
						break
					}
				}
			}

			// Normalize movement vector if > 1
			vlen := math.Sqrt(float64(mx*mx + my*my))
			if vlen > 1.0 {
				mx = float32(float64(mx) / vlen)
				my = float32(float64(my) / vlen)
			}

			resp := arena.InputMessage{Axes: map[string]float32{
				"move_x": mx, "move_y": my,
				"aim_x": ax, "aim_y": ay,
				"shoot": shoot, "dash": dash,
			}}
			out, _ := json.Marshal(resp)
			fmt.Println(string(out))

		case "end":
			log.Println("hunterbot: game over")
			return
		}
	}
}
