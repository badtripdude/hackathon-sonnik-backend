DOCKER_BUILDKIT=1 docker build \
  --network=host \
  -t ai \
  -f docker/Dockerfile \
  .

DOCKER_BUILDKIT=1 docker build \
  --network=host \
  -t ai_migrate \
  -f docker/Dockerfile \
  .

