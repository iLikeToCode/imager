{
  description = "Alpine TUI ISO builder";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    naersk = {
      url = "github:nix-community/naersk";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    fenix = {
      url = "github:nix-community/fenix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, flake-utils, naersk, fenix }:
  let
    system = "x86_64-linux";
    pkgs = import nixpkgs { inherit system; };
  in
  {
    devShells.${system}.default = pkgs.mkShell {
      buildInputs = with pkgs; [
        xorriso
        squashfsTools
        syslinux
        mtools
        dosfstools
        curl
        gzip
        qemu
        cpio
      ];
    };

    packages.${system}.iso = pkgs.stdenv.mkDerivation rec {
      pname = "imager";
      version = "0.1";

      src = ./.;

      alpineNetboot = pkgs.fetchurl {
        url = "https://dl-cdn.alpinelinux.org/alpine/v3.23/releases/x86_64/alpine-netboot-3.23.3-x86_64.tar.gz";
        sha256 = "sha256-U/tUZvdhLU/2Fr3g9jfwuM0mfX5SrtxwUiD0h+Qx8VA=";
      };

      imager = 
        let
          pkgs = nixpkgs.legacyPackages.${system};
          target = "x86_64-unknown-linux-musl";
          toolchain = with fenix.packages.${system}; combine [
            minimal.cargo
            minimal.rustc
            targets.${target}.latest.rust-std
          ];
        in

        (naersk.lib.${system}.override {
          cargo = toolchain;
          rustc = toolchain;
        }).buildPackage {
          src = ./imager;
          CARGO_BUILD_TARGET = target;
        };

      alpineRootfs = pkgs.fetchurl {
        url = "https://dl-cdn.alpinelinux.org/alpine/v3.22/releases/x86_64/alpine-minirootfs-3.22.0-x86_64.tar.gz";
        sha256 = "sha256-GIeYhONbBxjwF6UP+FteZWgnnpcjP8QoIiKVhf6y+k0=";
      };

      nativeBuildInputs = with pkgs; [
        xorriso
        squashfsTools
        syslinux
        mtools
        dosfstools
        curl
        gzip
        cpio
      ];

      buildPhase = ''
        mkdir -p work
        cd work

        # extract minimal Alpine rootfs
        mkdir rootfs
        tar -xzf $alpineRootfs -C rootfs
        
        # inject Rust TUI binary
        mkdir -p rootfs/usr/local/bin
        cp $imager/bin/imager rootfs/usr/local/bin/

        # simple init script for initramfs
        cat > rootfs/init <<'EOF'
#!/bin/sh
mount -t proc proc /proc
mount -t sysfs sys /sys
mount -t devtmpfs udev /dev
/usr/local/bin/imager
#poweroff -f
sh
EOF
        chmod +x rootfs/init
        mkdir -p iso/boot

        # create initramfs
        cd rootfs
        find . | cpio -o -H newc | gzip > ../iso/boot/initramfs.cpio.gz
        cd ..

        # copy kernel from Alpine netboot
        tar -xzf $alpineNetboot --strip-components=1 -C iso/boot boot/vmlinuz-lts

        # prepare ISO tree for Syslinux
        mkdir -p iso/boot/syslinux
      '';

      installPhase = ''
        mkdir -p $out

        # copy syslinux bootloader
        cp -r ${pkgs.syslinux}/share/syslinux/* iso/boot/syslinux/

        # generate minimal isolinux.cfg
        cat > iso/boot/syslinux/isolinux.cfg <<'EOF'
DEFAULT alpine
LABEL alpine
  KERNEL /boot/vmlinuz-lts
  INITRD /boot/initramfs.cpio.gz
  APPEND rw
EOF

        # build ISO
        xorriso -as mkisofs \
          -o $out/${pname}.iso \
          -b boot/syslinux/isolinux.bin \
          -c boot/syslinux/boot.cat \
          -no-emul-boot \
          -boot-load-size 4 \
          -boot-info-table \
          iso/
      '';
    };

    apps.${system}.default = {
      type = "app";
      program = "${pkgs.writeShellScript "run-iso" ''
        ISO=$(ls ${self.packages.${system}.iso}/*.iso | head -n1)
        exec ${pkgs.qemu}/bin/qemu-system-x86_64 \
          -enable-kvm \
          -cpu host \
          -m 1024 \
          -smp $(nproc) \
          -cdrom "$ISO"
      ''}";
    };

    defaultPackage.${system} = self.packages.${system}.iso;
  };
}