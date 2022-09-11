{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    nix-shell.url = "github:loophp/nix-shell";
  };

  outputs = { self, nixpkgs, flake-utils, nix-shell }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import nixpkgs {
            inherit system;
          };
        in
        {
          devShells = {
            default = pkgs.mkShellNoCC {
              name = "PHP project";

              buildInputs = [
                pkgs.postgresql_10
              ];
            };
          };
        }
      );
}