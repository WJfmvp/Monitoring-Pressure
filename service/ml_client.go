package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// MLPredictRequest 与 Python /api/v1/predict 期望的 JSON 字段对齐
type MLPredictRequest struct {
	StudentID          string  `json:"student_id"`
	StudyDuration      float64 `json:"study_duration"`
	HomeworkSubmitRate float64 `json:"homework_submit_rate"`
	StudyFrequency     int     `json:"study_frequency"`
	AbsenceCount       int     `json:"absence_count"`
	ExamScore          float64 `json:"exam_score"`
	ScoreFluctuation   float64 `json:"score_fluctuation"`
	ClassRank          int     `json:"class_rank"`
	AnxietyLevel       int     `json:"anxiety_level"`
	LearningMotivation int     `json:"learning_motivation"`
	EmotionalState     int     `json:"emotional_state"`
	SleepQuality       int     `json:"sleep_quality"`
	BehaviorScore      float64 `json:"behavior_score"`
	AcademicScore      float64 `json:"academic_score"`
	PsychologicalScore float64 `json:"psychological_score"`
}

// MLPredictData 与 Python 返回的 data 字段对齐
type MLPredictData struct {
	StudentID          string  `json:"student_id"`
	ModelPredictLevel  int     `json:"model_predict_level"`
	FinalLevel         int     `json:"final_level"`
	PredictProbability float64 `json:"predict_probability"`
	RiskIndex          float64 `json:"risk_index"`
	MentalRisk         float64 `json:"mental_risk"`
	BehaviorRisk       float64 `json:"behavior_risk"`
	HighAnxietyFlag    int     `json:"high_anxiety_flag"`
}

type MLPredictResponse struct {
	Code      int           `json:"code"`
	Message   string        `json:"message"`
	Data      MLPredictData `json:"data"`
	RequestID string        `json:"request_id"`
}

// CallMLPredict 向 Python 随机森林服务发送一次预测请求
func CallMLPredict(req MLPredictRequest) (*MLPredictData, error) {
	baseURL := strings.TrimRight(os.Getenv("ML_SERVICE_URL"), "/")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8000"
	}

	timeout := 10 * time.Second
	if t := os.Getenv("ML_SERVICE_TIMEOUT"); t != "" {
		if n, err := strconv.Atoi(t); err == nil && n > 0 {
			timeout = time.Duration(n) * time.Second
		}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化预测请求失败: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, baseURL+"/api/v1/predict", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("构造预测请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("调用预测服务失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取预测响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("预测服务返回非 200: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var parsed MLPredictResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("解析预测响应失败: %w body=%s", err, string(respBody))
	}

	if parsed.Code != 0 {
		return nil, fmt.Errorf("预测服务业务错误: code=%d message=%s", parsed.Code, parsed.Message)
	}

	return &parsed.Data, nil
}
