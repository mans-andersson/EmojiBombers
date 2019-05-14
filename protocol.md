# Protocol
All communication is handled by sending JSON string, ending with a semicolon (;).

## Client
The clients send very simple objects simply describing the action performed by the client. The possible actions are: up, down left, right and space. These are mapped to corresponding keypresses which are the arrow keys and the spacebar. This could of course be changed if you implemented your own client.

Example:

```json
{ "action":"up"};
```

The key is always "action" and "up" is just an example. It could be any of the possible actions described above.

## Server
The server also sends json obects ending with a semicolon. The server holds a state for the game and sends it continously to the clients. The state being sent has the following keys: "players", "blockades", "bombs", "snowflakes", "placed_bombs", "explosions", "winner".

Example:

```json
{
  "players": [
    {
      "id": 0,
      "x_pos": 100,
      "y_pos": 100,
      "spawned": true,
      "bomb_count": 0,
      "lives": 3,
      "damage_taken": false,
      "stunned": false,
      "dead": false
    },
    {
      "id": 1,
      "x_pos": 700,
      "y_pos": 700,
      "spawned": false,
      "bomb_count": 0,
      "lives": 3,
      "damage_taken": false,
      "stunned": false,
      "dead": false
    },
    {
      "id": 2,
      "x_pos": 100,
      "y_pos": 700,
      "spawned": false,
      "bomb_count": 0,
      "lives": 3,
      "damage_taken": false,
      "stunned": false,
      "dead": false
    },
    {
      "id": 3,
      "x_pos": 700,
      "y_pos": 100,
      "spawned": false,
      "bomb_count": 0,
      "lives": 3,
      "damage_taken": false,
      "stunned": false,
      "dead": false
    }
  ],
  "blockades": [
    {
      "x_pos": 0,
      "y_pos": 0
    },
    {
      "x_pos": 0,
      "y_pos": 200
    },
    {
      "x_pos": 0,
      "y_pos": 400
    },
    {
      "x_pos": 0,
      "y_pos": 600
    },
    {
      "x_pos": 0,
      "y_pos": 800
    },
    {
      "x_pos": 200,
      "y_pos": 0
    },
    {
      "x_pos": 200,
      "y_pos": 200
    },
    {
      "x_pos": 200,
      "y_pos": 400
    },
    {
      "x_pos": 200,
      "y_pos": 600
    },
    {
      "x_pos": 200,
      "y_pos": 800
    },
    {
      "x_pos": 400,
      "y_pos": 0
    },
    {
      "x_pos": 400,
      "y_pos": 200
    },
    {
      "x_pos": 400,
      "y_pos": 400
    },
    {
      "x_pos": 400,
      "y_pos": 600
    },
    {
      "x_pos": 400,
      "y_pos": 800
    },
    {
      "x_pos": 600,
      "y_pos": 0
    },
    {
      "x_pos": 600,
      "y_pos": 200
    },
    {
      "x_pos": 600,
      "y_pos": 400
    },
    {
      "x_pos": 600,
      "y_pos": 600
    },
    {
      "x_pos": 600,
      "y_pos": 800
    },
    {
      "x_pos": 800,
      "y_pos": 0
    },
    {
      "x_pos": 800,
      "y_pos": 200
    },
    {
      "x_pos": 800,
      "y_pos": 400
    },
    {
      "x_pos": 800,
      "y_pos": 600
    },
    {
      "x_pos": 800,
      "y_pos": 800
    }
  ],
  "bombs": [
    {
      "spawned": false,
      "x_pos": 100,
      "y_pos": 100
    },
    {
      "spawned": false,
      "x_pos": 100,
      "y_pos": 300
    },
    {
      "spawned": false,
      "x_pos": 100,
      "y_pos": 500
    },
    {
      "spawned": false,
      "x_pos": 100,
      "y_pos": 700
    },
    {
      "spawned": false,
      "x_pos": 300,
      "y_pos": 100
    },
    {
      "spawned": false,
      "x_pos": 300,
      "y_pos": 300
    },
    {
      "spawned": false,
      "x_pos": 300,
      "y_pos": 500
    },
    {
      "spawned": false,
      "x_pos": 300,
      "y_pos": 700
    },
    {
      "spawned": false,
      "x_pos": 500,
      "y_pos": 100
    },
    {
      "spawned": false,
      "x_pos": 500,
      "y_pos": 300
    },
    {
      "spawned": false,
      "x_pos": 500,
      "y_pos": 500
    },
    {
      "spawned": false,
      "x_pos": 500,
      "y_pos": 700
    },
    {
      "spawned": false,
      "x_pos": 700,
      "y_pos": 100
    },
    {
      "spawned": false,
      "x_pos": 700,
      "y_pos": 300
    },
    {
      "spawned": false,
      "x_pos": 700,
      "y_pos": 500
    },
    {
      "spawned": false,
      "x_pos": 700,
      "y_pos": 700
    }
  ],
  "snowflakes": [
    {
      "spawned": false,
      "x_pos": 100,
      "y_pos": 100
    },
    {
      "spawned": false,
      "x_pos": 100,
      "y_pos": 300
    },
    {
      "spawned": false,
      "x_pos": 100,
      "y_pos": 500
    },
    {
      "spawned": false,
      "x_pos": 100,
      "y_pos": 700
    },
    {
      "spawned": false,
      "x_pos": 300,
      "y_pos": 100
    },
    {
      "spawned": false,
      "x_pos": 300,
      "y_pos": 300
    },
    {
      "spawned": false,
      "x_pos": 300,
      "y_pos": 500
    },
    {
      "spawned": false,
      "x_pos": 300,
      "y_pos": 700
    },
    {
      "spawned": false,
      "x_pos": 500,
      "y_pos": 100
    },
    {
      "spawned": false,
      "x_pos": 500,
      "y_pos": 300
    },
    {
      "spawned": false,
      "x_pos": 500,
      "y_pos": 500
    },
    {
      "spawned": false,
      "x_pos": 500,
      "y_pos": 700
    },
    {
      "spawned": false,
      "x_pos": 700,
      "y_pos": 100
    },
    {
      "spawned": false,
      "x_pos": 700,
      "y_pos": 300
    },
    {
      "spawned": false,
      "x_pos": 700,
      "y_pos": 500
    },
    {
      "spawned": false,
      "x_pos": 700,
      "y_pos": 700
    }
  ],
  "placed_bombs": [
    {
      "spawned": false,
      "x_pos": 0,
      "y_pos": 0
    },
    {
      "spawned": false,
      "x_pos": 0,
      "y_pos": 0
    },
    {
      "spawned": false,
      "x_pos": 0,
      "y_pos": 0
    },
    {
      "spawned": false,
      "x_pos": 0,
      "y_pos": 0
    }
  ],
  "explosions": [
    {
      "spawned": false,
      "x_pos": 0,
      "y_pos": 0
    },
    {
      "spawned": false,
      "x_pos": 0,
      "y_pos": 0
    },
    {
      "spawned": false,
      "x_pos": 0,
      "y_pos": 0
    },
    {
      "spawned": false,
      "x_pos": 0,
      "y_pos": 0
    }
  ],
  "winner": -1
};
```

Each object in "players", "blockades", "bombs", "snowflakes", "placed_bombs" and "explosions" should be continously rendered at their corresponding "x_pos" and "y_pos" where "x_pos" and "y_pos" describe the center of the image. The players should be rendered with different images for different "id"s to be able to tell them apart. The possible ids are 0, 1, 2 and 3. The image size for players should be 100 x 100. The image size for blockades, bombs, snowflakes and placed_bombs should be 50 x 50. The image size for explosions should be 200 x 200. If you do not adhere to these recomendations when implementing your own client it will still work but the collision detection will look strange. You should render players with different images depending on their state. If "dead" is true you should have a special image to show that and the same goes for "stunned" and "damage_taken". You should only render "players", "blockades", "bombs", "snowlakes", "placed_bombs" and "explosions" if "spawned" is true. "blockades" should always be rendered. The "winner" field holds information on what player has won the current game. -1 means that there is no winner at the moment. Otherwise it will hold the id of the winner and then you should present a screen describing what player won the game. You can identify the player by rendering the image used for the player with that particular id if you used different images for players with different id's. After a delay "winner" will be set to -1 again and the game restarts so you should be able to render the game again.
