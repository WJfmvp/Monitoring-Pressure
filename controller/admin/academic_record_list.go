package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Monitoring-Pressure/service"
)

func GetAcademicRecordListHandle(c *gin.Context) {
	// 1. 分页参数
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

	// 2. 可选筛选参数
	var userIDPtr *int64
	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil || userID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "user_id 参数不合法",
			})
			return
		}
		userIDPtr = &userID
	}

	var importIDPtr *uint
	importIDStr := c.Query("import_id")
	if importIDStr != "" {
		importID64, err := strconv.ParseUint(importIDStr, 10, 64)
		if err != nil || importID64 == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "import_id 参数不合法",
			})
			return
		}
		importID := uint(importID64)
		importIDPtr = &importID
	}

	var sourceTypePtr *int
	sourceTypeStr := c.Query("source_type")
	if sourceTypeStr != "" {
		sourceType, err := strconv.Atoi(sourceTypeStr)
		if err != nil || (sourceType != 1 && sourceType != 2) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "source_type 参数不合法，只能是 1 或 2",
			})
			return
		}
		sourceTypePtr = &sourceType
	}

	var minScorePtr *float64
	minScoreStr := c.Query("min_score")
	if minScoreStr != "" {
		minScore, err := strconv.ParseFloat(minScoreStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "min_score 参数不合法",
			})
			return
		}
		minScorePtr = &minScore
	}

	var maxScorePtr *float64
	maxScoreStr := c.Query("max_score")
	if maxScoreStr != "" {
		maxScore, err := strconv.ParseFloat(maxScoreStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "max_score 参数不合法",
			})
			return
		}
		maxScorePtr = &maxScore
	}

	if minScorePtr != nil && maxScorePtr != nil && *minScorePtr > *maxScorePtr {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "min_score 不能大于 max_score",
		})
		return
	}

	examName := c.Query("exam_name")
	term := c.Query("term")

	// 3. 调用 service
	resp, err := service.GetAcademicRecordList(service.AcademicRecordListReq{
		Page:       page,
		PageSize:   pageSize,
		UserID:     userIDPtr,
		ImportID:   importIDPtr,
		ExamName:   examName,
		Term:       term,
		SourceType: sourceTypePtr,
		MinScore:   minScorePtr,
		MaxScore:   maxScorePtr,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询成绩列表失败: " + err.Error(),
		})
		return
	}

	// 4. 返回数据
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": resp,
	})
}
