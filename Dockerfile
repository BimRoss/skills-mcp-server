FROM golang:1.24.2-alpine AS builder
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/skills-mcp-server ./cmd/skills-mcp-server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/skills-mcp-server /usr/local/bin/skills-mcp-server
COPY examples /app/examples

ENV SKILLS_MCP_SERVER_PORT=8081
ENV SKILLS_MCP_SERVER_DIR=/app/skills
ENV SKILLS_EXAMPLES_DIR=/app/examples

RUN mkdir -p /app/skills
EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/skills-mcp-server"]
