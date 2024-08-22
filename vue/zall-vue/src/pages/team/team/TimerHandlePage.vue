<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode==='create'">创建定时任务</span>
        <span v-else-if="mode==='update'">编辑定时任务</span>
      </div>
      <div class="section" v-if="mode==='create'">
        <div class="section-title">选择环境</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            placeholder="选择环境"
            v-model:value="formState.selectedEnv"
            :options="envList"
          />
          <div class="input-desc">多环境选择, 选择其中一个环境</div>
        </div>
      </div>
      <div class="section" v-if="mode==='update'">
        <div class="section-title">已选环境</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section">
        <div class="section-title">名称</div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.name" />
          <div class="input-desc">描述定时任务的作用</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">cron表达式</div>
        <div class="section-body">
          <a-input
            style="width:100%"
            :value="formState.cronExp"
            @focus="cronInputFocus"
            ref="cronInput"
          />
          <div class="input-desc">cron表达式, 分钟级别</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">任务类型</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.taskType">
            <a-radio value="http">HTTP</a-radio>
            <div class="radio-option-desc">定时任务执行将会发送http请求</div>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.taskType === 'http'">
        <div class="section-title">HTTP请求配置</div>
        <div class="section-body">
          <ul class="input-ul">
            <li>
              <div class="input-name">Url</div>
              <a-input v-model:value="formState.url" />
              <div class="input-desc">必须以http://或https://开头</div>
            </li>
            <li>
              <div class="input-name">Method</div>
              <a-select
                v-model:value="formState.method"
                style="width:100%"
                :options="methodList.map(item=>{return {value:item}})"
              />
            </li>
            <li>
              <div class="input-name">Headers</div>
              <ul class="headers-ul">
                <li v-for="(item, index) in formState.headers" v-bind:key="`header_${index}`">
                  <a-input style="width: 40%" v-model:value="item.key" />
                  <div style="width: 10%;text-align:center">=</div>
                  <a-input style="width: 40%" v-model:value="item.value" />
                  <div style="width: 10%;">
                    <a-button
                      style="width: 80%;margin-left:20%"
                      type="primary"
                      danger
                      :icon="h(MinusOutlined)"
                      @click="deleteHeader(index)"
                    />
                  </div>
                </li>
              </ul>
              <div style="margin-top: 10px;">
                <a-button type="primary" @click="addHeader">新增一个header</a-button>
              </div>
            </li>
            <li v-if="formState.method === 'POST'">
              <div class="input-name">请求体</div>
              <a-textarea
                :auto-size="{ minRows: 5, maxRows: 10 }"
                v-model:value="formState.body"
                @keydown.tab="handleTab"
              />
            </li>
            <li v-if="formState.method === 'POST'">
              <div class="input-name">Content-Type</div>
              <a-input v-model:value="formState.contentType" />
            </li>
          </ul>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateTimerTask">立即保存</a-button>
      </div>
      <a-modal v-model:open="cronModalOpen" title="cron表达式" :width="800" @ok="handleCronModalOk">
        <ZCron v-model="addCronExp" />
      </a-modal>
    </div>
  </div>
</template>
<script setup>
import ZCron from "@/components/common/ZCron";
import { useRoute, useRouter } from "vue-router";
import { ref, reactive, h } from "vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import { createTimerRequest, updateTimerRequest } from "@/api/team/timerApi";
import { message } from "ant-design-vue";
import { timerTaskNameRegexp } from "@/utils/regexp";
import { MinusOutlined } from "@ant-design/icons-vue";
import { useTimerTaskStore } from "@/pinia/timerTaskStore";
const timerTaskStore = useTimerTaskStore();
const router = useRouter();
const addCronExp = ref("");
const formState = reactive({
  selectedEnv: null,
  name: "",
  cronExp: "",
  taskType: "http",
  url: "http://",
  method: "POST",
  contentType: "application/json;charset=utf-8",
  headers: [
    {
      key: "",
      value: ""
    }
  ],
  body: "{}"
});
const cronInput = ref();
const cronModalOpen = ref(false);
const envList = ref([]);
const route = useRoute();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const cronInputFocus = () => {
  cronModalOpen.value = true;
  cronInput.value.blur();
};
const getEnvCfg = () => {
  getEnvCfgRequest().then(res => {
    envList.value = res.data.map(item => {
      return {
        value: item,
        label: item
      };
    });
    if (route.query.env && res.data?.includes(route.query.env)) {
      formState.selectedEnv = route.query.env;
    } else if (res.data.length > 0) {
      formState.selectedEnv = res.data[0];
    }
  });
};
const handleTab = event => {
  event.preventDefault();
  let inputElement = event.target;
  let value = inputElement.value;
  let selectionStart = inputElement.selectionStart;
  let leftValue = value.substring(0, selectionStart);
  let rightValue = value.substring(selectionStart);
  inputElement.value = leftValue + "    " + rightValue;
  inputElement.selectionStart = selectionStart + 4;
  inputElement.selectionEnd = inputElement.selectionStart;
};
const methodList = ["GET", "POST"];
const handleCronModalOk = () => {
  formState.cronExp = addCronExp.value;
  cronModalOpen.value = false;
};
const addHeader = () => {
  formState.headers.push({
    key: "",
    value: ""
  });
};
const deleteHeader = index => {
  formState.headers.splice(index, 1);
};
const saveOrUpdateTimerTask = () => {
  if (!formState.selectedEnv) {
    message.warn("请选择环境");
    return;
  }
  if (!timerTaskNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!formState.cronExp) {
    message.warn("请配置cron");
    return;
  }
  let task = {
    taskType: formState.taskType
  };
  if (formState.taskType === "http") {
    if (!formState.url.startsWith("http")) {
      message.warn("请输入正确的http url");
      return;
    }
    let headers = {};
    formState.headers.forEach(item => {
      if (item.key && item.value) {
        headers[item.key] = item.value;
      }
    });
    let contentType = "";
    let body = "";
    if (formState.method === "POST") {
      contentType = formState.contentType;
      body = formState.body;
    }
    task.httpTask = {
      url: formState.url,
      method: formState.method,
      headers: headers,
      bodyStr: body,
      contentType: contentType
    };
  }
  if (mode === "create") {
    createTimerRequest({
      env: formState.selectedEnv,
      cronExp: formState.cronExp,
      teamId: parseInt(route.params.teamId),
      name: formState.name,
      task
    }).then(() => {
      message.success("创建成功");
      router.push(
        `/team/${route.params.teamId}/timer/list/${formState.selectedEnv}`
      );
    });
  } else if (mode === "update") {
    updateTimerRequest({
      cronExp: formState.cronExp,
      task,
      name: formState.name,
      id: timerTaskStore.id
    }).then(() => {
      message.success("编辑成功");
      router.push(
        `/team/${route.params.teamId}/timer/list/${formState.selectedEnv}`
      );
    });
  }
};
if (mode === "create") {
  getEnvCfg();
} else if (mode === "update") {
  if (timerTaskStore.id === 0) {
    router.push(`/team/${route.params.teamId}/timer/list`);
  } else {
    formState.selectedEnv = timerTaskStore.env;
    formState.cronExp = timerTaskStore.cronExp;
    addCronExp.value = timerTaskStore.cronExp;
    formState.name = timerTaskStore.name;
    formState.taskType = timerTaskStore.task.taskType;
    if (formState.taskType === "http") {
      formState.url = timerTaskStore.task.httpTask.url;
      formState.method = timerTaskStore.task.httpTask.method;
      formState.body = timerTaskStore.task.httpTask.bodyStr;
      formState.contentType = timerTaskStore.task.httpTask.contentType;
      let headers = timerTaskStore.task.httpTask.headers;
      let retHeaders = [];
      if (headers) {
        for (let key in headers) {
          retHeaders.push({
            key,
            value: headers[key]
          });
        }
        if (retHeaders.length > 0) {
          formState.headers = retHeaders;
        }
      }
    }
  }
}
</script>
<style scoped>
.headers-ul > li {
  height: 32px;
  line-height: 32px;
  width: 100%;
  display: flex;
  align-items: center;
  font-size: 14px;
}
.headers-ul > li + li {
  margin-top: 6px;
}
</style>