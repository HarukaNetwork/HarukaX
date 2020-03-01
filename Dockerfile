FROM golang

WORKDIR /go/src/github.com/HarukaNetwork/HarukaX

COPY . .

ENTRYPOINT ["go", "run", "."]
