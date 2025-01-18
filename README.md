這是一個使用GO 語言寫的類似LINE的功能有註冊和登入database ，還有即時通訊採用TCP跟Redis緩存

https://youtu.be/dz6nbNogw5Q

展示可以觀看上面YT 

前端部分都是找有人做好的模板和使用GPT協助完成，再去做參數調整

成果展示
![image](https://github.com/user-attachments/assets/ee4fdea8-7bd4-4b6d-924b-0f43a69e993b)
![image](https://github.com/user-attachments/assets/dc92e2b5-b284-438d-8567-42c8b9bb9050)
![image](https://github.com/user-attachments/assets/039a4500-5685-497f-8fae-b16a6525c97a)
這是跟好友對話的樣子，採用Goroutine一個負責送一個負責收和使用去redis讀取資料達成高效率，順便減輕database壓力

![image](https://github.com/user-attachments/assets/274e3409-9d54-4412-a4d5-47bd7faa6cb5)
上面是群組聊天樣子，一樣採用redis，並且都含有及時顯示功能

https://youtu.be/wTzc4GjK4Mo
群組聊天室未更新時候


# GO Second Project - 即時通訊與群組聊天應用

## 專案介紹
這是一個使用 Go 語言開發的即時通訊與群組聊天應用，提供實時消息傳遞和用戶互動功能。

## 技術棧
- 後端: Go (Golang)
- 數據庫: MySQL
- 快取: Redis
- 容器化: Docker
- 通訊協議: WebSocket

## 主要功能
- 用戶註冊與認證
- 即時消息傳遞
- 群組聊天
- 實時通知
- 用戶在線狀態

## 系統需求
- Go 1.23
- Docker
- Docker Compose

## 快速開始

### 克隆倉庫
```bash

git clone https://github.com/your-username/GO-second-project.git

cd GO-second-project
cd SecondProject
docker-compose up --build -d

```

暫停
```
docker-compose stop
``` 