FROM ghcr.io/geonet/base-images/fedora:41-aarch64
# Installing tools needed for rpmbuild
RUN dnf update -y && \
    dnf install -y \
    rpm-build rpmdevtools rpm-sign rpmlint git dnf-plugins-core nodejs npm make automake gcc \
    gcc-c++ kernel-devel python3 glibc systemd openssl fontconfig-devel freetype freetype-devel boost awscli

COPY rpmlintrc.toml /etc/xdg/rpmlint/rpmlintrc.toml
