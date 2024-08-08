<template>
  <div style="padding:10px">
    <div class="container" style="width:60%">
      <div class="title">SSH和GPG</div>
      <div class="section">
        <div class="section-title flex-between">
          <span>SSH</span>
          <span class="add-btn" @click="goto('/personalSetting/sshAndGpg/createSshKey')">新增SSH密钥</span>
        </div>
        <ul class="key-ul" v-if="sshList.length > 0">
          <li v-for="item in sshList" v-bind:key="item.id">
            <div style="padding-left:10px;width:90%">
              <div class="key-name no-wrap">{{item.name}}</div>
              <div class="key-sha no-wrap">{{item.fingerprint}}</div>
              <div class="key-extra">添加于 {{item.created}}</div>
              <div class="key-extra">最近使用于 {{item.lastOperated}}</div>
            </div>
            <div class="key-right">
              <a-button type="primary" danger @click="deleteSshKey(item)">删除</a-button>
            </div>
          </li>
        </ul>
        <div v-else class="no-data">没有上传SSH密钥, 请点击新增SSH密钥</div>
      </div>
      <div class="section">
        <div class="section-title flex-between">
          <span>GPG</span>
          <span class="add-btn" @click="goto('/personalSetting/sshAndGpg/createGpgKey')">新增GPG密钥</span>
        </div>
        <ul class="key-ul" v-if="gpgList.length > 0">
          <li v-for="item in gpgList" v-bind:key="item.id">
            <div style="padding-left:10px;width:90%">
              <div class="key-name no-wrap">{{item.name}}</div>
              <div class="key-sha no-wrap">{{item.sha}}</div>
              <div class="key-extra">添加于 {{item.created}}</div>
              <div class="key-extra">最近使用于 {{item.updated}}</div>
            </div>
            <div class="key-right">
              <a-button type="primary" danger>删除</a-button>
            </div>
          </li>
        </ul>
        <div v-else class="no-data">没有上传GPG密钥, 请点击新增GPG密钥</div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { useRouter } from "vue-router";
import { ref, createVNode } from "vue";
import {
  listAllSshKeyRequest,
  deleteSshKeyRequest
} from "@/api/user/sshKeyApi";
import { Modal, message } from "ant-design-vue";
import { ExclamationCircleOutlined } from "@ant-design/icons-vue";
const router = useRouter();
const sshList = ref([]);
const gpgList = ref([]);
// 页面跳转
const goto = href => {
  router.push(href);
};
// 获取ssh密钥列表
const listSshKey = () => {
  listAllSshKeyRequest().then(res => {
    sshList.value = res.data.map(item => {
      return {
        ...item
      };
    });
  });
};
// 删除ssh密钥
const deleteSshKey = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteSshKeyRequest(item.id).then(() => {
        message.success("删除成功");
        listSshKey();
      });
    }
  });
};
listSshKey();
</script>
<style scoped>
.key-ul {
  border-top: 1px solid #d9d9d9;
}
.key-ul > li {
  width: 100%;
  display: flex;
  align-items: center;
  padding: 10px;
  color: black;
}
.key-right {
  width: 10%;
  text-align: center;
}
.key-ul > li + li {
  border-top: 1px solid #d9d9d9;
}
.key-name {
  font-size: 14px;
  padding-bottom: 6px;
}
.key-sha {
  font-size: 14px;
  color: gray;
  padding-bottom: 6px;
}
.key-extra {
  font-size: 14px;
  padding-bottom: 4px;
  color: gray;
}
.no-data {
  border-top: 1px solid #d9d9d9;
  padding: 20px 10px;
  font-size: 16px;
  text-align: center;
  color: gray;
}
.add-btn:hover {
  cursor: pointer;
  color: #1677ff;
}
</style>