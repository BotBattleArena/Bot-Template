#include <iostream>
#include <string>
#include <map>
#include <vector>
#include <random>
#include <filesystem>
#include "json.hpp"

using json = nlohmann::json;

// ============================================================================
// 1. DATA STRUCTURES (Server -> Bot)
// ============================================================================

/// Bot represents a player in the arena
struct Bot {
    double x       = 0;
    double y       = 0;
    int    hp      = 0;
    int    kills   = 0;
    int    deaths  = 0;
    int    shoot_cd = 0;
    int    dash_cd  = 0;
    double aim_x   = 0;
    double aim_y   = 0;
};
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE_WITH_DEFAULT(Bot,
    x, y, hp, kills, deaths, shoot_cd, dash_cd, aim_x, aim_y)

/// SceneObject represents walls, bullets, and other geometry
struct SceneObject {
    double x  = 0;
    double y  = 0;
    double w  = 0;
    double h  = 0;
    double r  = 0;
    double dx = 0;
    double dy = 0;
    std::vector<std::array<double, 2>> points;
};
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE_WITH_DEFAULT(SceneObject,
    x, y, w, h, r, dx, dy, points)

struct SceneObjects {
    std::vector<SceneObject> rect;
    std::vector<SceneObject> circle;
    std::vector<SceneObject> poly;
    std::vector<SceneObject> bullets;
};
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE_WITH_DEFAULT(SceneObjects,
    rect, circle, poly, bullets)

/// Scene uses manual from_json/to_json because "static" is a C++ keyword.
/// The JSON field "static" is mapped to the struct field "static_objs".
struct Scene {
    double       mapw        = 0;
    double       maph        = 0;
    SceneObjects static_objs = {};
    SceneObjects dynamic     = {};
};

inline void from_json(const json& j, Scene& s) {
    s.mapw        = j.value("mapw", 0.0);
    s.maph        = j.value("maph", 0.0);
    s.static_objs = j.value("static", SceneObjects{});
    s.dynamic     = j.value("dynamic", SceneObjects{});
}

inline void to_json(json& j, const Scene& s) {
    j = json{
        {"mapw",    s.mapw},
        {"maph",    s.maph},
        {"static",  s.static_objs},
        {"dynamic", s.dynamic}
    };
}

struct Rules {
    double move_speed    = 0;
    double bullet_speed  = 0;
    int    damage        = 0;
    int    win_kills     = 0;
    int    respawn_ticks = 0;
    int    shoot_cd      = 0;
    int    dash_cd       = 0;
    int    max_hp        = 0;
};
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE_WITH_DEFAULT(Rules,
    move_speed, bullet_speed, damage, win_kills,
    respawn_ticks, shoot_cd, dash_cd, max_hp)

/// SetupFrame contains initial game rules and map
struct SetupFrame {
    std::string type             = "";
    int         tick_duration_ms = 0;
    int         timeout_ms       = 0;
    int         max_tick         = 0;
    int         countdown_sec    = 0;
    Scene       scene            = {};
    Rules       rules            = {};
};
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE_WITH_DEFAULT(SetupFrame,
    type, tick_duration_ms, timeout_ms, max_tick,
    countdown_sec, scene, rules)

/// TickFrame is the game state for a single frame
struct TickFrame {
    std::string                type    = "";
    int                        tick    = 0;
    std::map<std::string, Bot> players = {};
    Scene                      scene   = {};
};
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE_WITH_DEFAULT(TickFrame,
    type, tick, players, scene)

// ============================================================================
// 2. OUTPUT STRUCTURES (Bot -> Server)
// ============================================================================

struct InputMsg {
    float move_x = 0;
    float move_y = 0;
    float aim_x  = 0;
    float aim_y  = 0;
    float shoot  = 0;
    float dash   = 0;
};
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE_WITH_DEFAULT(InputMsg,
    move_x, move_y, aim_x, aim_y, shoot, dash)

// ============================================================================
// 3. GAME ENGINE
// ============================================================================

class Game {
public:
    SetupFrame setup;
    TickFrame  tick;
    std::string my_id;

    Game() {
        namespace fs = std::filesystem;
        fs::path exe = std::string(
            #ifdef _WIN32
                __argv[0]
            #else
                program_invocation_name
            #endif
        );
        my_id = exe.stem().string();
    }

    /// Send wraps input and encodes it to stdout
    void send(const InputMsg& msg) {
        json resp;
        resp["axes"] = msg;
        std::cout << resp.dump() << "\n";
        std::cout.flush();
    }

    /// Process determines the next bot action
    void process() {
        InputMsg input;

        auto it = tick.players.find(my_id);
        if (it == tick.players.end() || it->second.hp <= 0) {
            send(input);
            return;
        }
        const Bot& me = it->second;

        // Logic: Random movement and aiming
        input.move_x = random_float(-1.0f, 1.0f);
        input.move_y = random_float(-1.0f, 1.0f);
        input.aim_x  = random_float(-1.0f, 1.0f);
        input.aim_y  = random_float(-1.0f, 1.0f);

        if (me.shoot_cd == 0) {
            input.shoot = static_cast<float>(random_int(0, 1));
        }
        if (me.dash_cd == 0) {
            input.dash = static_cast<float>(random_int(0, 1));
        }

        send(input);
    }

private:
    std::mt19937 rng_{std::random_device{}()};

    float random_float(float lo, float hi) {
        std::uniform_real_distribution<float> dist(lo, hi);
        return dist(rng_);
    }

    int random_int(int lo, int hi) {
        std::uniform_int_distribution<int> dist(lo, hi);
        return dist(rng_);
    }
};

// ============================================================================
// 4. MAIN ENTRY POINT
// ============================================================================

int main() {
    Game game;

    // Log to stderr; stdout is reserved for the protocol
    std::cerr << "example_cpp_bot: started as " << game.my_id << std::endl;

    std::string line;
    while (std::getline(std::cin, line)) {
        json msg;
        try {
            msg = json::parse(line);
        } catch (...) {
            continue;
        }

        std::string msg_type = msg.value("type", "");

        if (msg_type == "start") {
            continue;
        } else if (msg_type == "end") {
            return 0;
        } else if (msg_type == "state") {
            if (!msg.contains("state")) continue;
            json state = msg["state"];

            std::string state_type = state.value("type", "");

            if (state_type == "setup") {
                game.setup = state.get<SetupFrame>();
                std::cerr << "example_cpp_bot: setup received" << std::endl;
            } else if (state_type == "tick") {
                game.tick = state.get<TickFrame>();
                game.process();
            }
        }
    }

    return 0;
}
