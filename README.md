# tic

A tic tac toe fight server.

# Spec

## `/tic` (GET)

A health check.

Response: empty.

## `/tic/version` (GET)

A health check.

Response example: `"1760301245"`

## `/tic/game` (POST)

Create a new game.

Body: ignored.

Response example: `"3ed8049c-a674-4e49-a3b0-230730178bd5"`

## `/tic/games` (GET)

Active games.

Response example:

```json
[
  {
    "id": "3ed8049c-a674-4e49-a3b0-230730178bd5",
    "turn": "jane",
    "players": ["jane","gwen"],
    "piece": "X",
    "state": [["_","X","O"],["X","O","_"],["O","_","X"]],
    "result": "", // "win", "draw", ""
    "winner": ""
  }
]
```
  
## `/tic/move` (POST)

Body example: 

```json
{
  "id": "3ed8049c-a674-4e49-a3b0-230730178bd5",
  "user": "jane",
  "row": 0, 
  "col": 0 
}
```

Reponse:

```json
3
```

Number codes:

- Win = 0
- Lose = 1
- Draw = 2
- Illegal = 3
- OK = 4
