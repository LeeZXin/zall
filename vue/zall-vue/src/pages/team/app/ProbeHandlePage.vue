<template>
  <div style="padding:14px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建探针</span>
        <span v-else-if="mode === 'update'">编辑探针</span>
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
          <a-input v-model:value="formState.name" />
          <div class="input-desc">标识探针作用</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">探针类型</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            v-model:value="formState.probeType"
            :options="probeTypeList"
          />
          <div class="input-desc">单选探针类型</div>
        </div>
      </div>
      <div class="section" v-if="formState.probeType === 'http'">
        <div class="section-title">
          <span>HTTP探针</span>
        </div>
        <div class="section-body">
          <ul class="input-ul">
            <li>
              <div class="input-name">URL</div>
              <a-input v-model:value="formState.httpUrl" />
              <div class="input-desc">必须以http为开头的url</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-if="formState.probeType === 'tcp'">
        <div class="section-title">
          <span>TCP探针</span>
        </div>
        <div class="section-body">
          <ul class="input-ul">
            <li>
              <div class="input-name">目标</div>
              <a-input v-model:value="formState.tcpAddr" />
              <div class="input-desc">ip:端口格式的字符串</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-if="formState.probeType === 'k8s'">
        <div class="section-title">
          <span>k8s</span>
        </div>
        <div class="section-body"></div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateProbe">立即保存</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import {
  probeNameRegexp,
  probeHttpUrlRegexp,
  probeTcpAddrRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import { createProbeRequest, updateProbeRequest } from "@/api/app/probeApi";
import { useRoute, useRouter } from "vue-router";
import { useProbeStore } from "@/pinia/probeStore";
const probeStore = useProbeStore();
const route = useRoute();
const router = useRouter();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const formState = reactive({
  name: "",
  content: "",
  selectedEnv: "",
  probeType: "http",
  httpUrl: "",
  tcpAddr: ""
});
const probeTypeList = [
  {
    value: "http",
    label: "HTTP探针"
  },
  {
    value: "tcp",
    label: "TCP探针"
  },
  {
    value: "k8s",
    label: "K8S"
  }
];
const envList = ref([]);
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
const saveOrUpdateProbe = () => {
  if (!probeNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  let config = {
    type: formState.probeType
  };
  switch (formState.probeType) {
    case "http":
      if (!probeHttpUrlRegexp.test(formState.httpUrl)) {
        message.warn("http url错误");
        return;
      }
      config["http"] = {
        url: formState.httpUrl
      };
      break;
    case "tcp":
      if (!probeTcpAddrRegexp.test(formState.tcpAddr)) {
        message.warn("tcp目标错误");
        return;
      }
      config["tcp"] = {
        addr: formState.tcpAddr
      };
      break;
    case "k8s":
      config["k8s"] = {};
      break;
  }
  if (mode === "create") {
    createProbeRequest(
      {
        env: formState.selectedEnv,
        appId: route.params.appId,
        config: config,
        name: formState.name
      },
      formState.selectedEnv
    ).then(() => {
      message.success("创建成功");
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/probe/list/${formState.selectedEnv}`
      );
    });
  } else if (mode === "update") {
    updateProbeRequest(
      {
        probeId: probeStore.id,
        config: config,
        name: formState.name
      },
      formState.selectedEnv
    ).then(() => {
      message.success("保存成功");
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/probe/list/${formState.selectedEnv}`
      );
    });
  }
};

if (mode === "create") {
  getEnvCfg();
} else if (mode === "update") {
  if (probeStore.id === 0) {
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/probe/list`
    );
  } else {
    formState.name = probeStore.name;
    formState.selectedEnv = probeStore.env;
    formState.probeType = probeStore.config.type;
    switch (probeStore.config.type) {
      case "http":
        formState.httpUrl = probeStore.config.http.url;
        break;
      case "tcp":
        formState.tcpAddr = probeStore.config.tcp.addr;
        break;
    }
  }
}
</script>
<style scoped>
.section-body > li {
  width: 100%;
}
</style>