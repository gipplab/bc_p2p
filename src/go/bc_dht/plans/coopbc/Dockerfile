ARG build_image
ARG run_image

FROM ${build_image} AS builder
WORKDIR /build
ENV CGO_ENABLED 0
COPY . .
RUN cd plan && go build -a -o /testplan

FROM ${run_image}
COPY --from=builder /testplan /testplan
COPY plan/test_doc.csv /
EXPOSE 80 443 6060
ENTRYPOINT [ "/testplan"]
