<template>
  <div>
    <div class="page-card">
      <h2 class="page-title">帮助中心</h2>
      <p class="hint">常见问题集合。找不到答案？联系系统管理员。</p>
    </div>

    <div class="page-card" style="margin-top: 16px">
      <el-input
        v-model="keyword"
        placeholder="搜索问题（如：测评、密码、导入）"
        clearable
        :prefix-icon="Search"
        style="margin-bottom: 16px"
      />

      <el-tabs v-model="tab">
        <el-tab-pane
          v-for="g in filteredGroups"
          :key="g.name"
          :label="`${g.name} (${g.items.length})`"
          :name="g.name"
        >
          <el-collapse accordion>
            <el-collapse-item
              v-for="(q, i) in g.items"
              :key="`${g.name}-${i}`"
              :title="q.q"
              :name="`${g.name}-${i}`"
            >
              <p class="answer">{{ q.a }}</p>
            </el-collapse-item>
          </el-collapse>
          <el-empty v-if="g.items.length === 0" description="该分类暂无匹配问题" />
        </el-tab-pane>
      </el-tabs>
    </div>

    <div class="page-card" style="margin-top: 16px">
      <h3 class="card-title">还没找到答案？</h3>
      <p>联系系统管理员，或参考 <router-link to="/about" style="color: #409eff">关于平台</router-link> 了解更多信息。</p>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { Search } from '@element-plus/icons-vue'

const tab = ref('账号')
const keyword = ref('')

const GROUPS = [
  {
    name: '账号',
    items: [
      { q: '如何注册账号？', a: '在登录页点击「注册新账号」，输入手机号、用户名、密码及手机验证码即可。注册后默认是学生角色，如需切换需联系管理员。' },
      { q: '忘记密码怎么办？', a: '在登录页点击「忘记密码」，使用注册手机号接收验证码后即可重置。' },
      { q: '验证码收不到？', a: '开发环境下短信网关可能不可达，验证码会打印到 Go 服务的控制台日志，可直接读取。' },
      { q: '如何完善个人资料？', a: '登录后点击顶栏头像 → 完善资料，可修改用户名、性别、邮箱。手机号绑定后不可修改。' },
    ],
  },
  {
    name: '学生',
    items: [
      { q: '如何提交心理测评？', a: '左侧菜单 →「心理测评」，填写焦虑/动力/情绪等指标后提交即可。建议每周或在压力波动时各做一次。' },
      { q: '压力等级是怎么计算的？', a: '系统从行为/学业/心理三个维度采集数据，先用随机森林模型给出基础等级，再用规则修正（如心理 ≥ 85 直接判为高压力）。' },
      { q: '我的成绩怎么没显示？', a: '成绩由管理员通过 Excel 批量导入。如果没看到记录，可能是本学期成绩还没入库，可联系管理员。' },
      { q: '收到干预建议后要怎么做？', a: '建议根据内容尝试调整，并将效果反馈给老师。后续可以在干预记录中查看历史。' },
    ],
  },
  {
    name: '教师',
    items: [
      { q: '怎么查看班级整体压力？', a: '左侧菜单 →「班级压力分布」，会展示等级分布饼图、维度均分柱状图，以及学生压力排行。' },
      { q: '怎么查看某个学生的详细情况？', a: '在班级压力分布的排行表中，点击预警学生右侧的「压力诊断书」按钮，可查看雷达图、趋势和干预建议。' },
      { q: '诊断书里的「建议干预方向」从哪里来？', a: '从管理员维护的「干预建议库」中按等级匹配，最多 6 条。如果对应等级没有模板，会显示空。' },
      { q: '能给学生推送建议吗？', a: '当前版本通过管理员后台统一推送，教师暂只读。' },
    ],
  },
  {
    name: '管理员',
    items: [
      { q: '如何批量导入学生成绩？', a: '管理员 →「Excel 导入」，上传 .xlsx / .xls 文件即可。导入操作人会从当前登录的 JWT 自动取，不需要手填。' },
      { q: '如何运行压力预测？', a: '管理员 →「运行预测」，填好 14 个特征后提交。系统会先建特征快照，再调用 Python 随机森林服务，最后写回评估结果表。' },
      { q: '需要启动哪些服务？', a: 'Go 后端（:8080）+ Python 随机森林服务（:8000）+ MySQL。Python 服务在「运行预测」时必须在线。' },
      { q: '怎么修改用户角色？', a: '管理员 →「用户管理」，在目标行点击「修改角色」按钮，可设为 学生 / 教师 / 管理员。' },
    ],
  },
  {
    name: '系统',
    items: [
      { q: '前后端如何联调？', a: '前端通过 Vite 反向代理把 /api 转发到 :8080。生产环境可改 .env.production 中的 VITE_BACKEND_URL。' },
      { q: '数据保存在哪里？', a: '业务数据在 MySQL，会话密钥（JWT）在前端 localStorage。验证码默认走内存，没装 Redis 也能正常用。' },
      { q: '想清空本机偏好设置？', a: '设置 →「偏好设置」→ 重置默认 即可恢复默认主题/字号。' },
    ],
  },
]

const filteredGroups = computed(() => {
  if (!keyword.value) return GROUPS
  const kw = keyword.value.toLowerCase()
  return GROUPS.map(g => ({
    ...g,
    items: g.items.filter(it => it.q.toLowerCase().includes(kw) || it.a.toLowerCase().includes(kw)),
  }))
})
</script>

<style scoped>
.hint {
  color: #909399;
  margin: 6px 0 0 0;
}
.card-title {
  margin: 0 0 8px 0;
}
.answer {
  color: #606266;
  line-height: 1.7;
  margin: 0;
}
</style>