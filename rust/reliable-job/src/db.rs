use sqlx::PgPool;
use uuid::Uuid;

pub async fn enqueue(
    pool: &PgPool,
    job_type: &str,
    payload: serde_json::Value,
) -> Result<(), sqlx::Error> {
    let id = Uuid::new_v4();

    sqlx::query!(
        "INSERT INTO jobs(id, job_type, payload, status) VALUES ($1, $2, $3, 'pending')",
        id,
        job_type,
        payload
    )
    .execute(pool)
    .await?;

    Ok(())
}
