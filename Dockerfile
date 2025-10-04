# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:3.20

LABEL maintainer="NoteSpark AI"
LABEL description="Lightweight Pandoc + XeLaTeX + Noto multilingual font support"

# Install basic packages and Pandoc
RUN apk add --no-cache \
    pandoc \
    texlive-xetex \
    texmf-dist-latexextra \
    texmf-dist-fontsextra \
    texmf-dist-langcjk \
    texmf-dist-langother \
    texmf-dist-latexrecommended \
    ghostscript \
    fontconfig \
    ttf-freefont \
    wget unzip

# Install Noto fonts (multilingual support)
RUN mkdir -p /usr/share/fonts/noto && cd /usr/share/fonts/noto && \
    wget -q https://noto-website-2.storage.googleapis.com/pkgs/Noto-unhinted.zip && \
    unzip -qq Noto-unhinted.zip && \
    fc-cache -fv && rm Noto-unhinted.zip

# Set XeLaTeX as default PDF engine environment (optional)
ENV PATH="/usr/bin:$PATH"
ENV LANG=C.UTF-8

COPY --from=builder /app/main /app/main
WORKDIR /app
CMD ["./main"]