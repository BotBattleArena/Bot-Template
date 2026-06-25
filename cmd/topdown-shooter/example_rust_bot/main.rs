use rand::Rng;
use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::collections::HashMap;
use std::io::{self, BufRead, Write};

// ============================================================================
// 1. DATA STRUCTURES (Server -> Bot)
// ============================================================================

/// ServerMessage is the generic envelope for all communication
#[derive(Deserialize, Default)]
struct ServerMessage {
    #[serde(default)]
    r#type: String,
    #[serde(default)]
    state: Value,
}

/// StateHeader is used to route "state" messages (setup vs tick)
#[derive(Deserialize, Default)]
struct StateHeader {
    #[serde(default)]
    r#type: String,
}

/// Bot represents a player in the arena
#[derive(Deserialize, Default)]
struct Bot {
    #[serde(default)]
    x: f64,
    #[serde(default)]
    y: f64,
    #[serde(default)]
    hp: i32,
    #[serde(default)]
    kills: i32,
    #[serde(default)]
    deaths: i32,
    #[serde(default)]
    shoot_cd: i32,
    #[serde(default)]
    dash_cd: i32,
    #[serde(default)]
    aim_x: f64,
    #[serde(default)]
    aim_y: f64,
}

/// SceneObject represents walls, bullets, and other geometry
#[derive(Deserialize, Default)]
struct SceneObject {
    #[serde(default)]
    x: f64,
    #[serde(default)]
    y: f64,
    #[serde(default)]
    w: f64,
    #[serde(default)]
    h: f64,
    #[serde(default)]
    r: f64,
    #[serde(default)]
    dx: f64,
    #[serde(default)]
    dy: f64,
    #[serde(default)]
    points: Vec<[f64; 2]>,
}

#[derive(Deserialize, Default)]
struct SceneObjects {
    #[serde(default)]
    rect: Vec<SceneObject>,
    #[serde(default)]
    circle: Vec<SceneObject>,
    #[serde(default)]
    poly: Vec<SceneObject>,
    #[serde(default)]
    bullets: Vec<SceneObject>,
}

#[derive(Deserialize, Default)]
struct Scene {
    #[serde(default)]
    mapw: f64,
    #[serde(default)]
    maph: f64,
    #[serde(default, rename = "static")]
    static_objs: SceneObjects,
    #[serde(default)]
    dynamic: SceneObjects,
}

#[derive(Deserialize, Default)]
struct Rules {
    #[serde(default)]
    move_speed: f64,
    #[serde(default)]
    bullet_speed: f64,
    #[serde(default)]
    damage: i32,
    #[serde(default)]
    win_kills: i32,
    #[serde(default)]
    respawn_ticks: i32,
    #[serde(default)]
    shoot_cd: i32,
    #[serde(default)]
    dash_cd: i32,
    #[serde(default)]
    max_hp: i32,
}

/// SetupFrame contains initial game rules and map
#[allow(dead_code)]
#[derive(Deserialize, Default)]
struct SetupFrame {
    #[serde(default)]
    r#type: String,
    #[serde(default)]
    tick_duration_ms: i32,
    #[serde(default)]
    timeout_ms: i32,
    #[serde(default)]
    max_tick: i32,
    #[serde(default)]
    countdown_sec: i32,
    #[serde(default)]
    scene: Scene,
    #[serde(default)]
    rules: Rules,
}

/// TickFrame is the game state for a single frame
#[derive(Deserialize, Default)]
struct TickFrame {
    #[serde(default)]
    r#type: String,
    #[serde(default)]
    tick: i32,
    #[serde(default)]
    players: HashMap<String, Bot>,
    #[serde(default)]
    scene: Scene,
}

// ============================================================================
// 2. OUTPUT STRUCTURES (Bot -> Server)
// ============================================================================

#[derive(Serialize, Default)]
struct InputMsg {
    move_x: f32,
    move_y: f32,
    aim_x: f32,
    aim_y: f32,
    shoot: f32,
    dash: f32,
}

#[derive(Serialize)]
struct Response<'a> {
    axes: &'a InputMsg,
}

// ============================================================================
// 3. GAME ENGINE
// ============================================================================

#[allow(dead_code)]
struct Game {
    setup: SetupFrame,
    tick: TickFrame,
    my_id: String,
}

impl Game {
    fn new() -> Self {
        let exe = std::env::current_exe().unwrap_or_default();
        let name = exe
            .file_stem()
            .unwrap_or_default()
            .to_string_lossy()
            .to_string();
        Game {
            setup: SetupFrame::default(),
            tick: TickFrame::default(),
            my_id: name,
        }
    }

    /// Send wraps input and encodes it to stdout
    fn send(&self, msg: &InputMsg) {
        let resp = Response { axes: msg };
        let out = serde_json::to_string(&resp).unwrap();
        let stdout = io::stdout();
        let mut handle = stdout.lock();
        let _ = writeln!(handle, "{}", out);
        let _ = handle.flush();
    }

    /// Process determines the next bot action
    fn process(&self) {
        let mut input = InputMsg::default();
        let mut rng = rand::thread_rng();

        let me = match self.tick.players.get(&self.my_id) {
            Some(bot) if bot.hp > 0 => bot,
            _ => {
                self.send(&input);
                return;
            }
        };

        // Logic: Random movement and aiming
        input.move_x = rng.gen_range(-1.0..=1.0);
        input.move_y = rng.gen_range(-1.0..=1.0);
        input.aim_x = rng.gen_range(-1.0..=1.0);
        input.aim_y = rng.gen_range(-1.0..=1.0);

        if me.shoot_cd == 0 {
            input.shoot = rng.gen_range(0..=1) as f32;
        }
        if me.dash_cd == 0 {
            input.dash = rng.gen_range(0..=1) as f32;
        }

        self.send(&input);
    }
}

// ============================================================================
// 4. MAIN ENTRY POINT
// ============================================================================

fn main() {
    let mut game = Game::new();

    // Log to stderr; stdout is reserved for the protocol
    eprintln!("example_rust_bot: started as {}", game.my_id);

    let stdin = io::stdin();
    for line in stdin.lock().lines() {
        let line = match line {
            Ok(l) => l,
            Err(_) => break,
        };

        let msg: ServerMessage = match serde_json::from_str(&line) {
            Ok(m) => m,
            Err(_) => continue,
        };

        match msg.r#type.as_str() {
            "start" => continue,
            "end" => return,
            "state" => {
                let header: StateHeader = match serde_json::from_value(msg.state.clone()) {
                    Ok(h) => h,
                    Err(_) => continue,
                };

                match header.r#type.as_str() {
                    "setup" => {
                        if let Ok(setup) = serde_json::from_value(msg.state) {
                            game.setup = setup;
                            eprintln!("example_rust_bot: setup received");
                        }
                    }
                    "tick" => {
                        if let Ok(tick) = serde_json::from_value(msg.state) {
                            game.tick = tick;
                            game.process();
                        }
                    }
                    _ => {}
                }
            }
            _ => {}
        }
    }
}
