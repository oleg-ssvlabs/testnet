FROM ghcr.io/foundry-rs/foundry:latest

USER root
WORKDIR /app

# Install Git and other dependencies
RUN apt-get update && apt-get install -y git curl bash jq

# Clone with --depth 1 and force pull to get latest
RUN git clone --branch foundry --depth 1 https://github.com/Zacholme7/ssv-network . 
RUN cd . && git fetch --all && git reset --hard origin/foundry

# Install dependencies via Foundry
RUN forge install

# Default command - keep container alive
CMD ["tail", "-f", "/dev/null"]
