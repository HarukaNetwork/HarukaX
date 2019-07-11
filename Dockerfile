FROM golang

WORKDIR /go/src/github.com/ATechnoHazard/ginko

COPY . .

ENTRYPOINT ["go", "run", "."]
