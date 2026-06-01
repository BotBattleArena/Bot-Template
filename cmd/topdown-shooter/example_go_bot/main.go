package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

// ============================================================================
// 1. DATA STRUCTURES (Server -> Bot)
// ============================================================================

// ServerMessage is the generic envelope for all communication
type ServerMessage struct {
	Type  string          `json:"type"`
	State json.RawMessage `json:"state"`
}

// StateHeader is used to route "state" messages (setup vs tick)
type StateHeader struct {
	Type string `json:"type"`
}

// Bot represents a player in the arena
type Bot struct {
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

// SceneObject represents walls, bullets, and other geometry
type SceneObject struct {
	X      float64      `json:"x"`
	Y      float64      `json:"y"`
	W      float64      `json:"w"`
	H      float64      `json:"h"`
	R      float64      `json:"r"`
	DX     float64      `json:"dx"`
	DY     float64      `json:"dy"`
	Points [][2]float64 `json:"points"`
}

type SceneObjects struct {
	Rect    []SceneObject `json:"rect"`
	Circle  []SceneObject `json:"circle"`
	Poly    []SceneObject `json:"poly"`
	Bullets []SceneObject `json:"bullets"`
}

type Scene struct {
	MapW    float64      `json:"mapw"`
	MapH    float64      `json:"maph"`
	Static  SceneObjects `json:"static"`
	Dynamic SceneObjects `json:"dynamic"`
}

type Rules struct {
	MoveSpeed    float64 `json:"move_speed"`
	BulletSpeed  float64 `json:"bullet_speed"`
	Damage       int     `json:"damage"`
	WinKills     int     `json:"win_kills"`
	RespawnTicks int     `json:"respawn_ticks"`
	ShootCD      int     `json:"shoot_cd"`
	DashCD       int     `json:"dash_cd"`
	MaxHP        int     `json:"max_hp"`
}

// TickFrame is the game state for a single frame
type TickFrame struct {
	Type    string          `json:"type"`
	Tick    int             `json:"tick"`
	Players map[string]*Bot `json:"players"`
	Scene   Scene           `json:"scene"`
}

// SetupFrame contains initial game rules and map
type SetupFrame struct {
	Type           string `json:"type"`
	TickDurationMs int    `json:"tick_duration_ms"`
	TimeoutMs      int    `json:"timeout_ms"`
	MaxTick        int    `json:"max_tick"`
	CountdownSec   int    `json:"countdown_sec"`
	Scene          Scene  `json:"scene"`
	Rules          Rules  `json:"rules"`
}

// ============================================================================
// 2. GAME ENGINE
// ============================================================================

type Game struct {
	Setup SetupFrame
	Tick  TickFrame
	MyID  string

	msg ServerMessage
	dec *json.Decoder
	enc *json.Encoder
}

func NewGame() *Game {
	name := filepath.Base(os.Args[0])
	return &Game{
		MyID: strings.TrimSuffix(name, filepath.Ext(name)),
		dec:  json.NewDecoder(os.Stdin),
		enc:  json.NewEncoder(os.Stdout),
	}
}

// Listen waits for the next message. Returns false on error/EOF.
func (game *Game) Listen() bool {
	return game.dec.Decode(&game.msg) == nil
}

// Send wraps input and encodes it to stdout
func (game *Game) Send(msg *InputMsg) {
	resp := struct {
		Axes *InputMsg `json:"axes"`
	}{Axes: msg}
	game.enc.Encode(resp)
}

// Process determines next bot action
func (game *Game) Process() {
	var input InputMsg

	me, exists := game.Tick.Players[game.MyID]
	if !exists || me.HP <= 0 {
		game.Send(&input)
		return
	}

	// Logic: Random movement and aiming
	input.MoveX = rand.Float32()*2 - 1
	input.MoveY = rand.Float32()*2 - 1
	input.AimX = rand.Float32()*2 - 1
	input.AimY = rand.Float32()*2 - 1

	if me.ShootCD == 0 {
		input.Shoot = float32(rand.Intn(2))
	}
	if me.DashCD == 0 {
		input.Dash = float32(rand.Intn(2))
	}

	game.Send(&input)
}

// ============================================================================
// 3. OUTPUT STRUCTURES (Bot -> Server)
// ============================================================================

type InputMsg struct {
	MoveX float32 `json:"move_x"`
	MoveY float32 `json:"move_y"`
	AimX  float32 `json:"aim_x"`
	AimY  float32 `json:"aim_y"`
	Shoot float32 `json:"shoot"`
	Dash  float32 `json:"dash"`
}

// ============================================================================
// 4. MAIN ENTRY POINT
// ============================================================================

func main() {
	game := NewGame()

	// Log to stderr; stdout is reserved for the protocol
	log.SetOutput(os.Stderr)
	log.Printf("example_bot: started as %s", game.MyID)

	for game.Listen() {
		switch game.msg.Type {
		case "start":
			continue
		case "end":
			return
		case "state":
			var header StateHeader
			json.Unmarshal(game.msg.State, &header)

			switch header.Type {
			case "setup":
				json.Unmarshal(game.msg.State, &game.Setup)
				log.Println("example_bot: setup received")
			case "tick":
				if err := json.Unmarshal(game.msg.State, &game.Tick); err == nil {
					game.Process()
				}
			}
		}
	}
}
