#!/usr/bin/env bash
# Copies the images of the NeuVector charts from the Prime registry, where the NeuVector
# releases push them, to the staging images registry, where the ship-family-staged action
# looks for them before it pushes a chart. When the Rancher release team releases the chart,
# cargueiro release moves each image on from staging to production by its digest, and leaves
# an image that production has already as it is.
#
# Each image is copied as it is: the index and every platform, with the same digests
# (skopeo copy --all --preserve-digests), and then what refers to it, the signature and the
# SLSA provenance that the NeuVector release made (cargueiro copy-image-signatures). So
# staging has the image that Prime has, signed as it is there, and a release candidate is
# tested with the image that ships. The Docker Hub copy is not used: it is another build, and
# it is not signed.
#
# For each image:
#   staging has the digest that Prime has   nothing is copied
#   staging has another digest, or none     it is copied: a tag that moves, such as the 6 of
#                                           the scanner, follows Prime
#   Prime does not have it                  an error, because the chart would name an image
#                                           that is not released
# Every image is looked at before anything is copied, and an error stops the copy of all.
# With DRY_RUN=true nothing is written: it says what would be copied.
#
# Usage: copy-images.sh <list>
#   <list>  a file with one image per line, as rancher/<name>:<tag>, such as
#           rancher/neuvector-controller:5.7.0
#
# Environment:
#   PRIME_REGISTRY, PRIME_REGISTRY_USERNAME, PRIME_REGISTRY_PASSWORD
#       the Prime registry, which the images are read from
#   STAGING_REGISTRY, STAGING_REGISTRY_USERNAME, STAGING_REGISTRY_PASSWORD
#       the staging images registry, which they are written to
#   DRY_RUN     true (the default) or false
#   PLAIN_HTTP  true to talk HTTP to both registries, only for local registries in a test
# A registry with no username is read anonymously.
#
# When GITHUB_OUTPUT is set, it writes pending=<the number of images that were copied, or
# would be>. When GITHUB_STEP_SUMMARY is set, it adds a table of the images to it.
#
# Requires skopeo, and for a real run cargueiro, on PATH.
set -euo pipefail

LIST=${1:?usage: $0 <file with one image per line, such as rancher/neuvector-controller:5.7.0>}
PRIME_REGISTRY=${PRIME_REGISTRY:?PRIME_REGISTRY is required: the host of the Prime registry}
STAGING_REGISTRY=${STAGING_REGISTRY:?STAGING_REGISTRY is required: the host of the staging images registry}
PRIME_REGISTRY_USERNAME=${PRIME_REGISTRY_USERNAME:-}
PRIME_REGISTRY_PASSWORD=${PRIME_REGISTRY_PASSWORD:-}
STAGING_REGISTRY_USERNAME=${STAGING_REGISTRY_USERNAME:-}
STAGING_REGISTRY_PASSWORD=${STAGING_REGISTRY_PASSWORD:-}
DRY_RUN=${DRY_RUN:-true}
PLAIN_HTTP=${PLAIN_HTTP:-false}

case "$DRY_RUN" in true | false) ;; *)
  echo "Error: DRY_RUN must be true or false, not '$DRY_RUN'" >&2
  exit 1
  ;;
esac
case "$PLAIN_HTTP" in
  true) TLS_VERIFY=false PLAIN_HTTP_FLAG=--plain-http ;;
  false) TLS_VERIFY=true PLAIN_HTTP_FLAG= ;;
  *)
    echo "Error: PLAIN_HTTP must be true or false, not '$PLAIN_HTTP'" >&2
    exit 1
    ;;
esac

# The images. Only rancher/<name>:<tag> is taken, which is what cargueiro release can move
# and what helm.sh/images holds: no registry, which comes from the environment, and no digest.
images=()
while IFS= read -r line || [ -n "$line" ]; do
  line=${line//[[:space:]]/}
  [ -z "$line" ] && continue
  if ! [[ "$line" =~ ^rancher/[a-z0-9]+([._-][a-z0-9]+)*:[A-Za-z0-9_][A-Za-z0-9._-]{0,127}$ ]]; then
    echo "Error: $line in $LIST must look like rancher/neuvector-controller:5.7.0, with no registry and no digest" >&2
    exit 1
  fi
  images+=("$line")
done <"$LIST"
if [ "${#images[@]}" -eq 0 ]; then
  echo "Error: there are no images in $LIST" >&2
  exit 1
fi

# skopeo keeps its logins in a file of its own, which is removed at the end.
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT # replaced below, once the summary can be written
export REGISTRY_AUTH_FILE="$WORK/auth.json"

login() { # <host> <username> <password>
  [ -z "$2" ] && return 0
  if ! printf '%s' "$3" | skopeo login --tls-verify="$TLS_VERIFY" --username "$2" --password-stdin "$1" >/dev/null; then
    echo "Error: cannot log in to $1 with the credentials that were given" >&2
    exit 1
  fi
}
login "$PRIME_REGISTRY" "$PRIME_REGISTRY_USERNAME" "$PRIME_REGISTRY_PASSWORD"
login "$STAGING_REGISTRY" "$STAGING_REGISTRY_USERNAME" "$STAGING_REGISTRY_PASSWORD"

# Prints the digest of the manifest that <reference> names: for a multi-platform image, the
# digest of its index. When it cannot be read, prints why and returns 1.
read_digest() { # <reference>
  local error
  if ! error=$(skopeo inspect --tls-verify="$TLS_VERIFY" --retry-times 3 --raw "docker://$1" 2>&1 >"$WORK/manifest"); then
    printf '%s\n' "$error" | tail -n 1
    return 1
  fi
  skopeo manifest-digest "$WORK/manifest"
}

# What became of each image: in staging, missing, pending (to be copied) or copied. The
# summary is written from it at the end, whatever the end is.
digests=()
states=()
pending=()
signatures=
missing=0

summary() {
  if [ -n "${GITHUB_OUTPUT:-}" ]; then
    echo "pending=${#pending[@]}" >>"$GITHUB_OUTPUT"
  fi
  [ -z "${GITHUB_STEP_SUMMARY:-}" ] && return 0
  local i result digest
  {
    echo "### The images, from Prime to staging"
    echo "| Image | Digest in Prime | |"
    echo "|---|---|---|"
    for i in "${!images[@]}"; do
      case "${states[$i]:-}" in
        "in staging") result="in staging" ;;
        missing) result="**missing from Prime**: release it first" ;;
        pending) if [ "$DRY_RUN" = true ]; then result="would be copied"; else result="**not copied**"; fi ;;
        copied) result="copied" ;;
        *) result="not looked at" ;;
      esac
      digest=${digests[$i]:-}
      if [ -n "$digest" ]; then digest="\`${digest:0:19}\`"; else digest="none"; fi
      echo "| \`${images[$i]}\` | $digest | $result |"
    done
    case "$signatures" in
      copied) echo "Signatures and provenance: copied, or in staging already." ;;
      failed) echo "Signatures and provenance: **not copied**, see the log. A new run copies them." ;;
    esac
  } >>"$GITHUB_STEP_SUMMARY"
}
on_exit() {
  local status=$?
  summary || true
  rm -rf "$WORK"
  exit "$status"
}
trap on_exit EXIT

for i in "${!images[@]}"; do
  image=${images[$i]}
  if ! prime=$(read_digest "$PRIME_REGISTRY/$image"); then
    echo "missing     $image is not in $PRIME_REGISTRY: $prime"
    digests[i]=
    states[i]=missing
    missing=$((missing + 1))
    continue
  fi
  digests[i]=$prime
  if staging=$(read_digest "$STAGING_REGISTRY/$image"); then
    if [ "$staging" = "$prime" ]; then
      echo "in staging  $image ($prime)"
      states[i]="in staging"
      continue
    fi
    reason="staging has $staging, Prime has $prime"
  else
    reason="not in staging ($staging)"
  fi
  echo "to copy     $image: $reason"
  states[i]=pending
  pending+=("$i")
done

if [ "$missing" -gt 0 ]; then
  echo "Error: $missing of the images are not in $PRIME_REGISTRY, so nothing was copied: release them first" >&2
  exit 1
fi
if [ "$DRY_RUN" = true ]; then
  echo "Dry run: ${#pending[@]} of the images would be copied from $PRIME_REGISTRY to $STAGING_REGISTRY, and nothing was."
  exit 0
fi

# The copy, by the digest that was read, so that a tag that moves in Prime in the meantime
# cannot bring another image. It is read back from staging.
for i in ${pending[@]+"${pending[@]}"}; do
  image=${images[$i]}
  digest=${digests[$i]}
  echo "::group::copy $image"
  skopeo copy --all --preserve-digests --retry-times 3 \
    --src-tls-verify="$TLS_VERIFY" --dest-tls-verify="$TLS_VERIFY" \
    "docker://$PRIME_REGISTRY/${image%:*}@$digest" "docker://$STAGING_REGISTRY/$image"
  echo "::endgroup::"
  if ! copied=$(read_digest "$STAGING_REGISTRY/$image") || [ "$copied" != "$digest" ]; then
    echo "Error: $image was copied, but $STAGING_REGISTRY gives back ${copied:-nothing} for it, not $digest" >&2
    exit 1
  fi
  states[i]=copied
  echo "copied      $image ($digest)"
done

# The signatures and the provenance, which refer to the image and are not part of it. Now
# that staging has every image with the digest of Prime, they are copied for all of them, not
# only for the ones copied now: cargueiro writes nothing for an image whose signatures
# staging has already, and so a run after one whose copy of the signatures failed copies
# them. cargueiro reads the source with CARGUEIRO_USERNAME and writes the destination with
# CARGUEIRO_PRODUCTION_USERNAME: here the destination is staging.
echo "::group::cargueiro copy-image-signatures"
signatures=failed
# shellcheck disable=SC2086 # PLAIN_HTTP_FLAG is empty or one flag
CARGUEIRO_USERNAME="$PRIME_REGISTRY_USERNAME" CARGUEIRO_PASSWORD="$PRIME_REGISTRY_PASSWORD" \
  CARGUEIRO_PRODUCTION_USERNAME="$STAGING_REGISTRY_USERNAME" CARGUEIRO_PRODUCTION_PASSWORD="$STAGING_REGISTRY_PASSWORD" \
  cargueiro copy-image-signatures "${images[@]}" --from "$PRIME_REGISTRY" --to "$STAGING_REGISTRY" $PLAIN_HTTP_FLAG
signatures=copied
echo "::endgroup::"

echo "Copied ${#pending[@]} of the images from $PRIME_REGISTRY to $STAGING_REGISTRY. Every image is there with its signatures."
