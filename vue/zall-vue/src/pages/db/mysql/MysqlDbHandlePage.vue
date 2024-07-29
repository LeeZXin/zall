<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建Mysql数据源</span>
        <span v-else-if="mode === 'update'">编辑Mysql数据源</span>
      </div>
      <div class="section">
        <div class="section-title">名称</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
          <div class="input-desc">标识数据源</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">写节点</div>
        <div class="section-body">
          <ul class="node-ul">
            <li>
              <div class="input-title">host</div>
              <a-input v-model:value="formState.writeHost" />
              <div class="input-desc">数据库host ip:port格式</div>
            </li>
            <li>
              <div class="input-title">账号</div>
              <a-input v-model:value="formState.writeUsername" />
              <div class="input-desc">数据库账号</div>
            </li>
            <li>
              <div class="input-title">密码</div>
              <a-input-password v-model:value="formState.writePassword" />
              <div class="input-desc">数据库密码</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="section">
        <div class="section-title">读节点</div>
        <div class="section-body">
          <ul class="node-ul">
            <li>
              <div class="input-title">host</div>
              <a-input v-model:value="formState.readHost" />
              <div class="input-desc">数据库host ip:port格式</div>
            </li>
            <li>
              <div class="input-title">账号</div>
              <a-input v-model:value="formState.readUsername" />
              <div class="input-desc">数据库账号</div>
            </li>
            <li>
              <div class="input-title">密码</div>
              <a-input-password v-model:value="formState.readPassword" />
              <div class="input-desc">数据库密码</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateDb">立即保存</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import { dbHostRegexp, dbNameRegexp, dbUsernameRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { createMysqlDbRequest, updateMysqlDbRequest } from "@/api/db/mysqlApi";
import { useRoute, useRouter } from "vue-router";
import { useMysqldbStore } from "@/pinia/mysqldbStore";
const dbStore = useMysqldbStore();
const route = useRoute();
const router = useRouter();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const formState = reactive({
  name: "",
  writeHost: "",
  writeUsername: "",
  writePassword: "",
  readHost: "",
  readUsername: "",
  readPassword: ""
});

const saveOrUpdateDb = () => {
  if (!dbNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!dbHostRegexp.test(formState.writeHost)) {
    message.warn("写节点host格式错误");
    return;
  }
  if (!dbUsernameRegexp.test(formState.writeUsername)) {
    message.warn("写节点账号格式错误");
    return;
  }
  if (!dbHostRegexp.test(formState.readHost)) {
    message.warn("读节点host格式错误");
    return;
  }
  if (!dbUsernameRegexp.test(formState.readUsername)) {
    message.warn("读节点账号格式错误");
    return;
  }
  if (mode === "create") {
    createMysqlDbRequest({
      name: formState.name,
      config: {
        writeNode: {
          host: formState.writeHost,
          username: formState.writeUsername,
          password: formState.writePassword
        },
        readNode: {
          host: formState.readHost,
          username: formState.readUsername,
          password: formState.readPassword
        }
      }
    }).then(() => {
      message.success("创建成功");
      router.push(`/db/mysqlDb/list`);
    });
  } else if (mode === "update") {
    updateMysqlDbRequest({
      dbId: dbStore.id,
      name: formState.name,
      config: {
        writeNode: {
          host: formState.writeHost,
          username: formState.writeUsername,
          password: formState.writePassword
        },
        readNode: {
          host: formState.readHost,
          username: formState.readUsername,
          password: formState.readPassword
        }
      }
    }).then(() => {
      message.success("保存成功");
      router.push(`/db/mysqlDb/list`);
    });
  }
};

if (mode === "update") {
  if (dbStore.id === 0) {
    router.push(`/db/mysqlDb/list`);
  } else {
    formState.name = dbStore.name;
    formState.writeHost = dbStore.writeHost;
    formState.writeUsername = dbStore.writeUsername;
    formState.writePassword = dbStore.writePassword;
    formState.readHost = dbStore.readHost;
    formState.readUsername = dbStore.readUsername;
    formState.readPassword = dbStore.readPassword;
  }
}
</script>
<style scoped>
.node-ul > li + li {
  margin-top: 20px;
}
</style>