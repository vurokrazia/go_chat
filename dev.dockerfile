FROM golang:1.15.2-buster
WORKDIR /app
RUN go mod init go_chat
COPY . .
#ENTRYPOINT ["reflex”, "-c", "reflex.conf"]
CMD bash