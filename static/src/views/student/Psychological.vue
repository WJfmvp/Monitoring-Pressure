<template>
  <div>
    <div class="page-card">
      <h2 class="page-title">提交心理测评</h2>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        style="max-width: 720px"
      >
        <el-form-item label="问卷 ID" prop="questionnaire_id">
          <el-input-number v-model="form.questionnaire_id" :min="1" />
        </el-form-item>
        <el-form-item label="评估日期" prop="assess_date">
          <el-date-picker
            v-model="form.assess_date"
            type="date"
            value-format="YYYY-MM-DD"
            placeholder="选择评估日期"
          />
        </el-form-item>

        <el-divider content-position="left">核心指标 (0-10)</el-divider>
        <el-form-item label="焦虑程度" prop="anxiety_level">
          <el-slider v-model="form.anxiety_level" :max="10" show-stops />
        </el-form-item>
        <el-form-item label="学习动力" prop="learning_motivation">
          <el-slider v-model="form.learning_motivation" :max="10" show-stops />
        </el-form-item>
        <el-form-item label="情绪状态" prop="emotional_state">
          <el-slider v-model="form.emotional_state" :max="10" show-stops />
        </el-form-item>

        <el-divider content-position="left">扩展指标</el-divider>
        <el-form-item label="自我压力感知">
          <el-slider v-model="form.stress_perception" :max="10" show-stops />
        </el-form-item>
        <el-form-item label="睡眠质量">
          <el-slider v-model="form.sleep_quality" :max="10" show-stops />
        </el-form-item>
        <el-form-item label="疲劳程度">
          <el-slider v-model="form.fatigue_level" :max="10" show-stops />
        </el-form-item>

        <el-form-item label="问卷分数">
          <el-input-number :model-value="questionnaireScore" :min="0" :max="100" :precision="2" disabled />
          <span class="hint">综合得分自动计算（0-100，越高压力越大）</span>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" maxlength="255" show-word-limit />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="onSubmit">提交</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="page-card" style="margin-top: 16px">
      <h3 class="card-title">历史问卷记录</h3>
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="questionnaire_id" label="问卷 ID" width="100" />
        <el-table-column prop="assess_date" label="评估日期" width="140">
          <template #default="{ row }">{{ formatDate(row.assess_date) }}</template>
        </el-table-column>
        <el-table-column prop="anxiety_level" label="焦虑" width="80" />
        <el-table-column prop="learning_motivation" label="动力" width="80" />
        <el-table-column prop="emotional_state" label="情绪" width="80" />
        <el-table-column prop="stress_perception" label="压力感知" width="100" />
        <el-table-column prop="sleep_quality" label="睡眠" width="80" />
        <el-table-column prop="fatigue_level" label="疲劳" width="80" />
        <el-table-column prop="questionnaire_score" label="总分" width="100" />
        <el-table-column prop="remark" label="备注" />
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { submitPsychological, getPsychologicalList } from '@/api/student'

const formRef = ref(null)
const loading = ref(false)
const list = ref([])

function today() {
  const d = new Date()
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

const defaults = () => ({
  questionnaire_id: 1,
  assess_date: today(),
  anxiety_level: 0,
  learning_motivation: 5,
  emotional_state: 5,
  stress_perception: 0,
  sleep_quality: 5,
  fatigue_level: 0,
  remark: '',
})

const form = ref(defaults())

const rules = {
  questionnaire_id: [{ required: true, message: '请输入问卷 ID', trigger: 'blur' }],
  assess_date: [{ required: true, message: '请选择评估日期', trigger: 'change' }],
}

// 6 项指标各 0-10，归一到 0-100 的压力分（越高表示压力越大）
function computeScore(f) {
  const clamp = (v) => Math.max(0, Math.min(10, Number(v) || 0))
  const positive = clamp(f.anxiety_level) + clamp(f.stress_perception) + clamp(f.fatigue_level)
  const inverted = (10 - clamp(f.learning_motivation)) + (10 - clamp(f.emotional_state)) + (10 - clamp(f.sleep_quality))
  const raw = positive + inverted // 0..60
  return Math.round((raw / 60) * 100 * 100) / 100
}

const questionnaireScore = computed(() => computeScore(form.value))

function toNumber(value, fallback = 0) {
  const n = Number(value)
  return Number.isFinite(n) ? n : fallback
}

function toScaleValue(value) {
  return Math.max(0, Math.min(10, Math.round(toNumber(value))))
}

function buildPayload() {
  return {
    questionnaire_id: Math.max(1, Math.round(toNumber(form.value.questionnaire_id, 1))),
    assess_date: form.value.assess_date,
    anxiety_level: toScaleValue(form.value.anxiety_level),
    learning_motivation: toScaleValue(form.value.learning_motivation),
    emotional_state: toScaleValue(form.value.emotional_state),
    stress_perception: toScaleValue(form.value.stress_perception),
    sleep_quality: toScaleValue(form.value.sleep_quality),
    fatigue_level: toScaleValue(form.value.fatigue_level),
    questionnaire_score: questionnaireScore.value,
    remark: String(form.value.remark || '').trim(),
  }
}

function formatDate(t) {
  if (!t) return ''
  return String(t).slice(0, 10)
}

async function loadList() {
  try {
    const data = await getPsychologicalList()
    list.value = Array.isArray(data) ? data : []
  } catch (_) {}
}

async function onSubmit() {
  await formRef.value.validate()
  loading.value = true
  try {
    await submitPsychological(buildPayload())
    ElMessage.success('提交成功')
    onReset()
    await loadList()
  } finally {
    loading.value = false
  }
}

function onReset() {
  form.value = defaults()
  formRef.value?.clearValidate?.()
}

onMounted(loadList)
</script>

<style scoped>
.hint {
  margin-left: 12px;
  color: #909399;
  font-size: 12px;
}
.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
}
</style>
