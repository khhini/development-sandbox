mod db;
mod jobs;
mod models;

use sqlx::postgres::PgPoolOptions;

#[tokio::main]
async fn main() -> Result<(), sqlx::Error> {
    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect("postgresql://postgres:postgres@localhost:5432/reliable-job")
        .await?;

    db::enqueue(
        &pool,
        "send_email",
        serde_json::json!({
            "to": "user@example.com"
        }),
    )
    .await?;

    jobs::run_worker(pool).await;

    Ok(())
}
