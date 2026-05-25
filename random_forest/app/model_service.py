from pathlib import Path
from threading import Lock

import joblib
import pandas as pd

from .config import settings
from .decision import final_decision
from .feature_engineering import FEATURE_COLUMNS, add_features
from .logger import logger


class ModelService:
    def __init__(self, model_path: Path):
        self.model_path = model_path
        self._model = None
        self._lock = Lock()

    def load_model(self) -> None:
        with self._lock:
            if self._model is not None:
                return

            if not self.model_path.exists():
                raise FileNotFoundError(f"模型文件不存在: {self.model_path}")

            self._model = joblib.load(self.model_path)
            logger.info("模型加载成功: %s", self.model_path)

    @property
    def model(self):
        if self._model is None:
            self.load_model()
        return self._model

    def predict_one(self, payload: dict) -> dict:
        df = pd.DataFrame([payload])

        df = add_features(df)

        x = df[FEATURE_COLUMNS]

        model_pred = self.model.predict(x)[0]

        probability = None
        if hasattr(self.model, "predict_proba"):
            proba = self.model.predict_proba(x)[0]
            probability = float(max(proba))
        else:
            probability = 0.0

        df["model_predict_level"] = int(model_pred)
        df["predict_probability"] = float(probability)
        df["final_level"] = df.apply(final_decision, axis=1)

        row = df.iloc[0]

        return {
            "student_id": str(payload["student_id"]),
            "model_predict_level": int(row["model_predict_level"]),
            "final_level": int(row["final_level"]),
            "predict_probability": round(float(row["predict_probability"]), 6),
            "risk_index": round(float(row["risk_index"]), 6),
            "mental_risk": round(float(row["mental_risk"]), 6),
            "behavior_risk": round(float(row["behavior_risk"]), 6),
            "high_anxiety_flag": int(row["high_anxiety_flag"]),
        }


model_service = ModelService(settings.model_path)