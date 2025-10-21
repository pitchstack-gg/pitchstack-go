{
  description = "Pitchstack Go Client";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    allSystems = [
      "x86_64-linux"
      "aarch64-linux"
      "x86_64-darwin"
      "aarch64-darwin"
    ];

    forAllSystems = f:
      nixpkgs.lib.genAttrs allSystems (system:
        f {
          pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
          };
          inherit system;
        });
  in {
    devShells = forAllSystems ({pkgs, ...}: {
      default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gotools
          go-mockery
          gopls
          gosec
          gofumpt
          golangci-lint
        ];
      };
    });
  };
}
