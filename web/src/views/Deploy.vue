<!--
  部署集群页面
  从已注册的主机中选择节点部署 SeaTunnel 集群
-->
<template>
  <div class="deploy">
    <!-- 步骤导航 -->
    <el-card class="steps-card">
      <el-steps :active="currentStep" align-center>
        <el-step title="选择主机" />
        <el-step title="配置参数" />
        <el-step title="环境检查" />
        <el-step title="插件选择" />
        <el-step title="开始部署" />
        <el-step title="完成" />
      </el-steps>
    </el-card>

    <!-- 步骤内容 -->
    <component :is="currentStepComponent" @next="handleNext" @prev="handlePrev" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import StepSelectHosts from './deploy/StepSelectHosts.vue'
import StepConfig from './install/StepConfig.vue'
import StepPrecheck from './install/StepPrecheck.vue'
import StepPlugins from './install/StepPlugins.vue'
import StepInstall from './install/StepInstall.vue'
import StepComplete from './install/StepComplete.vue'

const currentStep = ref(0)

const steps = [StepSelectHosts, StepConfig, StepPrecheck, StepPlugins, StepInstall, StepComplete]

const currentStepComponent = computed(() => steps[currentStep.value])

const handleNext = () => {
  if (currentStep.value < steps.length - 1) {
    currentStep.value++
  }
}

const handlePrev = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}
</script>

<style scoped>
.deploy {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

.steps-card {
  margin-bottom: 20px;
}
</style>
