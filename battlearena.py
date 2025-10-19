#!/usr/bin/env -S uv run --script
import sys
import subprocess

def run_bot(bot, player, board):
    try:
        output = subprocess.check_output([bot, player, board]).decode('utf-8').strip()
        return output
    except subprocess.CalledProcessError as e:
        print(f"Error running {bot}: {e}")
        sys.exit(1)

def check_winner(board, player):
    win_conditions = [
        (0, 1, 2),
        (3, 4, 5),
        (6, 7, 8),
        (0, 3, 6),
        (1, 4, 7),
        (2, 5, 8),
        (0, 4, 8),
        (2, 4, 6),
    ]

    for condition in win_conditions:
        if board[condition[0]] == player and board[condition[1]] == player and board[condition[2]] == player:
            return True

    return False

def check_draw(board):
    return '_' not in board

def main():
    if len(sys.argv) != 3:
        print("Usage: python battlearena.py <bot1.py> <bot2.py>")
        return

    bot1 = sys.argv[1]
    bot2 = sys.argv[2]

    board = '_' * 9
    players = ['X', 'O']
    bots = [bot1, bot2]

    current_player = 0
    while True:
        bot = bots[current_player]
        player = players[current_player]
        board = run_bot(bot, player, board)
        print(f"{bot},{player}")
        print(f"Board after {bot}:")
        for i in range(0, 9, 3):
            print(board[i:i+3])

        if check_winner(board, player):
            print(f"Winner: {player}")
            break

        if check_draw(board):
            print("It's a draw!")
            break

        current_player = (current_player + 1) % 2

if __name__ == "__main__":
    main()
