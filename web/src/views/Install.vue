<!--
  安装向导页面 - 参考 www/index.html
-->
<template>
  <div class="install">
    <!-- 步骤导航 -->
    <el-card class="steps-card">
      <el-steps :active="currentStep" align-center>
        <el-step :title="t('install.steps.config')" />
        <el-step :title="t('install.steps.precheck')" />
        <el-step :title="t('install.steps.plugins')" />
        <el-step :title="t('install.steps.install')" />
        <el-step :title="t('install.steps.complete')" />
      </el-steps>
    </el-card>

    <!-- 步骤内容 -->
    <component :is="currentStepComponent" @next="handleNext" @prev="handlePrev" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import StepConfig from './install/StepConfig.vue'
import StepPrecheck from './install/StepPrecheck.vue'
import StepPlugins from './install/StepPlugins.vue'
import StepInstall from './install/StepInstall.vue'
import StepComplete from './install/StepComplete.vue'
import { useI18n } from 'vue-i18n'

const currentStep = ref(0)
const { t } = useI18n()

const steps = [StepConfig, StepPrecheck, StepPlugins, StepInstall, StepComplete]

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
.install {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

.steps-card {
  margin-bottom: 20px;
}
</style>
