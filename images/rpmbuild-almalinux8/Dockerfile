FROM ghcr.io/geonet/base-images/almalinux:8.10
# Install prerequisites
RUN dnf module enable -y nodejs:20 \
  && dnf install -y epel-release 'dnf-command(config-manager)' \
  && dnf config-manager --set-enabled powertools
# Update system
RUN dnf update -y && \
  dnf install -y boost automake boost-filesystem boost-iostreams \
  boost-program-options boost-regex boost-signals boost-system \
  boost-thread cairo cairo-devel dnf-plugins-core fontconfig \
  fontconfig-devel freetype freetype-devel gcc gcc-c++ git glibc \
  kernel-devel libdbi libdbi-devel libgfortran libxml2 \
  libxml2-devel make ncurses ncurses-devel nodejs npm octave \
  openssl pango pango-devel perl-devel python3.11 qt5-qtwebengine \
  rpm-build rpmdevtools rpm-sign rpmlint shadow-utils systemd unzip \
  python3.11-devel python3.11 python3.11-rpm-macros

COPY rpmlintrc.toml /etc/xdg/rpmlint/rpmlintrc.toml
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && unzip awscliv2.zip
RUN ./aws/install
