FROM golang:1.18-bullseye

RUN mkdir /pub

WORKDIR /pub

COPY . .

RUN go build -o publisher .

EXPOSE 8080

CMD [ "/pub/publisher" ]