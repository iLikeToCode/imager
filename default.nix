{ pkgs ? import <nixpkgs> {} }:

rec {
  imager = pkgs.buildGo126Module {
    pname = "imager";
    version = "0.1";

    src = ./imager;

    vendorHash = null;
    
    ldflags = [ "-s" "-w" ];
    doCheck = false;
  };

  iso = pkgs.stdenv.mkDerivation rec {
    pname = "imager-iso";
    version = "0.1";

    src = ./.;

    alpineNetboot = pkgs.fetchurl {
      url = "https://dl-cdn.alpinelinux.org/alpine/v3.23/releases/x86_64/alpine-netboot-3.23.3-x86_64.tar.gz";
      sha256 = "sha256-U/tUZvdhLU/2Fr3g9jfwuM0mfX5SrtxwUiD0h+Qx8VA=";
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
      export HOME=$TMPDIR
      mkdir -p work
      cd work

      mkdir -p build
      cp ${pkgs.lib.getExe' imager "imager"} ./build/imager

      mkdir rootfs
      tar -xzf $alpineRootfs -C rootfs

      mkdir -p rootfs/usr/local/bin
      cp build/imager rootfs/usr/local/bin/

      cat > rootfs/init <<'EOF'
#!/bin/sh
mount -t proc proc /proc
mount -t sysfs sys /sys
mount -t devtmpfs udev /dev
/usr/local/bin/imager
sh
EOF
      chmod +x rootfs/init

      mkdir -p iso/boot

      cd rootfs
      find . | cpio -o -H newc | gzip > ../iso/boot/initramfs.cpio.gz
      cd ..

      tar -xzf $alpineNetboot --strip-components=1 -C iso/boot boot/vmlinuz-lts

      mkdir -p iso/boot/syslinux
    '';

    installPhase = ''
      mkdir -p $out

      cp -r ${pkgs.syslinux}/share/syslinux/* iso/boot/syslinux/

      cat > iso/boot/syslinux/isolinux.cfg <<'EOF'
DEFAULT alpine
LABEL alpine
  KERNEL /boot/vmlinuz-lts
  INITRD /boot/initramfs.cpio.gz
  APPEND rw quiet
EOF

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
  default = iso;
}