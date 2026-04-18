{
  description = "Alpine TUI OS Imager (Go embedded)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = import ./shell.nix { inherit pkgs; };

        packages = import ./default.nix { inherit pkgs; };

        apps.test-vm = {
          type = "app";
          meta = {
            description = "Run the test vm for the imager client";
          };
          program = "${pkgs.writeShellScript "run-iso" ''
            ISO=$(ls ${self.packages.${system}.iso}/*.iso | head -n1)
            exec ${pkgs.qemu}/bin/qemu-system-x86_64 \
              -enable-kvm \
              -cpu host \
              -m 4096 \
              -smp 2 \
              -cdrom "$ISO"
          ''}";
        };
      }
    );
}
