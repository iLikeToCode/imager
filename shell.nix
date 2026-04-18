{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
    buildInputs = with pkgs; [
        go_1_26
        gopls
        go-tools
        buf
        protoc-gen-go-grpc
        protoc-gen-go
        grpcurl
        
        xorriso
        syslinux
        mtools
        dosfstools
        curl
        gzip
        qemu
        cpio
    ];
}