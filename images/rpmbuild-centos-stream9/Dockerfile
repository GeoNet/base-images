FROM ghcr.io/geonet/base-images/centos:stream9
# Installing tools needed for rpmbuild

RUN dnf update -y && \
    dnf install -y \
    rpm-build rpmdevtools rpm-sign rpmlint git dnf-plugins-core nodejs npm make automake gcc \
    gcc-c++ kernel-devel python3 glibc systemd openssl fontconfig-devel freetype freetype-devel boost unzip

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install
