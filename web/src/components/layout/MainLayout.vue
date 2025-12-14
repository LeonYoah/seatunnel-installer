<!--
  主布局组件
  包含 Header、Sidebar、Main、Footer
-->
<template>
  <div class="main-layout">
    <Header class="layout-header" />
    <div class="layout-body">
      <Sidebar class="layout-sidebar" />
      <div class="layout-main">
        <div class="main-content">
          <router-view />
        </div>
        <Footer class="layout-footer" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Header from './Header.vue'
import Sidebar from './Sidebar.vue'
import Footer from './Footer.vue'
</script>

<style scoped>
.main-layout {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg);
}

.layout-header {
  flex-shrink: 0;
}

.layout-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.layout-sidebar {
  flex-shrink: 0;
}

.layout-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.main-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background: var(--bg);
}

.layout-footer {
  flex-shrink: 0;
}
</style>
