{
  description = "A tool to create tmux sessions from a configuration file";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      packages.default = pkgs.buildGoModule {
        pname = "sessionizer";
        version = "0.1.0";
        src = self;

        vendorHash = "sha256-tBIGkNzvwrUzLer8Wa4mntfyym4lnAOf/PEkPS95lgs=";

        goPackagePath = "github.com/salfel/sessionizer";

        nativeBuildInputs = [
          pkgs.fzf
        ];

        buildPhase = ''
          mkdir -p $out/bin
          go build -o $out/bin/sessionizer .
        '';
      };
    });
}
