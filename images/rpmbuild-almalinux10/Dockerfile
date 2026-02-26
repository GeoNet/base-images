FROM ghcr.io/geonet/base-images/almalinux:10.1

# Install prerequisites
RUN dnf install -y epel-release 'dnf-command(config-manager)'

# Update system
RUN dnf update -y && \
    dnf install -y \
    automake \
    awscli2 \
    boost-container \
    boost-filesystem \
    boost-iostreams \
    boost-program-options \
    boost-regex \
    boost-system \
    boost-thread \
    cairo \
    cairo-devel \
    dnf-plugins-core \
    fontconfig \
    fontconfig-devel \
    freetype \
    freetype-devel \
    gcc \
    gcc-c++ \
    git \
    glibc \
    kernel-devel \
    libdbi \
    libdbi-devel \
    libgfortran \
    libxml2 \
    libxml2-devel \
    make \
    ncurses \
    ncurses-devel \
    nodejs \
    npm \
    openssl \
    pango \
    pango-devel \
    perl-devel \
    python3 \
    python3-devel \
    python3-numpy \
    python3-pip \
    python3-rpm-macros \
    qt6-qtbase \
    qt6-qtbase-gui \
    qt6-qtsvg \
    rpm-build \
    rpmdevtools \
    rpm-sign \
    rpmlint \
    shadow-utils \
    systemd \
    unzip

COPY rpmlintrc.toml /etc/xdg/rpmlint/rpmlintrc.toml
