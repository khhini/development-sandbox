use clap::{Args, Parser, Subcommand, command};
use std::error::Error;

pub trait Command {
    fn run(&self) -> Result<(), Box<dyn Error>>;
    fn dry_run(&self) -> bool {
        false
    }
}

#[derive(Subcommand)]
enum Commands {
    Init(InitCommand),
    Deploy(DeployCommand),
}

#[derive(Parser)]
#[command(name = "mycli")]
struct Cli {
    #[command[subcommand]]
    command: Commands,
}

#[derive(Args, Debug)]
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

#[derive(Args, Debug)]
pub struct DeployCommand {
    pub target: String,

    #[arg(short('d'), long("debug"), action = clap::ArgAction::SetTrue)]
    debug: bool,
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
        Commands::Init(cmd) => Box::new(cmd),
        Commands::Deploy(cmd) => Box::new(cmd),
    }
}

fn main() -> Result<(), Box<dyn Error>> {
    let cli = Cli::parse();
    let command = dispatch_command(cli);
    command.run()
}
