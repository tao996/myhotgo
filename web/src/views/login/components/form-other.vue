<template>
  <n-form-item class="default-color">
    <div class="flex view-account-other">
      <div class="flex-initial mx-2">
        <a @click="handleLoginWechat">
          <n-icon size="24" color="rgb(24, 160, 88)">
            <LogoWechat />
          </n-icon>
        </a>
      </div>
      <div class="flex-initial mx-2">
        <a @click="handleLogoTiktok">
          <n-icon size="24" color="rgba(25, 28, 34, 0.88)">
            <LogoTiktok />
          </n-icon>
        </a>
      </div>
    </div>
  </n-form-item>
</template>

<script lang="ts" setup>
  import { LogoTiktok, LogoWechat } from '@vicons/ionicons5';
  import { useMessage } from 'naive-ui';

  interface Props {
    moduleKey: string;
    tag: string;
  }

  withDefaults(defineProps<Props>(), {
    moduleKey: 'register',
    tag: '注册账号',
  });

  const message = useMessage();
  const emit = defineEmits(['updateActiveModule']);

  function handleLogoTiktok() {
    console.log('handleLogoTiktok...');
    message.info('暂未开放');
  }

  function handleLoginWechat() {
    console.log('handleLoginWechat...');
    message.info('暂未开放');
  }
</script>

<style scoped></style>
