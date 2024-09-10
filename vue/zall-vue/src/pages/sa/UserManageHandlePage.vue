<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">创建用户</span>
        <span v-else-if="mode === 'update'">编辑用户</span>
      </div>
      <div class="section">
        <div class="section-title">帐号</div>
        <div class="section-body" v-if="mode === 'create'">
          <a-input v-model:value="formState.account" />
          <div class="input-desc">用户唯一标识, 长度为4-32</div>
        </div>
        <div class="section-body" v-else>{{formState.account}}</div>
      </div>
      <div class="section">
        <div class="section-title">名称</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
          <div class="input-desc">用户名称, 长度为1-32</div>
        </div>
      </div>
      <div class="section" v-if="mode === 'create'">
        <div class="section-title">密码</div>
        <div class="section-body">
          <a-input-password v-model:value="formState.password" />
          <div class="input-desc">帐号密码, 长度为6-255</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">邮箱</div>
        <div class="section-body">
          <a-input v-model:value="formState.email" />
          <div class="input-desc">邮箱格式</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">头像</div>
        <div class="section-body">
          <a-upload
            v-model:file-list="formState.avatar"
            action="/api/files/avatar/upload"
            list-type="picture"
            :maxCount="1"
            :before-upload="beforeUpload"
            @change="uploadChange"
          >
            <a-button :icon="h(UploadOutlined)">点击上传</a-button>
          </a-upload>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateUser">立即保存</a-button>
      </div>
    </div>
    <a-modal
      title="裁剪头像"
      v-model:open="cropModal.open"
      :width="448"
      @ok="handleCropModalOk"
      @cancel="handleCropModalCancel"
    >
      <VueCropper
        style="width:100%;height:400px"
        :img="cropModal.img"
        outputType="png"
        :autoCrop="true"
        :guides="false"
        :canMove="false"
        :autoCropWidth="200"
        :autoCropHeight="200"
        :centerBox="true"
        :round="true"
        :isCircleCropping="true"
        :fixed="true"
        :fixNumber="[1, 1]"
        ref="cropper"
      />
    </a-modal>
  </div>
</template>

<script setup>
import "vue-cropper/dist/index.css";
import { VueCropper } from "vue-cropper";
import { useRoute, useRouter } from "vue-router";
import { reactive, h, ref } from "vue";
import { UploadOutlined } from "@ant-design/icons-vue";
import { message } from "ant-design-vue";
import { createUserRequest, updateUserRequest } from "@/api/user/userApi";
import {
  usernameRegexp,
  accountRegexp,
  passwordRegexp,
  emailRegexp
} from "@/utils/regexp";
import { useUserManageStore } from "@/pinia/userManageStore";
import { useUserStore } from "@/pinia/userStore";
const cropper = ref(null);
const userStore = useUserStore();
const userManageStore = useUserManageStore();
const router = useRouter();
const route = useRoute();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
// 模式
const mode = getMode();
// 表单数据
const formState = reactive({
  account: "",
  name: "",
  password: "",
  email: "",
  avatar: [],
  oldAvatar: []
});
// 头像裁剪modal
const cropModal = reactive({
  open: false,
  img: "",
  resolve: null,
  reject: null
});
// 裁剪成功
const handleCropModalOk = () => {
  cropper.value.getCropBlob(blob => {
    if (cropModal.resolve) {
      cropModal.resolve(blob);
    }
    cropModal.open = false;
  });
};
// 取消裁剪
const handleCropModalCancel = () => {
  if (cropModal.reject) {
    cropModal.reject();
  }
  cropModal.open = false;
};
// 上传头像change
const uploadChange = info => {
  if (info.file.status === "done" && info.file.response.filePath) {
    message.success("上传成功");
    info.file.url = info.file.response.filePath;
  } else if (info.file.status === "error" && info.file.response.error) {
    message.error(info.file.response.error);
  }
  if (!info.file.status || info.file.status === "error") {
    formState.avatar = formState.oldAvatar;
  }
};
// 上传头像校验
const beforeUpload = file => {
  formState.oldAvatar = formState.avatar;
  return new Promise((resolve, reject) => {
    const isJpgOrPng = file.type === "image/jpeg" || file.type === "image/png";
    if (!isJpgOrPng) {
      message.error("不是图像");
      reject();
      return;
    }
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => {
      cropModal.img = reader.result;
      cropModal.open = true;
      cropModal.resolve = resolve;
      cropModal.reject = reject;
    };
  });
};
// 创建或编辑用户
const saveOrUpdateUser = () => {
  if (!accountRegexp.test(formState.account)) {
    message.warn("帐号格式错误");
    return;
  }
  if (!usernameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!emailRegexp.test(formState.email)) {
    message.warn("邮箱格式错误");
    return;
  }
  if (formState.avatar.length === 0 || !formState.avatar[0].url) {
    message.warn("请上传头像");
    return;
  }
  if (mode === "create") {
    if (!passwordRegexp.test(formState.password)) {
      message.warn("密码格式错误");
      return;
    }
    createUserRequest({
      account: formState.account,
      name: formState.name,
      password: formState.password,
      email: formState.email,
      avatarUrl: formState.avatar[0].url
    }).then(() => {
      message.success("创建成功");
      router.push("/sa/user/list");
    });
  } else {
    updateUserRequest({
      account: formState.account,
      name: formState.name,
      email: formState.email,
      avatarUrl: formState.avatar[0].url
    }).then(() => {
      if (formState.account === userStore.account) {
        message.info("你已编辑自己的信息, 将重新登录");
        setTimeout(() => {
          router.push("/login/login");
        }, 1000);
      } else {
        message.success("编辑成功");
        router.push("/sa/user/list");
      }
    });
  }
};
if (mode === "update") {
  if (userManageStore.account === "") {
    router.push("/sa/user/list");
  } else {
    formState.account = userManageStore.account;
    formState.name = userManageStore.name;
    formState.email = userManageStore.email;
    if (userManageStore.avatarUrl) {
      formState.avatar = [
        {
          uid: "1",
          name: userManageStore.account,
          status: "done",
          url: userManageStore.avatarUrl,
          thumbUrl: userManageStore.avatarUrl
        }
      ];
    }
  }
}
</script>

<style scoped>
</style>