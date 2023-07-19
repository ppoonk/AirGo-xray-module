<template>
	<div class="layout-footer pb15">
		<div class="layout-footer-warp">
      <el-link>Powered by AirGo</el-link>
      <span> | </span>
      <el-link href="http://qm.qq.com/cgi-bin/qm/qr?_wv=1027&k=o2OuceaPMla0734GE8U4Es9havy4hdni&authKey=yY%2BB3JSK52JuFXKe6s6Be6mkpk8nBRZRUy9yXDy9LNV9%2FrHDOFMu3e8DEs88ohJp&noverify=0&group_code=144545032">加入我们</el-link>
		</div>
	</div>
</template>

<script setup lang="ts" name="layoutFooter">
// 此处需有内容（注释也得），否则缓存将失败
</script>

<style scoped lang="scss">
.layout-footer {
	width: 100%;
	display: flex;
	&-warp {
		margin: auto;
		color: var(--el-text-color-secondary);
		text-align: center;
		animation: error-num 0.3s ease;
	}
}
</style>
