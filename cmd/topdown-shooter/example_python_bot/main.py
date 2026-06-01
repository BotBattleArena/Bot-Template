import sys
import json
import random
import os

# ============================================================================
# 1. GAME ENGINE
# ============================================================================

class Game:
    def __init__(self):
        # Derive MyID from the filename, similar to the Go version
        name = os.path.basename(sys.argv[0])
        self.my_id = os.path.splitext(name)[0]
        self.setup = {}
        self.tick = {}

    def listen(self):
        """Reads a single JSON message from stdin."""
        line = sys.stdin.readline()
        if not line:
            return None
        try:
            return json.loads(line)
        except json.JSONDecodeError:
            return None

    def send(self, input_msg):
        """Encodes and sends an input message to stdout."""
        resp = {"axes": input_msg}
        sys.stdout.write(json.dumps(resp) + "\n")
        sys.stdout.flush()

    def process(self):
        """Determines the next bot action based on current tick state."""
        players = self.tick.get("players", {})
        me = players.get(self.my_id)

        input_msg = {
            "move_x": 0.0,
            "move_y": 0.0,
            "aim_x": 0.0,
            "aim_y": 0.0,
            "shoot": 0.0,
            "dash": 0.0
        }

        # If I'm dead or not in the game, send empty input
        if not me or me.get("hp", 0) <= 0:
            self.send(input_msg)
            return

        # Logic: Random movement and aiming
        input_msg["move_x"] = random.uniform(-1, 1)
        input_msg["move_y"] = random.uniform(-1, 1)
        input_msg["aim_x"] = random.uniform(-1, 1)
        input_msg["aim_y"] = random.uniform(-1, 1)

        # Shoot and Dash randomly if off cooldown
        if me.get("shoot_cd", 0) == 0:
            input_msg["shoot"] = float(random.randint(0, 1))
        
        if me.get("dash_cd", 0) == 0:
            input_msg["dash"] = float(random.randint(0, 1))

        self.send(input_msg)

# ============================================================================
# 2. MAIN ENTRY POINT
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
                game.setup = state
                print("example_python_bot: setup received", file=sys.stderr)
            elif state_type == "tick":
                game.tick = state
                game.process()

if __name__ == "__main__":
    main()
