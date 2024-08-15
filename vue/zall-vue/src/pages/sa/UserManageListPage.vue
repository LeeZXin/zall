<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchUserKey"
        placeholder="搜索帐号"
        style="width:240px;margin-right:6px"
        @pressEnter="searchUser"
      >
        <template #suffix>
          <SearchOutlined />
        </template>
      </a-input>
      <a-button type="primary" @click="gotoCreatePage" :icon="h(PlusOutlined)">创建用户</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'avatarUrl'">
          <a-image :width="20" :height="20" :src="dataItem[dataIndex]" :fallback="fallbackAvatar" />
        </template>
        <span v-else-if="dataIndex === 'isDba'">{{dataItem[dataIndex] ? '是':'否'}}</span>
        <span v-else-if="dataIndex === 'isAdmin'">{{dataItem[dataIndex] ? '是':'否'}}</span>
        <span v-else-if="dataIndex === 'isProhibited'">{{dataItem[dataIndex] ? '是':'否'}}</span>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div
            class="op-icon"
            @click="deleteUser(dataItem)"
            v-if="dataItem['account'] !== userStore.account"
          >
            <a-tooltip placement="top">
              <template #title>
                <span>Delete User</span>
              </template>
              <DeleteOutlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="resetPassword(dataItem)">
                  <LockOutlined />
                  <span style="margin-left:4px">重置密码</span>
                </li>
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">编辑用户</span>
                </li>
                <li v-if="dataItem['isDba']" @click="setDba(dataItem, false)">
                  <UserSwitchOutlined />
                  <span style="margin-left:4px">取消DBA</span>
                </li>
                <li v-else @click="setDba(dataItem, true)">
                  <UserSwitchOutlined />
                  <span style="margin-left:4px">成为DBA</span>
                </li>
                <template v-if="dataItem['account'] !== userStore.account">
                  <li v-if="dataItem['isProhibited']" @click="setProhibited(dataItem, false)">
                    <UserSwitchOutlined />
                    <span style="margin-left:4px">取消禁用用户</span>
                  </li>
                  <li v-else @click="setProhibited(dataItem, true)">
                    <UserSwitchOutlined />
                    <span style="margin-left:4px">禁用用户</span>
                  </li>
                  <li v-if="dataItem['isAdmin']" @click="setAdmin(dataItem, false)">
                    <UserSwitchOutlined />
                    <span style="margin-left:4px">取消系统管理员</span>
                  </li>
                  <li v-else @click="setAdmin(dataItem, true)">
                    <UserSwitchOutlined />
                    <span style="margin-left:4px">成为系统管理员</span>
                  </li>
                </template>
              </ul>
            </template>
            <div class="op-icon">
              <EllipsisOutlined />
            </div>
          </a-popover>
        </div>
      </template>
    </ZTable>
    <a-pagination
      v-model:current="dataPage.current"
      :total="dataPage.totalCount"
      show-less-items
      :pageSize="dataPage.pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listUser()"
    />
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  UserSwitchOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined,
  LockOutlined,
  DeleteOutlined
} from "@ant-design/icons-vue";
import {
  listUserByAdminRequest,
  setDbaRequest,
  setAdminRequest,
  setProhibitedRequest,
  resetPasswordRequest,
  deleteUserRequest
} from "@/api/user/userApi";
import { ref, h, reactive, createVNode } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "@/pinia/userStore";
import { Modal, message } from "ant-design-vue";
import { useUserManageStore } from "@/pinia/userManageStore";
const userStore = useUserStore();
const userManageStore = useUserManageStore();
// 头像加载失败图像
const fallbackAvatar =
  "data:image/png;base64,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";
// 分页数据
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
// 搜索关键词
const searchUserKey = ref("");
const router = useRouter();
// 数据
const dataSource = ref([]);
// 数据项
const columns = [
  {
    title: "帐号",
    dataIndex: "account",
    key: "account"
  },
  {
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "邮箱",
    dataIndex: "email",
    key: "email"
  },
  {
    title: "头像",
    dataIndex: "avatarUrl",
    key: "avatarUrl"
  },
  {
    title: "是否是系统管理员",
    dataIndex: "isAdmin",
    key: "isAdmin"
  },
  {
    title: "是否是DBA",
    dataIndex: "isDba",
    key: "isDba"
  },
  {
    title: "是否禁用",
    dataIndex: "isProhibited",
    key: "isProhibited"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];
// 跳转创建用户界面
const gotoCreatePage = () => {
  router.push(`/sa/user/create`);
};
// 跳转编辑用户界面
const gotoUpdatePage = item => {
  userManageStore.account = item.account;
  userManageStore.email = item.email;
  userManageStore.name = item.name;
  userManageStore.avatarUrl = item.avatarUrl;
  router.push(`/sa/user/${item.account}/update`);
};
// 获取用户列表
const listUser = () => {
  listUserByAdminRequest({
    pageNum: dataPage.current,
    account: searchUserKey.value
  }).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.account,
        ...item
      };
    });
    dataPage.totalCount = res.totalCount;
  });
};
// 设置dba角色
const setDba = (item, isDba) => {
  let msg = isDba
    ? `是否使${item.account}成为DBA角色?`
    : `是否取消${item.account}DBA角色?`;
  Modal.confirm({
    title: msg,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      setDbaRequest({
        account: item.account,
        isDba
      }).then(() => {
        message.success("操作成功");
        item.isDba = isDba;
      });
    }
  });
};
// 设置系统管理员角色
const setAdmin = (item, isAdmin) => {
  let msg = isAdmin
    ? `是否使${item.account}成为系统管理员角色?`
    : `是否取消${item.account}系统管理员角色?`;
  Modal.confirm({
    title: msg,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      setAdminRequest({
        account: item.account,
        isAdmin
      }).then(() => {
        message.success("操作成功");
        item.isAdmin = isAdmin;
      });
    }
  });
};
// 设置禁用状态
const setProhibited = (item, isProhibited) => {
  let msg = isProhibited
    ? `是否禁用${item.account}?`
    : `是否取消禁用${item.account}?`;
  Modal.confirm({
    title: msg,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      setProhibitedRequest({
        account: item.account,
        isProhibited
      }).then(() => {
        message.success("操作成功");
        item.isProhibited = isProhibited;
      });
    }
  });
};
// 重置密码
const resetPassword = item => {
  Modal.confirm({
    title: `是否要重置${item.account}的密码?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      resetPasswordRequest(item.account).then(() => {
        message.success("操作成功");
      });
    }
  });
};
// 删除用户
const deleteUser = item => {
  Modal.confirm({
    title: `是否要删除${item.account}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteUserRequest(item.account).then(() => {
        message.success("操作成功");
        searchUser();
      });
    }
  });
};
// 搜索用户
const searchUser = () => {
  dataPage.current = 1;
  listUser();
};
listUser();
</script>
<style scoped>
</style>