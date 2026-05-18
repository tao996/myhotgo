<template>
  <n-tabs type="segment" justify-content="space-evenly">
    <n-tab-pane name="pwdLogin" tab="密码登录">
      <Form @updateActiveModule="updateActiveModule" mode="pwdLogin" />
    </n-tab-pane>
    <n-tab-pane name="mergeCode" tab="验证码登录">
      <Form @updateActiveModule="updateActiveModule" mode="mergeCode" />
    </n-tab-pane>
  </n-tabs>
</template>

<script lang="ts" setup>
  import Form from './form.vue';

  const emit = defineEmits(['updateActiveModule']);

  function updateActiveModule(key: string) {
    emit('updateActiveModule', key);
  }
</script>

<style lang="less" scoped></style>
