<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">新增SSH密钥</div>
      <div class="section">
        <div class="section-title">名称</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
          <div class="input-desc">密钥的名称, 长度为1-128</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">内容</div>
        <div class="section-body">
          <a-textarea
            style="width:100%"
            v-model:value="formState.content"
            :auto-size="{ minRows: 3, maxRows: 8 }"
          />
          <div
            class="input-desc"
          >以 'ssh-ed25519'、 'ssh-rsa'、 'ecdsa-sha2-nistp256'、'ecdsa-sha2-nistp384'、'ecdsa-sha2-nistp521'、 'sk-ecdsa-sha2-nistp256@openssh.com' 或 'sk-ssh-ed25519@openssh.com' 开头</div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveSshKey">立即新增</a-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive } from "vue";
import { createSshKeyRequest } from "@/api/user/sshKeyApi";
import { sshKeyNameRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { useRouter } from "vue-router";
const router = useRouter();
// 表单数据
const formState = reactive({
  name: "",
  content: ""
});
// 点击“立即新增”按钮
const saveSshKey = () => {
  if (!sshKeyNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!formState.content) {
    message.warn("密钥内容为空");
    return;
  }
  createSshKeyRequest({
    name: formState.name,
    pubKeyContent: formState.content
  }).then(() => {
    message.success("新增成功");
    router.push("/personalSetting/sshAndGpg/list");
  });
};
</script>

<style scoped>
</style>