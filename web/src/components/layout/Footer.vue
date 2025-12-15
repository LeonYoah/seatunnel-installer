<!--
  Footer 组件
  版权信息
-->
<template>
  <div class="footer">
    <p class="copyright">
      © {{ currentYear }} {{ t('footer.copyright') }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const currentYear = computed(() => new Date().getFullYear())
const { t } = useI18n()
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
