# ── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

# Dependências de build
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copia go.mod primeiro para aproveitar cache de dependências
COPY go.mod ./
# go.sum pode não existir se não houver dependências externas
COPY go.sum* ./

RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário — CGO desativado para imagem mínima
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/crypto-desert ./cmd/server

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
FROM alpine:3.19

# Certificados para chamadas HTTPS (CoinGecko)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copia o binário compilado
COPY --from=builder /app/crypto-desert .

# Copia o frontend estático
COPY --from=builder /app/web ./web

# Variáveis de ambiente com defaults
ENV PORT=8080
ENV WEB_DIR=./web

EXPOSE 8080

# Usuário não-root por segurança
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["./crypto-desert"]
