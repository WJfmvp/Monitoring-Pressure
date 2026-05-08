from pydantic import BaseModel, Field, ConfigDict


class PredictRequest(BaseModel):
    model_config = ConfigDict(extra="forbid")

    student_id: str = Field(..., min_length=1, max_length=64)

    study_duration: float = Field(..., ge=0, le=24)
    homework_submit_rate: float = Field(..., ge=0, le=1)
    study_frequency: int = Field(..., ge=0, le=7)
    absence_count: int = Field(..., ge=0, le=365)
    exam_score: float = Field(..., ge=0, le=100)
    score_fluctuation: float = Field(..., ge=0, le=100)
    class_rank: int = Field(..., ge=1, le=100000)
    anxiety_level: int = Field(..., ge=0, le=10)
    learning_motivation: int = Field(..., ge=0, le=10)
    emotional_state: int = Field(..., ge=0, le=10)
    sleep_quality: int = Field(..., ge=0, le=10)
    behavior_score: float = Field(..., ge=0, le=100)
    academic_score: float = Field(..., ge=0, le=100)
    psychological_score: float = Field(..., ge=0, le=100)


class PredictData(BaseModel):
    student_id: str
    model_predict_level: int
    final_level: int
    predict_probability: float
    risk_index: float
    mental_risk: float
    behavior_risk: float
    high_anxiety_flag: int


class PredictResponse(BaseModel):
    code: int
    message: str
    data: PredictData
    request_id: str