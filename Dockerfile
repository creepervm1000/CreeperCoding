# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.26-alpine3.23 AS frontend-build
RUN apk add --no-cache build-base git nodejs pnpm
WORKDIR /src

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .
RUN make frontend


FROM golang:1.26-alpine3.23 AS build-env

ARG CREEPERCODING_VERSION
ARG TAGS=""
ENV TAGS="bindata timetzdata $TAGS"

RUN apk add --no-cache build-base git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-build /src/public/assets public/assets

RUN mkdir -p .git

RUN make backend


FROM alpine:3.23

RUN apk add --no-cache \
    ca-certificates \
    git \
    sqlite

# create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=build-env /src/creepercoding /app/creepercoding

RUN mkdir -p /data && chown -R appuser:appgroup /data /app

ENV USER=appuser
ENV GITEA_CUSTOM=/data/creepercoding
ENV PORT=3000

EXPOSE 3000

RUN apk add sudo su-exec

CMD ["su-exec git /app/creepercoding"]
