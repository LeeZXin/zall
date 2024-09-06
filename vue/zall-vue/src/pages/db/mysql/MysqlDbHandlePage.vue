<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('mysqlSource.createSource')}}</span>
        <span v-else-if="mode === 'update'">{{t('mysqlSource.updateSource')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('mysqlSource.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('mysqlSource.writeHost')}}</div>
        <div class="section-body">
          <ul class="node-ul">
            <li>
              <div class="input-title">{{t('mysqlSource.host')}}</div>
              <a-input v-model:value="formState.writeHost" />
              <div class="input-desc">{{t('mysqlSource.hostFormat')}}</div>
            </li>
            <li>
              <div class="input-title">{{t('mysqlSource.username')}}</div>
              <a-input v-model:value="formState.writeUsername" />
            </li>
            <li>
              <div class="input-title">{{t('mysqlSource.password')}}</div>
              <a-input-password v-model:value="formState.writePassword" />
            </li>
          </ul>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('mysqlSource.readHost')}}</div>
        <div class="section-body">
          <ul class="node-ul">
            <li>
              <div class="input-title">{{t('mysqlSource.host')}}</div>
              <a-input v-model:value="formState.readHost" />
              <div class="input-desc">{{t('mysqlSource.hostFormat')}}</div>
            </li>
            <li>
              <div class="input-title">{{t('mysqlSource.username')}}</div>
              <a-input v-model:value="formState.readUsername" />
            </li>
            <li>
              <div class="input-title">{{t('mysqlSource.password')}}</div>
              <a-input-password v-model:value="formState.readPassword" />
            </li>
          </ul>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateDb">{{t('mysqlSource.save')}}</a-button>
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const dbStore = useMysqldbStore();
const route = useRoute();
const router = useRouter();
// 模式 根据页面 create/update 区分
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
// 表单数据
const formState = reactive({
  name: "",
  writeHost: "",
  writeUsername: "",
  writePassword: "",
  readHost: "",
  readUsername: "",
  readPassword: ""
});
// 新增或编辑
const saveOrUpdateDb = () => {
  if (!dbNameRegexp.test(formState.name)) {
    message.warn(t("mysqlSource.nameFormatErr"));
    return;
  }
  if (!dbHostRegexp.test(formState.writeHost)) {
    message.warn(t("mysqlSource.writeHostFormatErr"));
    return;
  }
  if (!dbUsernameRegexp.test(formState.writeUsername)) {
    message.warn(t("mysqlSource.writeUsernameFormatErr"));
    return;
  }
  if (!dbHostRegexp.test(formState.readHost)) {
    message.warn(t("mysqlSource.readHostFormatErr"));
    return;
  }
  if (!dbUsernameRegexp.test(formState.readUsername)) {
    message.warn(t("mysqlSource.readUsernameFormatErr"));
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
      message.success(t("operationSuccess"));
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
      message.success(t("operationSuccess"));
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
  margin-top: 10px;
}
</style>