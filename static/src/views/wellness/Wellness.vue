<template>
  <div>
    <div class="page-card hero">
      <div>
        <h2 class="page-title" style="margin: 0">放松一下</h2>
        <p class="hint" style="margin: 6px 0 0 0">
          压力大的时候，试试这些小工具，给自己一个喘息的间隙。
        </p>
      </div>
      <el-icon :size="56" class="hero-icon"><Sunny /></el-icon>
    </div>

    <el-row :gutter="16" style="margin-top: 16px">
      <!-- 心情打卡 -->
      <el-col :span="8">
        <div class="page-card tool-card">
          <h3 class="card-title">
            <el-icon><Comment /></el-icon>
            今日心情打卡
          </h3>
          <div class="mood-row">
            <div
              v-for="m in MOODS"
              :key="m.v"
              class="mood-item"
              :class="{ active: todayMood === m.v }"
              @click="setMood(m.v)"
            >
              <span class="emoji">{{ m.emoji }}</span>
              <span class="label">{{ m.label }}</span>
            </div>
          </div>
          <p class="mood-msg" v-if="moodMessage">{{ moodMessage }}</p>

          <h4 class="sub-title">近 7 天</h4>
          <div class="mood-history">
            <div v-for="d in last7Days" :key="d.date" class="mood-cell">
              <span class="cell-emoji">{{ d.emoji || '·' }}</span>
              <span class="cell-date">{{ d.label }}</span>
            </div>
          </div>
        </div>
      </el-col>

      <!-- 呼吸训练 -->
      <el-col :span="8">
        <div class="page-card tool-card">
          <h3 class="card-title">
            <el-icon><WindPower /></el-icon>
            478 呼吸训练
          </h3>
          <div class="breath-stage">
            <div
              class="breath-circle"
              :class="breathPhase"
              :style="{ animationDuration: breathDuration + 's' }"
            >
              <span class="breath-label">{{ breathPhaseLabel }}</span>
              <span class="breath-count">{{ breathRemaining }}s</span>
            </div>
          </div>
          <p class="hint" style="text-align: center">
            吸气 4 秒 → 屏息 7 秒 → 呼气 8 秒
          </p>
          <el-button
            :type="breathRunning ? 'danger' : 'primary'"
            class="full-btn"
            @click="toggleBreath"
          >
            {{ breathRunning ? '停止' : '开始训练' }}
          </el-button>
          <p class="counter">已完成 {{ breathCycles }} 个完整循环</p>
        </div>
      </el-col>

      <!-- 番茄钟 -->
      <el-col :span="8">
        <div class="page-card tool-card">
          <h3 class="card-title">
            <el-icon><Timer /></el-icon>
            学习番茄钟
          </h3>
          <el-input v-model="task" placeholder="此刻专注做什么..." maxlength="40" />
          <div class="timer-display" :class="pomoState">{{ timerDisplay }}</div>
          <p class="timer-state">{{ pomoStateLabel }}</p>
          <div class="timer-actions">
            <el-button type="primary" :disabled="pomoRunning" @click="startPomo">
              {{ pomoRemaining < pomoTotal ? '继续' : '开始' }}
            </el-button>
            <el-button :disabled="!pomoRunning" @click="pausePomo">暂停</el-button>
            <el-button @click="resetPomo">重置</el-button>
          </div>
          <p class="counter">今日已完成 {{ pomoCompletedToday }} 个番茄</p>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="14">
        <div class="page-card tip-card">
          <h3 class="card-title">
            <el-icon><Reading /></el-icon>
            心理调节小贴士
          </h3>
          <el-row :gutter="12">
            <el-col :span="12" v-for="t in TIPS" :key="t.title">
              <div class="tip-item">
                <el-icon :size="20" :color="t.color">
                  <component :is="t.icon" />
                </el-icon>
                <div class="tip-text">
                  <h4>{{ t.title }}</h4>
                  <p>{{ t.content }}</p>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
      </el-col>
      <el-col :span="10">
        <div class="page-card quote-card">
          <h3 class="card-title">每日一句</h3>
          <p class="quote">"{{ currentQuote.text }}"</p>
          <p class="quote-author">—— {{ currentQuote.author }}</p>
          <el-button link type="primary" @click="rotateQuote">
            换一句 <el-icon style="margin-left: 4px"><Refresh /></el-icon>
          </el-button>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

// =============== 心情打卡 ===============
const MOODS = [
  { v: 5, emoji: '😄', label: '很棒' },
  { v: 4, emoji: '🙂', label: '还行' },
  { v: 3, emoji: '😐', label: '一般' },
  { v: 2, emoji: '😞', label: '不太好' },
  { v: 1, emoji: '😢', label: '很糟' },
]

const MOOD_KEY = 'mp_mood_log'

function todayKey() {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const moodLog = ref({})
const todayMood = computed(() => moodLog.value[todayKey()] || 0)

function loadMood() {
  try {
    moodLog.value = JSON.parse(localStorage.getItem(MOOD_KEY) || '{}')
  } catch {
    moodLog.value = {}
  }
}
function setMood(v) {
  moodLog.value = { ...moodLog.value, [todayKey()]: v }
  localStorage.setItem(MOOD_KEY, JSON.stringify(moodLog.value))
}

const moodMessage = computed(() => {
  if (!todayMood.value) return '点一下表情，记录今天的心情吧。'
  const map = {
    5: '太棒了！今天的好心情记得保存哦 ✨',
    4: '保持这份从容，继续加油 💪',
    3: '平淡也是生活的一部分，慢慢来。',
    2: '辛苦了，记得给自己一个拥抱 🤗',
    1: '难过时不必硬撑，试试做几次深呼吸。',
  }
  return map[todayMood.value]
})

const last7Days = computed(() => {
  const out = []
  for (let i = 6; i >= 0; i--) {
    const d = new Date()
    d.setDate(d.getDate() - i)
    const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
    const v = moodLog.value[key]
    out.push({
      date: key,
      label: `${d.getMonth() + 1}/${d.getDate()}`,
      emoji: v ? MOODS.find(m => m.v === v)?.emoji : '',
    })
  }
  return out
})

// =============== 呼吸训练（478）===============
const PHASES = [
  { name: 'inhale', label: '吸气', sec: 4 },
  { name: 'hold', label: '屏息', sec: 7 },
  { name: 'exhale', label: '呼气', sec: 8 },
]
const breathRunning = ref(false)
const breathPhaseIdx = ref(0)
const breathRemaining = ref(4)
const breathCycles = ref(0)
let breathTimer = null

const breathPhase = computed(() => PHASES[breathPhaseIdx.value].name)
const breathPhaseLabel = computed(() => PHASES[breathPhaseIdx.value].label)
const breathDuration = computed(() => PHASES[breathPhaseIdx.value].sec)

function toggleBreath() {
  if (breathRunning.value) stopBreath()
  else startBreath()
}
function startBreath() {
  breathRunning.value = true
  breathPhaseIdx.value = 0
  breathRemaining.value = PHASES[0].sec
  breathTimer = setInterval(() => {
    breathRemaining.value -= 1
    if (breathRemaining.value <= 0) {
      const next = (breathPhaseIdx.value + 1) % PHASES.length
      if (next === 0) breathCycles.value += 1
      breathPhaseIdx.value = next
      breathRemaining.value = PHASES[next].sec
    }
  }, 1000)
}
function stopBreath() {
  breathRunning.value = false
  if (breathTimer) clearInterval(breathTimer)
  breathTimer = null
  breathPhaseIdx.value = 0
  breathRemaining.value = PHASES[0].sec
}

// =============== 番茄钟 ===============
const POMO_KEY = 'mp_pomo_log'
const POMO_WORK = 25 * 60
const POMO_BREAK = 5 * 60

const task = ref('')
const pomoState = ref('work') // 'work' | 'break'
const pomoTotal = ref(POMO_WORK)
const pomoRemaining = ref(POMO_WORK)
const pomoRunning = ref(false)
const pomoCompletedToday = ref(0)
let pomoTimer = null

const pomoStateLabel = computed(() => (pomoState.value === 'work' ? '专注时段 · 25 分钟' : '休息时段 · 5 分钟'))
const timerDisplay = computed(() => {
  const m = String(Math.floor(pomoRemaining.value / 60)).padStart(2, '0')
  const s = String(pomoRemaining.value % 60).padStart(2, '0')
  return `${m}:${s}`
})

function loadPomo() {
  try {
    const log = JSON.parse(localStorage.getItem(POMO_KEY) || '{}')
    pomoCompletedToday.value = log[todayKey()] || 0
  } catch {
    pomoCompletedToday.value = 0
  }
}
function incPomoToday() {
  let log = {}
  try {
    log = JSON.parse(localStorage.getItem(POMO_KEY) || '{}')
  } catch {}
  log[todayKey()] = (log[todayKey()] || 0) + 1
  localStorage.setItem(POMO_KEY, JSON.stringify(log))
  pomoCompletedToday.value = log[todayKey()]
}

function startPomo() {
  pomoRunning.value = true
  pomoTimer = setInterval(() => {
    pomoRemaining.value -= 1
    if (pomoRemaining.value <= 0) {
      clearInterval(pomoTimer)
      pomoTimer = null
      pomoRunning.value = false
      if (pomoState.value === 'work') {
        incPomoToday()
        pomoState.value = 'break'
        pomoTotal.value = POMO_BREAK
        pomoRemaining.value = POMO_BREAK
      } else {
        pomoState.value = 'work'
        pomoTotal.value = POMO_WORK
        pomoRemaining.value = POMO_WORK
      }
    }
  }, 1000)
}
function pausePomo() {
  pomoRunning.value = false
  if (pomoTimer) clearInterval(pomoTimer)
  pomoTimer = null
}
function resetPomo() {
  pausePomo()
  pomoState.value = 'work'
  pomoTotal.value = POMO_WORK
  pomoRemaining.value = POMO_WORK
}

// =============== 每日一句 ===============
const QUOTES = [
  { text: '把每一个清晨当作礼物，把每一次心跳当作奇迹。', author: '佚名' },
  { text: '我们必须接受失望，因为它是有限的，但千万不可失去希望，因为它是无穷的。', author: '马丁·路德·金' },
  { text: '比起结果，过程中的每一次小进步更值得珍惜。', author: '佚名' },
  { text: '允许自己有难过的瞬间，但别让它变成永远。', author: '佚名' },
  { text: '专注做好这一件事，其它的留给时间。', author: '佚名' },
  { text: '你不必完美，只要真实。', author: '佚名' },
  { text: '深呼吸，你已经走得很远了。', author: '佚名' },
  { text: '人生不是和别人比赛，而是和昨天的自己比赛。', author: '佚名' },
]
const quoteIdx = ref(Math.floor(Math.random() * QUOTES.length))
const currentQuote = computed(() => QUOTES[quoteIdx.value])
function rotateQuote() {
  let next = quoteIdx.value
  while (next === quoteIdx.value) next = Math.floor(Math.random() * QUOTES.length)
  quoteIdx.value = next
}

// =============== 调节小贴士 ===============
const TIPS = [
  { icon: 'MoonNight', color: '#6e8efb', title: '保证睡眠', content: '尽量在 23 点前入睡，连续 3 天就能看到注意力的明显改善。' },
  { icon: 'Sunrise', color: '#e6a23c', title: '晒太阳', content: '每天 15 分钟自然光，对调节情绪和昼夜节律都有帮助。' },
  { icon: 'Bicycle', color: '#67c23a', title: '规律运动', content: '每周 3 次有氧运动，焦虑水平可以显著下降。' },
  { icon: 'Coffee', color: '#a777e3', title: '减少咖啡因', content: '下午 3 点后避免摄入咖啡因，能更快入睡。' },
  { icon: 'ChatLineRound', color: '#f56c6c', title: '主动倾诉', content: '难受时找一个信任的人聊一聊，比独自硬扛更有效。' },
  { icon: 'Document', color: '#909399', title: '情绪记录', content: '把烦恼写下来，能让大脑跳出循环思维。' },
]

onMounted(() => {
  loadMood()
  loadPomo()
})
onBeforeUnmount(() => {
  if (breathTimer) clearInterval(breathTimer)
  if (pomoTimer) clearInterval(pomoTimer)
})
</script>

<style scoped>
.hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(120deg, #6e8efb 0%, #a777e3 100%);
  color: #fff;
}
.hero .page-title,
.hero .hint {
  color: #fff;
}
.hero-icon {
  opacity: 0.7;
}

.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.tool-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 340px;
}

/* 心情 */
.mood-row {
  display: flex;
  justify-content: space-between;
  gap: 6px;
}
.mood-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
}
.mood-item:hover {
  background: #f5f7fa;
}
.mood-item.active {
  background: #ecf5ff;
  outline: 1px solid #409eff;
}
.mood-item .emoji {
  font-size: 28px;
}
.mood-item .label {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}
.mood-msg {
  background: #f8f9fc;
  border-radius: 6px;
  padding: 8px 12px;
  margin: 0;
  color: #606266;
  font-size: 13px;
}
.sub-title {
  margin: 6px 0;
  color: #606266;
  font-size: 13px;
}
.mood-history {
  display: flex;
  gap: 6px;
  justify-content: space-between;
}
.mood-cell {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 6px 0;
  background: #f8f9fc;
  border-radius: 6px;
}
.mood-cell .cell-emoji {
  font-size: 20px;
}
.mood-cell .cell-date {
  font-size: 11px;
  color: #909399;
}

/* 呼吸 */
.breath-stage {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 160px;
}
.breath-circle {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: linear-gradient(135deg, #91d5ff, #1890ff);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: #fff;
  font-weight: 600;
  transform-origin: center;
}
.breath-circle.inhale {
  animation: breath-grow linear forwards;
}
.breath-circle.hold {
  transform: scale(1.4);
  background: linear-gradient(135deg, #b7eb8f, #52c41a);
}
.breath-circle.exhale {
  animation: breath-shrink linear forwards;
  background: linear-gradient(135deg, #ffd591, #fa8c16);
}
@keyframes breath-grow {
  from { transform: scale(1); }
  to { transform: scale(1.4); }
}
@keyframes breath-shrink {
  from { transform: scale(1.4); }
  to { transform: scale(1); }
}
.breath-label {
  font-size: 16px;
}
.breath-count {
  font-size: 12px;
  opacity: 0.8;
}
.full-btn {
  width: 100%;
}
.counter {
  text-align: center;
  color: #909399;
  font-size: 12px;
  margin: 6px 0 0 0;
}

/* 番茄 */
.timer-display {
  text-align: center;
  font-size: 42px;
  font-weight: 600;
  letter-spacing: 2px;
  font-family: 'JetBrains Mono', monospace;
  color: #303133;
}
.timer-display.break {
  color: #67c23a;
}
.timer-state {
  text-align: center;
  margin: 0;
  color: #909399;
  font-size: 13px;
}
.timer-actions {
  display: flex;
  justify-content: center;
  gap: 8px;
}

/* 贴士 */
.tip-card .tip-item {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  padding: 12px;
  background: #f8f9fc;
  border-radius: 8px;
  margin-bottom: 10px;
}
.tip-text h4 {
  margin: 0 0 4px 0;
  font-size: 14px;
}
.tip-text p {
  margin: 0;
  color: #606266;
  font-size: 12px;
  line-height: 1.6;
}

/* 金句 */
.quote-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
  background: linear-gradient(135deg, #fff7e6 0%, #fff 100%);
  min-height: 240px;
}
.quote {
  font-size: 18px;
  line-height: 1.7;
  color: #303133;
  margin: 12px 0;
}
.quote-author {
  text-align: right;
  color: #909399;
  margin: 0 0 12px 0;
}
.hint {
  color: #909399;
  font-size: 13px;
  margin: 0;
}
</style>