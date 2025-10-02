use std::{
    collections::HashMap,
    sync::{Arc, Mutex},
};

fn calcuate_hash(key: &str) -> usize {
    use std::hash::{Hash, Hasher};
    let mut hasher = std::collections::hash_map::DefaultHasher::new();
    key.hash(&mut hasher);
    hasher.finish() as usize
}

fn get_shard(key: &str, shards: &[Arc<Mutex<HashMap<String, String>>>]) -> usize {
    calcuate_hash(key) % shards.len()
}

fn main() {
    let shards: Vec<_> = (0..16)
        .map(|_| Arc::new(Mutex::new(HashMap::new())))
        .collect();

    let shards = Arc::new(shards);

    for _ in 0..8 {
        let shards = Arc::clone(&shards);
        std::thread::spawn(move || {
            for i in 0..10_000 {
                let key = format!("key-{i}");
                let shard_index = get_shard(&key, &shards);
                let mut shard = shards[shard_index].lock().unwrap();

                shard.insert(key, "value".to_string());
            }
        });
    }
}
