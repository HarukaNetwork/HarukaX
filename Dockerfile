FROM golang

WORKDIR /go/src/github.com/atechnohazard/ginko

COPY . .

ENTRYPOINT ["go", "run", "."]