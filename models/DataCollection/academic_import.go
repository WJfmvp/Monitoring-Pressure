package DataCollection

import "time"

// AcademicImportRecord 记录每次 Excel 导入
// 来源：管理员上传 Excel 后生成一条记录
type AcademicImportRecord struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	FileName string `json:"file_name" gorm:"type:varchar(255);not null"` // 文件名
	FilePath string `json:"file_path" gorm:"type:varchar(500)"`          // 文件存储路径

	ImportCount int `json:"import_count" gorm:"default:0"` // 导入条数

	OperatorID int64 `json:"operator_id" gorm:"not null;index"` // 操作人（管理员ID）

	// 状态：0失败 1成功
	Status int `json:"status" gorm:"not null;default:1"`

	ErrorMessage string `json:"error_message" gorm:"type:text"` // 导入失败时记录错误原因

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (AcademicImportRecord) TableName() string {
	return "academic_import_records"
}
