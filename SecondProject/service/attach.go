package service

import (
	"fmt"
	"github/qq900306ss/SecondProject/utils"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) { //上傳檔案
	w := c.Writer
	req := c.Request
	srcFile, head, err := req.FormFile("file") // src可以得知來源檔案，head可以得知檔案名稱
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	suffix := ".png"         //後綴名稱
	ofiName := head.Filename //名稱
	tem := strings.Split(ofiName, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1] //取得後綴名稱
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFile, err := os.Create("./asset/upload/" + fileName) //創造檔案並放那
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	_, err = io.Copy(dstFile, srcFile) //dest 目的地 , src來源 這裡就上傳完了
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	url := "./asset/upload/" + fileName //傳給他的檔案是哪個
	utils.RespOK(w, url, "上傳圖片成功")

}
