const readline = require('readline');
const path = require('path');

// ============================================================================
// 1. DATA STRUCTURES (Server -> Bot)
// ============================================================================

class Bot {
    constructor(data = {}) {
        this.x = data.x || 0;
        this.y = data.y || 0;
        this.hp = data.hp || 0;
        this.kills = data.kills || 0;
        this.deaths = data.deaths || 0;
        this.shootCD = data.shoot_cd || 0;
        this.dashCD = data.dash_cd || 0;
        this.aimX = data.aim_x || 0;
        this.aimY = data.aim_y || 0;
    }
}

class SceneObject {
    constructor(data = {}) {
        this.x = data.x || 0;
        this.y = data.y || 0;
        this.w = data.w || 0;
        this.h = data.h || 0;
        this.r = data.r || 0;
        this.dx = data.dx || 0;
        this.dy = data.dy || 0;
        this.points = data.points || [];
    }
}

class SceneObjects {
    constructor(data = {}) {
        this.rect = (data.rect || []).map(o => new SceneObject(o));
        this.circle = (data.circle || []).map(o => new SceneObject(o));
        this.poly = (data.poly || []).map(o => new SceneObject(o));
        this.bullets = (data.bullets || []).map(o => new SceneObject(o));
    }
}

class Scene {
    constructor(data = {}) {
        this.mapW = data.mapw || 0;
        this.mapH = data.maph || 0;
        // "static" is a reserved word, so access via bracket notation
        this.staticObjects = new SceneObjects(data["static"] || {});
        this.dynamic = new SceneObjects(data.dynamic || {});
    }
}

class Rules {
    constructor(data = {}) {
        this.moveSpeed = data.move_speed || 0;
        this.bulletSpeed = data.bullet_speed || 0;
        this.damage = data.damage || 0;
        this.winKills = data.win_kills || 0;
        this.respawnTicks = data.respawn_ticks || 0;
        this.shootCD = data.shoot_cd || 0;
        this.dashCD = data.dash_cd || 0;
        this.maxHP = data.max_hp || 0;
    }
}

class SetupFrame {
    constructor(data = {}) {
        this.type = data.type || '';
        this.tickDurationMs = data.tick_duration_ms || 0;
        this.timeoutMs = data.timeout_ms || 0;
        this.maxTick = data.max_tick || 0;
        this.countdownSec = data.countdown_sec || 0;
        this.scene = new Scene(data.scene || {});
        this.rules = new Rules(data.rules || {});
    }
}

class TickFrame {
    constructor(data = {}) {
        this.type = data.type || '';
        this.tick = data.tick || 0;
        this.players = {};
        this.scene = new Scene(data.scene || {});

        // Parse each player into a Bot instance
        const players = data.players || {};
        for (const [id, botData] of Object.entries(players)) {
            this.players[id] = new Bot(botData);
        }
    }
}

// ============================================================================
// 2. OUTPUT STRUCTURES (Bot -> Server)
// ============================================================================

class InputMsg {
    constructor() {
        this.move_x = 0;
        this.move_y = 0;
        this.aim_x = 0;
        this.aim_y = 0;
        this.shoot = 0;
        this.dash = 0;
    }
}

// ============================================================================
// 3. GAME ENGINE
// ============================================================================

class Game {
    constructor() {
        // Derive bot ID from executable/script name (without extension)
        const scriptName = path.basename(process.argv[1] || process.argv[0]);
        this.myID = scriptName.replace(path.extname(scriptName), '');
        this.setup = null;
        this.tick = null;
    }

    /** Send wraps input and encodes it to stdout */
    send(input) {
        const resp = { axes: input };
        process.stdout.write(JSON.stringify(resp) + '\n');
    }

    /** Process determines next bot action */
    process() {
        const input = new InputMsg();

        const me = this.tick ? this.tick.players[this.myID] : null;
        if (!me || me.hp <= 0) {
            this.send(input);
            return;
        }

        // Logic: Random movement and aiming
        input.move_x = Math.random() * 2 - 1;
        input.move_y = Math.random() * 2 - 1;
        input.aim_x = Math.random() * 2 - 1;
        input.aim_y = Math.random() * 2 - 1;

        if (me.shootCD === 0) {
            input.shoot = Math.random() < 0.5 ? 0 : 1;
        }
        if (me.dashCD === 0) {
            input.dash = Math.random() < 0.5 ? 0 : 1;
        }

        this.send(input);
    }
}

// ============================================================================
// 4. MAIN ENTRY POINT
// ============================================================================

function main() {
    const game = new Game();

    // Log to stderr; stdout is reserved for the protocol
    process.stderr.write(`example_js_bot: started as ${game.myID}\n`);

    const rl = readline.createInterface({
        input: process.stdin,
        terminal: false,
    });

    rl.on('line', (line) => {
        let msg;
        try {
            msg = JSON.parse(line);
        } catch {
            return;
        }

        switch (msg.type) {
            case 'start':
                return;
            case 'end':
                process.exit(0);
                return;
            case 'state': {
                const state = msg.state;
                if (!state) return;

                switch (state.type) {
                    case 'setup':
                        game.setup = new SetupFrame(state);
                        process.stderr.write('example_js_bot: setup received\n');
                        break;
                    case 'tick':
                        game.tick = new TickFrame(state);
                        game.process();
                        break;
                }
                break;
            }
        }
    });

    rl.on('close', () => {
        process.exit(0);
    });
}

main();
