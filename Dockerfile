FROM golang:1.21

WORKDIR /app

# Go 환경 설정
ENV GO111MODULE=on
ENV CGO_ENABLED=0

# 소스 코드 복사
COPY . .

# 의존성 설치 및 빌드
RUN go mod tidy && go build -o main .

EXPOSE 8080

# 서버 실행
CMD ["./main"] 