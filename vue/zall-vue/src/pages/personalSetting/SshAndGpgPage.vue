<template>
  <div style="padding:10px">
    <div class="container" style="width:70%">
      <div class="header">SSH&GPG</div>
      <div class="section">
        <div class="section-title flex-between">
          <span>SSH</span>
          <span
            class="add-btn"
            @click="goto('/personalSetting/sshAndGpg/createSshKey')"
          >{{t('sshGpg.createSsh')}}</span>
        </div>
        <ul class="key-ul" v-if="sshList.length > 0">
          <li v-for="item in sshList" v-bind:key="item.id">
            <div style="padding-left:10px;width:90%">
              <div class="key-name no-wrap">{{item.name}}</div>
              <div class="key-extra no-wrap">{{item.fingerprint}}</div>
              <div class="key-extra no-wrap">{{t('sshGpg.createdAt')}} {{item.created}}</div>
              <div class="key-extra no-wrap">{{t('sshGpg.recentlyUsedAt')}} {{item.lastOperated}}</div>
            </div>
            <div class="key-right">
              <a-button type="primary" danger @click="deleteSshKey(item)">{{t('sshGpg.delete')}}</a-button>
            </div>
          </li>
        </ul>
        <ZNoData v-else :unbordered="true" style="border-top: 1px solid #d9d9d9;"/>
      </div>
      <div class="section">
        <div class="section-title flex-between">
          <span>GPG</span>
          <span
            class="add-btn"
            @click="goto('/personalSetting/sshAndGpg/createGpgKey')"
          >{{t('sshGpg.createGpg')}}</span>
        </div>
        <ul class="key-ul" v-if="gpgList.length > 0">
          <li v-for="item in gpgList" v-bind:key="item.id">
            <div style="padding-left:10px;width:90%">
              <div class="key-name no-wrap">{{item.name}}</div>
              <div class="key-extra no-wrap">KeyID {{item.keyId}}</div>
              <div class="key-extra no-wrap">{{t('sshGpg.email')}} {{item.email}}</div>
              <div class="key-extra no-wrap">{{t('sshGpg.subKeys')}} {{item.subKeys}}</div>
              <div class="key-extra no-wrap">{{t('sshGpg.createdAt')}} {{item.created}}</div>
              <div class="key-extra no-wrap">{{t('sshGpg.expiredAt')}} {{item.expired}}</div>
            </div>
            <div class="key-right">
              <a-button type="primary" danger @click="deleteGpgKey(item)">{{t('sshGpg.delete')}}</a-button>
            </div>
          </li>
        </ul>
        <ZNoData v-else :unbordered="true" style="border-top: 1px solid #d9d9d9;"/>
      </div>
    </div>
  </div>
</template>
<script setup>
import ZNoData from "@/components/common/ZNoData";
import { useRouter } from "vue-router";
import { ref, createVNode } from "vue";
import {
  listAllSshKeyRequest,
  deleteSshKeyRequest
} from "@/api/user/sshKeyApi";
import {
  listAllGpgKeyRequest,
  deleteGpgKeyRequest
} from "@/api/user/gpgKeyApi";
import { Modal, message } from "ant-design-vue";
import { ExclamationCircleOutlined } from "@ant-design/icons-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
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
// 获取gpg密钥列表
const listGpgKey = () => {
  listAllGpgKeyRequest().then(res => {
    gpgList.value = res.data.map(item => {
      return {
        ...item
      };
    });
  });
};
// 删除ssh密钥
const deleteSshKey = item => {
  Modal.confirm({
    title: `${t("sshGpg.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteSshKeyRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listSshKey();
      });
    }
  });
};
// 删除gpg密钥
const deleteGpgKey = item => {
  Modal.confirm({
    title: `${t("sshGpg.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteGpgKeyRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listGpgKey();
      });
    }
  });
};
listGpgKey();
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
.key-extra {
  font-size: 14px;
  padding-bottom: 4px;
  color: gray;
}
.add-btn:hover {
  cursor: pointer;
  color: #1677ff;
}
</style>