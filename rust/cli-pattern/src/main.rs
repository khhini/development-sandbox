use clap::{Parser, Subcommand};
use std::error::Error;

pub trait Command {
    fn run(&self) -> Result<(), Box<dyn Error>>;
    fn dry_run(&self) -> bool {
        false
    }
}

#[derive(Subcommand)]
enum Commands {
    Init,
    Deploy { target: String },
}

#[derive(Parser)]
#[command(name = "mycli")]
struct Cli {
    #[command[subcommand]]
    command: Commands,
}

pub struct InitCommand;

impl Command for InitCommand {
    fn run(&self) -> Result<(), Box<dyn Error>> {
        println!("Initializing project...");
        Ok(())
    }

    fn dry_run(&self) -> bool {
        true
    }
}

pub struct DeployCommand {
    pub target: String,
}

impl Command for DeployCommand {
    fn run(&self) -> Result<(), Box<dyn Error>> {
        println!("Deploying to {}", self.target);
        Ok(())
    }

    fn dry_run(&self) -> bool {
        true
    }
}

fn dispatch_command(cli: Cli) -> Box<dyn Command> {
    match cli.command {
        Commands::Init => Box::new(InitCommand),
        Commands::Deploy { target } => Box::new(DeployCommand { target }),
    }
}

fn main() -> Result<(), Box<dyn Error>> {
    let cli = Cli::parse();
    let command = dispatch_command(cli);
    command.run()
}
