FROM golang:13
WORKDIR avatargen/
COPY go.* ./
RUN go get
COPY . ./
RUN go build -o out
CMD ["./out"]
