FROM golang:1.9

RUN mkdir -p /app

WORKDIR /app

ADD . /app

RUN curl https://glide.sh/get | sh

RUN glide install
RUN go build

CMD ["./rest"]