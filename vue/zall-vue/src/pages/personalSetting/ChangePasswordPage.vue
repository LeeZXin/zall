<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">{{t('changePassword.title')}}</div>
      <div class="item">
        <div class="input-title">{{t('changePassword.oldPwd')}}</div>
        <a-input-password v-model:value="formState.origin" />
      </div>
      <div class="item">
        <div class="input-title">{{t('changePassword.newPwd')}}</div>
        <a-input-password v-model:value="formState.password" />
      </div>
      <div class="item">
        <div class="input-title">{{t('changePassword.confirmPwd')}}</div>
        <a-input-password v-model:value="formState.confirm" />
      </div>
      <div>
        <a-button type="primary" @click="updatePassword">{{t('changePassword.save')}}</a-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { updatePasswordRequest } from "@/api/user/userApi";
import { reactive } from "vue";
import { passwordRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const formState = reactive({
  origin: "",
  password: "",
  confirm: ""
});
const updatePassword = () => {
  if (!passwordRegexp.test(formState.origin)) {
    message.warn(t("changePassword.oldPwdFormatErr"));
    return;
  }
  if (!passwordRegexp.test(formState.password)) {
    message.warn(t("changePassword.newPwdFormatErr"));
    return;
  }
  if (formState.confirm !== formState.password) {
    message.warn(t("changePassword.oldPwdNotEqualNewPwd"));
    return;
  }
  updatePasswordRequest({
    origin: formState.origin,
    password: formState.password
  }).then(() => {
    message.success(t("operationSuccess"));
    formState.origin = "";
    formState.password = "";
    formState.confirm = "";
  });
};
</script>

<style scoped>
.item {
  margin-bottom: 10px;
}
</style>