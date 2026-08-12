FROM golang:1.26
WORKDIR /app


RUN useradd --create-home app 
USER app 

COPY main.go go.mod ./

CMD ["go", "run", "main.go"]
