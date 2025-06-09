# FROM golang:1.23.3
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=admin
ENV DB_NAME=besidemom
ENV SSL_MODE=disable
ENV APP_HOST=0.0.0.0
ENV APP_PORT=5000
ENV BUCKET_NAME=Beside-Mom
ENV EMAIL_HOST=smtp.gmail.com
ENV EMAIL_PORT=587
ENV EMAIL_USER=kasianbot66@gmail.com
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
CMD ["./main"]