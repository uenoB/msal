{
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system: with nixpkgs.legacyPackages.${system}; rec {
        packages.msal = buildGoModule {
          pname = "msal";
          version = "1.0.0";
          src = ./.;
          vendorHash = "sha256-oeJWRMrGAM/qdd+9DzHRMWSQy95Aq4tKboq47cQEQIM=";
        };
        packages.msal-simple = buildGoModule {
          inherit (packages.msal) version src vendorHash;
          pname = "msal-simple";
          GOFLAGS = [ "-tags=simple" ];
        };
        defaultPackage = packages.msal;
      }
    );
}
