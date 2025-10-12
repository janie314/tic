# tic

A tic tac toe fight server.

# Spec

- `/tic` (GET)
- `/tic/game` (POST)
- `/tic/games` (GET)
- `/tic/move` (POST)
- `/tic/results` (GET)
- `/tic/version` (GET)

# Example 

`curl -X POST http://192.168.1.82/tic/move -H 'Content-Type: application/json' -d '{ "id": "d0f124a5-47eb-4153-bd23-6ff1659d24f3", "user": "gwen", "row": 0, "col": 1 }'`