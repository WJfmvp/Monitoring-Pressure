package DataCollection

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/DataCollection"
	"gorm.io/gorm"
	"time"
)

// LearningBehaviorRepo 学习行为表数据访问层
type LearningBehaviorRepo struct {
	db *gorm.DB
}

// NewLearningBehaviorRepo 创建学习行为Repo实例
func NewLearningBehaviorRepo() *LearningBehaviorRepo {
	return &LearningBehaviorRepo{db: db.DB}
}

// Create 新增单条学习行为数据
func (repo *LearningBehaviorRepo) Create(behavior *DataCollection.LearningBehavior) error {
	// 补充默认时间
	behavior.CreateTime = time.Now()
	behavior.UpdateTime = time.Now()
	return repo.db.Create(behavior).Error
}

// CreateBatch 批量新增学习行为数据
func (repo *LearningBehaviorRepo) CreateBatch(behaviors []*DataCollection.LearningBehavior) error {
	// 批量补充时间
	for _, b := range behaviors {
		b.CreateTime = time.Now()
		b.UpdateTime = time.Now()
	}
	return repo.db.CreateInBatches(behaviors, 100).Error // 每100条批量插入
}

// GetByID 根据ID查询单条数据
func (repo *LearningBehaviorRepo) GetByID(id int64) (*DataCollection.LearningBehavior, error) {
	var behavior DataCollection.LearningBehavior
	err := repo.db.Where("id = ?", id).First(&behavior).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // 无数据返回nil，避免上层处理空指针
	}
	return &behavior, err
}

// ListByCondition 条件查询（支持学生ID、日期范围、数据来源）
func (repo *LearningBehaviorRepo) ListByCondition(studentID string, startDate, endDate time.Time, dataSource int8) ([]*DataCollection.LearningBehavior, int64, error) {
	var behaviors []*DataCollection.LearningBehavior
	var count int64

	// 构建查询条件
	query := repo.db.Model(&DataCollection.LearningBehavior{})
	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("learning_date BETWEEN ? AND ?", startDate, endDate)
	}
	if dataSource != 0 {
		query = query.Where("data_source = ?", dataSource)
	}

	// 先查总数
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 查列表（支持分页，这里示例默认查全部，可扩展page/pageSize）
	err = query.Find(&behaviors).Error
	return behaviors, count, err
}

// Update 更新学习行为数据（只更新非空字段）
func (repo *LearningBehaviorRepo) Update(behavior *DataCollection.LearningBehavior) error {
	behavior.UpdateTime = time.Now()
	// 使用Select指定更新字段，避免更新主键/创建时间等字段
	return repo.db.Model(behavior).Select(
		"course_name", "learning_duration", "homework_status",
		"homework_quality", "data_source", "update_time",
	).Where("id = ?", behavior.ID).Updates(behavior).Error
}

// Delete 软删除（物理删除需用Unscoped()）
func (repo *LearningBehaviorRepo) Delete(id int64) error {
	return repo.db.Delete(&DataCollection.LearningBehavior{}, id).Error
}

// DeleteBatch 批量软删除
func (repo *LearningBehaviorRepo) DeleteBatch(ids []int64) error {
	return repo.db.Delete(&DataCollection.LearningBehavior{}, ids).Error
}
