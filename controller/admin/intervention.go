package admin

import (
	"Monitoring-Pressure/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// CreateInterventionSuggestionHandle 创建干预建议
func CreateInterventionSuggestionHandle(c *gin.Context) {

	var req service.CreateInterventionSuggestionReq

	// ========= 1️⃣ 绑定参数 =========
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	// ========= 2️⃣ 参数校验 =========

	// 标题
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title 不能为空"})
		return
	}

	// 内容
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content 不能为空"})
		return
	}

	// 等级
	if req.Level < 0 || req.Level > 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "level 只能是 0/1/2"})
		return
	}

	// 分类校验（可选）
	if req.Category != "" {
		switch req.Category {
		case "学习", "心理", "生活":
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "category 必须是 学习/心理/生活"})
			return
		}
	}

	// 分数校验
	if req.MinScore < 0 || req.MaxScore < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分数不能为负"})
		return
	}

	if req.MinScore > req.MaxScore {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_score 不能大于 max_score"})
		return
	}

	// ========= 3️⃣ 调用 service =========
	err := service.CreateInterventionSuggestion(req)
	if err != nil {
		log.Printf("创建干预建议失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ========= 4️⃣ 返回 =========
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
	})
}

func GetInterventionSuggestionListHandle(c *gin.Context) {

	// ========= 1️⃣ 分页 =========
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// ========= 2️⃣ 筛选 =========
	var level *int
	if l := c.Query("level"); l != "" {
		v, err := strconv.Atoi(l)
		if err != nil || v < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "level 参数错误"})
			return
		}
		level = &v
	}

	category := c.Query("category")

	var minScore *float64
	if ms := c.Query("min_score"); ms != "" {
		v, err := strconv.ParseFloat(ms, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "min_score 参数错误"})
			return
		}
		minScore = &v
	}

	var maxScore *float64
	if ms := c.Query("max_score"); ms != "" {
		v, err := strconv.ParseFloat(ms, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "max_score 参数错误"})
			return
		}
		maxScore = &v
	}

	// ========= 3️⃣ 调用 service =========
	result, err := service.GetInterventionSuggestionList(service.InterventionQuery{
		Level:    level,
		Category: category,
		MinScore: minScore,
		MaxScore: maxScore,
		Page:     page,
		PageSize: pageSize,
	})

	if err != nil {
		log.Printf("获取建议列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取失败",
		})
		return
	}

	// ========= 4️⃣ 返回 =========
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    result,
	})
}
