package admin

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"Monitoring-Pressure/service"
)

func ImportAcademicExcelHandle(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请上传Excel文件",
		})
		return
	}

	operatorIDStr := c.PostForm("operator_id")
	operatorID, err := strconv.ParseInt(operatorIDStr, 10, 64)
	if err != nil || operatorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "operator_id 不合法",
		})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".xlsx" && ext != ".xls" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "只支持 .xlsx 或 .xls 文件",
		})
		return
	}

	saveDir := "static/uploads/academic"
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建上传目录失败",
		})
		return
	}

	saveName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	savePath := filepath.Join(saveDir, saveName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "保存上传文件失败",
		})
		return
	}

	err = service.ImportAcademicExcel(savePath, file.Filename, operatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "导入失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成绩导入成功",
	})
}
