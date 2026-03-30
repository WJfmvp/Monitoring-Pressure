package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Monitoring-Pressure/service"
)

func GetAcademicImportRecordListHandle(c *gin.Context) {
	// 1. 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "page 参数不合法",
		})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "page_size 参数不合法",
		})
		return
	}

	if pageSize > 100 {
		pageSize = 100
	}

	// 2. 获取筛选参数
	var statusPtr *int
	statusStr := c.Query("status")
	if statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err != nil || (status != 0 && status != 1) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "status 参数不合法，只能是 0 或 1",
			})
			return
		}
		statusPtr = &status
	}

	var operatorIDPtr *int64
	operatorIDStr := c.Query("operator_id")
	if operatorIDStr != "" {
		operatorID, err := strconv.ParseInt(operatorIDStr, 10, 64)
		if err != nil || operatorID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "operator_id 参数不合法",
			})
			return
		}
		operatorIDPtr = &operatorID
	}

	fileName := c.Query("file_name")

	// 3. 调 service
	resp, err := service.GetAcademicImportRecordList(service.AcademicImportRecordListReq{
		Page:       page,
		PageSize:   pageSize,
		Status:     statusPtr,
		FileName:   fileName,
		OperatorID: operatorIDPtr,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询导入记录失败: " + err.Error(),
		})
		return
	}

	// 4. 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": resp,
	})
}
