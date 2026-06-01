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
	log.Printf("shooterbot: started as %s", myID)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*256), 1024*256)

	moveAngle := rand.Float64() * math.Pi * 2
	changeTick := 0

	for scanner.Scan() {
		var msg arena.ServerMessage
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			log.Printf("shooterbot: parse error: %v", err)
			continue
		}

		switch msg.Type {
		case "start":
			for _, ax := range msg.Axes {
				log.Printf("shooterbot: axis %s (%.2f)", ax.Name, ax.Value)
			}

		case "state":
			var base map[string]interface{}
			if err := json.Unmarshal(msg.State, &base); err != nil {
				continue
			}

			frameType, _ := base["type"].(string)

			if frameType == "setup" {
				log.Println("shooterbot: received setup")
				var setup struct {
					Scene struct {
						MapW   float64     `json:"mapw"`
						MapH   float64     `json:"maph"`
						Static StaticScene `json:"static"`
					} `json:"scene"`
				}
				json.Unmarshal(msg.State, &setup)
				continue
			}

			if frameType != "tick" {
				continue
			}

			var frame BotFrame
			if err := json.Unmarshal(msg.State, &frame); err != nil {
				log.Printf("shooterbot: state error: %v", err)
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

			changeTick++
			if changeTick > 60+rand.Intn(120) {
				moveAngle = rand.Float64() * math.Pi * 2
				changeTick = 0
			}

			mx := float32(math.Cos(moveAngle))
			my := float32(math.Sin(moveAngle))

			var aimX, aimY float32
			shoot := float32(0)
			dash := float32(0)

			closest := math.MaxFloat64
			var targetX, targetY float64
			hasTarget := false

			for id, p := range frame.Players {
				if id == myID || p.HP <= 0 {
					continue
				}
				dx, dy := p.X-me.X, p.Y-me.Y
				d := math.Sqrt(dx*dx + dy*dy)
				if d < closest {
					closest = d
					targetX, targetY = p.X, p.Y
					hasTarget = true
				}
			}

			if hasTarget {
				dx, dy := targetX-me.X, targetY-me.Y
				d := math.Sqrt(dx*dx + dy*dy)
				if d > 0.1 {
					aimX = float32(dx / d)
					aimY = float32(dy / d)
				}
				if d < 500 && me.ShootCD == 0 {
					shoot = 1
				}
				if d < 150 && me.DashCD == 0 && rand.Float64() < 0.3 {
					perpAngle := math.Atan2(float64(aimY), float64(aimX)) + math.Pi/2
					if rand.Float64() < 0.5 {
						perpAngle = -perpAngle
					}
					aimX = float32(math.Cos(perpAngle))
					aimY = float32(math.Sin(perpAngle))
					dash = 1
				}
			} else {
				aimX = mx
				aimY = my
				if rand.Float64() < 0.05 {
					shoot = 1
				}
			}

			// EVASION LOGIC (Dodge incoming bullets)
			for _, b := range frame.Scene.Dynamic.Bullets {
				dx, dy := b.X-me.X, b.Y-me.Y
				d := math.Sqrt(dx*dx + dy*dy)
				if d < 120 {
					dot := dx*b.DX + dy*b.DY
					// moving towards us
					if dot < 0 {
						evadeX := -b.DY
						evadeY := b.DX
						if rand.Float64() < 0.5 {
							evadeX, evadeY = -evadeX, -evadeY
						}
						mx = float32(evadeX)
						my = float32(evadeY)
						if d < 60 && me.DashCD == 0 {
							aimX = float32(evadeX)
							aimY = float32(evadeY)
							dash = 1
						}
						break
					}
				}
			}

			// Normalize
			vlen := math.Sqrt(float64(mx*mx + my*my))
			if vlen > 1.0 {
				mx = float32(float64(mx) / vlen)
				my = float32(float64(my) / vlen)
			}

			resp := arena.InputMessage{Axes: map[string]float32{
				"move_x": mx, "move_y": my,
				"aim_x": aimX, "aim_y": aimY,
				"shoot": shoot, "dash": dash,
			}}
			out, _ := json.Marshal(resp)
			fmt.Println(string(out))

		case "end":
			log.Println("shooterbot: game over")
			return
		}
	}
}
