FROM rust AS builder

COPY src/ src/
COPY Cargo.toml .
COPY Cargo.lock .

RUN cargo build --release 

# CMD ["/target/release/bc_p2p"]

FROM debian:buster-slim
RUN apt-get update && apt-get install -y libssl-dev
COPY --from=builder /target/release/bc_p2p /bin/bc_p2p

CMD ["bc_p2p"]

# run detached with: 
# docker run -it 7811 /bin/bash

# or

# docker run -it 7811 /target/release/bc_p2p

# build with:
# docker build . -t ihlec/p2p