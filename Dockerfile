FROM alpine:3.21
RUN apk add --no-cache ca-certificates && \
    adduser -D -h /home/appuser appuser
COPY zotero-mcp-go-server /usr/local/bin/zotero-mcp-go-server
USER appuser
ENTRYPOINT ["zotero-mcp-go-server"]
