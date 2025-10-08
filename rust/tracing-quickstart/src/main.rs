use tracing::{debug, instrument};
use tracing_subscriber::FmtSubscriber;

#[tokio::main]
async fn main() {
    let sub = FmtSubscriber::builder()
        .with_max_level(tracing::Level::DEBUG)
        .finish();
    tracing::subscriber::set_global_default(sub).unwrap();

    work().await;
}

#[instrument]
async fn work() {
    let number = 2;
    debug!(name: "competed", number, "test number is {:}", number);
    dbg!(number);

    //
}
