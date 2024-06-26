<template>
  <div style="padding:14px">
    <ul class="info-list">
      <li>
        <div class="info-name">名称</div>
        <div class="info-value">cnm</div>
      </li>
      <li>
        <div class="info-name">状态</div>
        <div class="info-value">
          <PlanStatusTag :status="1" />
        </div>
      </li>
      <li>
        <div class="info-name">创建人</div>
        <div class="info-value">zxjcli3</div>
      </li>
    </ul>
    <ul class="info-list">
      <li>
        <div class="info-name">创建时间</div>
        <div class="info-value">2022-09-09 00:00:00</div>
      </li>
      <li>
        <div class="info-name">制品号</div>
        <div class="info-value">202342342342342.uat</div>
      </li>
      <li>
        <div class="info-name">配置</div>
        <div class="info-value">查看配置</div>
      </li>
    </ul>
    <div>
      <a-button type="primary" :icon="h(PlayCircleOutlined)">开始发布计划</a-button>
      <a-button type="primary" :icon="h(CloseCircleOutlined)" danger style="margin-left: 6px">关闭发布计划</a-button>
    </div>
    <div class="service-name no-wrap">华南地区发布-集群1</div>
    <div class="flex" style="padding-left:calc(50% - 470px)">
      <ul class="stages">
        <li v-for="(item, index) in stageList" v-bind:key="item.id">
          <div class="stage-card">
            <div class="stage-info" @click="selectStage(item)">
              <div class="stage-name">{{item.name}}</div>
              <div class="stage-progress">
                <a-progress
                  :percent="item.percent"
                  size="small"
                  style="padding:9px"
                  :status="item.status"
                />
                <div style="text-align:center;line-height:32px;">
                  <LoadingOutlined />
                  <span style="margin-left:8px">等待交互</span>
                </div>
              </div>
            </div>
            <div class="arrow-down" v-if="index < stageList.length - 1">
              <ArrowDownOutlined />
            </div>
          </div>
          <div class="stage-line"></div>
        </li>
      </ul>
      <div class="stage-detail">
        <div class="stage-name">灰度第一批</div>
        <ul class="detail-list">
          <li>
            <div class="flex-between">
              <span class="no-wrap">agent1</span>
              <span>
                <CheckCircleFilled style="color:#52c41a" />
                <span style="margin-left:8px">成功</span>
              </span>
            </div>
          </li>
          <li>
            <div class="flex-between">
              <span class="no-wrap">agent1</span>
              <span>
                <CloseCircleFilled style="color:#ff4d4f" />
                <span style="margin-left:8px">失败</span>
              </span>
            </div>
          </li>
          <li>
            <div class="flex-between">
              <span class="no-wrap">agent1</span>
              <span>
                <CheckCircleFilled style="color:#52c41a" />
                <span style="margin-left:8px">成功</span>
              </span>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import PlanStatusTag from "@/components/app/PlanStatusTag";
import {
  ArrowDownOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  PlayCircleOutlined,
  CloseCircleOutlined,
  LoadingOutlined
} from "@ant-design/icons-vue";
import { ref, h } from "vue";
const selectedStageId = ref(1);
const stageList = [
  {
    id: 1,
    name: "自动节点-灰度发布-第一批",
    percent: 20,
    status: "exception"
  },
  {
    id: 2,
    name: "交互节点-灰度发布-第一批",
    percent: 100,
    status: ""
  },
  {
    id: 3,
    name: "交互节点-灰度发布-第一批",
    percent: 100,
    status: ""
  },
  {
    id: 4,
    name: "交互节点-灰度发布-第一批",
    percent: 100,
    status: ""
  },
  {
    id: 5,
    name: "交互节点-灰度发布-第一批",
    percent: 100,
    status: ""
  },
  {
    id: 6,
    name: "交互节点-灰度发布-第一批",
    percent: 100,
    status: ""
  }
];

const selectStage = item => {
  selectedStageId.value = item.id;
};
</script>

<style scoped>
.service-name {
  font-size: 16px;
  font-weight: bold;
  line-height: 42px;
  padding-left: 10px;
  margin-bottom: 20px;
}

.stages {
  margin-right: 40px;
}

.stages > li {
  display: flex;
}

.stage-card {
  width: 400px;
}

.stage-info {
  width: 100%;
  cursor: pointer;
}

.stage-progress {
  border-bottom: 1px solid #d9d9d9;
  border-left: 1px solid #d9d9d9;
  border-right: 1px solid #d9d9d9;
  border-bottom-left-radius: 4px;
  border-bottom-right-radius: 4px;
}

.stage-name {
  width: 100%;
  font-size: 14px;
  word-break: break-all;
  line-height: 40px;
  padding: 0 9px;
  line-height: 36px;
  color: white;
  border-top-left-radius: 4px;
  border-top-right-radius: 4px;
  background-color: #1677ff;
}

.arrow-down {
  width: 100%;
  text-align: center;
  font-size: 22px;
  color: gray;
  padding: 6px 0px;
}

.stage-detail {
  width: 500px;
  overflow: scroll;
  max-height: 600px;
}

.detail-list {
  border-bottom-left-radius: 4px;
  border-bottom-right-radius: 4px;
  border-bottom: 1px solid #d9d9d9;
  border-left: 1px solid #d9d9d9;
  border-right: 1px solid #d9d9d9;
}

.detail-list > li {
  padding: 10px;
  font-size: 14px;
}

.detail-list > li + li {
  border-top: 1px solid #d9d9d9;
}

.detail-list > li:hover {
  cursor: pointer;
  background-color: #f0f0f0;
}

.info-list > li {
  width: 33.33%;
}
</style>