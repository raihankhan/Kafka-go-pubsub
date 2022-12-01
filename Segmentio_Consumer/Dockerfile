FROM golang:1.18-bullseye

RUN mkdir /sub

WORKDIR /sub

COPY . .

RUN go build -o consumer .

EXPOSE 8080

CMD [ "/sub/consumer" ]