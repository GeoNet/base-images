FROM ghcr.io/geonet/base-images/centos:stream9
# Installing tools needed for rpmbuild

RUN dnf module enable -y nodejs:20 && \
    dnf update -y && \
    dnf install -y \
    boost automake boost-filesystem boost-iostreams \
    boost-program-options boost-regex boost-system \
    boost-thread cairo cairo-devel dnf-plugins-core fontconfig \
    fontconfig-devel freetype freetype-devel gcc gcc-c++ git glibc \
    kernel-devel libgfortran libxml2 libxml2-devel make ncurses \
    ncurses-devel nodejs npm python3.11-devel python3.11 \
    openssl pango pango-devel perl-devel python3.11 qt5-qtbase \
    rpm-build rpmdevtools rpm-sign rpmlint shadow-utils systemd unzip

COPY rpmlintrc.toml /etc/xdg/rpmlint/rpmlintrc.toml
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install
