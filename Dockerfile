FROM alpine:3.23
RUN apk add --no-cache ca-certificates && \
    adduser -D -h /home/appuser appuser
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/zotero-mcp-go-server /usr/local/bin/zotero-mcp-go-server
USER appuser
ENTRYPOINT ["zotero-mcp-go-server"]
