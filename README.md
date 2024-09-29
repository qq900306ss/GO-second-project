"# second-project" 
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
