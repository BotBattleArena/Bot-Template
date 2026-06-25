import sys
import json
import random
import os
from dataclasses import dataclass, field
from typing import Dict, List, Optional

# ============================================================================
# 1. DATA STRUCTURES (Server -> Bot)
# ============================================================================


@dataclass
class Bot:
    """Represents a player in the arena."""
    x: float = 0.0
    y: float = 0.0
    hp: int = 0
    kills: int = 0
    deaths: int = 0
    shoot_cd: int = 0
    dash_cd: int = 0
    aim_x: float = 0.0
    aim_y: float = 0.0

    @classmethod
    def from_dict(cls, d: dict) -> "Bot":
        return cls(
            x=d.get("x", 0.0),
            y=d.get("y", 0.0),
            hp=d.get("hp", 0),
            kills=d.get("kills", 0),
            deaths=d.get("deaths", 0),
            shoot_cd=d.get("shoot_cd", 0),
            dash_cd=d.get("dash_cd", 0),
            aim_x=d.get("aim_x", 0.0),
            aim_y=d.get("aim_y", 0.0),
        )


@dataclass
class SceneObject:
    """Represents walls, bullets, and other geometry."""
    x: float = 0.0
    y: float = 0.0
    w: float = 0.0
    h: float = 0.0
    r: float = 0.0
    dx: float = 0.0
    dy: float = 0.0
    points: List[List[float]] = field(default_factory=list)

    @classmethod
    def from_dict(cls, d: dict) -> "SceneObject":
        return cls(
            x=d.get("x", 0.0),
            y=d.get("y", 0.0),
            w=d.get("w", 0.0),
            h=d.get("h", 0.0),
            r=d.get("r", 0.0),
            dx=d.get("dx", 0.0),
            dy=d.get("dy", 0.0),
            points=d.get("points", []),
        )


@dataclass
class SceneObjects:
    """Container for categorised scene geometry."""
    rect: List[SceneObject] = field(default_factory=list)
    circle: List[SceneObject] = field(default_factory=list)
    poly: List[SceneObject] = field(default_factory=list)
    bullets: List[SceneObject] = field(default_factory=list)

    @classmethod
    def from_dict(cls, d: dict) -> "SceneObjects":
        return cls(
            rect=[SceneObject.from_dict(o) for o in d.get("rect", []) or []],
            circle=[SceneObject.from_dict(o) for o in d.get("circle", []) or []],
            poly=[SceneObject.from_dict(o) for o in d.get("poly", []) or []],
            bullets=[SceneObject.from_dict(o) for o in d.get("bullets", []) or []],
        )


@dataclass
class Scene:
    """Describes the arena map and its objects."""
    mapw: float = 0.0
    maph: float = 0.0
    static_objs: SceneObjects = field(default_factory=SceneObjects)
    dynamic: SceneObjects = field(default_factory=SceneObjects)

    @classmethod
    def from_dict(cls, d: dict) -> "Scene":
        return cls(
            mapw=d.get("mapw", 0.0),
            maph=d.get("maph", 0.0),
            static_objs=SceneObjects.from_dict(d.get("static", {})),
            dynamic=SceneObjects.from_dict(d.get("dynamic", {})),
        )


@dataclass
class Rules:
    """Game balance parameters."""
    move_speed: float = 0.0
    bullet_speed: float = 0.0
    damage: int = 0
    win_kills: int = 0
    respawn_ticks: int = 0
    shoot_cd: int = 0
    dash_cd: int = 0
    max_hp: int = 0

    @classmethod
    def from_dict(cls, d: dict) -> "Rules":
        return cls(
            move_speed=d.get("move_speed", 0.0),
            bullet_speed=d.get("bullet_speed", 0.0),
            damage=d.get("damage", 0),
            win_kills=d.get("win_kills", 0),
            respawn_ticks=d.get("respawn_ticks", 0),
            shoot_cd=d.get("shoot_cd", 0),
            dash_cd=d.get("dash_cd", 0),
            max_hp=d.get("max_hp", 0),
        )


@dataclass
class SetupFrame:
    """Initial game configuration sent once at start."""
    tick_duration_ms: int = 0
    timeout_ms: int = 0
    max_tick: int = 0
    countdown_sec: int = 0
    scene: Scene = field(default_factory=Scene)
    rules: Rules = field(default_factory=Rules)

    @classmethod
    def from_dict(cls, d: dict) -> "SetupFrame":
        return cls(
            tick_duration_ms=d.get("tick_duration_ms", 0),
            timeout_ms=d.get("timeout_ms", 0),
            max_tick=d.get("max_tick", 0),
            countdown_sec=d.get("countdown_sec", 0),
            scene=Scene.from_dict(d.get("scene", {})),
            rules=Rules.from_dict(d.get("rules", {})),
        )


@dataclass
class TickFrame:
    """Game state for a single frame."""
    tick: int = 0
    players: Dict[str, Bot] = field(default_factory=dict)
    scene: Scene = field(default_factory=Scene)

    @classmethod
    def from_dict(cls, d: dict) -> "TickFrame":
        return cls(
            tick=d.get("tick", 0),
            players={
                pid: Bot.from_dict(pdata)
                for pid, pdata in (d.get("players") or {}).items()
            },
            scene=Scene.from_dict(d.get("scene", {})),
        )


# ============================================================================
# 2. OUTPUT STRUCTURES (Bot -> Server)
# ============================================================================


@dataclass
class InputMsg:
    """Controls sent back to the server each tick."""
    move_x: float = 0.0
    move_y: float = 0.0
    aim_x: float = 0.0
    aim_y: float = 0.0
    shoot: float = 0.0
    dash: float = 0.0

    def to_dict(self) -> dict:
        return {
            "move_x": self.move_x,
            "move_y": self.move_y,
            "aim_x": self.aim_x,
            "aim_y": self.aim_y,
            "shoot": self.shoot,
            "dash": self.dash,
        }


# ============================================================================
# 3. GAME ENGINE
# ============================================================================


class Game:
    """Handles the stdin/stdout JSON protocol."""

    def __init__(self):
        # Derive MyID from the executable filename, same as the Go version
        name = os.path.basename(sys.argv[0])
        self.my_id: str = os.path.splitext(name)[0]
        self.setup: Optional[SetupFrame] = None
        self.tick: Optional[TickFrame] = None

    def listen(self) -> Optional[dict]:
        """Reads a single JSON message from stdin."""
        line = sys.stdin.readline()
        if not line:
            return None
        try:
            return json.loads(line)
        except json.JSONDecodeError:
            return None

    def send(self, msg: InputMsg) -> None:
        """Encodes and sends an input message to stdout."""
        resp = {"axes": msg.to_dict()}
        sys.stdout.write(json.dumps(resp) + "\n")
        sys.stdout.flush()

    def process(self) -> None:
        """Determines the next bot action based on current tick state."""
        inp = InputMsg()

        if self.tick is None:
            self.send(inp)
            return

        me = self.tick.players.get(self.my_id)

        # If dead or not in the game, send empty input
        if me is None or me.hp <= 0:
            self.send(inp)
            return

        # Logic: Random movement and aiming
        inp.move_x = random.uniform(-1, 1)
        inp.move_y = random.uniform(-1, 1)
        inp.aim_x = random.uniform(-1, 1)
        inp.aim_y = random.uniform(-1, 1)

        # Shoot and Dash randomly if off cooldown
        if me.shoot_cd == 0:
            inp.shoot = float(random.randint(0, 1))
        if me.dash_cd == 0:
            inp.dash = float(random.randint(0, 1))

        self.send(inp)


# ============================================================================
# 4. MAIN ENTRY POINT
# ============================================================================


def main():
    game = Game()

    # Log to stderr; stdout is reserved for the protocol
    print(f"example_python_bot: started as {game.my_id}", file=sys.stderr)

    while True:
        msg = game.listen()
        if msg is None:
            break

        msg_type = msg.get("type")
        if msg_type == "start":
            continue
        elif msg_type == "end":
            break
        elif msg_type == "state":
            state = msg.get("state")
            if not state:
                continue

            state_type = state.get("type")
            if state_type == "setup":
                game.setup = SetupFrame.from_dict(state)
                print("example_python_bot: setup received", file=sys.stderr)
            elif state_type == "tick":
                game.tick = TickFrame.from_dict(state)
                game.process()


if __name__ == "__main__":
    main()
