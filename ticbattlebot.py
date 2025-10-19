#!/usr/bin/env -S uv run --script
import sys

def print_board(board):
    for i in range(0, 9, 3):
        print(board[i:i+3])

def make_move(board, player):
    opponent = 'O' if player == 'X' else 'X'

    # Check if we can win in the next move
    for i in range(9):
        if board[i] == '_':
            new_board = list(board)
            new_board[i] = player
            new_board = ''.join(new_board)
            if check_winner(new_board, player):
                return new_board

    # Check if opponent can win in the next move, block them
    for i in range(9):
        if board[i] == '_':
            new_board = list(board)
            new_board[i] = opponent
            if check_winner(''.join(new_board), opponent):
                new_board[i] = player
                return ''.join(new_board)

    # Try to control the center
    if board[4] == '_':
        new_board = list(board)
        new_board[4] = player
        return ''.join(new_board)

    # Try to control the corners
    for i in [0, 2, 6, 8]:
        if board[i] == '_':
            new_board = list(board)
            new_board[i] = player
            return ''.join(new_board)

    # Play a random move
    for i in range(9):
        if board[i] == '_':
            new_board = list(board)
            new_board[i] = player
            return ''.join(new_board)

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

def main():
    if len(sys.argv) != 3:
        print("Usage: ticbattlebot.py <player> <board>")
        return

    player = sys.argv[1]
    board = sys.argv[2]

    new_board = make_move(board, player)
    print(new_board)

if __name__ == "__main__":
    main()
