FROM golang:1.21-alpine AS development

WORKDIR /app

# 필요한 빌드 도구들 설치
RUN apk add --no-cache git gcc musl-dev

# air 설치 및 PATH 설정
RUN go install github.com/cosmtrek/air@v1.49.0
ENV PATH="/go/bin:${PATH}"

# 소스 코드 복사
COPY go.* ./

# 의존성 설치
RUN go get github.com/gin-gonic/gin@v1.9.1 && \
    go get github.com/go-sql-driver/mysql@v1.8.0 && \
    go mod tidy && \
    go mod download && \
    go mod verify

COPY . .

# 개발 환경에서는 air를 사용하여 자동 리로드
CMD ["air", "-c", ".air.toml"]

FROM golang:1.21-alpine AS production

WORKDIR /app

# 소스 코드 복사
COPY . .

# 빌드
RUN go get github.com/gin-gonic/gin@v1.9.1 && \
    go get github.com/go-sql-driver/mysql@v1.8.0 && \
    go mod tidy && \
    go mod download && \
    go mod verify && \
    go build -o main .

# 실행
CMD ["./main"] 