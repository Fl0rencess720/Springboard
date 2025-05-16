<template>
  <view class="preview-page">
    <!-- Custom Navigation Bar Placeholder - Actual implementation might vary -->
    <!-- <view class="custom-navbar">
      <view class="back-icon" @click="goBack">&lt;</view>
      <view class="page-title">内页预览</view>
      <view class="actions">
        <text class="icon">...</text>
        <text class="icon">◎</text>
      </view>
    </view> -->

    <view class="content-area">
      <!-- Loading State -->
      <view v-if="loading">加载预览中...</view>
      <!-- Error State -->
      <view v-else-if="error && !templateDetails">{{ error }}</view> <!-- Show main error only if details failed -->
      <!-- Preview Content -->
      <template v-else-if="templateDetails">
        <!-- Display URL fetch error if it occurred but details loaded -->
        <view v-if="error && templateDetails" style="color: red; margin-bottom: 20rpx; text-align: center;">{{ error }}</view>
        <view class="page-section">
          <view class="section-title">
            <text class="dot"></text>封面
          </view>
          <!-- Use the signed URL ref -->
          <image class="preview-image" :src="coverImageSignedUrl" mode="widthFix"></image>
        </view>

        <view class="page-section">
          <view class="section-title">
            <text class="dot"></text>目录页
          </view>
          <!-- Use the signed URL ref -->
          <image class="preview-image" :src="tocImageSignedUrl" mode="widthFix"></image>
        </view>

        <view class="page-section">
          <view class="section-title">
            <text class="dot"></text>内容页
          </view>
          <!-- Use the signed URL ref -->
          <image class="preview-image" :src="contentImageSignedUrl" mode="widthFix"></image>
        </view>
      </template>
    </view>

    <!-- Confirm Button -->
    <!-- Allow confirmation even if some image URLs failed, adjust if needed -->
    <button class="confirm-button" @click="confirmPreview" :disabled="loading || !templateDetails">
      <image src="/static/icons/confirm.svg" class="button-icon" mode="aspectFit"></image> <!-- Updated icon -->
      确定 <text class="arrow">&gt;</text>
    </button>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'; // Import watch
import { onLoad } from '@dcloudio/uni-app';
import type { Template as TemplateDetails } from '../../types/api'; // Import the detailed Template type
import { fetchTemplateDetails } from '../../api/template'; // Import the API function
import { getPreviewUrl } from '../../api/oss'; // Import the OSS function

const templateId = ref<string | null>(null);
const templateDetails = ref<TemplateDetails | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);

// Refs for signed image URLs
const coverImageSignedUrl = ref<string>('/static/images/template1.png'); // Default placeholder
const tocImageSignedUrl = ref<string>('/static/images/template2.png'); // Default placeholder
const contentImageSignedUrl = ref<string>('/static/images/template3.png'); // Default placeholder

// Function to fetch signed URLs
const fetchSignedUrls = async () => {
  if (!templateDetails.value || !templateDetails.value.pages) return;

  try {
    // Cover Page (assuming index 0)
    const coverKey = templateDetails.value.pages[0]?.oss_key;
    if (coverKey) {
      coverImageSignedUrl.value = await getPreviewUrl(coverKey);
    } else {
      console.warn('Cover page preview_oss_key not found.');
      // Keep default or set specific placeholder if needed
    }

    // TOC Page (assuming index 1)
    const tocKey = templateDetails.value.pages[1]?.oss_key;
    if (tocKey) {
      tocImageSignedUrl.value = await getPreviewUrl(tocKey);
    } else {
      console.warn('TOC page preview_oss_key not found.');
      // Keep default or set specific placeholder if needed
    }

    // Content Page (find first content page)
    const contentPage = templateDetails.value.pages.find(p => p.is_content_page);
    const contentKey = contentPage?.oss_key;
    if (contentKey) {
      contentImageSignedUrl.value = await getPreviewUrl(contentKey);
    } else {
      console.warn('Content page preview_oss_key not found.');
      // Keep default or set specific placeholder if needed
    }
    console.log('Signed URLs fetched:', { cover: coverImageSignedUrl.value, toc: tocImageSignedUrl.value, content: contentImageSignedUrl.value });
  } catch (urlError: any) {
    console.error('Failed to fetch signed URLs:', urlError);
    // Handle URL fetching errors, maybe show a specific message or fallback images
    error.value = `无法加载部分预览图片: ${urlError.message || 'OSS错误'}`;
    uni.showToast({ title: error.value, icon: 'none' });
    // Decide if this error should prevent confirmation - currently it doesn't block the button
  }
};

// Computed properties for preview images based on fetched template details
// Assumes specific page indices for cover, TOC, content. Adjust if needed.
// --- These computed properties are no longer directly used for src, but kept for potential logic ---
const coverImageUrl = computed(() => templateDetails.value?.pages?.[0]?.preview_oss_key || '');
const tocImageUrl = computed(() => templateDetails.value?.pages?.[1]?.preview_oss_key || '');
const contentImageUrl = computed(() => {
  // Find the first content page
  const contentPage = templateDetails.value?.pages?.find(p => p.is_content_page);
  return contentPage?.preview_oss_key || '';
});

onLoad(async (options) => {
  if (options && options.templateId) {
    templateId.value = options.templateId;
    console.log('Previewing template ID:', templateId.value);
    loading.value = true;
    error.value = null;
    try {
      // Fetch actual template details
      // TODO: Implement fetchTemplateDetails in src/api/template.ts
      const details = await fetchTemplateDetails(templateId.value);
      templateDetails.value = details;
      console.log('Fetched template details:', details);
      if (!details || !details.pages || details.pages.length === 0) {
        throw new Error('模板数据不完整');
      }
      // Fetch signed URLs after getting details
      await fetchSignedUrls();
    } catch (err: any) {
      console.error('Failed to fetch template details:', err);
      error.value = `无法加载模板预览: ${err.message || '未知错误'}`;
      uni.showToast({ title: error.value, icon: 'none' });
      // Optionally navigate back
      // uni.navigateBack();
    } finally {
      loading.value = false;
    }
  } else {
    console.error('Template ID not provided to preview page.');
    error.value = '模板信息错误';
    uni.showToast({ title: error.value, icon: 'none' });
    uni.navigateBack();
  }
});

const goBack = () => {
  uni.navigateBack();
};

const confirmPreview = () => {
  if (!templateId.value || loading.value || error.value) {
    uni.showToast({ title: '无法确认模板', icon: 'none' });
    return;
  }
  console.log('Confirming preview, navigating to group management...');
  // Navigate to the next step: Group Management
  uni.navigateTo({
    // Pass necessary info like templateId
    url: `/pages/edit/group?templateId=${templateId.value}`
  });
};
</script>

<style scoped>
.preview-page {
  background-color: #f8f8f8;
  min-height: 100vh;
  padding-bottom: 140rpx; /* Space for the fixed button */
  box-sizing: border-box;
}

/* Placeholder for custom navbar styling if needed */
/* .custom-navbar { ... } */

.content-area {
  padding: 30rpx;
}

.page-section {
  margin-bottom: 40rpx;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 20rpx;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
}

.section-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 20rpx;
  display: flex;
  align-items: center;
}

.dot {
  width: 12rpx;
  height: 12rpx;
  background-color: #57DD00; /* Green dot */
  border-radius: 50%;
  margin-right: 15rpx;
}

.preview-image {
  width: 100%;
  display: block;
  border-radius: 8rpx; /* Slightly rounded corners for image inside section */
  /* box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.08); */
  background-color: #eee; /* Placeholder bg */
}

.confirm-button {
  position: fixed;
  bottom: 40rpx; /* Adjust as needed */
  left: 50%;
  transform: translateX(-50%);
  width: calc(100% - 80rpx); /* Full width minus padding */
  max-width: 600rpx; /* Max width for larger screens */
  background-color: #333; /* Dark background */
  color: #ffffff;
  border: none;
  border-radius: 40rpx; /* Pill shape */
  padding: 20rpx 0;
  font-size: 30rpx;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.1);
}

.confirm-button:disabled {
  background-color: #ccc;
  color: #666;
}

.button-icon {
  width: 36rpx; /* Adjusted size */
  height: 36rpx; /* Adjusted size */
  margin-right: 12rpx; /* Adjusted spacing */
  /* filter: brightness(0) invert(1); */ /* Uncomment if SVG needs to be white */
}

.arrow {
  margin-left: 10rpx;
  font-weight: bold;
}

</style>
