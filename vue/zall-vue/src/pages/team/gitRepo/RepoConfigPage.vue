<template>
  <div style="padding:14px">
    <div class="container">
      <div class="section">
        <div class="section-title">
          <span>仓库优化</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <div class="input-title">仓库瘦身</div>
            <a-button type="primary">触发仓库GC</a-button>
            <div class="input-desc">执行git gc命令瘦身仓库</div>
          </div>
          <div class="input-item">
            <div class="input-title">重建hook</div>
            <a-button type="primary">重建hook</a-button>
            <div class="input-desc">将原有的hook脚本重新初始化</div>
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
            <a-input />
            <div class="input-desc"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
</script>
<style scoped>
</style>