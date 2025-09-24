FROM ghcr.io/geonet/base-images/python:3.11-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates tzdata \
    build-essential \
    python3-dev libffi-dev libssl-dev libgdal-dev \
    proj-bin libproj-dev libkrb5-dev libgeos-dev \
    python3-pip curl cmake autoconf automake libtool libdw-dev \
    rustc cargo \
 && rm -rf /var/lib/apt/lists/*

# Installed Cargo (and Rust) is required by one of the indirect dependencies of Poetry.
# Need to add it to PATH.
ENV PATH="${PATH}:/root/.cargo/bin"

ENV PROJ_DIR=/usr
ENV POETRY_REQUESTS_TIMEOUT=30

# Install Poetry via pipx.
RUN pip install pipx==1.7.1
RUN pipx install poetry==2.1.3
RUN pipx inject poetry poetry-plugin-bundle==1.7.0

# Install arcgis with Poetry.
COPY pyproject.toml .
ENV PATH="${PATH}:/root/.local/bin"
RUN poetry install --no-root --only main --no-directory

# Delete Poetry files from image.
RUN rm pyproject.toml
RUN rm poetry.lock