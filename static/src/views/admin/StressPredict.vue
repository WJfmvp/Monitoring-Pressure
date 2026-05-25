<template>
  <div>
    <div class="page-card">
      <h2 class="page-title">运行压力预测</h2>
      <p class="hint">字段范围与 Python 随机森林服务的 schema 一致；提交后会保存特征快照与评估结果。</p>

      <el-form ref="formRef" :model="form" :rules="rules" label-width="140px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="学生 ID" prop="user_id">
              <el-input-number v-model="form.user_id" :min="1" controls-position="right" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="快照日期">
              <el-date-picker
                v-model="form.snapshot_date"
                type="date"
                value-format="YYYY-MM-DD"
                placeholder="缺省取今天"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">学习行为</el-divider>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="学习时长 (h/天)" prop="study_duration">
              <el-input-number v-model="form.study_duration" :min="0" :max="24" :precision="1" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="作业提交率 (0-1)" prop="homework_submit_rate">
              <el-input-number
                v-model="form.homework_submit_rate"
                :min="0"
                :max="1"
                :step="0.05"
                :precision="2"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="学习频率 (天/周)" prop="study_frequency">
              <el-input-number v-model="form.study_frequency" :min="0" :max="7" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="缺勤次数" prop="absence_count">
              <el-input-number v-model="form.absence_count" :min="0" :max="365" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">学业成绩</el-divider>
        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="考试成绩" prop="exam_score">
              <el-input-number v-model="form.exam_score" :min="0" :max="100" :precision="2" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="成绩波动" prop="score_fluctuation">
              <el-input-number
                v-model="form.score_fluctuation"
                :min="0"
                :max="100"
                :precision="2"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="班级排名" prop="class_rank">
              <el-input-number v-model="form.class_rank" :min="1" :max="100000" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">心理特征 (0-10)</el-divider>
        <el-row :gutter="16">
          <el-col :span="6">
            <el-form-item label="焦虑">
              <el-input-number v-model="form.anxiety_level" :min="0" :max="10" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="学习动力">
              <el-input-number v-model="form.learning_motivation" :min="0" :max="10" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="情绪状态">
              <el-input-number v-model="form.emotional_state" :min="0" :max="10" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="睡眠质量">
              <el-input-number v-model="form.sleep_quality" :min="0" :max="10" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">维度综合得分 (0-100)</el-divider>
        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="行为得分">
              <el-input-number
                v-model="form.behavior_score"
                :min="0"
                :max="100"
                :precision="2"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="学业得分">
              <el-input-number
                v-model="form.academic_score"
                :min="0"
                :max="100"
                :precision="2"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="心理得分">
              <el-input-number
                v-model="form.psychological_score"
                :min="0"
                :max="100"
                :precision="2"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="onSubmit">运行预测</el-button>
          <el-button @click="reset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div v-if="result" class="page-card" style="margin-top: 16px">
      <h3 class="card-title">预测结果</h3>
      <el-descriptions :column="3" border>
        <el-descriptions-item label="评估 ID">{{ result.result_id }}</el-descriptions-item>
        <el-descriptions-item label="快照 ID">{{ result.snapshot_id }}</el-descriptions-item>
        <el-descriptions-item label="生成时间">{{ result.created_at }}</el-descriptions-item>
        <el-descriptions-item label="模型预测等级">
          <el-tag :type="LEVEL_TAG[result.model_predict_level]">
            {{ LEVEL_LABEL[result.model_predict_level] }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="综合预警等级">
          <el-tag :type="LEVEL_TAG[result.warning_status]">
            {{ LEVEL_LABEL[result.warning_status] }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="预测概率">
          {{ ((result.predict_probability || 0) * 100).toFixed(2) }}%
        </el-descriptions-item>
        <el-descriptions-item label="综合得分">{{ result.total_score?.toFixed?.(2) }}</el-descriptions-item>
        <el-descriptions-item label="心理风险">{{ result.mental_risk?.toFixed?.(2) }}</el-descriptions-item>
        <el-descriptions-item label="行为风险">{{ result.behavior_risk?.toFixed?.(2) }}</el-descriptions-item>
        <el-descriptions-item label="高焦虑标记">
          {{ result.high_anxiety_flag ? '是' : '否' }}
        </el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { predictStress } from '@/api/admin'

const LEVEL_LABEL = ['正常', '中等', '高压力']
const LEVEL_TAG = ['success', 'warning', 'danger']

const formRef = ref(null)
const loading = ref(false)
const result = ref(null)

const defaults = () => ({
  user_id: 1,
  snapshot_date: '',
  study_duration: 6,
  homework_submit_rate: 0.85,
  study_frequency: 5,
  absence_count: 0,
  exam_score: 75,
  score_fluctuation: 5,
  class_rank: 10,
  anxiety_level: 5,
  learning_motivation: 6,
  emotional_state: 6,
  sleep_quality: 6,
  behavior_score: 75,
  academic_score: 70,
  psychological_score: 65,
})

const form = ref(defaults())

const rules = {
  user_id: [{ required: true, message: '请输入学生 ID', trigger: 'blur' }],
}

async function onSubmit() {
  await formRef.value.validate()
  loading.value = true
  try {
    result.value = await predictStress(form.value)
    ElMessage.success('预测完成')
  } finally {
    loading.value = false
  }
}

function reset() {
  form.value = defaults()
  result.value = null
}
</script>

<style scoped>
.hint {
  color: #909399;
  margin: 0 0 12px 0;
}
.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
}
</style>