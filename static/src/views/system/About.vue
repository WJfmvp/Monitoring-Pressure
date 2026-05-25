<template>
  <div>
    <div class="page-card hero">
      <el-icon :size="48" color="#fff"><DataAnalysis /></el-icon>
      <div>
        <h1 class="hero-title">学业压力监测平台</h1>
        <p class="hero-sub">基于机器学习的学生学业压力检测与早期干预系统</p>
      </div>
    </div>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="14">
        <div class="page-card">
          <h3 class="card-title">平台简介</h3>
          <p class="paragraph">
            本平台面向中学场景，整合学生的<strong>学习行为</strong>、<strong>学业成绩</strong>
            和<strong>心理状态</strong>三类数据，通过随机森林模型量化压力等级，
            为学生提供自我察觉、为教师提供班级洞察、为管理员提供干预闭环。
          </p>
          <p class="paragraph">
            模型在 Python 端独立服务化，可单独迭代；Go 端负责账户、权限、数据采集与最终决策；
            前端基于 Vue 3 + Element Plus 提供分角色看板。
          </p>
        </div>
      </el-col>
      <el-col :span="10">
        <div class="page-card">
          <h3 class="card-title">版本信息</h3>
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="平台版本">v1.0.0</el-descriptions-item>
            <el-descriptions-item label="发布日期">2026-05</el-descriptions-item>
            <el-descriptions-item label="构建模式">{{ buildMode }}</el-descriptions-item>
            <el-descriptions-item label="后端地址">{{ apiBase }}</el-descriptions-item>
          </el-descriptions>
        </div>
      </el-col>
    </el-row>

    <div class="page-card" style="margin-top: 16px">
      <h3 class="card-title">技术栈</h3>
      <el-row :gutter="16">
        <el-col :span="8" v-for="g in STACK" :key="g.title">
          <div class="stack-card">
            <h4>
              <el-icon :color="g.color"><component :is="g.icon" /></el-icon>
              {{ g.title }}
            </h4>
            <el-tag
              v-for="t in g.items"
              :key="t"
              size="small"
              effect="plain"
              style="margin: 4px 6px 4px 0"
            >
              {{ t }}
            </el-tag>
          </div>
        </el-col>
      </el-row>
    </div>

    <div class="page-card" style="margin-top: 16px">
      <h3 class="card-title">核心模块</h3>
      <el-timeline>
        <el-timeline-item
          v-for="m in MODULES"
          :key="m.title"
          :type="m.color"
          :timestamp="m.tag"
          placement="top"
        >
          <h4>{{ m.title }}</h4>
          <p class="paragraph">{{ m.desc }}</p>
        </el-timeline-item>
      </el-timeline>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const buildMode = computed(() => (import.meta.env.PROD ? '生产' : '开发'))
const apiBase = computed(() => import.meta.env.VITE_API_BASE || '/api')

const STACK = [
  {
    title: '前端',
    color: '#42b883',
    icon: 'Monitor',
    items: ['Vue 3', 'Vite 5', 'Element Plus', 'Pinia', 'Vue Router', 'ECharts', 'Axios'],
  },
  {
    title: '后端',
    color: '#00add8',
    icon: 'Coin',
    items: ['Go 1.25', 'Gin', 'GORM', 'JWT', 'MySQL', 'Sonyflake'],
  },
  {
    title: '算法',
    color: '#ffd43b',
    icon: 'MagicStick',
    items: ['Python', 'FastAPI', 'scikit-learn', 'pandas', 'Random Forest'],
  },
]

const MODULES = [
  {
    title: '数据采集',
    tag: '基础',
    color: 'primary',
    desc: '心理问卷、学习行为、Excel 成绩导入三路汇聚到特征快照表。',
  },
  {
    title: '压力检测',
    tag: '核心',
    color: 'success',
    desc: '随机森林模型预测基础等级 + 规则后置修正，最终判定 0/1/2 三档预警。',
  },
  {
    title: '班级洞察',
    tag: '教师视角',
    color: 'warning',
    desc: '等级分布、维度均分、压力排行表，对高压力学生提供「压力诊断书」入口。',
  },
  {
    title: '干预闭环',
    tag: '管理员维度',
    color: 'danger',
    desc: '管理员维护干预建议库，按等级匹配学生，记录推送与反馈状态。',
  },
]
</script>

<style scoped>
.hero {
  display: flex;
  align-items: center;
  gap: 18px;
  background: linear-gradient(135deg, #6e8efb 0%, #a777e3 100%);
  color: #fff;
  padding: 28px 32px;
}
.hero-title {
  margin: 0 0 4px 0;
  color: #fff;
  font-size: 24px;
}
.hero-sub {
  margin: 0;
  color: rgba(255, 255, 255, 0.85);
  font-size: 14px;
}
.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
}
.paragraph {
  color: #606266;
  line-height: 1.8;
  margin: 0 0 12px 0;
}
.stack-card {
  background: #f8f9fc;
  border-radius: 8px;
  padding: 14px 16px;
}
.stack-card h4 {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0 0 10px 0;
}
</style>