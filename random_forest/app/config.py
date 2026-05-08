from pathlib import Path
from pydantic import BaseModel

class Settings(BaseModel):
    app_name: str = "stress-detection-service"
    app_version: str = "1.0.0"
    model_path: Path = Path(__file__).resolve().parent.parent / "stress_model.pkl"


settings = Settings()