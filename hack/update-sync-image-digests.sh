#!/bin/bash

cd "$(git rev-parse --show-toplevel)" || exit 1

for IMAGE in $(< ./config.yaml yq e '.. comments="" |= . | .sync[] | select(.auto-update-mutable-tag-digest==true) | .source' -P -); do
    IMAGE_BASE="$(echo "$IMAGE" | cut -d'@' -f1)"
    TAG="$(echo "$IMAGE" | cut -d: -f2 | cut -d@ -f1)"
    DATE="$(date +%Y-%m-%d)"
    NEW_IMAGE_DIGEST="$(crane digest "$IMAGE_BASE")"
    NEW_IMAGE="$IMAGE_BASE@$NEW_IMAGE_DIGEST"
    if [ "$IMAGE" = "$NEW_IMAGE" ]; then
        echo "notice: image '$IMAGE_BASE' is already up to date"
        continue
    fi
    echo "updating: image '$IMAGE_BASE'"
    export FROM="$IMAGE" IMAGE="$NEW_IMAGE # $TAG for $DATE"
    yq e -i 'with(.sync[] | select(.source==env(FROM)); .source = env(IMAGE))' ./config.yaml
done
