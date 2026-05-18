<template>
  <n-form
    ref="formRef"
    label-placement="left"
    size="large"
    :model="mode === 'pwdLogin' ? formPwdData : formCodeData"
    :rules="mode === 'pwdLogin' ? rules : mergeAccountRules"
  >
    <template v-if="mode === 'pwdLogin'">
      <n-form-item path="username">
        <n-input
          @keyup.enter="debouncePwdSubmit"
          v-model:value="formPwdData.username"
          placeholder="请输入邮箱/手机号"
        >
          <template #prefix>
            <n-icon size="18" color="#808695">
              <PersonOutline />
            </n-icon>
          </template>
        </n-input>
      </n-form-item>
      <n-form-item path="pass">
        <n-input
          @keyup.enter="debouncePwdSubmit"
          v-model:value="formPwdData.pass"
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
        <n-button type="primary" size="large" :block="true" @click="handlePwdLoginCaptcha">
          登录
        </n-button>
      </n-space>

      <n-modal
        v-model:show="showLoginCaptcha"
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
              @keyup.enter="debouncePwdSubmit"
              v-model:value="formPwdData.captcha"
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
            @click="handleLogin"
          >
            确定
          </n-button>
        </n-space>
      </n-modal>
    </template>
    <!-- 登录码登录-->
    <template v-if="mode === 'mergeCode'">
      <n-form-item path="mergeCode">
        <n-input
          @keyup.enter="handleCodeSubmit"
          v-model:value="formCodeData.account"
          placeholder="请输入账号/邮箱/手机号"
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
            @keyup.enter="handleCodeSubmit"
            v-model:value="formCodeData.code"
            placeholder="请输入验证码"
          >
            <template #prefix>
              <n-icon size="18" color="#808695" :component="SafetyCertificateOutlined" />
            </template>
          </n-input>
          <n-button
            type="primary"
            ghost
            @click="handleSendMergeAccountCode"
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
          <n-button :text="true" @click="handleResetPassword">忘记密码？</n-button>
        </div>
        <n-button type="primary" size="large" :block="true" :loading="loading" @click="handleLogin">
          登录
        </n-button>
      </n-space>

      <n-modal
        v-model:show="showSendCodeCaptcha"
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
              @keyup.enter="sendMergeAccountCode"
              v-model:value="formCodeData.captcha"
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
            @click="sendMergeAccountCode"
          >
            确定
          </n-button>
        </n-space>
      </n-modal>
    </template>

    <DemoAccount @login="handleDemoAccountLogin" />
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

  interface Props {
    mode: string;
  }

  const props = withDefaults(defineProps<Props>(), {
    mode: 'pwdLogin',
  });

  interface FormPwdState {
    username: string;
    pass: string;
    cid: string;
    captcha: string;
    password: string;
  }

  interface FormCodeState {
    account: string;
    code: string;
    cid: string;
    captcha: string;
  }

  const formRef = ref();
  const message = useMessage();
  const loading = ref(false);
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
  const debouncePwdSubmit = useDebounceFn((e) => {
    handleSubmit(e);
  }, 500);
  const formPwdData = ref<FormPwdState>({
    username: '',
    pass: '',
    password: '',
    cid: '',
    captcha: '',
  });

  const formCodeData = ref<FormCodeState>({
    account: '',
    code: '',
    cid: '',
    captcha: '',
  });

  const rules = {
    username: { required: true, message: '请输入用户名', trigger: 'blur' },
    pass: { required: true, message: '请输入密码', trigger: 'blur' },
  };

  const mergeAccountRules = {
    account: { required: true, message: '请输入合法的账号', trigger: 'blur' },
    code: { required: true, message: '请输入验证码', trigger: 'blur' },
  };
  // 判断当前登录类型
  const isPwdLogin = computed(() => {
    return props.mode === 'pwdLogin';
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    formRef.value.validate(async (errors) => {
      if (!errors) {
        if (userStore.loginConfig?.loginCaptchaSwitch === 1 && formPwdData.value.captcha === '') {
          message.error('请输入验证码');
          return;
        }

        const params = {
          username: formPwdData.value.username,
          password: aesEcb.encrypt(formPwdData.value.pass),
          cid: formPwdData.value.cid,
          captcha: formPwdData.value.captcha,
        };
        await handleLoginResp(userStore.login(params));
      } else {
        message.error('请填写完整信息，并且进行验证码校验');
      }
    });
  };

  // 刷新图形验证码
  async function refreshCode() {
    if (userStore.loginConfig?.loginCaptchaSwitch !== 1) {
      return;
    }
    loadingBar.start();
    const data = await GetCaptcha();
    codeBase64.value = data.base64;
    if (isPwdLogin.value) {
      formPwdData.value.cid = data.cid;
      formPwdData.value.captcha = '';
    } else {
      formCodeData.value.cid = data.cid;
      formCodeData.value.captcha = '';
    }
    loadingBar.finish();
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

  // 验证码登录
  const handleCodeSubmit = (e) => {
    e.preventDefault();
    formRef.value.validate(async (errors) => {
      if (!errors) {
        const params = {
          account: formCodeData.value.account,
          cid: formCodeData.value.cid,
          captcha: formCodeData.value.captcha,
          code: formCodeData.value.code,
        };
        await handleLoginResp(userStore.login(params));
      } else {
        message.error('请填写完整信息，并且进行验证码校验');
      }
    });
  };
  // 是否弹出手机验证码图片对话框
  const showSendCodeCaptcha = ref(false);

  //  弹出发送联合账号验证码的图片验证码
  function handleSendMergeAccountCode() {
    console.log('弹出发送联合账号验证码的图片验证码', formCodeData.value);
    validate.mergeAccount(
      mergeAccountRules.account,
      formCodeData.value.account,
      function (error?: Error) {
        if (error === undefined) {
          showSendCodeCaptcha.value = true;
          return;
        }
        message.error(error.message);
      }
    );
  }

  // 发送验证码
  function sendMergeAccountCode() {
    validate.mergeAccount(
      mergeAccountRules.account,
      formCodeData.value.account,
      function (error?: Error) {
        showSendCodeCaptcha.value = false;
        if (error === undefined) {
          activateSend(
            SendLoginCode({
              account: formCodeData.value.account,
              cid: formCodeData.value.cid,
              captcha: formCodeData.value.captcha,
              event: 'login',
            })
          );
          return;
        }
        message.error(error.message);
      }
    );
  }

  function handleResetPassword() {
    message.info('如果你忘记了密码，请联系管理员找回');
  }

  // 密码登录
  const showLoginCaptcha = ref(false);

  function handlePwdLoginCaptcha() {
    const username = formPwdData.value.username;
    if (username.trim() == '') {
      message.error('请输入你的账号');
      return;
    }
    const password = formPwdData.value.pass;
    if (password.trim() == '') {
      message.error('请输入你的密码');
      return;
    }
    showLoginCaptcha.value = true;
  }

  function handleLogin(e) {
    if (isPwdLogin.value) {
      showLoginCaptcha.value = false;
      debouncePwdSubmit(e);
      return;
    }
    showSendCodeCaptcha.value = false;
    handleCodeSubmit(e);
  }

  async function handleLoginResp(request: Promise<any>) {
    message.loading('登录中...');
    loading.value = true;
    try {
      const { code, message: msg } = await request;
      message.destroyAll();
      if (code == ResultEnum.SUCCESS) {
        const toPath = decodeURIComponent((route.query?.redirect || '/') as string);
        message.success('登录成功，即将进入系统');
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

  onMounted(() => {});
</script>
