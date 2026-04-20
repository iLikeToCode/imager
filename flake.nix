{
  description = "Alpine TUI OS Imager (Go embedded)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        lib = pkgs.lib;

        imager = pkgs.buildGoModule {
          pname = "imager";
          version = "0.1";

          src = ./imager;

          vendorHash = null;

          env = {
            CGO_ENABLED = 0;
          };

          ldflags = [
            "-s"
            "-w"
          ];

          meta.mainProgram = "client";
        };
      in
      {
        packages = {
          inherit imager;
          default = imager;
          iso = self.nixosConfigurations.iso.config.system.build.isoImage;
        };

        devShells.default = pkgs.mkShell {
          packages = [ pkgs.go ];
        };

        apps.test-vm = {
          type = "app";
          program = "${pkgs.writeShellScript "run-vm" ''
            set -euo pipefail

            ISO=$(ls ${self.packages.${system}.iso}/iso/*.iso)

            exec ${pkgs.qemu}/bin/qemu-system-x86_64 \
              -enable-kvm \
              -cpu host \
              -m 4096 \
              -smp 2 \
              -cdrom "$ISO" \
              -boot d \
              -display gtk \
              -netdev bridge,id=net0,br=vmbr0 \
              -device virtio-net-pci,netdev=net0
          ''}";
        };
      }) // {
      nixosConfigurations.iso = nixpkgs.lib.nixosSystem {
        system = "x86_64-linux";

        modules = [
          "${nixpkgs}/nixos/modules/installer/cd-dvd/installation-cd-minimal.nix"

          ({ pkgs, lib, ... }: {
            environment.systemPackages = [ self.packages.x86_64-linux.imager ];

            systemd.services."imager" = {
              enable = true;
              wantedBy = [ "multi-user.target" ];
              after = [ "multi-user.target" ];
              serviceConfig = {
                ExecStart = "${lib.getExe self.packages.x86_64-linux.imager}";
                StandardInput = "tty";
                StandardOutput = "tty";
                TTYPath = "/dev/tty1";
              };
            };
            systemd.services."autovt@tty1".enable = false;
          })
        ];
      };
    };
}