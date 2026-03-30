package service

import (
	"Monitoring-Pressure/dao/db"
	"strings"

	DataCollection "Monitoring-Pressure/models/data_collection"
)

type AcademicImportRecordListReq struct {
	Page       int
	PageSize   int
	Status     *int
	FileName   string
	OperatorID *int64
}

type AcademicImportRecordListItem struct {
	ID           uint   `json:"id"`
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	ImportCount  int    `json:"import_count"`
	OperatorID   int64  `json:"operator_id"`
	Status       int    `json:"status"`
	ErrorMessage string `json:"error_message"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type AcademicImportRecordListResp struct {
	List     []AcademicImportRecordListItem `json:"list"`
	Total    int64                          `json:"total"`
	Page     int                            `json:"page"`
	PageSize int                            `json:"page_size"`
}

func GetAcademicImportRecordList(req AcademicImportRecordListReq) (*AcademicImportRecordListResp, error) {
	// 分页参数兜底
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	db := db.DB.Model(&DataCollection.AcademicImportRecord{})

	// 条件筛选
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if strings.TrimSpace(req.FileName) != "" {
		db = db.Where("file_name LIKE ?", "%"+strings.TrimSpace(req.FileName)+"%")
	}

	if req.OperatorID != nil && *req.OperatorID > 0 {
		db = db.Where("operator_id = ?", *req.OperatorID)
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 查询数据
	var records []DataCollection.AcademicImportRecord
	offset := (req.Page - 1) * req.PageSize

	if err := db.Order("id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&records).Error; err != nil {
		return nil, err
	}

	// 转返回结构
	list := make([]AcademicImportRecordListItem, 0, len(records))
	for _, item := range records {
		list = append(list, AcademicImportRecordListItem{
			ID:           item.ID,
			FileName:     item.FileName,
			FilePath:     item.FilePath,
			ImportCount:  item.ImportCount,
			OperatorID:   item.OperatorID,
			Status:       item.Status,
			ErrorMessage: item.ErrorMessage,
			CreatedAt:    item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &AcademicImportRecordListResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
