package service

import (
	"Monitoring-Pressure/dao/db"
	"strings"

	DataCollection "Monitoring-Pressure/models/data_collection"
)

type AcademicRecordListReq struct {
	Page       int
	PageSize   int
	UserID     *int64
	ImportID   *uint
	ExamName   string
	Term       string
	SourceType *int
	MinScore   *float64
	MaxScore   *float64
}

type AcademicRecordListItem struct {
	ID               uint    `json:"id"`
	UserID           int64   `json:"user_id"`
	ImportID         uint    `json:"import_id"`
	ExamName         string  `json:"exam_name"`
	Term             string  `json:"term"`
	ExamDate         string  `json:"exam_date"`
	RawScore         float64 `json:"raw_score"`
	ExamScore        float64 `json:"exam_score"`
	AverageScore     float64 `json:"average_score"`
	ScoreFluctuation float64 `json:"score_fluctuation"`
	ClassRank        int     `json:"class_rank"`
	GradeRank        int     `json:"grade_rank"`
	FailCount        int     `json:"fail_count"`
	ScoreChangeRate  float64 `json:"score_change_rate"`
	SourceType       int     `json:"source_type"`
	Remark           string  `json:"remark"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type AcademicRecordListResp struct {
	List     []AcademicRecordListItem `json:"list"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
}

func GetAcademicRecordList(req AcademicRecordListReq) (*AcademicRecordListResp, error) {
	// 分页兜底
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	db := db.DB.Model(&DataCollection.AcademicRecord{})

	// 条件筛选
	if req.UserID != nil && *req.UserID > 0 {
		db = db.Where("user_id = ?", *req.UserID)
	}

	if req.ImportID != nil && *req.ImportID > 0 {
		db = db.Where("import_id = ?", *req.ImportID)
	}

	if strings.TrimSpace(req.ExamName) != "" {
		db = db.Where("exam_name LIKE ?", "%"+strings.TrimSpace(req.ExamName)+"%")
	}

	if strings.TrimSpace(req.Term) != "" {
		db = db.Where("term = ?", strings.TrimSpace(req.Term))
	}

	if req.SourceType != nil {
		db = db.Where("source_type = ?", *req.SourceType)
	}

	if req.MinScore != nil {
		db = db.Where("exam_score >= ?", *req.MinScore)
	}

	if req.MaxScore != nil {
		db = db.Where("exam_score <= ?", *req.MaxScore)
	}

	// 统计总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 查询分页数据
	var records []DataCollection.AcademicRecord
	offset := (req.Page - 1) * req.PageSize

	if err := db.Order("exam_date DESC, id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&records).Error; err != nil {
		return nil, err
	}

	// 组装返回数据
	list := make([]AcademicRecordListItem, 0, len(records))
	for _, item := range records {
		examDate := ""
		if !item.ExamDate.IsZero() {
			examDate = item.ExamDate.Format("2006-01-02")
		}

		list = append(list, AcademicRecordListItem{
			ID:               item.ID,
			UserID:           item.UserID,
			ImportID:         item.ImportID,
			ExamName:         item.ExamName,
			Term:             item.Term,
			ExamDate:         examDate,
			RawScore:         item.RawScore,
			ExamScore:        item.ExamScore,
			AverageScore:     item.AverageScore,
			ScoreFluctuation: item.ScoreFluctuation,
			ClassRank:        item.ClassRank,
			GradeRank:        item.GradeRank,
			FailCount:        item.FailCount,
			ScoreChangeRate:  item.ScoreChangeRate,
			SourceType:       item.SourceType,
			Remark:           item.Remark,
			CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:        item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &AcademicRecordListResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
