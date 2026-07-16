use std::{
    collections::HashMap,
    sync::{Arc, RwLock},
};

use axum::{
    Json, Router,
    extract::{Path, State},
    http::StatusCode,
    response::{IntoResponse, Redirect},
    routing::{get, post},
};
use serde::{Deserialize, Serialize};
use tokio::net::TcpListener;
use uuid::Uuid;

#[derive(Clone)]
struct AppState {
    urls: Arc<RwLock<HashMap<String, String>>>,
}

#[derive(Deserialize)]
struct ShortenRequest {
    url: String,
}

#[derive(Serialize)]
struct ShortenResponse {
    short_url: String,
}

#[tokio::main]
async fn main() {
    let state = AppState {
        urls: Arc::new(RwLock::new(HashMap::new())),
    };

    let app = Router::new()
        .route("/shorten", post(shorten))
        .route("/{code}", get(redirect))
        .with_state(state);

    let listener = TcpListener::bind("127.0.0.1:8080").await.unwrap();
    println!("Listening on http://127.0.0.1:8080");

    axum::serve(listener, app).await.unwrap();
}

async fn shorten(
    State(state): State<AppState>,
    Json(payload): Json<ShortenRequest>,
) -> impl IntoResponse {
    let code = &Uuid::new_v4().to_string()[..8];

    state
        .urls
        .write()
        .unwrap()
        .insert(code.to_string(), payload.url);

    let short_url = format!("http://127.0.0.1:8080/{}", code);

    (StatusCode::CREATED, Json(ShortenResponse { short_url }))
}

async fn redirect(State(state): State<AppState>, Path(code): Path<String>) -> impl IntoResponse {
    let urls = state.urls.read().unwrap();
    match urls.get(&code) {
        Some(long_url) => Redirect::to(long_url).into_response(),
        None => (StatusCode::NOT_FOUND, "404 Not Found: No URL for this code").into_response(),
    }
}
