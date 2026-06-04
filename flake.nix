{
  description = "Bitbucket MCP server (Go)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # Go toolchain
            go
            gopls
            golangci-lint
            gotools        # goimports, godoc, etc.
            delve          # debugger

            # Python (for agent scripts and tooling)
            python3
            python3Packages.pip

            # General dev utilities
            git
            curl
            jq
            gnumake
          ];

          shellHook = ''
            export GOPATH="$HOME/go"
            export PATH="$GOPATH/bin:$PATH"
            echo "Go $(go version | awk '{print $3}') ready"
          '';
        };
      });
}
