use std::ops::DerefMut;

use sqlx::{PgPool, Postgres, Transaction};
use tokio::time::{Duration, sleep};

use crate::models::Job;

pub async fn run_worker(pool: PgPool) {
    loop {
        match fetch_and_lock_job(&pool).await {
            Ok(Some(job)) => {
                if let Err(e) = process_job(&pool, job).await {
                    eprintln!("Job failed: {e:?}");
                }
            }
            Ok(None) => {
                sleep(Duration::from_secs(2)).await;
            }
            Err(e) => {
                eprintln!("DB error: {e:?}");
                sleep(Duration::from_secs(5)).await;
            }
        }
    }
}

async fn fetch_and_lock_job(pool: &PgPool) -> Result<Option<Job>, sqlx::Error> {
    let mut tx: Transaction<'_, Postgres> = pool.begin().await?;

    let job = sqlx::query_as!(
        Job,
        "SELECT * FROM jobs WHERE status = 'pending' ORDER BY created_at LIMIT 1 FOR UPDATE SKIP LOCKED"
    )
    .fetch_optional(tx.deref_mut())
    .await?;

    if let Some(ref job) = job {
        sqlx::query!(
            "UPDATE jobs SET status = 'processing', updated_at = now() WHERE id = $1",
            job.id
        )
        .execute(tx.deref_mut())
        .await?;
    }

    tx.commit().await?;
    Ok(None)
}

async fn process_job(pool: &PgPool, job: Job) -> Result<(), sqlx::Error> {
    println!("Processing job {} of type {}", job.id, job.job_type);
    tokio::time::sleep(Duration::from_secs(1)).await;

    sqlx::query!(
        "UPDATE jobs SET status = 'done', updated_at = now() WHERE id = $1",
        job.id
    )
    .execute(pool)
    .await?;
    Ok(())
}
