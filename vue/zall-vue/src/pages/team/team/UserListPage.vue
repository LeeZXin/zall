<template>
  <div style="padding:14px">
    <div style="margin-bottom:10px">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="showAddModal">添加成员</a-button>
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
              <a-tooltip placement="top">
                <template #title>
                  <span>删除成员</span>
                </template>
                <delete-outlined />
              </a-tooltip>
            </div>
            <a-popover placement="bottomRight" trigger="hover">
              <template #content>
                <ul class="op-list">
                  <li @click="showChangeRoleModal(dataItem)">
                    <edit-outlined />
                    <span style="margin-left:4px">更换角色</span>
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
    <a-modal v-model:open="addUserModalOpen" title="添加成员" @ok="handleAddModalOk">
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">选择成员</div>
        <a-select
          v-model:value="addFormState.userList"
          style="width:100%;"
          :options="userList"
          show-search
          mode="multiple"
          :filter-option="filterUserListOption"
          placeholder="选择成员"
        />
      </div>
      <div>
        <div style="font-size:12px;margin-bottom:4px">选择角色</div>
        <a-select
          v-model:value="addFormState.roleId"
          style="width:100%;"
          :options="roleList"
          placeholder="选择角色"
        />
      </div>
    </a-modal>
    <a-modal v-model:open="changeRoleModalOpen" title="更换角色" @ok="handleChangeModalOk">
      <div style="margin-bottom:10px">
        <div style="font-size:12px;margin-bottom:4px">选择成员</div>
        <div style="font-size:14px;font-weight:bold">{{changeRoleFormState.account}}</div>
      </div>
      <div>
        <div style="font-size:12px;margin-bottom:4px">选择角色</div>
        <a-select
          v-model:value="changeRoleFormState.roleId"
          style="width:100%;"
          :options="roleList"
          placeholder="选择角色"
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
const addFormState = reactive({
  userList: [],
  roleId: undefined
});
const changeRoleModalOpen = ref(false);
const changeRoleFormState = reactive({
  id: 0,
  roleId: undefined,
  account: ""
});
const route = useRoute();
const dataSource = ref([]);
const addUserModalOpen = ref(false);
const allUserList = ref([]);
const userList = ref([]);
const columns = ref([
  {
    title: "帐号",
    dataIndex: "account",
    key: "account"
  },
  {
    title: "姓名",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "角色",
    dataIndex: "roleName",
    key: "roleName"
  },

  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);
const userStore = useUserStore();
const showChangeRoleModal = item => {
  changeRoleModalOpen.value = true;
  changeRoleFormState.roleId = item.roleId;
  changeRoleFormState.account = item.account;
  changeRoleFormState.id = item.id;
  loadRoleList();
};
const roleList = ref([]);
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
const handleChangeModalOk = () => {
  if (!changeRoleFormState.roleId) {
    message.warn("请选择角色");
    return;
  }
  changeRoleRequest({
    relationId: changeRoleFormState.id,
    roleId: changeRoleFormState.roleId
  }).then(() => {
    message.success("操作成功");
    changeRoleModalOpen.value = false;
    listTeamUsers();
  });
};
const deleteTeamUser = item => {
  Modal.confirm({
    title: `你确定要删除${item.account}(${item.name})吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteTeamUserRequest(item.id).then(() => {
        message.success("删除成功");
        listTeamUsers();
      });
    },
    onCancel() {}
  });
};
const showAddModal = () => {
  addUserModalOpen.value = true;
  loadRoleList();
};
const handleAddModalOk = () => {
  if (addFormState.userList.length === 0) {
    message.warn("请选择成员");
    return;
  }
  if (!addFormState.roleId) {
    message.warn("请选择角色");
    return;
  }
  createTeamUserRequest({
    roleId: addFormState.roleId,
    accounts: addFormState.userList
  }).then(() => {
    message.success("添加成功");
    addUserModalOpen.value = false;
    addFormState.userList = [];
    addFormState.roleId = undefined;
    listTeamUsers();
  });
};
const filterUserListOption = (input, option) => {
  return option.value.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
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