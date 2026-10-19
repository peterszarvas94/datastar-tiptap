FROM node:24-bookworm-slim AS assets
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npx tailwindcss -i css/input.css -o static/generated.css \
    && npm run bundle:js-beautify \
    && npm run bundle:tiptap

FROM golang:1.26.9-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY src/ src/
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/app ./src

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/app ./app
COPY --from=assets /src/static/ ./static/
COPY templates/ ./templates/
USER 10001:10001
EXPOSE 3000
CMD ["./app"]
