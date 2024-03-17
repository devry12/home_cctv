# Gunakan image resmi dari Golang sebagai base image
FROM golang:latest

# Install FFmpeg
RUN apt-get update && apt-get install -y ffmpeg

# Kerjaan direktori default di dalam kontainer
WORKDIR /go/src/app

# Salin kode Golang Anda ke dalam kontainer
COPY . .

# Kompilasi kode Golang
RUN go build -o app .


# Jalankan aplikasi ketika kontainer dimulai
CMD ["./app"]
