package service

import (
	"Monitoring-Pressure/dao/db"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	DataCollection "Monitoring-Pressure/models/data_collection"
	"Monitoring-Pressure/models/users"
)

type academicExcelRow struct {
	Username     string
	ExamName     string
	Term         string
	ExamDate     time.Time
	RawScore     float64
	AverageScore float64
	ClassRank    int
	GradeRank    int
	FailCount    int
}

func ImportAcademicExcel(filePath string, originalFileName string, operatorID int64) error {
	// 1. 先创建导入记录
	importRecord := DataCollection.AcademicImportRecord{
		FileName:    originalFileName,
		FilePath:    filePath,
		ImportCount: 0,
		OperatorID:  operatorID,
		Status:      0, // 先标记处理中/失败，成功后再改成1
	}

	if err := db.DB.Create(&importRecord).Error; err != nil {
		return fmt.Errorf("创建导入记录失败: %w", err)
	}

	// 2. 打开 Excel
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		_ = db.DB.Model(&importRecord).Updates(map[string]interface{}{
			"status":        0,
			"error_message": "打开Excel失败: " + err.Error(),
		}).Error
		return fmt.Errorf("打开Excel失败: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	// 3. 取第一个 sheet
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		_ = db.DB.Model(&importRecord).Updates(map[string]interface{}{
			"status":        0,
			"error_message": "Excel中没有工作表",
		}).Error
		return errors.New("Excel中没有工作表")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		_ = db.DB.Model(&importRecord).Updates(map[string]interface{}{
			"status":        0,
			"error_message": "读取工作表失败: " + err.Error(),
		}).Error
		return fmt.Errorf("读取工作表失败: %w", err)
	}

	if len(rows) <= 1 {
		_ = db.DB.Model(&importRecord).Updates(map[string]interface{}{
			"status":        0,
			"error_message": "Excel没有数据行",
		}).Error
		return errors.New("Excel没有数据行")
	}

	importCount := 0

	// 4. 开事务，保证导入要么都成功，要么都失败
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		for i := 1; i < len(rows); i++ { // 跳过表头，第2行开始
			row := rows[i]

			if isEmptyRow(row) {
				continue
			}

			rowData, err := parseAcademicRow(row)
			if err != nil {
				return fmt.Errorf("第 %d 行解析失败: %w", i+1, err)
			}

			// 5. 根据用户名查找用户
			userID, err := findUserIDByUsername(tx, rowData.Username)
			if err != nil {
				return fmt.Errorf("第 %d 行用户匹配失败: %w", i+1, err)
			}

			// 6. 计算成绩波动、成绩变化率
			scoreFluctuation, scoreChangeRate, err := calcScoreChange(tx, userID, rowData.ExamScore())
			if err != nil {
				return fmt.Errorf("第 %d 行计算成绩变化失败: %w", i+1, err)
			}

			record := DataCollection.AcademicRecord{
				UserID:           userID,
				ImportID:         importRecord.ID,
				ExamName:         rowData.ExamName,
				Term:             rowData.Term,
				ExamDate:         rowData.ExamDate,
				RawScore:         rowData.RawScore,
				ExamScore:        rowData.ExamScore(),
				AverageScore:     rowData.AverageScore,
				ScoreFluctuation: scoreFluctuation,
				ClassRank:        rowData.ClassRank,
				GradeRank:        rowData.GradeRank,
				FailCount:        rowData.FailCount,
				ScoreChangeRate:  scoreChangeRate,
				SourceType:       2, // Excel导入
			}

			if err := tx.Create(&record).Error; err != nil {
				return fmt.Errorf("第 %d 行保存失败: %w", i+1, err)
			}

			importCount++
		}

		return nil
	})

	if err != nil {
		_ = db.DB.Model(&importRecord).Updates(map[string]interface{}{
			"status":        0,
			"error_message": err.Error(),
		}).Error
		return err
	}

	// 7. 更新导入成功状态
	if err := db.DB.Model(&importRecord).Updates(map[string]interface{}{
		"status":        1,
		"import_count":  importCount,
		"error_message": "",
	}).Error; err != nil {
		return fmt.Errorf("更新导入记录失败: %w", err)
	}

	return nil
}

func parseAcademicRow(row []string) (*academicExcelRow, error) {
	// 约定列顺序：
	// A 用户名
	// B 考试名称
	// C 学期
	// D 考试日期
	// E 原始成绩
	// F 平均分
	// G 班级排名
	// H 年级排名
	// I 挂科数

	if len(row) < 9 {
		return nil, fmt.Errorf("列数不足，至少需要9列，当前只有%d列", len(row))
	}

	username := strings.TrimSpace(row[0])
	examName := strings.TrimSpace(row[1])
	term := strings.TrimSpace(row[2])
	examDateStr := strings.TrimSpace(row[3])

	if username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if examName == "" {
		return nil, errors.New("考试名称不能为空")
	}

	examDate, err := parseDate(examDateStr)
	if err != nil {
		return nil, fmt.Errorf("考试日期格式错误: %w", err)
	}

	rawScore, err := parseFloat(row[4])
	if err != nil {
		return nil, fmt.Errorf("原始成绩解析失败: %w", err)
	}

	averageScore, err := parseFloat(row[5])
	if err != nil {
		return nil, fmt.Errorf("平均分解析失败: %w", err)
	}

	classRank, err := parseInt(row[6])
	if err != nil {
		return nil, fmt.Errorf("班级排名解析失败: %w", err)
	}

	gradeRank, err := parseInt(row[7])
	if err != nil {
		return nil, fmt.Errorf("年级排名解析失败: %w", err)
	}

	failCount, err := parseInt(row[8])
	if err != nil {
		return nil, fmt.Errorf("挂科数解析失败: %w", err)
	}

	return &academicExcelRow{
		Username:     username,
		ExamName:     examName,
		Term:         term,
		ExamDate:     examDate,
		RawScore:     rawScore,
		AverageScore: averageScore,
		ClassRank:    classRank,
		GradeRank:    gradeRank,
		FailCount:    failCount,
	}, nil
}

func (r *academicExcelRow) ExamScore() float64 {
	// 这里先直接等于原始成绩
	// 如果你后续要做成绩清洗，可以在这里改
	return r.RawScore
}

func findUserIDByUsername(tx *gorm.DB, username string) (int64, error) {
	var user users.User
	err := tx.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("未找到用户名为 %s 的用户", username)
		}
		return 0, err
	}
	return user.UserID, nil
}

func calcScoreChange(tx *gorm.DB, userID int64, currentScore float64) (float64, float64, error) {
	var lastRecord DataCollection.AcademicRecord

	err := tx.Where("user_id = ?", userID).
		Order("exam_date desc, id desc").
		First(&lastRecord).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, nil
		}
		return 0, 0, err
	}

	scoreFluctuation := math.Abs(currentScore - lastRecord.ExamScore)

	scoreChangeRate := 0.0
	if lastRecord.ExamScore != 0 {
		scoreChangeRate = (currentScore - lastRecord.ExamScore) / lastRecord.ExamScore * 100
	}

	return scoreFluctuation, scoreChangeRate, nil
}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

func parseInt(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("日期不能为空")
	}

	layouts := []string{
		"2006-01-02",
		"2006/01/02",
		"2006.01.02",
		"2006-1-2",
		"2006/1/2",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("不支持的日期格式: %s", s)
}

func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
