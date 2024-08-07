<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">修改密码</div>
      <div class="item">
        <div class="input-title">原密码</div>
        <a-input-password v-model:value="formState.origin" />
      </div>
      <div class="item">
        <div class="input-title">新密码</div>
        <a-input-password v-model:value="formState.password" />
      </div>
      <div class="item">
        <div class="input-title">确认密码</div>
        <a-input-password v-model:value="formState.confirm" />
      </div>
      <div>
        <a-button type="primary" @click="updatePassword">立即修改</a-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { updatePasswordRequest } from "@/api/user/userApi";
import { reactive } from "vue";
import { passwordRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
const formState = reactive({
  origin: "",
  password: "",
  confirm: ""
});
const updatePassword = () => {
  if (!passwordRegexp.test(formState.origin)) {
    message.warn("原密码格式错误");
    return;
  }
  if (!passwordRegexp.test(formState.passwordRegexp)) {
    message.warn("新密码格式错误");
    return;
  }
  if (formState.confirm !== formState.password) {
    message.warn("确认密码与新密码不一致");
    return;
  }
  updatePasswordRequest({
    origin: formState.origin,
    password: formState.password
  }).then(() => {
    message.success("修改成功");
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