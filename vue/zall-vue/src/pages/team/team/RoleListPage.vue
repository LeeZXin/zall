<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建角色</a-button>
      <a-button
        type="primary"
        :icon="h(UserOutlined)"
        @click="gotoUserPage"
        style="margin-left:6px"
      >成员列表</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <div class="op-icon" v-if="!dataItem['isAdmin']">
            <a-tooltip placement="top">
              <template #title>
                <span>删除角色</span>
              </template>
              <delete-outlined @click="deleteRole(dataItem)" />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <edit-outlined />
                  <span style="margin-left:4px">编辑角色</span>
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
const teamRoleStore = useTeamRoleStore();
const router = useRouter();
const route = useRoute();
const dataSource = ref([]);
const columns = [
  {
    i18nTitle: "roleListPage.roleName",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];
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
const deleteRole = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteRoleRequest(item.roleId).then(() => {
        message.success("删除成功");
        listRoles();
      });
    },
    onCancel() {}
  });
};
const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/role/create`);
};
const gotoUpdatePage = item => {
  teamRoleStore.roleId = item.roleId;
  teamRoleStore.teamId = item.teamId;
  teamRoleStore.name = item.name;
  teamRoleStore.teamPerm = item.perm.teamPerm;
  teamRoleStore.defaultRepoPerm = item.perm.defaultRepoPerm;
  teamRoleStore.repoPermList = item.perm.repoPermList;
  teamRoleStore.developAppList = item.perm.developAppList;
  router.push(`/team/${route.params.teamId}/role/${item.roleId}/update`);
};
const gotoUserPage = () => {
  router.push(`/team/${route.params.teamId}/role/user/list`);
};
listRoles();
</script>
<style scoped>
</style>