# Use the official Golang image as the base image for building
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./server

FROM golang:1.21.1-alpine

# Install required packages
RUN apk --no-cache add python3 python3-dev ca-certificates postgresql-client
RUN apk --no-cache add py3-scikit-learn py3-pandas py3-joblib
RUN apk --no-cache add py3-numpy

WORKDIR /root/

COPY --from=builder /app/backend .
COPY --from=builder /app/web ./web
COPY --from=builder /app/python/model_class.py .
COPY --from=builder /app/python/model_reg.py .
# RUN cp /root/.venv/bin/python /usr/local/bin/python3

#RUN chmod -R 777 ./downloads
#RUN chmod -R 777 ./uploads

EXPOSE 8080

CMD ["./backend"]
