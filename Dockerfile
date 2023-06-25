FROM golang AS builder
WORKDIR /traefik-plugin-manual-access-control
COPY go.mod /traefik-plugin-manual-access-control/go.mod
COPY go.sum /traefik-plugin-manual-access-control/go.sum
RUN go mod download
COPY . /traefik-plugin-manual-access-control
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM node as builder-ui
WORKDIR /traefik-plugin-manual-access-control/ui
COPY ui/package.json /traefik-plugin-manual-access-control/ui/package.json
COPY ui/package-lock.json /traefik-plugin-manual-access-control/ui/package-lock.json
RUN npm install
COPY ui /traefik-plugin-manual-access-control/ui
RUN npm run build

FROM alpine
WORKDIR /traefik-plugin-manual-access-control
COPY --from=builder /traefik-plugin-manual-access-control/traefik-plugin-manual-access-control /traefik-plugin-manual-access-control/traefik-plugin-manual-access-control
COPY --from=builder-ui /traefik-plugin-manual-access-control/ui/dist /traefik-plugin-manual-access-control/ui/dist
ENTRYPOINT ["/traefik-plugin-manual-access-control/traefik-plugin-manual-access-control"]
EXPOSE 9502