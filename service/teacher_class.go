package service

import (
	"errors"
	"fmt"

	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/data_collection"
	"Monitoring-Pressure/models/users"

	"gorm.io/gorm"
)

// ===================== 班级压力总览 =====================

// ClassStressItem 班级压力排行/分布列表中的一项学生数据
type ClassStressItem struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Sex       int    `json:"sex"`
	Telephone string `json:"telephone"`

	HasResult bool   `json:"has_result"`
	ResultID  uint   `json:"result_id"`
	Assessed  string `json:"assessed_at"`

	BehaviorScore      float64 `json:"behavior_score"`
	AcademicScore      float64 `json:"academic_score"`
	PsychologicalScore float64 `json:"psychological_score"`
	TotalScore         float64 `json:"total_score"`

	ModelPredictLevel  int     `json:"model_predict_level"`
	PredictProbability float64 `json:"predict_probability"`
	WarningStatus      int     `json:"warning_status"`
}

// ClassStressSummary 班级整体统计
type ClassStressSummary struct {
	TotalStudents int `json:"total_students"`
	AssessedCount int `json:"assessed_count"`
	WarnedCount   int `json:"warned_count"`    // warning_status >= 1
	HighRiskCount int `json:"high_risk_count"` // warning_status == 2

	LevelDistribution []LevelBucket `json:"level_distribution"` // 0/1/2 + unassessed

	AvgBehavior      float64 `json:"avg_behavior"`
	AvgAcademic      float64 `json:"avg_academic"`
	AvgPsychological float64 `json:"avg_psychological"`
	AvgTotal         float64 `json:"avg_total"`
}

type LevelBucket struct {
	Level int    `json:"level"` // 0/1/2  -1 表示未评估
	Label string `json:"label"`
	Count int    `json:"count"`
}

type ClassStressOverview struct {
	Summary ClassStressSummary `json:"summary"`
	Items   []ClassStressItem  `json:"items"`
}

// classStressRow 用于扫 JOIN 结果，不暴露给外部
type classStressRow struct {
	UserID    int64
	Username  string
	Sex       int
	Telephone string

	ResultID *uint
	Assessed *string

	BehaviorScore      *float64
	AcademicScore      *float64
	PsychologicalScore *float64
	TotalScore         *float64

	ModelPredictLevel  *int
	PredictProbability *float64
	WarningStatus      *int
}

// GetClassStressOverview 返回班级所有学生的最新一次压力评估 + 总体统计
func GetClassStressOverview() (*ClassStressOverview, error) {
	var rows []classStressRow

	// 取每个学生最新的一条 stress_assessment_result (id 最大)
	sql := `
SELECT
    u.user_id, u.username, u.sex, u.telephone,
    r.id AS result_id,
    DATE_FORMAT(r.created_at, '%Y-%m-%d %H:%i:%s') AS assessed,
    r.behavior_score, r.academic_score, r.psychological_score, r.total_score,
    r.model_predict_level, r.predict_probability, r.warning_status
FROM user u
LEFT JOIN (
    SELECT r1.*
    FROM stress_assessment_results r1
    INNER JOIN (
        SELECT user_id, MAX(id) AS max_id
        FROM stress_assessment_results
        GROUP BY user_id
    ) lt ON r1.id = lt.max_id
) r ON u.user_id = r.user_id
WHERE u.role = ?
ORDER BY r.total_score IS NULL, r.total_score DESC, u.user_id ASC
`
	if err := db.DB.Raw(sql, string(users.RoleStudent)).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("查询班级压力总览失败: %w", err)
	}

	items := make([]ClassStressItem, 0, len(rows))
	var (
		assessed                         int
		warned                           int
		highRisk                         int
		levelBuckets                     = map[int]int{-1: 0, 0: 0, 1: 0, 2: 0}
		sumBeh, sumAca, sumPsy, sumTotal float64
	)

	for _, r := range rows {
		item := ClassStressItem{
			UserID:    r.UserID,
			Username:  r.Username,
			Sex:       r.Sex,
			Telephone: r.Telephone,
		}
		if r.ResultID != nil {
			item.HasResult = true
			item.ResultID = *r.ResultID
			if r.Assessed != nil {
				item.Assessed = *r.Assessed
			}
			if r.BehaviorScore != nil {
				item.BehaviorScore = *r.BehaviorScore
			}
			if r.AcademicScore != nil {
				item.AcademicScore = *r.AcademicScore
			}
			if r.PsychologicalScore != nil {
				item.PsychologicalScore = *r.PsychologicalScore
			}
			if r.TotalScore != nil {
				item.TotalScore = *r.TotalScore
			}
			if r.ModelPredictLevel != nil {
				item.ModelPredictLevel = *r.ModelPredictLevel
			}
			if r.PredictProbability != nil {
				item.PredictProbability = *r.PredictProbability
			}
			if r.WarningStatus != nil {
				item.WarningStatus = *r.WarningStatus
			}

			assessed++
			levelBuckets[item.WarningStatus]++
			if item.WarningStatus >= 1 {
				warned++
			}
			if item.WarningStatus >= 2 {
				highRisk++
			}
			sumBeh += item.BehaviorScore
			sumAca += item.AcademicScore
			sumPsy += item.PsychologicalScore
			sumTotal += item.TotalScore
		} else {
			levelBuckets[-1]++
		}
		items = append(items, item)
	}

	summary := ClassStressSummary{
		TotalStudents: len(rows),
		AssessedCount: assessed,
		WarnedCount:   warned,
		HighRiskCount: highRisk,
		LevelDistribution: []LevelBucket{
			{Level: 0, Label: "正常", Count: levelBuckets[0]},
			{Level: 1, Label: "中等", Count: levelBuckets[1]},
			{Level: 2, Label: "高压力", Count: levelBuckets[2]},
			{Level: -1, Label: "未评估", Count: levelBuckets[-1]},
		},
	}
	if assessed > 0 {
		summary.AvgBehavior = round2(sumBeh / float64(assessed))
		summary.AvgAcademic = round2(sumAca / float64(assessed))
		summary.AvgPsychological = round2(sumPsy / float64(assessed))
		summary.AvgTotal = round2(sumTotal / float64(assessed))
	}

	return &ClassStressOverview{Summary: summary, Items: items}, nil
}

func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}

// ===================== 学生压力诊断书 =====================

type DiagnosisStudent struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Sex       int    `json:"sex"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
}

type DiagnosisHistoryPoint struct {
	ResultID           uint    `json:"result_id"`
	AssessedAt         string  `json:"assessed_at"`
	BehaviorScore      float64 `json:"behavior_score"`
	AcademicScore      float64 `json:"academic_score"`
	PsychologicalScore float64 `json:"psychological_score"`
	TotalScore         float64 `json:"total_score"`
	WarningStatus      int     `json:"warning_status"`
}

type DiagnosisSuggestion struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Level    int    `json:"level"`
}

type StudentDiagnosis struct {
	Student      DiagnosisStudent        `json:"student"`
	Latest       *ClassStressItem        `json:"latest"`
	History      []DiagnosisHistoryPoint `json:"history"`
	Explanations []string                `json:"explanations"`
	Suggestions  []DiagnosisSuggestion   `json:"suggestions"`
}

// GetStudentDiagnosis 单个学生的"压力诊断书"
func GetStudentDiagnosis(studentID int64) (*StudentDiagnosis, error) {
	if studentID <= 0 {
		return nil, errors.New("student_id 不合法")
	}

	// 1. 基本信息
	var u users.User
	if err := db.DB.Where("user_id = ?", studentID).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("学生不存在")
		}
		return nil, fmt.Errorf("查询学生失败: %w", err)
	}

	diag := &StudentDiagnosis{
		Student: DiagnosisStudent{
			UserID:    u.UserID,
			Username:  u.Username,
			Sex:       u.Sex,
			Telephone: u.Telephone,
			Email:     u.Email,
		},
	}

	// 2. 历史评估（最多 12 次）
	var results []data_collection.StressAssessmentResult
	if err := db.DB.
		Where("user_id = ?", studentID).
		Order("id DESC").
		Limit(12).
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("查询历史评估失败: %w", err)
	}

	for i := len(results) - 1; i >= 0; i-- {
		r := results[i]
		diag.History = append(diag.History, DiagnosisHistoryPoint{
			ResultID:           r.ID,
			AssessedAt:         r.CreatedAt.Format("2006-01-02 15:04:05"),
			BehaviorScore:      r.BehaviorScore,
			AcademicScore:      r.AcademicScore,
			PsychologicalScore: r.PsychologicalScore,
			TotalScore:         r.TotalScore,
			WarningStatus:      r.WarningStatus,
		})
	}

	// 3. 最新一次结果（用于卡片 + 雷达 + 风险解读）
	if len(results) > 0 {
		latest := results[0]
		assessedAt := latest.CreatedAt.Format("2006-01-02 15:04:05")
		diag.Latest = &ClassStressItem{
			UserID:             u.UserID,
			Username:           u.Username,
			Sex:                u.Sex,
			Telephone:          u.Telephone,
			HasResult:          true,
			ResultID:           latest.ID,
			Assessed:           assessedAt,
			BehaviorScore:      latest.BehaviorScore,
			AcademicScore:      latest.AcademicScore,
			PsychologicalScore: latest.PsychologicalScore,
			TotalScore:         latest.TotalScore,
			ModelPredictLevel:  latest.ModelPredictLevel,
			PredictProbability: latest.PredictProbability,
			WarningStatus:      latest.WarningStatus,
		}
		diag.Explanations = buildExplanations(latest)
		diag.Suggestions = matchSuggestions(latest)
	} else {
		diag.Explanations = []string{"该学生暂无压力评估记录，建议尽快完成一次评估。"}
	}

	return diag, nil
}

// buildExplanations 根据维度得分生成风险解读
// 规则与 Python decision.py 保持一致
func buildExplanations(r data_collection.StressAssessmentResult) []string {
	out := []string{}

	switch r.WarningStatus {
	case 2:
		out = append(out, "⚠ 综合预警等级为「高压力」，需要立即介入干预。")
	case 1:
		out = append(out, "⚠ 综合预警等级为「中等压力」，建议关注并适度干预。")
	default:
		out = append(out, "✓ 综合预警等级为「正常」，整体状态良好。")
	}

	if r.PsychologicalScore >= 85 {
		out = append(out, fmt.Sprintf("心理维度得分严重偏高 (%.1f / 100)，存在明显焦虑/情绪困扰风险。", r.PsychologicalScore))
	} else if r.PsychologicalScore >= 70 {
		out = append(out, fmt.Sprintf("心理维度得分偏高 (%.1f / 100)，建议关注情绪状态与睡眠质量。", r.PsychologicalScore))
	}

	if r.AcademicScore >= 75 {
		out = append(out, fmt.Sprintf("学业维度得分偏高 (%.1f / 100)，可能存在成绩波动或排名压力。", r.AcademicScore))
	}

	if r.BehaviorScore >= 75 {
		out = append(out, fmt.Sprintf("行为维度得分偏高 (%.1f / 100)，可能伴随缺勤、作业提交不稳定等情况。", r.BehaviorScore))
	}

	if r.TotalScore >= 75 {
		out = append(out, fmt.Sprintf("综合风险指数为 %.1f / 100，多维度压力叠加，优先级最高。", r.TotalScore))
	} else if r.TotalScore >= 55 {
		out = append(out, fmt.Sprintf("综合风险指数为 %.1f / 100，存在中等强度压力。", r.TotalScore))
	}

	if r.ModelPredictLevel != r.WarningStatus {
		out = append(out, fmt.Sprintf("说明：模型预测等级为 %d，规则修正后的最终等级为 %d。", r.ModelPredictLevel, r.WarningStatus))
	}

	return out
}

// matchSuggestions 按等级匹配建议模板，最多返回 6 条
func matchSuggestions(r data_collection.StressAssessmentResult) []DiagnosisSuggestion {
	out := []DiagnosisSuggestion{}

	var suggestions []data_collection.InterventionSuggestion
	if err := db.DB.
		Where("level = ?", r.WarningStatus).
		Order("id ASC").
		Limit(6).
		Find(&suggestions).Error; err != nil {
		return out
	}

	for _, s := range suggestions {
		out = append(out, DiagnosisSuggestion{
			ID:       s.ID,
			Title:    s.Title,
			Content:  s.Content,
			Category: s.Category,
			Level:    s.Level,
		})
	}
	return out
}
