FROM golang:1.18-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY . .
RUN go build -o youtubeapi
EXPOSE 5001

CMD [ "./youtubeapi" ]
