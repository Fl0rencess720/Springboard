<template>
  <!-- Main container with padding for status bar -->
  <view class="custom-navbar" :style="{ paddingTop: statusBarHeight + 'px' }">
    <!-- Content area with fixed height -->
    <view class="navbar-content" :style="{ height: navBarHeight + 'px' }">
      <image class="logo" src="/static/icons/head-logo.png" mode="aspectFit" />
      <view class="navbar-actions">
        <!-- Use image tag for icons in uni-app -->
        <image class="icon" src="/static/icons/tips.svg" mode="aspectFit" @click="showTips" />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue';

// Get system info to determine heights
const systemInfo = uni.getSystemInfoSync();
const statusBarHeight = ref(systemInfo.statusBarHeight || 0); // Height of status bar
const navBarHeight = ref(44); // Standard height for navbar content (adjust if needed)

const showTips = () => {
  // Use uni.showToast for consistency
  uni.showToast({ title: 'Here are some tips!', icon: 'none' });
};

</script>

<style scoped>
.custom-navbar {
  background-color: #fff; /* Or your desired background */
  box-sizing: border-box;
  width: 100%;
  /* The padding-top is set dynamically */
}

.navbar-content {
  display: flex;
  align-items: center;
  /* Change justification to start */
  justify-content: flex-start;
  padding: 0 30rpx 0 0; /* Keep right padding for the actions */
  box-sizing: border-box;
  /* The height is set dynamically */
}

.logo {
  height: 36px; /* Keep fixed height */
  width: 200rpx; /* Set an explicit width in rpx */
  flex-shrink: 0; /* Keep this to prevent shrinking */
}

.navbar-actions {
  display: flex;
  align-items: center;
  /* Push actions to the right */
  margin-left: auto;
}

.icon {
  width: 40rpx;
  height: 40rpx;
  /* Add margin if needed */
  margin-left: 20rpx;
}
</style>