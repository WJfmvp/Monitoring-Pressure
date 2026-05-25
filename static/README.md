# 学业压力监测平台 · 前端

Vue 3 + Vite + Element Plus + ECharts，对接 Go 后端（`:8080`）和 Python 随机森林预测服务（`:8000`）。

## 目录结构

```
static/
├── index.html
├── package.json
├── vite.config.js              # 含 /api → http://127.0.0.1:8080 反向代理
├── .env.development            # 开发环境变量
├── .env.production             # 生产环境变量（自行修改）
└── src/
    ├── main.js                 # 入口：注册 Element Plus / Pinia / Router
    ├── App.vue
    ├── api/                    # 按角色拆分的接口模块
    │   ├── request.js          # axios 实例 + 拦截器（统一响应/错误/JWT）
    │   ├── account.js
    │   ├── student.js
    │   ├── teacher.js
    │   └── admin.js
    ├── router/index.js         # 路由 + 守卫（按角色分发主页 + 校验权限）
    ├── stores/user.js          # Pinia 用户/Token 状态
    ├── utils/auth.js           # localStorage 读写
    ├── layouts/BasicLayout.vue # 侧边栏 + 顶栏 + 面包屑（按角色生成菜单）
    ├── styles/index.css        # 全局样式
    └── views/
        ├── account/            # Login / Register / ResetPassword
        ├── profile/            # Profile / CompleteInfo
        ├── student/            # Home(仪表盘) / Psychological / Academic / StressResult / Intervention
        ├── teacher/            # Home / 学生成绩 / 学生压力 / 学生干预
        └── admin/              # Home / 用户 / Excel 导入 / 导入记录 / 成绩 / 压力评估 / 运行预测 / 干预建议
```

## 启动步骤

### 1. 启动后端（两端都要在线）

```bash
# Go 后端 :8080
cd F:/Monitoring-Pressure
go run main.go

# Python 随机森林服务 :8000（运行预测时必须）
cd F:/Monitoring-Pressure/random_forest
python run.py
```

### 2. 安装前端依赖（首次）

```bash
cd F:/Monitoring-Pressure/static
npm install
```

> 镜像建议：国内可用 `npm install --registry=https://registry.npmmirror.com`

### 3. 开发模式

```bash
npm run dev
```

打开 http://localhost:5173 ，登录/注册后会按 JWT 中的 `role` 自动跳转到对应主页：
- `student` → `/student/home`
- `teacher` → `/teacher/home`
- `admin`   → `/admin/home`

### 4. 生产构建

```bash
npm run build         # 输出到 static/dist
npm run preview       # 本地预览构建结果
```

## 前后端联调

`vite.config.js` 已配置反向代理：所有 `/api/*` 请求会被转发到 `http://127.0.0.1:8080/*`（去掉 `/api` 前缀）。

axios 的 `baseURL = /api`（来自 `.env.development` 中的 `VITE_API_BASE`），生产环境部署时按需修改：
- 同源部署：保留 `/api`，由 nginx 反代
- 跨域部署：改成完整地址，并要求 Go 端开启 CORS

## 接口适配说明

后端返回结构有多种形态，前端在 `src/api/request.js` 中统一兼容：

| 后端返回 | 处理 |
|---|---|
| `{code:0,msg,data}` (`util.ResponseSuccess`) | 解出 `data` |
| `{code:200,msg,data}` (手写) | 解出 `data` |
| `{code:1007/1009/1010,msg}` | 清空 token 并跳登录 |
| `{message,data}` | 解出 `data` |
| `{error}` | 弹错误并 reject |

## 预置角色入口

- 注册接口默认创建 `student` 角色，可登录 `admin` 后在 **用户管理** 页面手动改角色
- 登录后顶栏头像 → 退出 / 个人中心 / 完善资料

## 已对接的后端路由总览

| 模块 | 路由 |
|---|---|
| 登录注册 | `/account/sendVerifyCode` `/account/login` `/account/login/code` `/account/register` `/account/resetPassword` |
| 个人 | `/account/profile` `/account/completeInformation` |
| 学生 | `/student/home` `/student/academic` `/student/psychological/submit` `/student/psychological/list` `/student/stress/result/latest` `/student/intervention/list` |
| 教师 | `/teacher/home` `/teacher/student/academic/list` `/teacher/student/stress/list` `/teacher/student/intervention/list` |
| 管理员 | `/admin/home` `/admin/user/list` `/admin/user/updateRole` `/admin/academic/import` `/admin/academic/import/list` `/admin/academic/list` `/admin/stress/result/list` `/admin/stress/predict` `/admin/intervention/create` `/admin/intervention/list` |