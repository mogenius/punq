docker buildx create \
  --name amd64_builder \
  --node linux_amd64_builder \
  --platform linux/amd64 \
  ssh://build@v2202112163193172205.supersrv.de