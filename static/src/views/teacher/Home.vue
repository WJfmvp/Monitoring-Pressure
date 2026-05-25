<template>
  <div>
    <div class="page-card">
      <h2 class="page-title">教师工作台</h2>
      <p class="hint">
        推荐先进入「班级压力分布」获取整体视图，对预警学生点击查看「压力诊断书」。
      </p>
    </div>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="8">
        <el-card shadow="hover" class="entry primary">
          <h3>班级压力分布</h3>
          <p>等级分布、维度均分、预警排行</p>
          <el-button type="primary" @click="$router.push('/teacher/class/overview')">
            前往班级总览
          </el-button>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="entry">
          <h3>学生数据查询</h3>
          <p>输入学生 ID 查询其学业 / 压力 / 干预记录</p>
          <el-button @click="go('/teacher/student/academic')">学业成绩</el-button>
          <el-button @click="go('/teacher/student/stress')">压力评估</el-button>
          <el-button @click="go('/teacher/student/intervention')">干预记录</el-button>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="entry">
          <h3>快捷诊断</h3>
          <p>直接输入学号查看压力诊断书</p>
          <div class="quick-input">
            <el-input v-model="studentId" placeholder="学生 ID" />
            <el-button type="primary" :disabled="!studentId" @click="goDiagnosis">
              查看诊断书
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const studentId = ref('')

function go(path) {
  router.push({ path, query: studentId.value ? { student_id: studentId.value } : {} })
}

function goDiagnosis() {
  router.push({ path: '/teacher/student/diagnosis', query: { student_id: studentId.value } })
}
</script>

<style scoped>
.hint {
  color: #909399;
  margin: 0;
}
.entry h3 {
  margin: 0 0 6px 0;
}
.entry p {
  color: #909399;
  margin: 0 0 12px 0;
}
.entry.primary {
  border-top: 3px solid #409eff;
}
.quick-input {
  display: flex;
  gap: 8px;
}
.quick-input .el-input {
  flex: 1;
}
</style>