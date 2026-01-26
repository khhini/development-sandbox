use bevy::{
    DefaultPlugins,
    app::{App, PluginGroup, Startup},
    asset::AssetPlugin,
    camera::{Camera2d, ClearColor},
    color::Color,
    ecs::system::Commands,
    utils::default,
};
use bevy_game::player::PlayerPlugin;

fn main() {
    App::new()
        .insert_resource(ClearColor(Color::WHITE))
        .add_plugins(DefaultPlugins.set(AssetPlugin {
            file_path: "src/assets".into(),
            ..default()
        }))
        .add_systems(Startup, setup_camera)
        .add_plugins(PlayerPlugin)
        .run();
}

fn setup_camera(mut commands: Commands) {
    commands.spawn(Camera2d);
}
