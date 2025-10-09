use std::{
    fs::File,
    io::{BufWriter, Write},
};

use flate2::{Compression, write::GzEncoder};
use serde::Serialize;

#[derive(Serialize)]
struct Row {
    id: u64,
    value: String,
}

fn generate_row(id: u64) -> Row {
    Row {
        id,
        value: format!("value_{}", id),
    }
}

struct RotatingWriter {
    base_path: String,
    max_rows_per_file: usize,
    current_file_index: usize,
    current_row_count: usize,
    writer: Option<BufWriter<GzEncoder<File>>>,
}

impl RotatingWriter {
    pub fn new(base_path: String, max_rows_per_file: usize) -> Self {
        Self {
            base_path,
            max_rows_per_file,
            current_file_index: 0,
            current_row_count: 0,
            writer: None,
        }
    }

    fn create_new_writer(&mut self) -> std::io::Result<()> {
        let filename = format!(
            "{}_part_{}.json.gz",
            self.base_path, self.current_file_index
        );
        let file = File::create(filename)?;
        let encoder = GzEncoder::new(file, Compression::default());

        self.writer = Some(BufWriter::new(encoder));
        self.current_row_count = 0;

        Ok(())
    }

    pub fn write_row(&mut self, row: &Row) -> std::io::Result<()> {
        if self.writer.is_none() || self.current_row_count >= self.max_rows_per_file {
            if let Some(w) = self.writer.take() {
                let encoder = w.into_inner()?;
                encoder.finish()?;
            }

            self.current_file_index += 1;
            self.create_new_writer()?;
        }

        if let Some(w) = self.writer.as_mut() {
            let json = serde_json::to_string(row)?;
            w.write_all(json.as_bytes())?;
            w.write_all(b"\n")?;
            self.current_row_count += 1;
        }

        Ok(())
    }

    pub fn finalize(mut self) -> std::io::Result<()> {
        if let Some(w) = self.writer.take() {
            let encoder = w.into_inner()?;
            encoder.finish()?;
        }
        Ok(())
    }
}

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let mut writer = RotatingWriter::new("output/data".to_string(), 1_000);

    for id in 0..10_000 {
        let row = generate_row(id);
        writer.write_row(&row)?;
    }

    writer.finalize()?;
    println!("Export completed.");

    Ok(())
}
