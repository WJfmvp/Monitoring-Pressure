import uuid
from contextlib import asynccontextmanager

from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import JSONResponse

from .config import settings
from .logger import logger
from .model_service import model_service
from .schemas import PredictRequest, PredictResponse


@asynccontextmanager
async def lifespan(app: FastAPI):
    model_service.load_model()
    logger.info("服务启动完成")
    yield


app = FastAPI(
    title="Academic Stress Detection Service",
    version=settings.app_version,
    docs_url="/docs",
    redoc_url="/redoc",
    lifespan=lifespan,
)


@app.middleware("http")
async def add_request_id(request: Request, call_next):
    request_id = uuid.uuid4().hex[:12]
    request.state.request_id = request_id

    try:
        response = await call_next(request)
        response.headers["X-Request-ID"] = request_id
        return response
    except Exception as exc:
        logger.exception("未处理异常, request_id=%s, error=%s", request_id, str(exc))
        return JSONResponse(
            status_code=500,
            content={
                "code": 500,
                "message": "internal server error",
                "data": None,
                "request_id": request_id,
            },
        )


@app.get("/health")
def health():
    return {
        "code": 0,
        "message": "ok",
        "data": {
            "service": settings.app_name,
            "version": settings.app_version,
            "model_loaded": model_service.model is not None,
        },
    }


@app.post("/api/v1/predict", response_model=PredictResponse)
def predict(request: PredictRequest, raw_request: Request):
    request_id = getattr(raw_request.state, "request_id", uuid.uuid4().hex[:12])

    try:
        result = model_service.predict_one(request.model_dump())
        return {
            "code": 0,
            "message": "success",
            "data": result,
            "request_id": request_id,
        }
    except FileNotFoundError as exc:
        logger.error("模型文件错误, request_id=%s, error=%s", request_id, str(exc))
        raise HTTPException(status_code=500, detail="model file not found")
    except Exception as exc:
        logger.exception("预测失败, request_id=%s, error=%s", request_id, str(exc))
        raise HTTPException(status_code=500, detail="prediction failed")
