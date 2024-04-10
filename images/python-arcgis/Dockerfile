FROM ghcr.io/geonet/base-images/python:3.11.9-alpine3.19

RUN apk add --update \
    ca-certificates tzdata \
    gcc libc-dev g++ make \
    python3-dev libffi-dev openssl-dev gdal-dev \
    proj-util proj-dev krb5-dev geos-dev \
    py3-pip curl cmake autoconf automake libtool elfutils-dev \
    rust cargo

# Installed Cargo (and Rust) is required by one of the indirect dependencies of Poetry.
# Need to add it to PATH.
ENV PATH="${PATH}:/root/.cargo/bin"

ENV PROJ_DIR=/usr
ENV POETRY_REQUESTS_TIMEOUT=30

# Install Poetry via pipx.
RUN pip install pipx
RUN pipx install poetry==1.7.1
RUN pipx inject poetry poetry-plugin-bundle==1.3.0

# Install arcgis with Poetry.
COPY pyproject.toml .
ENV PATH="${PATH}:/root/.local/bin"
RUN poetry install --no-root --only main --no-directory

# Delete Poetry files from image.
RUN rm pyproject.toml
RUN rm poetry.lock