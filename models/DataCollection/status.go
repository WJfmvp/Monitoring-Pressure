package DataCollection

// 数据来源枚举
const (
	DataSourceManual       int8 = 1 // 手动录入
	DataSourceEducationSys int8 = 2 // 教务系统导入
)

// 作业提交状态枚举
const (
	HomeworkStatusUnsubmitted int8 = 0 // 未提交
	HomeworkStatusSubmitted   int8 = 1 // 已提交
)

// 作业质量枚举
const (
	HomeworkQualityExcellent int8 = 1 // 优
	HomeworkQualityGood      int8 = 2 // 良
	HomeworkQualityMedium    int8 = 3 // 中
	HomeworkQualityPoor      int8 = 4 // 差
)

// 问卷必填状态枚举
const (
	QuestionnaireNotRequired int8 = 0 // 非必填
	QuestionnaireRequired    int8 = 1 // 必填
)

// 问卷同步状态枚举
const (
	SyncStatusPending    int8 = 0 // 待同步
	SyncStatusProcessing int8 = 1 // 同步中
	SyncStatusSuccess    int8 = 2 // 同步成功
	SyncStatusFailed     int8 = 3 // 同步失败
)
