import pandas as pd

def final_decision(row: pd.Series) -> int:
    level = int(row["model_predict_level"])

    if row["psychological_score"] >= 85:
        return 2

    if row["psychological_score"] >= 70:
        level = max(level, 1)

    if row["risk_index"] >= 75:
        return 2

    if row["risk_index"] >= 55:
        level = max(level, 1)

    return int(level)