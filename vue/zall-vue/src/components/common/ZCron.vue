<template>
  <div :style="props.style">
    <ul class="cron-input-ul">
      <li>
        <div class="cron-input">
          <a-input :style="istyle" :value="minute.value" @focus="inputFocus('minute')" />
        </div>
        <div class="cron-item">分</div>
      </li>
      <li>
        <div class="cron-input">
          <a-input :style="istyle" :value="hour.value" @focus="inputFocus('hour')" />
        </div>
        <div class="cron-item">时</div>
      </li>
      <li>
        <div class="cron-input">
          <a-input :style="istyle" :value="dayOfMonth.value" @focus="inputFocus('dayOfMonth')" />
        </div>
        <div class="cron-item">日</div>
      </li>
      <li>
        <div class="cron-input">
          <a-input :style="istyle" :value="month.value" @focus="inputFocus('month')" />
        </div>
        <div class="cron-item">月</div>
      </li>
      <li>
        <div class="cron-input">
          <a-input :style="istyle" :value="dayOfWeek.value" @focus="inputFocus('dayOfWeek')" />
        </div>
        <div class="cron-item">周</div>
      </li>
    </ul>
    <div style="font-size:14px;padding:6px" v-show="tab === 'minute'">
      <a-radio-group v-model:value="minute.opt" @change="minuteRadioChange">
        <a-radio :value="1" :style="radioStyle">每分钟</a-radio>
        <a-radio :value="2" :style="radioStyle">
          <span>从</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="0"
            :max="59"
            size="small"
            :disabled="minute.opt !== 2"
            v-model:value="minute.periodStart"
            @change="minutePeriodChange"
          />
          <span>-</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="0"
            :max="59"
            size="small"
            :disabled="minute.opt !== 2"
            v-model:value="minute.periodEnd"
            @change="minutePeriodChange"
          />
          <span>每分钟执行一次</span>
        </a-radio>
        <a-radio :value="3" :style="radioStyle">
          <span>从</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="0"
            :max="59"
            size="small"
            :disabled="minute.opt !== 3"
            v-model:value="minute.loopStart"
            @change="minuteLoopChange"
          />
          <span>开始, 每</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="59"
            size="small"
            :disabled="minute.opt !== 3"
            v-model:value="minute.loopEnd"
            @change="minuteLoopChange"
          />
          <span>分钟执行一次</span>
        </a-radio>
        <a-radio :value="4" :style="radioStyle">
          <span>指定在</span>
          <a-select
            style="width: 240px;margin:0 4px"
            size="small"
            :disabled="minute.opt !== 4"
            v-model:value="minute.fix"
            :options="minuteOptions"
            mode="multiple"
            @click="preventDefault"
            :max-tag-count="3"
            @change="minuteFixChange"
          />
          <span>分钟执行</span>
        </a-radio>
      </a-radio-group>
    </div>
    <div style="font-size:14px;padding:6px" v-show="tab === 'hour'">
      <a-radio-group v-model:value="hour.opt" @change="hourRadioChange">
        <a-radio :value="1" :style="radioStyle">每小时</a-radio>
        <a-radio :value="2" :style="radioStyle">
          <span>从</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="0"
            :max="23"
            size="small"
            :disabled="hour.opt !== 2"
            v-model:value="hour.periodStart"
            @change="hourPeriodChange"
          />
          <span>-</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="0"
            :max="23"
            size="small"
            :disabled="hour.opt !== 2"
            v-model:value="hour.periodEnd"
            @change="hourPeriodChange"
          />
          <span>每小时执行一次</span>
        </a-radio>
        <a-radio :value="3" :style="radioStyle">
          <span>从</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="0"
            :max="23"
            size="small"
            :disabled="hour.opt !== 3"
            v-model:value="hour.loopStart"
            @change="hourLoopChange"
          />
          <span>开始, 每</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="23"
            size="small"
            :disabled="hour.opt !== 3"
            v-model:value="hour.loopEnd"
            @change="hourLoopChange"
          />
          <span>小时执行一次</span>
        </a-radio>
        <a-radio :value="4" :style="radioStyle">
          <span>指定在</span>
          <a-select
            style="width: 240px;margin:0 4px"
            size="small"
            :disabled="hour.opt !== 4"
            v-model:value="hour.fix"
            :options="hourOptions"
            mode="multiple"
            @click="preventDefault"
            :max-tag-count="3"
            @change="hourFixChange"
          />
          <span>小时执行</span>
        </a-radio>
      </a-radio-group>
    </div>
    <div style="font-size:14px;padding:6px" v-show="tab === 'dayOfMonth'">
      <a-radio-group v-model:value="dayOfMonth.opt" @change="dayOfMonthRadioChange">
        <a-radio :value="1" :style="radioStyle">每日</a-radio>
        <a-radio :value="2" :style="radioStyle">
          <span>从</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="31"
            size="small"
            :disabled="dayOfMonth.opt !== 2"
            v-model:value="dayOfMonth.periodStart"
            @change="dayOfMonthPeriodChange"
          />
          <span>-</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="31"
            size="small"
            :disabled="dayOfMonth.opt !== 2"
            v-model:value="dayOfMonth.periodEnd"
            @change="dayOfMonthPeriodChange"
          />
          <span>每日执行一次</span>
        </a-radio>
        <a-radio :value="3" :style="radioStyle">
          <span>从</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="31"
            size="small"
            :disabled="dayOfMonth.opt !== 3"
            v-model:value="dayOfMonth.loopStart"
            @change="dayOfMonthLoopChange"
          />
          <span>开始, 每</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="31"
            size="small"
            :disabled="dayOfMonth.opt !== 3"
            v-model:value="dayOfMonth.loopEnd"
            @change="dayOfMonthLoopChange"
          />
          <span>日执行一次</span>
        </a-radio>
        <a-radio :value="4" :style="radioStyle">
          <span>指定在</span>
          <a-select
            style="width: 240px;margin:0 4px"
            size="small"
            :disabled="dayOfMonth.opt !== 4"
            v-model:value="dayOfMonth.fix"
            :options="dayOfMonthOptions"
            mode="multiple"
            @click="preventDefault"
            :max-tag-count="3"
            @change="dayOfMonthFixChange"
          />
          <span>日执行</span>
        </a-radio>
        <a-radio :value="5" :style="radioStyle">
          <span>不指定</span>
        </a-radio>
      </a-radio-group>
    </div>
    <div style="font-size:14px;padding:6px" v-show="tab === 'month'">
      <a-radio-group v-model:value="month.opt" @change="monthRadioChange">
        <a-radio :value="1" :style="radioStyle">每月</a-radio>
        <a-radio :value="2" :style="radioStyle">
          <span>从</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="12"
            size="small"
            :disabled="month.opt !== 2"
            v-model:value="month.periodStart"
            @change="monthPeriodChange"
          />
          <span>-</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="12"
            size="small"
            :disabled="month.opt !== 2"
            v-model:value="month.periodEnd"
            @change="monthPeriodChange"
          />
          <span>每月执行一次</span>
        </a-radio>
        <a-radio :value="3" :style="radioStyle">
          <span>从</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="12"
            size="small"
            :disabled="month.opt !== 3"
            v-model:value="month.loopStart"
            @change="monthLoopChange"
          />
          <span>开始, 每</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="12"
            size="small"
            :disabled="month.opt !== 3"
            v-model:value="month.loopEnd"
            @change="monthLoopChange"
          />
          <span>月执行一次</span>
        </a-radio>
        <a-radio :value="4" :style="radioStyle">
          <span>指定在</span>
          <a-select
            style="width: 240px;margin:0 4px"
            size="small"
            :disabled="month.opt !== 4"
            v-model:value="month.fix"
            :options="monthOptions"
            mode="multiple"
            @click="preventDefault"
            :max-tag-count="3"
            @change="monthFixChange"
          />
          <span>月执行</span>
        </a-radio>
      </a-radio-group>
    </div>
    <div style="font-size:14px;padding:6px" v-show="tab === 'dayOfWeek'">
      <a-radio-group v-model:value="dayOfWeek.opt" @change="dayOfWeekRadioChange">
        <a-radio :value="5" :style="radioStyle">不指定</a-radio>
        <a-radio :value="2" :style="radioStyle">
          <span>从周</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="7"
            size="small"
            :disabled="dayOfWeek.opt !== 2"
            v-model:value="dayOfWeek.periodStart"
            @change="dayOfWeekPeriodChange"
          />
          <span>-周</span>
          <a-input-number
            style="width: 40px;margin:0 4px"
            :controls="false"
            :min="1"
            :max="7"
            size="small"
            :disabled="dayOfWeek.opt !== 2"
            v-model:value="dayOfWeek.periodEnd"
            @change="dayOfWeekPeriodChange"
          />
          <span>每天执行一次</span>
        </a-radio>
        <a-radio :value="4" :style="radioStyle">
          <span>指定在周</span>
          <a-select
            style="width: 240px;margin:0 4px"
            size="small"
            :disabled="dayOfWeek.opt !== 4"
            v-model:value="dayOfWeek.fix"
            :options="dayOfWeekOptions"
            mode="multiple"
            @click="preventDefault"
            :max-tag-count="3"
            @change="dayOfWeekFixChange"
          />
          <span>执行</span>
        </a-radio>
      </a-radio-group>
    </div>
  </div>
</template>
<script setup>
import { ref, reactive, defineEmits, defineProps } from "vue";
/*
  Cron表达式获取
*/
const props = defineProps(["style", "modelValue"]);
const emit = defineEmits(["update:modelValue"]);
// input style
const istyle = {
  width: "80%"
};
// 单选项style
const radioStyle = {
  display: "flex",
  height: "40px",
  lineHeight: "40px",
  alignItems: "center"
};
// 分钟下拉框选项 0-59
const minuteOptions = ref([]);
for (let i = 0; i <= 59; i++) {
  minuteOptions.value.push({
    value: `${i}`
  });
}
// 小时下拉框选项 0-23
const hourOptions = ref([]);
for (let i = 0; i <= 23; i++) {
  hourOptions.value.push({
    value: `${i}`
  });
}
// 天下拉框选项 1-31
const dayOfMonthOptions = ref([]);
for (let i = 1; i <= 31; i++) {
  dayOfMonthOptions.value.push({
    value: `${i}`
  });
}
// 月下拉框选项 1-12
const monthOptions = ref([]);
for (let i = 1; i <= 12; i++) {
  monthOptions.value.push({
    value: `${i}`
  });
}
// 周x 下拉框 1-7
const dayOfWeekOptions = ref([]);
for (let i = 1; i <= 7; i++) {
  dayOfWeekOptions.value.push({
    value: `${i}`
  });
}
const minute = reactive({
  value: "*",
  opt: 1,
  periodStart: 0,
  periodEnd: 1,
  loopStart: 0,
  loopEnd: 1,
  fix: ["0"]
});
const hour = reactive({
  value: "*",
  opt: 1,
  periodStart: 0,
  periodEnd: 1,
  loopStart: 0,
  loopEnd: 1,
  fix: ["0"]
});
const dayOfMonth = reactive({
  value: "*",
  opt: 1,
  periodStart: 1,
  periodEnd: 1,
  loopStart: 1,
  loopEnd: 1,
  fix: ["1"]
});
const month = reactive({
  value: "*",
  opt: 1,
  periodStart: 1,
  periodEnd: 1,
  loopStart: 1,
  loopEnd: 1,
  fix: ["1"]
});
const dayOfWeek = reactive({
  value: "?",
  opt: 5,
  periodStart: 1,
  periodEnd: 1,
  fix: ["1"]
});
const tab = ref("minute");
const inputFocus = t => {
  tab.value = t;
};

const minuteRadioChange = e => {
  let value = e.target.value;
  switch (value) {
    case 1:
      minute.value = "*";
      break;
    case 2:
      minute.value = `${minute.periodStart}-${minute.periodEnd}`;
      break;
    case 3:
      minute.value = `${minute.loopStart}/${minute.loopEnd}`;
      break;
    case 4:
      minute.value = minute.fix.join(",");
      break;
  }
  emitResult();
};

const monthRadioChange = e => {
  let value = e.target.value;
  switch (value) {
    case 1:
      month.value = "*";
      break;
    case 2:
      month.value = `${month.periodStart}-${month.periodEnd}`;
      break;
    case 3:
      month.value = `${month.loopStart}/${month.loopEnd}`;
      break;
    case 4:
      month.value = month.fix.join(",");
      break;
  }
  emitResult();
};

const hourRadioChange = e => {
  let value = e.target.value;
  switch (value) {
    case 1:
      hour.value = "*";
      break;
    case 2:
      hour.value = `${hour.periodStart}-${hour.periodEnd}`;
      break;
    case 3:
      hour.value = `${hour.loopStart}/${hour.loopEnd}`;
      break;
    case 4:
      hour.value = hour.fix.join(",");
      break;
  }
  emitResult();
};

const dayOfMonthRadioChange = e => {
  let value = e.target.value;
  switch (value) {
    case 1:
      dayOfMonth.value = "*";
      break;
    case 2:
      dayOfMonth.value = `${dayOfMonth.periodStart}-${dayOfMonth.periodEnd}`;
      break;
    case 3:
      dayOfMonth.value = `${dayOfMonth.loopStart}/${dayOfMonth.loopEnd}`;
      break;
    case 4:
      dayOfMonth.value = dayOfMonth.fix.join(",");
      break;
    case 5:
      dayOfMonth.value = "?";
      break;
  }
  emitResult();
};

const dayOfWeekRadioChange = e => {
  let value = e.target.value;
  switch (value) {
    case 5:
      dayOfWeek.value = "?";
      break;
    case 2:
      dayOfWeek.value = `${dayOfWeek.periodStart}-${dayOfWeek.periodEnd}`;
      break;
    case 4:
      dayOfWeek.value = dayOfWeek.fix.join(",");
      break;
  }
  emitResult();
};

const minutePeriodChange = () => {
  if (!minute.periodStart) {
    minute.periodStart = 0;
  }
  if (!minute.periodEnd) {
    minute.periodEnd = 1;
  }
  if (minute.periodStart > minute.periodEnd) {
    minute.periodEnd = minute.periodStart;
  }
  minute.value = `${minute.periodStart}-${minute.periodEnd}`;
  emitResult();
};

const minuteLoopChange = () => {
  if (!minute.loopStart) {
    minute.loopStart = 0;
  }
  if (!minute.loopEnd) {
    minute.loopEnd = 1;
  }
  minute.value = `${minute.loopStart}/${minute.loopEnd}`;
  emitResult();
};

const minuteFixChange = () => {
  if (!minute.fix || minute.fix.length === 0) {
    minute.fix = ["0"];
  }
  minute.fix.sort((a, b) => a - b);
  minute.value = minute.fix.join(",");
  emitResult();
};

const hourPeriodChange = () => {
  if (!hour.periodStart) {
    hour.periodStart = 0;
  }
  if (!hour.periodEnd) {
    hour.periodEnd = 1;
  }
  if (hour.periodStart > hour.periodEnd) {
    hour.periodEnd = hour.periodStart;
  }
  hour.value = `${hour.periodStart}-${hour.periodEnd}`;
  emitResult();
};

const hourLoopChange = () => {
  if (!hour.loopStart) {
    hour.loopStart = 0;
  }
  if (!hour.loopEnd) {
    hour.loopEnd = 1;
  }
  hour.value = `${hour.loopStart}/${hour.loopEnd}`;
  emitResult();
};

const hourFixChange = () => {
  if (!hour.fix || hour.fix.length === 0) {
    hour.fix = ["0"];
  }
  hour.fix.sort((a, b) => a - b);
  hour.value = hour.fix.join(",");
  emitResult();
};

const dayOfMonthPeriodChange = () => {
  if (!dayOfMonth.periodStart) {
    hour.periodStart = 0;
  }
  if (!dayOfMonth.periodEnd) {
    hour.periodEnd = 1;
  }
  if (dayOfMonth.periodStart > dayOfMonth.periodEnd) {
    dayOfMonth.periodEnd = dayOfMonth.periodStart;
  }
  dayOfMonth.value = `${dayOfMonth.periodStart}-${dayOfMonth.periodEnd}`;
  emitResult();
};

const dayOfMonthLoopChange = () => {
  if (!dayOfMonth.loopStart) {
    dayOfMonth.loopStart = 0;
  }
  if (!dayOfMonth.loopEnd) {
    dayOfMonth.loopEnd = 1;
  }
  dayOfMonth.value = `${dayOfMonth.loopStart}/${dayOfMonth.loopEnd}`;
  emitResult();
};

const dayOfMonthFixChange = () => {
  if (!dayOfMonth.fix || dayOfMonth.fix.length === 0) {
    dayOfMonth.fix = ["1"];
  }
  dayOfMonth.fix.sort((a, b) => a - b);
  dayOfMonth.value = dayOfMonth.fix.join(",");
  emitResult();
};

const monthPeriodChange = () => {
  if (!month.periodStart) {
    month.periodStart = 1;
  }
  if (!month.periodEnd) {
    month.periodEnd = 1;
  }
  if (month.periodStart > month.periodEnd) {
    month.periodEnd = month.periodStart;
  }
  month.value = `${month.periodStart}-${month.periodEnd}`;
  emitResult();
};

const monthLoopChange = () => {
  if (!month.loopStart) {
    month.loopStart = 1;
  }
  if (!month.loopEnd) {
    month.loopEnd = 1;
  }
  month.value = `${month.loopStart}/${month.loopEnd}`;
  emitResult();
};

const monthFixChange = () => {
  if (!month.fix || month.fix.length === 0) {
    month.fix = ["1"];
  }
  month.fix.sort((a, b) => a - b);
  month.value = month.fix.join(",");
  emitResult();
};

const dayOfWeekPeriodChange = () => {
  if (!dayOfWeek.periodStart) {
    dayOfWeek.periodStart = 1;
  }
  if (!dayOfWeek.periodEnd) {
    dayOfWeek.periodEnd = 1;
  }
  if (dayOfWeek.periodStart > dayOfWeek.periodEnd) {
    dayOfWeek.periodEnd = dayOfWeek.periodStart;
  }
  dayOfWeek.value = `${dayOfWeek.periodStart}-${dayOfWeek.periodEnd}`;
  emitResult();
};

const dayOfWeekFixChange = () => {
  if (!dayOfWeek.fix || dayOfWeek.fix.length === 0) {
    dayOfWeek.fix = ["1"];
  }
  dayOfWeek.fix.sort((a, b) => a - b);
  dayOfWeek.value = dayOfWeek.fix.join(",");
  emitResult();
};
// 阻止默认事件
const preventDefault = e => {
  e.preventDefault();
  e.stopPropagation();
};
// 修改modelValue
const emitResult = () => {
  let result = `${minute.value} ${hour.value} ${dayOfMonth.value} ${month.value} ${dayOfWeek.value}`;
  emit("update:modelValue", result);
};
// 解析cron表达式
const parseCron = () => {
  // 获取传进来的值
  let cron = props.modelValue;
  if (!cron) {
    return;
  }
  let fields = cron.split(" ");
  if (fields.length !== 5) {
    return;
  }
  let parseResult = parseItem(fields[0]);
  minute.value = parseResult.value;
  minute.opt = parseResult.opt;
  switch (parseResult.opt) {
    case 2:
      minute.periodStart = parseResult.start;
      minute.periodEnd = parseResult.end;
      break;
    case 3:
      minute.loopStart = parseResult.start;
      minute.loopEnd = parseResult.end;
      break;
    case 4:
      minute.fix = parseResult.fix;
      break;
  }
  parseResult = parseItem(fields[1]);
  hour.value = parseResult.value;
  hour.opt = parseResult.opt;
  switch (parseResult.opt) {
    case 2:
      hour.periodStart = parseResult.start;
      hour.periodEnd = parseResult.end;
      break;
    case 3:
      hour.loopStart = parseResult.start;
      hour.loopEnd = parseResult.end;
      break;
    case 4:
      hour.fix = parseResult.fix;
      break;
  }
  parseResult = parseItem(fields[2]);
  dayOfMonth.value = parseResult.value;
  dayOfMonth.opt = parseResult.opt;
  switch (parseResult.opt) {
    case 2:
      dayOfMonth.periodStart = parseResult.start;
      dayOfMonth.periodEnd = parseResult.end;
      break;
    case 3:
      dayOfMonth.loopStart = parseResult.start;
      dayOfMonth.loopEnd = parseResult.end;
      break;
    case 4:
      dayOfMonth.fix = parseResult.fix;
      break;
  }
  parseResult = parseItem(fields[3]);
  month.value = parseResult.value;
  month.opt = parseResult.opt;
  switch (parseResult.opt) {
    case 2:
      month.periodStart = parseResult.start;
      month.periodEnd = parseResult.end;
      break;
    case 3:
      month.loopStart = parseResult.start;
      month.loopEnd = parseResult.end;
      break;
    case 4:
      month.fix = parseResult.fix;
      break;
  }
  parseResult = parseItem(fields[4]);
  dayOfWeek.value = parseResult.value;
  dayOfWeek.opt = parseResult.opt;
  switch (parseResult.opt) {
    case 2:
      dayOfWeek.periodStart = parseResult.start;
      dayOfWeek.periodEnd = parseResult.end;
      break;
    case 4:
      dayOfWeek.fix = parseResult.fix;
      break;
  }
};

const parseItem = item => {
  if (item === "*") {
    return {
      value: item,
      opt: 1
    };
  }
  if (item === "?") {
    return {
      value: item,
      opt: 5
    };
  }
  let s1 = item.split("-");
  if (s1.length === 2) {
    return {
      value: item,
      opt: 2,
      start: s1[0],
      end: s1[1]
    };
  }
  s1 = item.split("/");
  if (s1.length === 2) {
    return {
      value: item,
      opt: 3,
      start: s1[0],
      end: s1[1]
    };
  }
  s1 = item.split(",");
  if (s1.length > 0) {
    return {
      value: item,
      opt: 4,
      fix: s1
    };
  }
  // 未知
  return {
    opt: 6
  };
};
parseCron();
emitResult();
</script>
<style scoped>
.cron-input {
  text-align: center;
}
.cron-input-ul {
  width: 100%;
  display: flex;
  align-items: center;
  margin-top: 20px;
}
.cron-input-ul > li {
  width: 20%;
}
.cron-item {
  margin-top: 4px;
  font-size: 14px;
  width: 100%;
  text-align: center;
}
</style>