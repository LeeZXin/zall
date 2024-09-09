<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-button
        type="primary"
        :icon="h(PlusOutlined)"
        @click="gotoCreatePage"
      >{{t('teamRole.createRole')}}</a-button>
      <a-button
        type="primary"
        :icon="h(UserOutlined)"
        @click="gotoMembersPage"
        style="margin-left:6px"
      >{{t('teamRole.members')}}</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <div class="op-icon" v-if="!dataItem['isAdmin']" @click="deleteRole(dataItem)">
            <DeleteOutlined />
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">{{t('teamRole.updateRole')}}</span>
                </li>
              </ul>
            </template>
            <div class="op-icon" v-if="!dataItem['isAdmin']">
              <EllipsisOutlined />
            </div>
          </a-popover>
        </template>
      </template>
    </ZTable>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { ref, createVNode, h } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  DeleteOutlined,
  EditOutlined,
  ExclamationCircleOutlined,
  EllipsisOutlined,
  PlusOutlined,
  UserOutlined
} from "@ant-design/icons-vue";
import { listRolesRequest, deleteRoleRequest } from "@/api/team/teamApi";
import { Modal, message } from "ant-design-vue";
import { useTeamRoleStore } from "@/pinia/teamRoleStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const teamRoleStore = useTeamRoleStore();
const router = useRouter();
const route = useRoute();
const dataSource = ref([]);
// 数据项
const columns = [
  {
    i18nTitle: "teamRole.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "teamRole.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 获取角色列表
const listRoles = () => {
  listRolesRequest(route.params.teamId).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.roleId,
        ...item
      };
    });
  });
};
// 删除角色
const deleteRole = item => {
  Modal.confirm({
    title: `${t("teamRole.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteRoleRequest(item.roleId).then(() => {
        message.success(t("operationSuccess"));
        listRoles();
      });
    },
    onCancel() {}
  });
};
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/role/create`);
};
// 跳转编辑页面
const gotoUpdatePage = item => {
  teamRoleStore.roleId = item.roleId;
  teamRoleStore.teamId = item.teamId;
  teamRoleStore.name = item.name;
  teamRoleStore.teamPerm = item.perm.teamPerm;
  teamRoleStore.defaultRepoPerm = item.perm.defaultRepoPerm;
  teamRoleStore.repoPermList = item.perm.repoPermList;
  teamRoleStore.defaultAppPerm = item.perm.defaultAppPerm;
  teamRoleStore.appPermList = item.perm.appPermList;
  router.push(`/team/${route.params.teamId}/role/${item.roleId}/update`);
};
// 成员列表
const gotoMembersPage = () => {
  router.push(`/team/${route.params.teamId}/role/user/list`);
};
listRoles();
</script>
<style scoped>
</style>