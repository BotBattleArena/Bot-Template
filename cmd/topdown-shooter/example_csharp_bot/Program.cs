using System;
using System.Collections.Generic;
using System.IO;
using System.Text.Json;
using System.Text.Json.Serialization;

// ============================================================================
// 1. DATA STRUCTURES (Server -> Bot)
// ============================================================================

/// <summary>ServerMessage is the generic envelope for all communication</summary>
public class ServerMessage
{
    [JsonPropertyName("type")]
    public string Type { get; set; } = "";

    [JsonPropertyName("state")]
    public JsonElement? State { get; set; }
}

/// <summary>Bot represents a player in the arena</summary>
public class Bot
{
    [JsonPropertyName("x")]
    public double X { get; set; }

    [JsonPropertyName("y")]
    public double Y { get; set; }

    [JsonPropertyName("hp")]
    public int HP { get; set; }

    [JsonPropertyName("kills")]
    public int Kills { get; set; }

    [JsonPropertyName("deaths")]
    public int Deaths { get; set; }

    [JsonPropertyName("shoot_cd")]
    public int ShootCD { get; set; }

    [JsonPropertyName("dash_cd")]
    public int DashCD { get; set; }

    [JsonPropertyName("aim_x")]
    public double AimX { get; set; }

    [JsonPropertyName("aim_y")]
    public double AimY { get; set; }
}

/// <summary>SceneObject represents walls, bullets, and other geometry</summary>
public class SceneObject
{
    [JsonPropertyName("x")]
    public double X { get; set; }

    [JsonPropertyName("y")]
    public double Y { get; set; }

    [JsonPropertyName("w")]
    public double W { get; set; }

    [JsonPropertyName("h")]
    public double H { get; set; }

    [JsonPropertyName("r")]
    public double R { get; set; }

    [JsonPropertyName("dx")]
    public double DX { get; set; }

    [JsonPropertyName("dy")]
    public double DY { get; set; }

    [JsonPropertyName("points")]
    public double[][]? Points { get; set; }
}

public class SceneObjects
{
    [JsonPropertyName("rect")]
    public List<SceneObject>? Rect { get; set; }

    [JsonPropertyName("circle")]
    public List<SceneObject>? Circle { get; set; }

    [JsonPropertyName("poly")]
    public List<SceneObject>? Poly { get; set; }

    [JsonPropertyName("bullets")]
    public List<SceneObject>? Bullets { get; set; }
}

public class Scene
{
    [JsonPropertyName("mapw")]
    public double MapW { get; set; }

    [JsonPropertyName("maph")]
    public double MapH { get; set; }

    [JsonPropertyName("static")]
    public SceneObjects? StaticObjects { get; set; }

    [JsonPropertyName("dynamic")]
    public SceneObjects? Dynamic { get; set; }
}

public class Rules
{
    [JsonPropertyName("move_speed")]
    public double MoveSpeed { get; set; }

    [JsonPropertyName("bullet_speed")]
    public double BulletSpeed { get; set; }

    [JsonPropertyName("damage")]
    public int Damage { get; set; }

    [JsonPropertyName("win_kills")]
    public int WinKills { get; set; }

    [JsonPropertyName("respawn_ticks")]
    public int RespawnTicks { get; set; }

    [JsonPropertyName("shoot_cd")]
    public int ShootCD { get; set; }

    [JsonPropertyName("dash_cd")]
    public int DashCD { get; set; }

    [JsonPropertyName("max_hp")]
    public int MaxHP { get; set; }
}

/// <summary>SetupFrame contains initial game rules and map</summary>
public class SetupFrame
{
    [JsonPropertyName("type")]
    public string Type { get; set; } = "";

    [JsonPropertyName("tick_duration_ms")]
    public int TickDurationMs { get; set; }

    [JsonPropertyName("timeout_ms")]
    public int TimeoutMs { get; set; }

    [JsonPropertyName("max_tick")]
    public int MaxTick { get; set; }

    [JsonPropertyName("countdown_sec")]
    public int CountdownSec { get; set; }

    [JsonPropertyName("scene")]
    public Scene? Scene { get; set; }

    [JsonPropertyName("rules")]
    public Rules? Rules { get; set; }
}

/// <summary>TickFrame is the game state for a single frame</summary>
public class TickFrame
{
    [JsonPropertyName("type")]
    public string Type { get; set; } = "";

    [JsonPropertyName("tick")]
    public int Tick { get; set; }

    [JsonPropertyName("players")]
    public Dictionary<string, Bot>? Players { get; set; }

    [JsonPropertyName("scene")]
    public Scene? Scene { get; set; }
}

// ============================================================================
// 2. GAME ENGINE
// ============================================================================

public class Game
{
    public SetupFrame Setup { get; set; } = new();
    public TickFrame TickState { get; set; } = new();
    public string MyID { get; }

    private static readonly Random Rng = new();

    public Game()
    {
        // Derive bot ID from executable name (without extension)
        var processPath = Environment.ProcessPath;
        MyID = Path.GetFileNameWithoutExtension(processPath ?? "example_csharp_bot");
    }

    /// <summary>Send wraps input and encodes it to stdout</summary>
    public void Send(InputMsg input)
    {
        var resp = new { axes = input };
        string json = JsonSerializer.Serialize(resp);
        Console.Out.WriteLine(json);
        Console.Out.Flush();
    }

    /// <summary>Process determines next bot action</summary>
    public void Process()
    {
        var input = new InputMsg();

        Bot? me = null;
        TickState.Players?.TryGetValue(MyID, out me);

        if (me == null || me.HP <= 0)
        {
            Send(input);
            return;
        }

        // Logic: Random movement and aiming
        input.MoveX = (float)(Rng.NextDouble() * 2 - 1);
        input.MoveY = (float)(Rng.NextDouble() * 2 - 1);
        input.AimX = (float)(Rng.NextDouble() * 2 - 1);
        input.AimY = (float)(Rng.NextDouble() * 2 - 1);

        if (me.ShootCD == 0)
        {
            input.Shoot = Rng.Next(2);
        }
        if (me.DashCD == 0)
        {
            input.Dash = Rng.Next(2);
        }

        Send(input);
    }
}

// ============================================================================
// 3. OUTPUT STRUCTURES (Bot -> Server)
// ============================================================================

public class InputMsg
{
    [JsonPropertyName("move_x")]
    public float MoveX { get; set; }

    [JsonPropertyName("move_y")]
    public float MoveY { get; set; }

    [JsonPropertyName("aim_x")]
    public float AimX { get; set; }

    [JsonPropertyName("aim_y")]
    public float AimY { get; set; }

    [JsonPropertyName("shoot")]
    public float Shoot { get; set; }

    [JsonPropertyName("dash")]
    public float Dash { get; set; }
}

// ============================================================================
// 4. MAIN ENTRY POINT
// ============================================================================

public class Program
{
    public static void Main(string[] args)
    {
        var game = new Game();

        // Log to stderr; stdout is reserved for the protocol
        Console.Error.WriteLine($"example_csharp_bot: started as {game.MyID}");

        string? line;
        while ((line = Console.ReadLine()) != null)
        {
            ServerMessage? msg;
            try
            {
                msg = JsonSerializer.Deserialize<ServerMessage>(line);
            }
            catch
            {
                continue;
            }

            if (msg == null) continue;

            switch (msg.Type)
            {
                case "start":
                    continue;
                case "end":
                    return;
                case "state":
                    if (msg.State == null) continue;

                    var stateJson = msg.State.Value.GetRawText();

                    // Read the "type" field to route setup vs tick
                    using (var doc = JsonDocument.Parse(stateJson))
                    {
                        var stateType = doc.RootElement.GetProperty("type").GetString();

                        switch (stateType)
                        {
                            case "setup":
                                game.Setup = JsonSerializer.Deserialize<SetupFrame>(stateJson)!;
                                Console.Error.WriteLine("example_csharp_bot: setup received");
                                break;
                            case "tick":
                                var tick = JsonSerializer.Deserialize<TickFrame>(stateJson);
                                if (tick != null)
                                {
                                    game.TickState = tick;
                                    game.Process();
                                }
                                break;
                        }
                    }
                    break;
            }
        }
    }
}
