FROM ghcr.io/geonet/base-images/centos:stream9
# Installing tools needed for rpmbuild
RUN dnf update -y && \
    dnf install -y \
    rpm-build rpmdevtools rpm-sign rpmlint git dnf-plugins-core nodejs npm make automake gcc \
    gcc-c++ kernel-devel python3 glibc
