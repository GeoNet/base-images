FROM ghcr.io/geonet/base-images/centos:centos7
# Install prerequisites
RUN curl -O https://nodejs.org/dist/latest-v20.x/node-v20.10.0-linux-arm64.tar.xz \
  && tar --strip-components 1 -xvf node-v* -C /usr/local \
  && yum install -y epel-release

# Update system
RUN yum update -y && \
  yum install -y autoconf automake boost boost-filesystem boost-iostreams \
  boost-program-options boost-regex boost-signals boost-system boost-thread \
  cairo cairo-devel createrepo dnf dnf-plugins-core fontconfig fontconfig-devel \
  freetype freetype-devel gcc gcc-c++ git libdbi libdbi-devel libgfortran \
  libxml2 libxml2-devel m2crypto make ncurses ncurses-devel numpy octave \
  openssl openssl11 pango pango-devel perl-devel python3 qtwebkit rpm-build \
  rpmdevtools rpmlint rpm-sign shadow-utils systemd tar unzip

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install
