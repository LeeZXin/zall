<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">{{t('profile.title')}}</div>
      <div class="section">
        <div class="section-title">{{t('profile.account')}}</div>
        <div class="section-body">{{userStore.account}}</div>
      </div>
      <div class="section">
        <div class="section-title">{{t('profile.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('profile.email')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.email" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('profile.avatar')}}</div>
        <div class="section-body">
          <a-upload
            v-model:file-list="formState.avatar"
            action="/api/files/avatar/upload"
            list-type="picture"
            :maxCount="1"
            :before-upload="beforeUpload"
            @change="uploadChange"
          >
            <a-button :icon="h(UploadOutlined)">{{t('profile.clickUpload')}}</a-button>
          </a-upload>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveProfile">{{t('profile.save')}}</a-button>
      </div>
    </div>
    <a-modal
      :title="t('profile.avatar')"
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
import { reactive, h, ref } from "vue";
import { UploadOutlined } from "@ant-design/icons-vue";
import { message } from "ant-design-vue";
import { updateUserRequest } from "@/api/user/userApi";
import { usernameRegexp, emailRegexp } from "@/utils/regexp";
import { useUserStore } from "@/pinia/userStore";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
const userStore = useUserStore();
const cropper = ref(null);
// 表单数据
const formState = reactive({
  name: "",
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
    message.success(t("profile.uploadFileSuccess"));
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
      message.error(t("profile.notImage"));
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
const saveProfile = () => {
  if (!usernameRegexp.test(formState.name)) {
    message.warn(t("profile.nameFormatErr"));
    return;
  }
  if (!emailRegexp.test(formState.email)) {
    message.warn(t("profile.emailFormatErr"));
    return;
  }
  if (formState.avatar.length === 0 || !formState.avatar[0].url) {
    message.warn(t("profile.pleaseUploadAvatar"));
    return;
  }
  updateUserRequest({
    account: userStore.account,
    name: formState.name,
    email: formState.email,
    avatarUrl: formState.avatar[0].url
  }).then(() => {
    message.success("你已编辑自己的信息, 将重新登录");
    setTimeout(() => {
      router.push("/login/login");
    }, 1000);
  });
};
const init = () => {
  formState.name = userStore.name;
  formState.email = userStore.email;
  if (userStore.avatarUrl) {
    formState.avatar = [
      {
        uid: "1",
        name: userStore.account,
        status: "done",
        url: userStore.avatarUrl,
        thumbUrl: userStore.avatarUrl
      }
    ];
  }
};
init();
</script>

<style scoped>
</style>