FROM golang:1.23-alpine AS builder 

# 設定容器內的工作目錄為 /app
WORKDIR /app

# 複製 go.mod 和 go.sum 檔案到容器的當前工作目錄
COPY go.mod go.sum ./

# 下載所有專案依賴套件
RUN go mod download

# 將專案的所有原始碼複製到容器的工作目錄
COPY . .

# 複製配置文件和視圖模板
COPY config /root/config
COPY views /app/views

EXPOSE 8000

# 編譯 GO 程式
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 執行階段映像
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /root/config /root/config
COPY --from=builder /app/views /app/views
EXPOSE 8000
CMD ["./main"]