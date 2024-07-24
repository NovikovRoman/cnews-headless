FROM golang:1.22-alpine3.20 as build
LABEL stage=builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.* ./
RUN apk --no-cache add make && go mod download

COPY . .
RUN make

FROM chromedp/headless-shell:128.0.6601.2
RUN apt-get update; apt install ca-certificates tzdata dumb-init -y

ENV TZ="Europe/Moscow"

ENTRYPOINT ["dumb-init", "--"]

WORKDIR /app
COPY --from=build /app/bin .
EXPOSE 4444
CMD [ "./app"]
