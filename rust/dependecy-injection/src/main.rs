pub trait Logger {
    fn log(&self, message: &str);
}

pub struct ConsoleLogger;

impl Logger for ConsoleLogger {
    fn log(&self, message: &str) {
        println!("[LOG]: {}", message)
    }
}

pub struct MockLogger;
impl Logger for MockLogger {
    fn log(&self, message: &str) {
        println!("[MOCK LOG]: {}", message);
    }
}

pub struct UserService<'a> {
    logger: &'a dyn Logger,
}

impl<'a> UserService<'a> {
    pub fn new(logger: &'a dyn Logger) -> Self {
        Self { logger }
    }

    pub fn create_user(&self, username: &str) {
        self.logger.log(&format!("Created user: {}", username))
    }
}

#[test]
fn test_user_creation_logs() {
    let logger = MockLogger;
    let service = UserService::new(&logger);
    service.create_user("test_user");
}

fn main() {
    let logger = MockLogger;
    let service = UserService::new(&logger);
    service.create_user("khhini");
}
