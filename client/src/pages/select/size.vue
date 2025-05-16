<template>
  <view class="size-setting-page">
    <view class="title">尺寸设置:</view>

    <view class="setting-card">
      <view class="setting-item">
        <text class="label">预设:</text>
        <picker mode="selector" :range="presetOptions" :value="selectedPresetIndex" @change="onPresetChange">
          <view class="picker-view">
            {{ presetOptions[selectedPresetIndex] }}
            <text class="arrow">▼</text>
          </view>
        </picker>
      </view>

      <view class="setting-item">
        <text class="label">出血线设置:</text>
        <input class="input-field" type="text" v-model="bleedValue" placeholder="例如: 3mm" />
      </view>

      <button class="confirm-button" @click="confirmSettings">确定></button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

const presetOptions = ref(['A4', 'A3', 'B5', '自定义']); // Example presets
const selectedPresetIndex = ref(0); // Default to A4
const bleedValue = ref('3mm'); // Default bleed value

const selectedPreset = computed(() => presetOptions.value[selectedPresetIndex.value]);

const onPresetChange = (event: any) => {
  selectedPresetIndex.value = parseInt(event.detail.value, 10);
  // Add logic if '自定义' is selected, e.g., show custom dimension inputs
  console.log('Preset changed to:', selectedPreset.value);
};

const confirmSettings = () => {
  // TODO: Save the settings (e.g., in store or pass via navigation)
  console.log('Selected Preset:', selectedPreset.value);
  console.log('Bleed Value:', bleedValue.value);

  // Navigate to the next step (e.g., select template page)
  uni.navigateTo({
    url: '/pages/select/select' // Adjust the path as needed
  });
};
</script>

<style scoped>
.size-setting-page {
  padding: 40rpx 30rpx;
  background-color: #f8f8f8;
  min-height: 100vh;
  box-sizing: border-box;
}

.title {
  font-size: 40rpx;
  font-weight: bold;
  margin-bottom: 40rpx;
  color: #333;
}

.setting-card {
  background-color: #ffffff;
  padding: 40rpx 30rpx;
  border-radius: 16rpx;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
  height: fit-content;
}

.setting-item {
  margin-bottom: 40rpx;
}

.label {
  display: block;
  font-size: 30rpx;
  color: #555;
  margin-bottom: 15rpx;
}

.picker-view {
  background-color: #f0f0f0;
  padding: 18rpx 25rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
  color: #333;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.arrow {
  color: #999;
}

.input-field {
  background-color: #f0f0f0;
  padding: 18rpx 25rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
  color: #333;
  width: auto; /* Let it take available width */
  box-sizing: border-box;
}

.confirm-button {
  background-color: #ffffff;
  color: #57DD00; /* Green color from tabbar */
  border: none;
  border-radius: 10rpx;
  /* padding: 20rpx 0; */
  font-size: 30rpx;
  text-align: right;
  margin-top: 20rpx; /* Add some space above the button */
  /* Simulate the > arrow */
  position: relative;
  padding-right: 40rpx; /* Make space for the pseudo-element */
}

.confirm-button::after {
    content: '>';
    position: absolute;
    right: 30rpx; /* Adjust position */
    top: 50%;
    transform: translateY(-50%);
    font-weight: bold;
}


.confirm-button:active {
  background-color: #4cc200; /* Darker shade on press */
}
</style>
