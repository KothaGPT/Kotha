use chrono::Utc;
use serde::{Deserialize, Serialize};
use serde_json::json;
use std::io::{self, BufRead, Write};
use std::thread;
use std::time::Duration;

// Temporarily stub out rdev functionality since the private dependency is not accessible
// This allows the build to succeed but disables keyboard event handling

// mod key_codes; // Commented out since rdev is not available

#[cfg(target_os = "macos")]
use cocoa::base::{id, nil};
#[cfg(target_os = "macos")]
use cocoa::foundation::{NSProcessInfo, NSString};
#[cfg(target_os = "macos")]
use objc::{msg_send, sel, sel_impl};

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
struct HotkeyCombo {
    keys: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(tag = "command")]
enum Command {
    #[serde(rename = "register_hotkeys")]
    RegisterHotkeys { hotkeys: Vec<HotkeyCombo> },
}

// Global state for registered hotkeys and currently pressed keys
static mut REGISTERED_HOTKEYS: Vec<HotkeyCombo> = Vec::new();
static mut CURRENTLY_PRESSED: Vec<String> = Vec::new();

// Global state for tracking modifier keys to detect Cmd+C/Ctrl+C combinations
static mut CMD_PRESSED: bool = false;
static mut CTRL_PRESSED: bool = false;
static mut COPY_IN_PROGRESS: bool = false;

/// Prevents macOS App Nap from suspending this process.
/// Returns an activity token that must be retained for the entire process lifetime.
/// On non-macOS platforms, returns a dummy value.
#[cfg(target_os = "macos")]
fn prevent_app_nap() -> id {
    unsafe {
        let process_info = NSProcessInfo::processInfo(nil);
        let reason = NSString::alloc(nil)
            .init_str("Keyboard event monitoring requires continuous operation");

        // NSActivityOptions flags:
        // NSActivityUserInitiated = 0x00FFFFFF (includes all protective flags)
        // This prevents App Nap and idle system sleep
        let options: u64 = 0x00FFFFFF;

        let activity: id = msg_send![process_info, beginActivityWithOptions:options reason:reason];

        eprintln!("macOS App Nap prevention enabled for keyboard listener process");
        activity
    }
}

#[cfg(not(target_os = "macos"))]
fn prevent_app_nap() -> () {
    // No-op on non-macOS platforms
}

fn main() {
    // Prevent macOS App Nap from suspending this process
    // Must retain this for the entire process lifetime
    let _activity = prevent_app_nap();

    // Spawn a thread to read commands from stdin
    thread::spawn(|| {
        let stdin = io::stdin();
        for line in stdin.lock().lines() {
            if let Ok(line) = line {
                match serde_json::from_str::<Command>(&line) {
                    Ok(command) => handle_command(command),
                    Err(e) => eprintln!("Error parsing command: {}", e),
                }
            }
        }
    });

    // Spawn heartbeat thread
    thread::spawn(|| {
        let mut heartbeat_id = 0u64;
        loop {
            thread::sleep(Duration::from_secs(10)); // Send heartbeat every 10 seconds

            heartbeat_id += 1;
            let heartbeat_json = json!({
                "type": "heartbeat_ping",
                "id": heartbeat_id.to_string(),
                "timestamp": Utc::now().to_rfc3339()
            });

            println!("{}", heartbeat_json);
            io::stdout().flush().unwrap();
        }
    });

    // Since rdev is not available, we'll just run an infinite loop
    // that sends periodic heartbeat messages
    eprintln!("Warning: Keyboard event handling is disabled due to missing rdev dependency");
    loop {
        thread::sleep(Duration::from_secs(60)); // Sleep for 1 minute
    }
}

fn handle_command(command: Command) {
    match command {
        Command::RegisterHotkeys { hotkeys } => unsafe {
            REGISTERED_HOTKEYS = hotkeys.clone();
            eprintln!("Registered {} hotkeys", REGISTERED_HOTKEYS.len());
        },
    }
    io::stdout().flush().unwrap();
}
