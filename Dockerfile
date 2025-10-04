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

# Install Pandoc and minimal dependencies
RUN apk add --no-cache \
    pandoc \
    perl \
    wget \
    fontconfig \
    freetype

# Install TinyTeX (lightweight LaTeX distribution ~200MB vs 2GB+ for full TeX Live)
RUN wget -qO- "https://yihui.org/tinytex/install-bin-unix.sh" | sh && \
    /root/.TinyTeX/bin/*/tlmgr install \
    xetex \
    fontspec \
    unicode-math \
    xecjk \
    geometry \
    && /root/.TinyTeX/bin/*/tlmgr path add

ENV PATH="/root/.TinyTeX/bin/x86_64-linuxmusl:$PATH"

# Install minimal Noto fonts (only what you need)
RUN mkdir -p /usr/share/fonts/noto && cd /usr/share/fonts/noto && \
    wget -q https://github.com/notofonts/noto-cjk/releases/download/Sans2.004/08_NotoSansCJKjp.zip && \
    wget -q https://github.com/notofonts/noto-fonts/releases/download/NotoSans-v2.013/NotoSans-v2.013.zip && \
    wget -q https://github.com/notofonts/arabic/releases/download/NotoSansArabic-v2.010/NotoSansArabic-v2.010.zip && \
    wget -q https://github.com/notofonts/devanagari/releases/download/NotoSansDevanagari-v2.004/NotoSansDevanagari-v2.004.zip && \
    wget -q https://github.com/notofonts/thai/releases/download/NotoSansThai-v2.002/NotoSansThai-v2.002.zip && \
    unzip -qq "*.zip" && \
    fc-cache -fv && \
    rm -rf *.zip && \
    apk del wget

# Set XeLaTeX as default PDF engine environment (optional)
ENV PATH="/usr/bin:$PATH"
ENV LANG=C.UTF-8

COPY --from=builder /app/main /app/main
WORKDIR /app
CMD ["./main"]