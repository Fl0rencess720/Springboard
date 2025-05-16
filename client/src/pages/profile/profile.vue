<template>
  <view class="profile-page">
    <!-- Header -->
    <view class="header">
      <view class="user-info">
        <image class="avatar" src="/static/icons/app-logo.png" mode="aspectFit"></image>
        <view class="user-details">
          <text class="username">Springbroad</text>
          <!-- TODO: Fetch and display actual user info -->
          <text class="uid">uid: loading...</text>
        </view>
      </view>
      <view class="actions">
        <!-- Placeholder for icons -->
        <img class="icon" src="/static/icons/moon.svg" mode="aspectFit" @click="switchTheme" />
        <img class="icon" src="/static/icons/question.svg" mode="aspectFit" @click="feedback" />
      </view>
    </view>

    <!-- My Portfolio Section -->
    <view class="section">
      <text class="section-title">我的作品集</text>
      <view v-if="loadingPortfolios" class="loading-state">加载中...</view>
      <view v-else class="grid">
        <view class="grid-item" v-for="item in portfolioItems" :key="item.uid">
          <!-- Display preview image from the first page of the template -->
          <image class="item-image" :src="item.template?.pages?.[0]?.preview_oss_key || '/static/images/template1.png'" mode="aspectFill"></image>
          <text class="portfolio-title">{{ item.title }}</text>
        </view>
        <!-- Add placeholder if needed -->
        <view v-if="portfolioItems.length % 2 !== 0" class="grid-item placeholder"></view>
      </view>
      <view v-if="!loadingPortfolios && portfolioItems.length === 0" class="empty-state">
        暂无作品集
      </view>
    </view>

    <!-- History Section -->
    <view class="section">
      <text class="section-title">历史使用</text>
      <view v-if="loadingHistory" class="loading-state">加载中...</view>
      <view v-else class="grid">
        <view class="grid-item" v-for="item in historyItems" :key="item.uid">
          <!-- History items use oss_key directly as per GetHistoryTemplatesResponseItem -->
          <image class="item-image" :src="item.oss_key || '/static/images/template1.png'" mode="aspectFill"></image>
          <!-- Optional: Display template name -->
          <!-- <text class="template-name">{{ item.name }}</text> -->
        </view>
        <!-- Add placeholders if needed -->
        <view v-for="i in (2 - (historyItems.length % 2)) % 2" :key="'history-placeholder-' + i" class="grid-item placeholder"></view>
      </view>
      <view v-if="!loadingHistory && historyItems.length === 0" class="empty-state">
        暂无历史记录
      </view>
    </view>

    <tabbar current-page="/pages/profile/profile" />
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Tabbar from '../../components/tabbar/tabbar.vue';
import { fetchMyPortfolios, fetchHistoryTemplates } from '../../api/portfolio';
import type { GetMyPortfolioResponseItem, GetHistoryTemplatesResponseItem } from '../../types/api';
import { login } from '../../api/auth'; // Import login function for testing

const portfolioItems = ref<GetMyPortfolioResponseItem[]>([]);
const historyItems = ref<GetHistoryTemplatesResponseItem[]>([]);
const loadingPortfolios = ref(true);
const loadingHistory = ref(true);

onMounted(async () => {
  try {
    await login();
    console.log('Login successful or token already valid.');

    // Fetch data in parallel
    loadingPortfolios.value = true;
    loadingHistory.value = true;
    const [portfolioData, historyData] = await Promise.all([
      fetchMyPortfolios(),
      fetchHistoryTemplates()
    ]);

    portfolioItems.value = portfolioData;
    historyItems.value = historyData;
    console.log('Fetched portfolios:', portfolioItems.value);
    console.log('Fetched history templates:', historyItems.value);

  } catch (error) {
    console.error('Failed to fetch profile data or login:', error);
    uni.showToast({ title: '加载数据失败', icon: 'none' });
  } finally {
    loadingPortfolios.value = false;
    loadingHistory.value = false;
  }
});

const switchTheme = () => {
  // TODO: Implement theme switching logic
  console.log('Switch theme clicked');
};

const feedback = () => {
  // TODO: Navigate to feedback page or open feedback modal
  console.log('Feedback clicked');
  // Example navigation:
  // uni.navigateTo({ url: '/pages/feedback/feedback' });
};

</script>

<style lang="scss" scoped>
.profile-page {
  padding: 20rpx;
  background-color: #f8f8f8; // Adjust background color as needed
  min-height: 100vh;
  /* Add padding to the bottom to prevent content from being hidden by the fixed tabbar */
  padding-bottom: 50px; /* Match the height of the tabbar */
  /* Add safe area padding if not already handled globally */
  padding-bottom: calc(50px + constant(safe-area-inset-bottom));
  padding-bottom: calc(50px + env(safe-area-inset-bottom));
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #ffffff;
  padding: 30rpx 20rpx;
  border-radius: 16rpx;
  margin-bottom: 30rpx;
  box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.user-info {
  display: flex;
  align-items: center;
}

.avatar {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  margin-right: 20rpx;
  background-color: #eee; // Placeholder background
}

.user-details {
  display: flex;
  flex-direction: column;
}

.username {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
}

.uid {
  font-size: 24rpx;
  color: #999;
}

.actions {
  display: flex;
  align-items: center;
  height: 100rpx;
}

.icon {
  width: 50rpx;
  height: 50rpx;
  margin-left: 30rpx;
}

.section {
  margin-bottom: 40rpx;
}

.section-title {
  display: block;
  font-size: 32rpx;
  font-weight: bold;
  margin-bottom: 20rpx;
  color: #333;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
}

.grid-item {
  background-color: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);
  position: relative;
  // Aspect ratio box trick
  padding-top: 60%; // Adjust this percentage based on desired aspect ratio

  &.placeholder {
    background-color: #f0f0f0; // Lighter background for placeholder
    box-shadow: none;
  }

  display: flex; /* Use flexbox for content alignment */
  flex-direction: column; /* Stack image and title vertically */
  justify-content: flex-start; /* Align content to the top */
  padding: 0; /* Remove padding if aspect ratio trick is used */
}

.item-image {
  // If not using aspect ratio trick, define size explicitly
  // width: 100%;
  // height: 200rpx; /* Example fixed height */
  // object-fit: cover;

  // If using aspect ratio trick, keep position absolute
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.portfolio-title {
  position: absolute; /* Position title over the image */
  bottom: 0;
  left: 0;
  width: 100%;
  background-color: rgba(0, 0, 0, 0.5); /* Semi-transparent background */
  color: white;
  font-size: 24rpx;
  padding: 8rpx 12rpx;
  box-sizing: border-box;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-state {
  text-align: center;
  color: #999;
  font-size: 28rpx;
  padding: 40rpx 0;
}

.loading-state {
  text-align: center;
  color: #999;
  font-size: 28rpx;
  padding: 40rpx 0;
}

.template-name {
    position: absolute; /* Position name over the image */
    bottom: 0;
    left: 0;
    width: 100%;
    background-color: rgba(0, 0, 0, 0.5); /* Semi-transparent background */
    color: white;
    font-size: 24rpx;
    padding: 8rpx 12rpx;
    box-sizing: border-box;
    text-align: center;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

</style>