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
        version = "v0.1.7";
        src = self;

        vendorHash = "sha256-5s5c4gwtc1thtUUkl8+ntwz1JwJ+UHJGm2voyLJVEeQ=";

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
