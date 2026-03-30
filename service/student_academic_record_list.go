package service

import (
	"Monitoring-Pressure/dao/db"
	"strings"

	DataCollection "Monitoring-Pressure/models/data_collection"
)

type StudentAcademicRecordListReq struct {
	UserID   int64
	Page     int
	PageSize int
	ExamName string
	Term     string
	MinScore *float64
	MaxScore *float64
}

type StudentAcademicRecordListItem struct {
	ID               uint    `json:"id"`
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

type StudentAcademicRecordListResp struct {
	List     []StudentAcademicRecordListItem `json:"list"`
	Total    int64                           `json:"total"`
	Page     int                             `json:"page"`
	PageSize int                             `json:"page_size"`
}

func GetStudentAcademicRecordList(req StudentAcademicRecordListReq) (*StudentAcademicRecordListResp, error) {
	// 1. 基础校验
	if req.UserID <= 0 {
		return &StudentAcademicRecordListResp{
			List:     []StudentAcademicRecordListItem{},
			Total:    0,
			Page:     1,
			PageSize: 10,
		}, nil
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	db := db.DB.Model(&DataCollection.AcademicRecord{}).Where("user_id = ?", req.UserID)

	// 2. 可选筛选条件
	if strings.TrimSpace(req.ExamName) != "" {
		db = db.Where("exam_name LIKE ?", "%"+strings.TrimSpace(req.ExamName)+"%")
	}

	if strings.TrimSpace(req.Term) != "" {
		db = db.Where("term = ?", strings.TrimSpace(req.Term))
	}

	if req.MinScore != nil {
		db = db.Where("exam_score >= ?", *req.MinScore)
	}

	if req.MaxScore != nil {
		db = db.Where("exam_score <= ?", *req.MaxScore)
	}

	// 3. 总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 4. 分页查询
	var records []DataCollection.AcademicRecord
	offset := (req.Page - 1) * req.PageSize

	if err := db.Order("exam_date DESC, id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&records).Error; err != nil {
		return nil, err
	}

	// 5. 组装返回
	list := make([]StudentAcademicRecordListItem, 0, len(records))
	for _, item := range records {
		examDate := ""
		if !item.ExamDate.IsZero() {
			examDate = item.ExamDate.Format("2006-01-02")
		}

		list = append(list, StudentAcademicRecordListItem{
			ID:               item.ID,
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

	return &StudentAcademicRecordListResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
