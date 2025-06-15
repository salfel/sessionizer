# Sessionizer

Sessionizer is a tool to search through your projects and create custom tmux sessions for each project.

## Features
Sessionizer has the following features:

- Search through your projects and create custom tmux sessions for each project.
- Keep track of the most used projects and order them accordingly.
- Supports per-project configuration files.

## Installation

You can install Sessionizer by running the following command:
```bash
go install github.com/salfel/sessionizer@latest
```

Because Sessionizer uses fzf to select the project, you will have to install fzf through your package manager.

### Nix

If you are using NixOS or Home-Manager, you can include the sessionizer flake as a input of your configuration and add the binary to you packages like so:

```nix flake.nix
{
  inputs = {
    sessionizer.url = "github:salfel/sessionizer";
    sessionizer.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, sessionizer, ...
```

```nix configuration.nix
environment.systemPackages = [
    inputs.sessionizer.packages.${system}.default
];

```

## Configuration

To get started, you need to create the configuration file at `~/.config/sessionizer/config.toml`.

Below is an example configuration file:

```toml
search_paths = ["/home/felix/Projects"]
active = "Terminal"

[[windows]]
name = "Editor"
cmd = ["nvim"]

[[windows]]
name = "Git"
cmd = ["lazygit"]

[[windows]]
name = "Terminal"
path = "test"
```

The configuration file is divided into two sections:

- `search_paths`: A list of paths to search for projects.
- `windows`: A list of windows to create for each project.
    - `name`: The name of the window.
    - `path`: The path to the directory to start the window in.
    - `cmd`: A list of commands to run in the window.
    - `active`: The name of the window to activate after starting the session. If not specified, the first window will be activated.

## Local Configuration

You can also add a `sessionizer.toml` file to the project directory to override the global configuration.
That configuration however, does not contain the `search_paths` field as this is a global option.

### Usage

To start a session, run the following command:

```bash
sessionizer
```

This will open a fzf window listing all directories containing a git repository inside them.

When selecting a directory, it will check if the directory contains a `sessionizer.toml` file. If it does, it will override the global window configuration with the local one.
It will then launch a tmux session, set its name to the name of the directory selected with the specified windows and will focus the first window.

If a tmux session with the same name already exists, it will switch to that session.
