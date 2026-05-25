package service

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/data_collection"
	"Monitoring-Pressure/models/users"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// GetInterventionRecordsByUserID 获取用户的干预建议记录
type InterventionRecordResp struct {
	ID             uint                                   `json:"id"`
	UserID         int64                                  `json:"user_id"`
	ResultID       uint                                   `json:"result_id"`
	SuggestionID   uint                                   `json:"suggestion_id"`
	Suggestion     data_collection.InterventionSuggestion `json:"suggestion"`
	Result         data_collection.StressAssessmentResult `json:"result"`
	PushTime       time.Time                              `json:"push_time"`
	Feedback       string                                 `json:"feedback"`
	FeedbackStatus int                                    `json:"feedback_status"`
	CreatedAt      time.Time                              `json:"created_at"`
	UpdatedAt      time.Time                              `json:"updated_at"`
}

func toInterventionRecordResp(record data_collection.StudentInterventionRecord) InterventionRecordResp {
	return InterventionRecordResp{
		ID:             record.ID,
		UserID:         record.UserID,
		ResultID:       record.ResultID,
		SuggestionID:   record.SuggestionID,
		Suggestion:     record.Suggestion,
		Result:         record.Result,
		PushTime:       record.PushTime,
		Feedback:       record.Feedback,
		FeedbackStatus: record.FeedbackStatus,
		CreatedAt:      record.CreatedAt,
		UpdatedAt:      record.UpdatedAt,
	}
}

func GetInterventionRecordsByUserID(userID int64) ([]InterventionRecordResp, error) {
	var interventions []data_collection.StudentInterventionRecord

	err := db.DB.
		Preload("Suggestion").
		Preload("Result").
		Where("user_id = ?", userID).
		Order("push_time DESC").
		Find(&interventions).Error
	if err != nil {
		return nil, fmt.Errorf("查询干预建议记录失败: %w", err)
	}

	result := make([]InterventionRecordResp, 0, len(interventions))
	for _, intervention := range interventions {
		result = append(result, toInterventionRecordResp(intervention))
	}
	return result, nil
}

// GenerateInterventionRecordsForResult 根据压力评估结果生成学生干预记录。
func GenerateInterventionRecordsForResult(tx *gorm.DB, result data_collection.StressAssessmentResult) error {
	if tx == nil {
		tx = db.DB
	}
	if result.ID == 0 || result.UserID <= 0 {
		return fmt.Errorf("压力评估结果不完整，无法生成干预记录")
	}
	if result.WarningStatus < 0 || result.WarningStatus > 2 {
		return nil
	}

	matchScore := result.TotalScore
	// 高压预警可能由心理分触发，综合分未必落入高压建议区间。
	if result.WarningStatus == 2 && result.PsychologicalScore > matchScore {
		matchScore = result.PsychologicalScore
	}

	var suggestions []data_collection.InterventionSuggestion
	if err := tx.
		Where("level = ?", result.WarningStatus).
		Where("min_score <= ? AND max_score >= ?", matchScore, matchScore).
		Order("id ASC").
		Find(&suggestions).Error; err != nil {
		return fmt.Errorf("匹配干预建议失败: %w", err)
	}
	if len(suggestions) == 0 {
		return nil
	}

	pushTime := result.CreatedAt
	if result.WarningTime != nil {
		pushTime = *result.WarningTime
	}
	if pushTime.IsZero() {
		pushTime = time.Now()
	}

	for _, suggestion := range suggestions {
		record := data_collection.StudentInterventionRecord{
			UserID:       result.UserID,
			ResultID:     result.ID,
			SuggestionID: suggestion.ID,
		}
		if err := tx.
			Where(&record).
			Attrs(data_collection.StudentInterventionRecord{
				PushTime:       pushTime,
				FeedbackStatus: 0,
			}).
			FirstOrCreate(&record).Error; err != nil {
			return fmt.Errorf("生成学生干预记录失败: %w", err)
		}
	}

	return nil
}

// GetStudentInterventionRecordsAdvanced 获取干预记录（高级版）
func GetStudentInterventionRecordsAdvanced(
	studentID int64,
	page int,
	pageSize int,
	startDate string,
	endDate string,
	feedbackStatusStr string,
) (*PageResult, error) {

	var (
		list  []data_collection.StudentInterventionRecord
		total int64
	)

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		// ========= 1️⃣ 校验用户 =========
		var count int64
		if err := tx.Model(&users.User{}).
			Where("user_id = ?", studentID).
			Count(&count).Error; err != nil {
			return fmt.Errorf("查询用户失败: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("学生不存在")
		}

		// ========= 2️⃣ 构建查询 =========
		query := tx.Model(&data_collection.StudentInterventionRecord{}).
			Where("user_id = ?", studentID)

		// ========= 3️⃣ 时间范围 =========
		if startDate != "" {
			t, err := time.Parse("2006-01-02", startDate)
			if err != nil {
				return fmt.Errorf("start_date 格式错误，应为 yyyy-MM-dd")
			}
			query = query.Where("push_time >= ?", t)
		}

		if endDate != "" {
			t, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				return fmt.Errorf("end_date 格式错误，应为 yyyy-MM-dd")
			}
			query = query.Where("push_time <= ?", t)
		}

		// ========= 4️⃣ 反馈状态 =========
		if feedbackStatusStr != "" {
			status, err := strconv.Atoi(feedbackStatusStr)
			if err != nil || status < 0 {
				return fmt.Errorf("feedback_status 参数错误")
			}
			query = query.Where("feedback_status = ?", status)
		}

		// ========= 5️⃣ count =========
		if err := query.Count(&total).Error; err != nil {
			return fmt.Errorf("统计失败: %w", err)
		}

		// ========= 6️⃣ 分页 + 关联 =========
		offset := (page - 1) * pageSize

		if err := query.
			Preload("Suggestion").
			Preload("Result").
			Order("push_time DESC").
			Limit(pageSize).
			Offset(offset).
			Find(&list).Error; err != nil {
			return fmt.Errorf("查询失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

type CreateInterventionSuggestionReq struct {
	Level    int     `json:"level"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Category string  `json:"category"`
	MinScore float64 `json:"min_score"`
	MaxScore float64 `json:"max_score"`
}

// CreateInterventionSuggestion 创建干预建议
func CreateInterventionSuggestion(req CreateInterventionSuggestionReq) error {

	return db.DB.Transaction(func(tx *gorm.DB) error {

		// ========= 1️⃣ 防重复（标题）=========
		var count int64
		if err := tx.Model(&data_collection.InterventionSuggestion{}).
			Where("title = ?", req.Title).
			Count(&count).Error; err != nil {
			return fmt.Errorf("检查重复失败: %w", err)
		}

		if count > 0 {
			return fmt.Errorf("该标题已存在")
		}

		// ========= 2️⃣ 区间冲突校验（关键 ⭐⭐⭐）=========
		var overlap int64
		if err := tx.Model(&data_collection.InterventionSuggestion{}).
			Where("level = ?", req.Level).
			Where("NOT (max_score < ? OR min_score > ?)", req.MinScore, req.MaxScore).
			Count(&overlap).Error; err != nil {
			return fmt.Errorf("区间检查失败: %w", err)
		}

		if overlap > 0 {
			return fmt.Errorf("存在重叠的分数区间，请检查")
		}

		// ========= 3️⃣ 构造数据 =========
		suggestion := data_collection.InterventionSuggestion{
			Level:    req.Level,
			Title:    req.Title,
			Content:  req.Content,
			Category: req.Category,
			MinScore: req.MinScore,
			MaxScore: req.MaxScore,
		}

		// ========= 4️⃣ 入库 =========
		if err := tx.Create(&suggestion).Error; err != nil {
			return fmt.Errorf("创建失败: %w", err)
		}

		return nil
	})
}

type InterventionQuery struct {
	Level    *int
	Category string
	MinScore *float64
	MaxScore *float64
	Page     int
	PageSize int
}

func GetInterventionSuggestionList(queryParam InterventionQuery) (*PageResult, error) {

	var (
		list  []data_collection.InterventionSuggestion
		total int64
	)

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		query := tx.Model(&data_collection.InterventionSuggestion{})

		// ========= 1️⃣ 条件拼接 =========

		if queryParam.Level != nil {
			query = query.Where("level = ?", *queryParam.Level)
		}

		if queryParam.Category != "" {
			query = query.Where("category = ?", queryParam.Category)
		}

		// 分数区间筛选（重点 ⭐⭐⭐）
		if queryParam.MinScore != nil {
			query = query.Where("min_score >= ?", *queryParam.MinScore)
		}

		if queryParam.MaxScore != nil {
			query = query.Where("max_score <= ?", *queryParam.MaxScore)
		}

		// ========= 2️⃣ count =========
		if err := query.Count(&total).Error; err != nil {
			return fmt.Errorf("统计失败: %w", err)
		}

		// ========= 3️⃣ 分页 =========
		offset := (queryParam.Page - 1) * queryParam.PageSize

		if err := query.
			Order("created_at DESC").
			Limit(queryParam.PageSize).
			Offset(offset).
			Find(&list).Error; err != nil {
			return fmt.Errorf("查询失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &PageResult{
		List:     list,
		Total:    total,
		Page:     queryParam.Page,
		PageSize: queryParam.PageSize,
	}, nil
}
