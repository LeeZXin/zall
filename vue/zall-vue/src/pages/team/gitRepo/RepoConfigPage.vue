<template>
  <div style="padding:10px">
    <div class="container">
      <div class="section">
        <div class="section-title">
          <span>仓库大小</span>
        </div>
        <div class="section-body">
          <div class="input-item">代码大小: {{readableVolumeSize(repoInfo.gitSize)}}</div>
          <div class="input-item">lfs大小: {{readableVolumeSize(repoInfo.lfsSize)}}</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>仓库优化</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <a-button type="primary" @click="triggerGc">触发仓库GC</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>仓库设置</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <div class="input-title">仓库描述</div>
            <a-input v-model:value="repoInfo.repoDesc" />
            <div class="input-desc">用简短的话来描述仓库的作用</div>
          </div>
          <div class="input-item">
            <a-checkbox v-model:checked="repoInfo.disableLfs">禁用LFS</a-checkbox>
            <div class="checkbox-option-desc">禁用LFS将不允许用户上传LFS大文件, 对已有配置了LFS的仓库仍能下载原有的文件</div>
          </div>
          <div class="input-item">
            <div class="input-title">代码仓库大小限制</div>
            <a-input-number
              style="width:100%"
              v-model:value="repoInfo.gitLimitSize"
              :min="0"
              :precision="0"
            >
              <template #addonAfter>
                <a-select style="width: 100px" v-model:value="repoInfo.gitSizeUnit">
                  <a-select-option
                    v-for="item in volumeList"
                    v-bind:key="item"
                    :value="item"
                  >{{item}}</a-select-option>
                </a-select>
              </template>
            </a-input-number>
            <div class="input-desc">不包含LFS大小, 指所有代码文件大小总和, 0代表不限制, 超过限制将不允许push代码</div>
          </div>
          <div class="input-item">
            <div class="input-title">LFS大小限制</div>
            <a-input-number
              style="width:100%"
              v-model:value="repoInfo.lfsLimitSize"
              :min="0"
              :precision="0"
            >
              <template #addonAfter>
                <a-select style="width: 100px" v-model:value="repoInfo.lfsSizeUnit">
                  <a-select-option
                    v-for="item in volumeList"
                    v-bind:key="item"
                    :value="item"
                  >{{item}}</a-select-option>
                </a-select>
              </template>
            </a-input-number>
            <div class="input-desc">所有被git lfs track的文件大小总和, 0代表不限制, 超过限制将不允许push代码</div>
          </div>
          <div class="input-item">
            <a-button type="primary" @click="updateRepo">保存仓库配置</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">危险操作</div>
        <div class="section-body">
          <div class="input-item">
            <a-button type="primary" danger @click="deleteRepo">删除代码仓库</a-button>
            <div class="input-desc">将临时删除代码仓库, 在回收站保持30天, 30天后若没重置, 则将永久删除代码文件、lfs文件、合并请求、工作流等记录</div>
          </div>
          <div class="input-item" v-if="repoInfo.isArchived">
            <a-button type="primary" danger @click="setArchivedStatus(false)">取消归档代码仓库</a-button>
            <div class="input-desc">将代码仓库置从归档状态变成正常状态, 代码可读可写</div>
          </div>
          <div class="input-item" v-else>
            <a-button type="primary" danger @click="setArchivedStatus(true)">归档代码仓库</a-button>
            <div class="input-desc">将代码仓库置为归档状态且后续代码仅可读, 不可被推送</div>
          </div>
          <div class="input-item">
            <a-button type="primary" danger>迁移仓库至其他团队</a-button>
            <div class="input-desc">将代码仓库迁移至其他团队, 该仓库原有的团队配置将失效</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { readableVolumeSize, calcUnit, Unit } from "@/utils/size";
import { reactive, createVNode } from "vue";
import {
  getRepoRequest,
  gcRequest,
  updateRepoRequest,
  setArchivedRequest,
  setUnArchivedRequest,
  deleteRepoRequest
} from "@/api/git/repoApi";
import { ExclamationCircleOutlined } from "@ant-design/icons-vue";
import { useRoute, useRouter } from "vue-router";
import { message, Modal } from "ant-design-vue";
import { useRepoStore } from "@/pinia/repoStore";
const route = useRoute();
const router = useRouter();
const repoId = parseInt(route.params.repoId);
const volumeList = ["KB", "MB", "GB", "TB"];
const repoInfo = reactive({
  gitSize: 0,
  lfsSize: 0,
  disableLfs: false,
  lfsLimitSize: 0,
  gitLimitSize: 0,
  gitSizeUnit: "KB",
  lfsSizeUnit: "KB",
  isArchived: false,
  repoDesc: "",
  loaded: false
});
const getRepo = () => {
  getRepoRequest(repoId).then(res => {
    repoInfo.gitSize = res.data.gitSize;
    repoInfo.lfsSize = res.data.lfsSize;
    if (!repoInfo.loaded) {
      repoInfo.repoDesc = res.data.repoDesc;
      repoInfo.loaded = true;
      repoInfo.disableLfs = res.data.disableLfs;
      let lfs = calcUnit(res.data.lfsLimitSize);
      repoInfo.lfsLimitSize = lfs.size;
      repoInfo.lfsSizeUnit = lfs.unit.unit;
      let git = calcUnit(res.data.gitLimitSize);
      repoInfo.gitLimitSize = git.size;
      repoInfo.gitSizeUnit = git.unit.unit;
    }
    repoInfo.isArchived = res.data.isArchived;
  });
};
const updateRepo = () => {
  updateRepoRequest({
    repoId,
    desc: repoInfo.repoDesc,
    disableLfs: repoInfo.disableLfs,
    gitLimitSize: new Unit(repoInfo.gitSizeUnit).toNumber(
      repoInfo.gitLimitSize
    ),
    lfsLimitSize: new Unit(repoInfo.lfsSizeUnit).toNumber(repoInfo.lfsLimitSize)
  }).then(() => {
    message.success("编辑成功");
  });
};
const triggerGc = () => {
  gcRequest(repoId).then(() => {
    message.success("gc成功");
    getRepo();
  });
};
const setArchivedStatus = isArchived => {
  let warnMsg;
  let request;
  if (isArchived) {
    warnMsg = "你确定要归档该仓库吗?";
    request = setArchivedRequest;
  } else {
    warnMsg = "你确定要取消归档仓库吗?";
    request = setUnArchivedRequest;
  }
  Modal.confirm({
    title: warnMsg,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      request(repoId).then(()=>{
          message.success("操作成功");
          getRepo();
      })
    },
    onCancel() {}
  });
};
const deleteRepo = () => {
  Modal.confirm({
    title: "你确定要删除该仓库吗?",
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteRepoRequest(repoId).then(()=>{
          message.success("删除成功");
          router.push(`/team/${useRepoStore().teamId}/gitRepo/list`)
      })
    },
    onCancel() {}
  });
};
getRepo();
</script>
<style scoped>
</style>