<template>
  <view class="content index-page">
    <!-- Add the custom navbar component here -->
    <custom-navbar />

    <!-- Header -->
    <!-- <view class="header">
      <image class="logo" src="/static/icons/head-logo.png" mode="aspectFit"></image>
      <view class="header-actions">
      </view>
    </view> -->

    <!-- Hot Templates Section -->
    <view class="section hot-section">
      <view class="section-header">
        <text class="section-title">热门模板</text>
        <text class="section-subtitle">更多人选择的优秀模板</text>
      </view>
      <scroll-view class="template-scroll" scroll-x="true" show-scrollbar="false">
        <view class="template-list">
          <view
            v-for="template in hotTemplates"
            :key="template.uid"
            class="template-item hot-item"
            @click="navigateToPreview(template.uid)"
          >
            <image class="template-image hot-image" :src="template.oss_key" mode="aspectFill"></image>
            <!-- Optional: Add template name or other info -->
          </view>
        </view>
      </scroll-view>
      <!-- Optional: Add pagination dots here -->
    </view>

    <!-- Recommended Templates Section -->
    <view class="section recommended-section">
       <view class="section-header">
        <text class="section-title">推荐模板</text>
        <text class="section-subtitle">Springbroad精品模板</text>
      </view>
      <view class="template-grid">
        <view
          v-for="template in recommendedTemplates"
          :key="template.uid"
          class="template-item recommended-item"
          @click="navigateToPreview(template.uid)"
        >
          <image class="template-image recommended-image" :src="template.oss_key" mode="aspectFill"></image>
           <!-- Optional: Add template name or other info -->
           <!-- <text class="template-name">{{ template.name }}</text> -->
        </view>
        <!-- Add placeholders for grid alignment -->
        <view v-for="i in (2 - (recommendedTemplates.length % 2)) % 2" :key="'placeholder-' + i" class="template-item recommended-item placeholder"></view>
      </view>
    </view>

    <!-- Add the custom tab bar -->
    <tabbar current-page="/pages/index/index" />
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
// Import the custom tab bar component
import Tabbar from '../../components/tabbar/tabbar.vue'; // Adjust path if necessary
// Import the custom navbar component
import CustomNavbar from '../../components/custom-navbar/custom-navbar.vue'; // Adjust path if necessary
import type { Template } from '../../types/api'; // Import Template type
// Import API functions
import { fetchHotTemplates, fetchAllTemplates } from '../../api/template';
import { login } from '../../api/auth'; // Import login function for testing

// Template data refs
const hotTemplates = ref<Template[]>([]);
const recommendedTemplates = ref<Template[]>([]);

onMounted(async () => {
  try {
    await login();
    console.log('Login successful or token already valid.');

    // Fetch templates in parallel
    const [hotData, allData] = await Promise.all([
      fetchHotTemplates(),
      fetchAllTemplates()
    ]);

    hotTemplates.value = hotData;
    // Use all templates for recommended for now, maybe slice or filter later
    recommendedTemplates.value = allData;

    console.log('Fetched hot templates:', hotTemplates.value);
    console.log('Fetched recommended (all) templates:', recommendedTemplates.value);

  } catch (error) {
    console.error('Failed to fetch templates or login:', error);
    uni.showToast({ title: '加载模板失败', icon: 'none' });
  }
});

const navigateToPreview = (templateId: string) => {
  console.log('Navigating to preview for template ID:', templateId);
  uni.navigateTo({
    url: `/pages/select/preview?templateId=${templateId}`
  });
};

// Placeholder for potential future actions
const handleMoreActions = () => {
  console.log('More actions clicked');
};

const handleViewModeToggle = () => {
  console.log('View mode toggle clicked');
};

</script>

<style scoped>
.index-page {
  background-color: #ffffff; /* White background for the page */
  min-height: 100vh;
  padding: 20rpx 30rpx; /* Add padding */
  /* Ensure content doesn't hide behind tabbar */
  padding-bottom: 50px; /* Match the height of the tabbar */
  padding-bottom: calc(50px + constant(safe-area-inset-bottom)); /* iOS safe area */
  padding-bottom: calc(50px + env(safe-area-inset-bottom)); /* iOS safe area */
  box-sizing: border-box;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 40rpx;
  padding-top: 20rpx; /* Add some space at the top */
}

.logo {
  height: 50rpx; /* Adjust logo size */
  width: 250rpx; /* Adjust logo width */
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 30rpx;
}

.icon {
  font-size: 40rpx;
  color: #333;
  /* Basic styling for icons, replace with actual icons/images if needed */
}
.more-icon {
    font-weight: bold;
    letter-spacing: -4rpx; /* Make dots closer */
}
.view-icon {
    font-weight: bold;
    /* Style for the circle icon */
}


.section {
  margin-bottom: 50rpx;
}

.section-header {
    margin-bottom: 20rpx;
}

.section-title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
  display: block; /* Ensure it takes full width */
}

.section-subtitle {
    font-size: 26rpx;
    color: #999;
    margin-top: 5rpx;
    display: block;
}

/* Hot Templates Scroll */
.template-scroll {
  white-space: nowrap;
}

.template-list {
  display: inline-flex; /* Use inline-flex for horizontal layout */
  gap: 25rpx;
}

.template-item {
  display: inline-block; /* Ensure items stay inline */
  border-radius: 16rpx;
  overflow: hidden;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: transform 0.2s ease-in-out;
}
.template-item:active {
    transform: scale(0.98);
}

.hot-item {
  width: 450rpx; /* Adjust width for horizontal scroll items */
}

.template-image {
  display: block;
  width: 100%;
  background-color: #eee; /* Placeholder bg */
}

.hot-image {
  height: 250rpx; /* Adjust height for hot templates */
}

/* Recommended Templates Grid */
.template-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 25rpx;
}

.recommended-image {
  height: 200rpx; /* Adjust height for recommended templates */
}

.recommended-item.placeholder {
  background-color: transparent;
  box-shadow: none;
}

/* Remove default content styles if they conflict */
.content {
  display: block; /* Override flex settings from original code */
  align-items: initial;
  justify-content: initial;
}
.text-area, .title {
    display: none; /* Hide the old title */
}

</style>
