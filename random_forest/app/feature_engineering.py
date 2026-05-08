import pandas as pd


FEATURE_COLUMNS = [
    "study_duration",
    "homework_submit_rate",
    "study_frequency",
    "absence_count",
    "exam_score",
    "score_fluctuation",
    "class_rank",
    "anxiety_level",
    "learning_motivation",
    "emotional_state",
    "sleep_quality",
    "behavior_score",
    "academic_score",
    "psychological_score",
    "risk_index",
    "mental_risk",
    "behavior_risk",
    "high_anxiety_flag",
]


def add_features(df: pd.DataFrame) -> pd.DataFrame:
    result = df.copy()

    result["risk_index"] = (
            result["behavior_score"] * 0.2
            + result["academic_score"] * 0.3
            + result["psychological_score"] * 0.5
    )

    result["mental_risk"] = (
            result["anxiety_level"] * 2
            - result["learning_motivation"]
    )

    result["behavior_risk"] = (
            result["absence_count"] * 2
            + (5 - result["study_frequency"])
    )

    result["high_anxiety_flag"] = (result["anxiety_level"] >= 8).astype(int)

    return result