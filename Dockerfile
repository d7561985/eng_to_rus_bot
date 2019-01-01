FROM golang:1.11.4-alpine3.8 as build

COPY . src/github.com/d7561985/eng_to_rus_bot/
WORKDIR src/github.com/d7561985/eng_to_rus_bot/
RUN go build -o /bot cmd/main.go

FROM alpine:3.8

COPY --from=build /bot /bot
COPY --from=build /go/src/github.com/d7561985/eng_to_rus_bot/assets/ /assets/

# init certificates for https connection
RUN apk add --no-cache libstdc++ \
	ca-certificates

# add heroku user
RUN adduser -D -u 1000 heroku
USER heroku

CMD /bot