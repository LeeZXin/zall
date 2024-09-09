<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-button
        type="primary"
        :icon="h(PlusOutlined)"
        @click="showAddModal"
      >{{t('teamRole.addUser')}}</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'created'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <template v-if="userStore.account !== dataItem['account']">
            <div class="op-icon" @click="deleteTeamUser(dataItem)">
              <DeleteOutlined />
            </div>
            <a-popover placement="bottomRight" trigger="hover">
              <template #content>
                <ul class="op-list">
                  <li @click="showChangeRoleModal(dataItem)">
                    <EditOutlined />
                    <span style="margin-left:4px">{{t('teamRole.changeRole')}}</span>
                  </li>
                </ul>
              </template>
              <div class="op-icon">
                <EllipsisOutlined />
              </div>
            </a-popover>
          </template>
        </template>
      </template>
    </ZTable>
    <a-modal v-model:open="addUserModalOpen" :title="t('teamRole.addUser')" @ok="handleAddModalOk">
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">{{t('teamRole.selectUser')}}</div>
        <a-select
          v-model:value="addFormState.userList"
          style="width:100%;"
          :options="userList"
          show-search
          mode="multiple"
          :filter-option="filterUserListOption"
        />
      </div>
      <div>
        <div style="font-size:12px;margin-bottom:4px">{{t('teamRole.selectRole')}}</div>
        <a-select v-model:value="addFormState.roleId" style="width:100%;" :options="roleList" />
      </div>
    </a-modal>
    <a-modal
      v-model:open="changeRoleFormState.open"
      :title="t('teamRole.changeRole')"
      @ok="handleChangeModalOk"
    >
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">{{t('teamRole.selectUser')}}</div>
        <div style="font-size:14px;font-weight:bold">{{changeRoleFormState.account}}</div>
      </div>
      <div>
        <div style="font-size:12px;margin-bottom:4px">{{t('teamRole.selectRole')}}</div>
        <a-select
          v-model:value="changeRoleFormState.roleId"
          style="width:100%;"
          :options="roleList"
        />
      </div>
    </a-modal>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { ref, createVNode, h, reactive } from "vue";
import { useRoute } from "vue-router";
import {
  DeleteOutlined,
  ExclamationCircleOutlined,
  PlusOutlined,
  EditOutlined,
  EllipsisOutlined
} from "@ant-design/icons-vue";
import {
  listRoleUserRequest,
  deleteTeamUserRequest,
  createTeamUserRequest,
  listRolesRequest,
  changeRoleRequest
} from "@/api/team/teamApi";
import { listAllUserRequest } from "@/api/user/userApi";
import { Modal, message } from "ant-design-vue";
import { useUserStore } from "@/pinia/userStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
// 添加成员表单
const addFormState = reactive({
  userList: [],
  roleId: undefined
});
// 修改角色modal
const changeRoleFormState = reactive({
  open: false,
  id: 0,
  roleId: undefined,
  account: ""
});
const route = useRoute();
const dataSource = ref([]);
const addUserModalOpen = ref(false);
const allUserList = ref([]);
const userList = ref([]);
// 数据项
const columns = [
  {
    i18nTitle: "teamRole.account",
    dataIndex: "account",
    key: "account"
  },
  {
    i18nTitle: "teamRole.userName",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "teamRole.roleName",
    dataIndex: "roleName",
    key: "roleName"
  },

  {
    i18nTitle: "teamRole.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
const userStore = useUserStore();
// 展示修改角色modal
const showChangeRoleModal = item => {
  changeRoleFormState.open = true;
  changeRoleFormState.roleId = item.roleId;
  changeRoleFormState.account = item.account;
  changeRoleFormState.id = item.id;
  loadRoleList();
};
// 角色列表
const roleList = ref([]);
// 获取成员
const listTeamUsers = () => {
  listRoleUserRequest(route.params.teamId).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.account,
        ...item
      };
    });
    listAllUsers();
  });
};
// 加载角色列表
const loadRoleList = () => {
  if (roleList.value.length === 0) {
    listRolesRequest(route.params.teamId).then(res => {
      roleList.value = res.data.map(item => {
        return {
          value: item.roleId,
          label: item.name
        };
      });
    });
  }
};
// 修改角色ok
const handleChangeModalOk = () => {
  if (!changeRoleFormState.roleId) {
    message.warn(t('teamRole.pleaseSelectRole'));
    return;
  }
  changeRoleRequest({
    relationId: changeRoleFormState.id,
    roleId: changeRoleFormState.roleId
  }).then(() => {
    message.success(t("operationSuccess"));
    changeRoleFormState.open = false;
    listTeamUsers();
  });
};
// 删除成员
const deleteTeamUser = item => {
  Modal.confirm({
    title: `${t('teamRole.confirmDelete')} ${item.account}(${item.name})?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteTeamUserRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listTeamUsers();
      });
    },
    onCancel() {}
  });
};
// 展示添加成员modal
const showAddModal = () => {
  addUserModalOpen.value = true;
  loadRoleList();
};
// 添加成员
const handleAddModalOk = () => {
  if (addFormState.userList.length === 0) {
    message.warn(t('teamRole.pleaseSelectUser'));
    return;
  }
  if (!addFormState.roleId) {
    message.warn(t('teamRole.pleaseSelectRole'));
    return;
  }
  createTeamUserRequest({
    roleId: addFormState.roleId,
    accounts: addFormState.userList
  }).then(() => {
    message.success(t("operationSuccess"));
    addUserModalOpen.value = false;
    addFormState.userList = [];
    addFormState.roleId = undefined;
    listTeamUsers();
  });
};
// 下拉框过滤用户
const filterUserListOption = (input, option) => {
  return option.value.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
// 加载所有用户
const listAllUsers = () => {
  listAllUserRequest().then(res => {
    allUserList.value = res.data.map(item => {
      return {
        value: item.account,
        label: `${item.account}(${item.name})`
      };
    });
    userList.value = allUserList.value.filter(item => {
      return !dataSource.value.find(data => data.account === item.value);
    });
  });
};
listTeamUsers();
</script>
<style scoped>
</style>