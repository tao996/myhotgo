<template>
  <div>
    <n-spin :show="show" description="请稍候...">
      <n-form :label-width="80" :model="formValue" :rules="rules" ref="formRef">
        <n-layout has-sider>
          <n-layout-sider width="120">
            <FileChooser v-model:value="formValue.contactWorkQrcode" file-type="image" />
          </n-layout-sider>
          <n-layout>
            <n-form-item label="企业微信" path="contactWork" :show-feedback="false">
              <n-input v-model:value="formValue.contactWork" placeholder="请输入企业微信" />
            </n-form-item>
            <n-form-item label="二维码地址" path="contactWorkQrcode" :show-feedback="false">
              <n-input
                v-model:value="formValue.contactWorkQrcode"
                placeholder="请上传或者输入企业微信二维码地址"
              />
            </n-form-item>
          </n-layout>
        </n-layout>

        <n-layout has-sider>
          <n-layout-sider width="120">
            <FileChooser v-model:value="formValue.contactWechatQrcode" file-type="image" />
          </n-layout-sider>
          <n-layout>
            <n-form-item label="公众号" path="contactWechat" :show-feedback="false">
              <n-input v-model:value="formValue.contactWechat" placeholder="请输入微信公众号" />
            </n-form-item>
            <n-form-item label="二维码地址" path="contactWechatQrcode" :show-feedback="false">
              <n-input
                v-model:value="formValue.contactWechatQrcode"
                placeholder="请上传或者输入公众号二维码地址"
              />
            </n-form-item>
          </n-layout>
        </n-layout>

        <n-form-item
          label="联系邮箱"
          path="contactEmail"
          :show-feedback="false"
          label-placement="left"
          label-align="left"
        >
          <n-input v-model:value="formValue.contactEmail" placeholder="请输入电子邮箱地址" />
        </n-form-item>
        <n-form-item
          label="联系电话"
          path="contactTel"
          :show-feedback="false"
          label-placement="left"
          label-align="left"
        >
          <n-input v-model:value="formValue.contactTel" placeholder="请输入联系电话" />
        </n-form-item>
        <n-form-item
          label="传真号码"
          path="contactFax"
          :show-feedback="false"
          label-placement="left"
          label-align="left"
        >
          <n-input v-model:value="formValue.contactFax" placeholder="请输入传真号码" />
        </n-form-item>

        <n-form-item
          label="QQ"
          path="contactQQ"
          :show-feedback="false"
          label-placement="left"
          label-align="left"
        >
          <n-input
            v-model:value="formValue.contactQQ"
            placeholder="请输入QQ群，多个Q号使用空格分割"
          />
        </n-form-item>

        <n-form-item
          label="微博"
          path="contactWeibo"
          :show-feedback="false"
          label-placement="left"
          label-align="left"
        >
          <n-input v-model:value="formValue.contactWeibo" placeholder="请输入微博账号" />
        </n-form-item>

        <n-form-item
          label="抖音"
          path="contactDouyin"
          :show-feedback="false"
          label-placement="left"
          label-align="left"
        >
          <n-input v-model:value="formValue.contactDouyin" placeholder="请输入抖音账号" />
        </n-form-item>

        <n-form-item
          label="Youtube"
          path="contactYoutube"
          :show-feedback="false"
          label-placement="left"
          label-align="left"
        >
          <n-input v-model:value="formValue.contactYoutube" placeholder="请输入 Youtube 账号" />
        </n-form-item>

        <n-form-item
          label="Facebook"
          path="contactFacebook"
          :show-feedback="false"
          label-placement="left"
          label-align="left"
        >
          <n-input v-model:value="formValue.contactFacebook" placeholder="请输入 Facebook" />
        </n-form-item>

        <n-form-item
          label="Twitter"
          path="contactTwitter"
          :show-feedback="false"
          label-placement="left"
          label-align="left"
        >
          <n-input v-model:value="formValue.contactTwitter" placeholder="请输入 Twitter" />
        </n-form-item>

        <n-form-item
          label="Telegram"
          path="contactTelegram"
          :show-feedback="false"
          label-placement="left"
          label-align="left"
        >
          <n-input v-model:value="formValue.contactTelegram" placeholder="请输入 Telegram" />
        </n-form-item>

        <n-form-item label="联系方式" path="loginProtocol">
          <Editor
            style="height: 320px"
            v-model:value="formValue.contactText"
            id="contactProtocol"
          />
        </n-form-item>

        <div>
          <n-space>
            <n-button type="primary" @click="formSubmit">保存更新</n-button>
          </n-space>
        </div>
      </n-form>
    </n-spin>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { useMessage } from 'naive-ui';
  import { getConfig, updateConfig } from '@/api/sys/config';
  import FileChooser from '@/components/FileChooser/index.vue';
  import Editor from '@/components/Editor/editor.vue';

  const group = ref('contact');
  const show = ref(false);
  const formRef: any = ref(null);
  const message = useMessage();

  const formValue = ref({
    contactEmail: '',
    contactTel: '',
    contactFax: '',
    contactWechat: '',
    contactWechatQrcode: '',
    contactWork: '',
    contactWorkQrcode: '',
    contactQQ: '',
    contactWeibo: '',
    contactDouyin: '',
    contactYoutube: '',
    contactFacebook: '',
    contactTwitter: '',
    contactWhatsapp: '',
    contactTelegram: '',
    contactText: '',
  });

  const rules = {};

  function formSubmit() {
    formRef.value.validate((errors) => {
      if (!errors) {
        updateConfig({ group: group.value, list: formValue.value }).then((_res) => {
          message.success('更新成功');
          load();
        });
      } else {
        message.error('验证失败，请填写完整信息');
      }
    });
  }

  onMounted(() => {
    load();
  });

  function load() {
    show.value = true;
    new Promise((_resolve, _reject) => {
      getConfig({ group: group.value })
        .then((res) => {
          console.log('ContactSetting reponse', res);
          if (res.list) {
            formValue.value = res.list;
          }
        })
        .finally(() => {
          show.value = false;
        });
    });
  }
</script>
