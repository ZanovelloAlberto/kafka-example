{
  description = "An ascii webcam in your console";

  inputs = {
    # go.url = ./overlays/go.nix;
    nixpkgs = {
      url = "github:NixOS/nixpkgs";
      # overlays = [
      #   (import ./overlays/go.nix)
      # ];
    };
    flake-utils.url = "github:numtide/flake-utils";

  };

  outputs = { nixpkgs, flake-utils, ... }:

    flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          # res = pkgs.goBuildModule {
          # };
        in
        {

          overlays = [ import ./overlays/go.nix ];
          packages = rec {
            # hello = pkgs.stdenv.mkDerivation {



            #   name = "ll";
            #   buildInput = [ pkgs.go_1_21 ];
            #   src = ./.;

            # };
            hello = pkgs.buildGoModule {
              
              nativeBuildInput = [ pkgs.go];
              vendorHash = "sha256-JOP7hcCOwVZ0hb2UXHHdxpKxpZqs6a8AjOFbrs711ps=";
              name = "cuolo";
              src = ./.;
            };

            default = hello;
          };

        }
      );

}
