FROM golang:1.19


WORKDIR /app 

COPY . . 

#RUN go get 

RUN go build main.go

EXPOSE 8000 

CMD ["./main"]

