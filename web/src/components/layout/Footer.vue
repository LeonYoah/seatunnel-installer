<!--
  Footer 组件
  版权信息
-->
<template>
  <div class="footer">
    <p class="copyright">
      © {{ currentYear }} SeaTunnel 企业级管理平台 | Apache SeaTunnel
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const currentYear = computed(() => new Date().getFullYear())
</script>

<style scoped>
.footer {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--surface);
  border-top: 1px solid var(--border);
}

.copyright {
  margin: 0;
  font-size: 14px;
  color: var(--muted);
}
</style>
