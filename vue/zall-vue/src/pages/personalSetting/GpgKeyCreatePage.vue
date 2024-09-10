<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">{{t('sshGpg.createGpg')}}</div>
      <div class="section">
        <div class="section-title">{{t('sshGpg.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('sshGpg.content')}}</div>
        <div class="section-body">
          <a-textarea
            style="width:100%"
            v-model:value="formState.content"
            :auto-size="{ minRows: 3, maxRows: 8 }"
          />
          <div class="input-desc">{{t('sshGpg.gpgContentDesc')}}</div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveGpgKey">{{t('sshGpg.save')}}</a-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive } from "vue";
import { createGpgKeyRequest } from "@/api/user/gpgKeyApi";
import { gpgKeyNameRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
// 表单数据
const formState = reactive({
  name: "",
  content: ""
});
// 点击“立即新增”按钮
const saveGpgKey = () => {
  if (!gpgKeyNameRegexp.test(formState.name)) {
    message.warn(t("sshGpg.nameFormatErr"));
    return;
  }
  if (!formState.content) {
    message.warn(t("sshGpg.pleaseFillContent"));
    return;
  }
  createGpgKeyRequest({
    name: formState.name,
    content: formState.content
  }).then(() => {
    message.success(t("operationSuccess"));
    router.push("/personalSetting/sshAndGpg/list");
  });
};
</script>

<style scoped>
</style>