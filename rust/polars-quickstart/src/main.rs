use chrono::NaiveDate;
use polars::{df, frame::DataFrame, prelude::*};

fn create_dataframe() -> DataFrame {
    df!(
    "name" => ["Alice Archer", "Ben Brown", "Chloe Cooper", "Daniel Donovan"],
    "birthdate" => [
        NaiveDate::from_ymd_opt(1997, 1, 10).unwrap(),
        NaiveDate::from_ymd_opt(1985, 2, 15).unwrap(),
        NaiveDate::from_ymd_opt(1983, 3, 22).unwrap(),
        NaiveDate::from_ymd_opt(1981, 4, 30).unwrap(),
    ],
    "weight" => [57.9, 72.5, 53.6, 83.1],  // (kg)
    "height" => [1.56, 1.77, 1.65, 1.75],  // (m)
    )
    .unwrap()
}

fn select_data(data: &DataFrame) -> Option<DataFrame> {
    data.clone()
        .lazy()
        .select([
            col("name"),
            col("birthdate").dt().year().alias("birth_year"),
            (col("weight") / col("height").pow(2)).alias("bmi"),
        ])
        .collect()
        .ok()
}

fn filter_data(data: &DataFrame) -> Option<DataFrame> {
    data.clone()
        .lazy()
        .filter(col("birthdate").dt().year().lt(lit(1990)))
        .collect()
        .ok()
}

fn main() {
    let data = create_dataframe();
    println!("{data}");

    let selected_data = select_data(&data).unwrap();

    println!("{selected_data}");

    let filtered_data = filter_data(&data).unwrap();

    println!("{filtered_data}");
}
