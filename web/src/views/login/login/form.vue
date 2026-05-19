<template>
  <n-form
    ref="formRef"
    label-placement="left"
    size="large"
    :model="mode === 'pwdLogin' ? formInPassword : formCode"
    :rules="mode === 'pwdLogin' ? passwordRules : codeRules"
  >
    <!-- 密码登录 -->
    <template v-if="mode === 'pwdLogin'">
      <n-form-item path="username">
        <n-input
          @keyup.enter="debouncePasswordSubmit"
          v-model:value="formInPassword.username"
          placeholder="请输入邮箱/手机号"
        >
          <template #prefix>
            <n-icon size="18" color="#808695">
              <PersonOutline />
            </n-icon>
          </template>
        </n-input>
      </n-form-item>
      <n-form-item path="password">
        <n-input
          @keyup.enter="debouncePasswordSubmit"
          v-model:value="formInPassword.password"
          type="password"
          show-password-on="click"
          placeholder="请输入登录密码"
        >
          <template #prefix>
            <n-icon size="18" color="#808695">
              <LockClosedOutline />
            </n-icon>
          </template>
        </n-input>
      </n-form-item>

      <n-space :vertical="true" :size="24">
        <div class="flex-y-center justify-between">
          <n-checkbox v-model:checked="autoLogin">自动登录</n-checkbox>
          <n-button :text="true" @click="handleResetPassword">忘记密码？</n-button>
        </div>
        <n-button type="primary" size="large" :block="true" @click="handlePasswordLoginModal">
          登录
        </n-button>
      </n-space>

      <n-modal
        v-model:show="showPasswordLoginModal"
        :show-icon="false"
        preset="dialog"
        title="验证码"
        style="max-width: 350px; bottom: 100px"
      >
        <n-form-item path="code" v-show="codeBase64 !== ''">
          <n-input-group>
            <n-input
              :style="{ width: '100%' }"
              placeholder="请输入验证码"
              @keyup.enter="debouncePasswordSubmit"
              v-model:value="formInPassword.captcha"
            >
              <template #prefix>
                <n-icon size="18" color="#808695" :component="SafetyCertificateOutlined" />
              </template>
              <template #suffix></template>
            </n-input>

            <n-loading-bar-provider :to="loadingBarTargetRef" container-style="position: absolute;">
              <img
                ref="loadingBarTargetRef"
                style="width: 100px"
                :src="codeBase64"
                @click="refreshCode"
                loading="lazy"
                alt="点击获取"
              />
              <loading-bar-trigger />
            </n-loading-bar-provider>
          </n-input-group>
        </n-form-item>
        <n-space :vertical="true" :size="24">
          <n-button
            type="primary"
            size="large"
            :block="true"
            :loading="loading"
            @click="handlePasswordSubmit"
          >
            确定
          </n-button>
        </n-space>
      </n-modal>
    </template>
    <!-- 验证码登录-->
    <template v-if="mode === 'mergeCode'">
      <n-form-item path="mergeCode">
        <n-input
          @keyup.enter="debounceCodeSubmit"
          v-model:value="formCode.account"
          placeholder="请输入邮箱/手机号"
        >
          <template #prefix>
            <n-icon size="18" color="#808695">
              <PersonOutline />
            </n-icon>
          </template>
        </n-input>
      </n-form-item>

      <n-form-item path="code">
        <n-input-group>
          <n-input
            @keyup.enter="debounceCodeSubmit"
            v-model:value="formCode.code"
            placeholder="请输入验证码"
          >
            <template #prefix>
              <n-icon size="18" color="#808695" :component="SafetyCertificateOutlined" />
            </template>
          </n-input>
          <n-button
            type="primary"
            ghost
            @click="handleCodeLoginModal"
            :disabled="isCounting"
            :loading="sendLoading"
          >
            {{ sendLabel }}
          </n-button>
        </n-input-group>
      </n-form-item>
      <n-space :vertical="true" :size="24">
        <div class="flex-y-center justify-between">
          <n-checkbox v-model:checked="autoLogin">自动登录</n-checkbox>
          <n-checkbox v-if="isRuningDev" v-model:checked="mockAccountCode">MockCode</n-checkbox>
          <n-button :text="true" @click="handleResetPassword">忘记密码？</n-button>
        </div>
        <n-button
          type="primary"
          size="large"
          :block="true"
          :loading="loading"
          @click="handleCodeSubmit"
        >
          登录
        </n-button>
      </n-space>

      <n-modal
        v-model:show="showCodeLoginModal"
        :show-icon="false"
        preset="dialog"
        title="验证码"
        style="max-width: 350px; bottom: 100px"
      >
        <n-form-item path="code" v-show="codeBase64 !== ''">
          <n-input-group>
            <n-input
              :style="{ width: '100%' }"
              placeholder="请输入验证码"
              @keyup.enter="sendAccountLoginCode"
              v-model:value="formCode.captcha"
            >
              <template #prefix>
                <n-icon size="18" color="#808695" :component="SafetyCertificateOutlined" />
              </template>
              <template #suffix></template>
            </n-input>

            <n-loading-bar-provider :to="loadingBarTargetRef" container-style="position: absolute;">
              <img
                ref="loadingBarTargetRef"
                style="width: 100px"
                :src="codeBase64"
                @click="refreshCode"
                loading="lazy"
                alt="点击获取"
              />
              <loading-bar-trigger />
            </n-loading-bar-provider>
          </n-input-group>
        </n-form-item>
        <n-space :vertical="true" :size="24">
          <n-button
            type="primary"
            size="large"
            :block="true"
            :loading="loading"
            @click="sendAccountLoginCode"
          >
            发送验证码
          </n-button>
        </n-space>
      </n-modal>
    </template>
    <!-- 演示角色登录 -->
    <DemoAccount v-if="isRuningDev" @login="handleDemoAccountLogin" />
  </n-form>
</template>

<script lang="ts" setup>
  import '../components/style.less';
  import { ref, onMounted, computed } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { useUserStore } from '@/store/modules/user';
  import { useMessage, useLoadingBar } from 'naive-ui';
  import { ResultEnum } from '@/enums/httpEnum';
  import { PersonOutline, LockClosedOutline } from '@vicons/ionicons5';
  import { PageEnum } from '@/enums/pageEnum';
  import { SafetyCertificateOutlined } from '@vicons/antd';
  import { GetCaptcha } from '@/api/base';
  import { aesEcb } from '@/utils/encrypt';
  import DemoAccount from './demo-account.vue';
  import { useSendCode } from '@/hooks/common';
  import { SendLoginCode } from '@/api/system/user';
  import { validate } from '@/utils/validateUtil';
  import { useDebounceFn } from '@vueuse/core';
  import { createSocket } from '@/utils/websocket';
  import { isDevMode } from '@/utils/env';

  interface Props {
    mode: string;
  }

  const props = withDefaults(defineProps<Props>(), {
    mode: 'pwdLogin',
  });

  interface FormState {
    username: string;
    password: string;
    cid: string;
    captcha: string;
  }

  interface FormCodeLogicState {
    account: string;
    code: string;
    cid: string;
    captcha: string;
  }

  const formRef = ref();
  const message = useMessage();
  const loading = ref(false);
  const isRuningDev = ref(isDevMode());
  const mockAccountCode = ref(true);
  const autoLogin = ref(true);
  const codeBase64 = ref('');
  const loadingBar = useLoadingBar();
  const loadingBarTargetRef = ref<undefined | HTMLElement>(undefined);
  const userStore = useUserStore();
  const router = useRouter();
  const route = useRoute();
  const { sendLabel, isCounting, loading: sendLoading, activateSend } = useSendCode();
  const emit = defineEmits(['updateActiveModule']);
  const LOGIN_NAME = PageEnum.BASE_LOGIN_NAME;

  // 判断当前登录类型
  const isPwdLogin = computed(() => {
    return props.mode === 'pwdLogin';
  });

  // 刷新图形验证码
  async function refreshCode() {
    if (userStore.loginConfig?.loginCaptchaSwitch !== 1) {
      return;
    }
    loadingBar.start();
    const data = await GetCaptcha();
    codeBase64.value = data.base64;
    if (isPwdLogin.value) {
      formInPassword.value.cid = data.cid;
      formInPassword.value.captcha = '';
    } else {
      formCode.value.cid = data.cid;
      formCode.value.captcha = '';
    }
    loadingBar.finish();
  }

  function handleResetPassword() {
    message.info('如果你忘记了密码，请使用验证码登录');
  }

  // 演示角色登录
  async function handleDemoAccountLogin(user: {
    username: string;
    password: string;
    captcha: string;
  }) {
    const params = {
      username: user.username,
      captcha: user.captcha,
      password: aesEcb.encrypt(user.password),
      isLock: true,
    };
    await handleLoginResp(userStore.login(params));
  }

  // 密码登录表单
  const formInPassword = ref<FormState>({
    username: '',
    password: '',
    cid: '',
    captcha: '',
  });
  const passwordRules = {
    username: { required: true, message: '请输入用户名', trigger: 'blur' },
    password: { required: true, message: '请输入密码', trigger: 'blur' },
    captcha: { required: true, message: '请输入图片验证码' },
  };

  // 密码登录
  const showPasswordLoginModal = ref(false);

  function handlePasswordLoginModal() {
    const username = formInPassword.value.username;
    if (username.trim() == '') {
      message.error('请输入你的账号');
      return;
    }
    const password = formInPassword.value.password;
    if (password.trim() == '') {
      message.error('请输入你的密码');
      return;
    }
    refreshCode();
    showPasswordLoginModal.value = true;
  }

  const handlePasswordSubmit = (e) => {
    e.preventDefault();
    showPasswordLoginModal.value = false;
    formRef.value.validate(async (errors) => {
      if (!errors) {
        if (
          userStore.loginConfig?.loginCaptchaSwitch === 1 &&
          formInPassword.value.captcha === ''
        ) {
          message.error('请输入验证码');
          return;
        }

        const params = {
          username: formInPassword.value.username,
          password: aesEcb.encrypt(formInPassword.value.password),
          cid: formInPassword.value.cid,
          captcha: formInPassword.value.captcha,
        };
        await handleLoginResp(userStore.login(params));
      } else {
        message.error('请填写完整信息，并且进行验证码校验');
      }
    });
  };
  const debouncePasswordSubmit = useDebounceFn((e) => {
    handlePasswordSubmit(e);
  }, 500);

  // 验证码登录表单
  const formCode = ref<FormCodeLogicState>({
    account: '',
    code: '',
    cid: '',
    captcha: '',
  });
  const codeRules = {
    account: { required: true, message: '请输入电子邮箱/手机号码', trigger: 'blur' },
    code: { required: true, message: '请输入验证码', trigger: 'blur' },
    captcha: { required: true, message: '请输入图片验证码' },
  };
  const showCodeLoginModal = ref(false);

  function handleCodeLoginModal() {
    refreshCode();
    validate.mergeAccount(codeRules.account, formCode.value.account, function (error?: Error) {
      if (error === undefined) {
        showCodeLoginModal.value = true;
        return;
      }
      message.error(error.message);
    });
  }

  // 发送账号验证码
  function sendAccountLoginCode() {
    validate.mergeAccount(codeRules.account, formCode.value.account, function (error?: Error) {
      showCodeLoginModal.value = false;
      if (error === undefined) {
        activateSend(
          SendLoginCode({
            account: formCode.value.account,
            cid: formCode.value.cid,
            captcha: formCode.value.captcha,
            event: 'login',
          })
        );
        return;
      }
      message.error(error.message);
    });
  }

  // 验证码登录
  const handleCodeSubmit = (e) => {
    e.preventDefault();
    showCodeLoginModal.value = false;
    formRef.value.validate(async (errors) => {
      if (!errors) {
        const params = {
          account: formCode.value.account,
          cid: formCode.value.cid,
          captcha: formCode.value.captcha,
          code: formCode.value.code,
        };
        await handleLoginResp(userStore.login(params));
      } else {
        message.error('请填写完整信息，并且进行验证码校验');
      }
    });
  };
  const debounceCodeSubmit = useDebounceFn((e) => {
    handleCodeSubmit(e);
  }, 500);

  async function handleLoginResp(request: Promise<any>) {
    message.loading('登录中...');
    loading.value = true;
    try {
      const { code, message: msg } = await request;
      message.destroyAll();
      if (code == ResultEnum.SUCCESS) {
        const toPath = decodeURIComponent((route.query?.redirect || '/') as string);
        message.success('登录成功，即将进入系统');
        createSocket();
        if (route.name === LOGIN_NAME) {
          await router.replace('/');
        } else {
          await router.replace(toPath);
        }
      } else {
        message.destroyAll();
        message.error(msg || '登录失败');
        await refreshCode();
      }
    } finally {
      loading.value = false;
    }
  }

  onMounted(() => {
    // setTimeout(function () {
    //   refreshCode();
    // });
  });
</script>
