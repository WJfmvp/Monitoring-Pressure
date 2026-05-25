package admin

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"Monitoring-Pressure/service"
	"Monitoring-Pressure/util"
)

func ImportAcademicExcelHandle(c *gin.Context) {
	// operator_id 必须从 JWT 取，避免被前端伪造
	operatorIDValue, exists := c.Get("user_id")
	if !exists {
		util.ResponseError(c, util.ErrCodeNeedLogin)
		return
	}
	operatorID, ok := operatorIDValue.(int64)
	if !ok || operatorID <= 0 {
		util.ResponseError(c, util.ErrCodeNeedLogin)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请上传Excel文件",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
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

	if err := service.ImportAcademicExcel(savePath, file.Filename, operatorID); err != nil {
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
