# 使用最新的 Go 鏡像作為構建階段
FROM golang:1.23-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 複製 go mod 和 go sum 文件
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製整個項目
COPY . .

# 編譯應用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 第二階段：使用輕量級 Alpine 鏡像 上個階段比較像是商品製作都含在 第二階段則是把用不到的東西刪掉 精美的給使用者
FROM alpine:latest

# 設置工作目錄
WORKDIR /root/

# 複製編譯好的二進制文件
COPY --from=builder /app/main .

# 複製配置文件
COPY --from=builder /app/config ./config

# 複製視圖模板
COPY --from=builder /app/views ./views

# 複製靜態資源
COPY --from=builder /app/asset ./asset

# 複製 HTML 文件
COPY --from=builder /app/index.html .

# 暴露端口
EXPOSE 3000

# 運行應用
CMD ["./main"]