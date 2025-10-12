use axum::{
    Json, Router,
    extract::State,
    response::Html,
    routing::{get, post},
};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::net::SocketAddr;
use std::sync::Arc;
use tokio::sync::RwLock;
use std::fs;

#[derive(Debug, Serialize, Deserialize, Clone)]
struct Game {
    id: String,
    turn: String,
    piece: String,
    state: Vec<Vec<String>>,
}

#[derive(Debug, Serialize, Deserialize)]
struct MoveRequest {
    id: String,
    user: String,
    state: Position,
}

#[derive(Debug, Serialize, Deserialize)]
struct Position {
    row: usize,
    col: usize,
}

type GameState = Arc<RwLock<HashMap<String, Game>>>;

#[tokio::main]
async fn main() {
    // Initialize game state
    let game_state = Arc::new(RwLock::new(HashMap::new()));

    // Create a sample game for testing
    let sample_game = Game {
        id: "3ed8049c-a674-4e49-a3b0-230730178bd5".to_string(),
        turn: "jane".to_string(),
        piece: "X".to_string(),
        state: vec![
            vec!["_".to_string(), "X".to_string(), "O".to_string()],
            vec!["X".to_string(), "O".to_string(), "_".to_string()],
            vec!["O".to_string(), "_".to_string(), "X".to_string()],
        ],
    };

    game_state
        .write()
        .await
        .insert(sample_game.id.clone(), sample_game);

    // Build our application with routes
    let app = Router::new()
        .route("/", get(serve_html))
        .route("/tic/challenges", get(get_challenges))
        .route("/tic/move", post(make_move))
        .with_state(game_state);

    // Run the server
    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));
    println!("Server running on http://localhost:3000");
    axum::serve(
        tokio::net::TcpListener::bind(addr).await.unwrap(),
        app.into_make_service(),
    )
    .await
    .unwrap();
}

async fn get_challenges(State(state): State<GameState>) -> Json<Vec<Game>> {
    let games = state.read().await;
    Json(games.values().cloned().collect())
}

async fn make_move(
    State(state): State<GameState>,
    Json(move_request): Json<MoveRequest>,
) -> Json<Game> {
    let mut games = state.write().await;

    if let Some(game) = games.get_mut(&move_request.id) {
        // Verify it's the user's turn
        if game.turn == move_request.user {
            // Verify the move is valid
            if game.state[move_request.state.row][move_request.state.col] == "_" {
                // Make the move
                game.state[move_request.state.row][move_request.state.col] = game.piece.clone();

                // Switch turns and pieces
                game.turn = if game.turn == "jane" {
                    "opponent".to_string()
                } else {
                    "jane".to_string()
                };
                game.piece = if game.piece == "X" {
                    "O".to_string()
                } else {
                    "X".to_string()
                };
            }
        }

        Json(game.clone())
    } else {
        // Return empty game if not found
        Json(Game {
            id: move_request.id,
            turn: "".to_string(),
            piece: "".to_string(),
            state: vec![vec!["".to_string(); 3]; 3],
        })
    }
}

async fn serve_html() -> Html<String> {
    let html_content = fs::read_to_string("index.html")
        .unwrap_or_else(|_| "<h1>Error loading page</h1>".to_string());
    Html(html_content)
}
