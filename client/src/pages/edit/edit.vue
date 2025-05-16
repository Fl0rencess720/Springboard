<template>
  <view class="edit-page">
    <!-- Top Navigation Placeholder -->
    <!-- <view class="top-nav"> ... </view> -->

    <!-- Group Selection Row -->
    <scroll-view class="group-scroll" scroll-x="true" show-scrollbar="false">
      <view class="group-list">
        <view
          v-for="(project, index) in projects"
          :key="project.uid"
          :class="['group-item', { active: selectedGroupIndex === index }]"
          @click="selectGroup(index)"
        >
          <!-- Placeholder for group thumbnail -->
          <view class="thumbnail-placeholder group-thumb"></view>
          <text class="group-name">{{ project.name }}</text>
        </view>
      </view>
    </scroll-view>

    <!-- Work Selection Row -->
    <scroll-view class="work-scroll" scroll-x="true" show-scrollbar="false">
      <view class="work-list">
        <view
          v-for="(work, index) in currentWorks"
          :key="work.oss_key || index"
          :class="['work-item', { active: selectedWorkIndex === index }]"
          @click="selectWork(index)"
        >
          <!-- Placeholder for work thumbnail -->
          <view class="thumbnail-placeholder work-thumb"></view>
          <text class="work-page">P{{ work.page }}</text>
        </view>
         <!-- Add Work Button within scroll -->
        <view class="work-item add-work-thumb" @click="addWork">
          <view class="thumbnail-placeholder work-thumb add-placeholder">
            <text class="add-icon-small">+</text>
          </view>
           <text class="work-page">添加</text>
        </view>
      </view>
    </scroll-view>

    <!-- Toolbar -->
    <view class="toolbar">
      <view class="tool-left">
        <text class="tool-icon">T</text> <!-- Text Tool -->
        <image class="tool-icon image-icon" src="/static/icons/image.svg" mode="aspectFit"></image> <!-- Image Tool -->
      </view>
      <view class="tool-right">
        <text class="bleed-label">出血线</text>
        <switch :checked="bleedEnabled" @change="toggleBleed" color="#57DD00" style="transform:scale(0.7)" />
      </view>
    </view>

    <!-- Editor Area -->
    <view class="editor-container">
       <EditorComponent
         class="editor-component-wrapper"
         v-if="selectedWork"
         :key="selectedWork.oss_key" 
         :work-data="selectedWork"
         :bleed-enabled="bleedEnabled"
         :image-url="editorImageUrl"
         @update:workData="handleEditorUpdate"
       />
       <view v-else class="editor-placeholder">
         <text>请选择或添加页面</text>
       </view>
    </view>

    <!-- Action Buttons -->
    <view class="action-buttons">
      <button class="action-btn add-btn" @click="addWork" :disabled="isUploading">
        <text class="btn-icon">+</text>
      </button>
      <button class="action-btn save-btn" @click="savePortfolio" :disabled="isUploading">
         <image class="btn-icon save-icon" src="/static/icons/icon.svg#upload" mode="aspectFit"></image>
      </button>
    </view>

  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { onLoad } from '@dcloudio/uni-app';
import EditorComponent from '../../components/editor/editor.vue'; // Import the editor component
import type { Project, Work } from '../../types/api'; // Import types
// Import OSS functions
import { uploadWorkFile, getFileSuffix, getPreviewUrl } from '../../api/oss';
import { savePortfolio as apiSavePortfolio } from '../../api/portfolio'; // Import the actual save function

// Define the structure for groups passed from the previous page
interface PassedGroup {
  name: string;
  order: number;
  // Add uid if it's generated/passed from group page, otherwise generate here
  uid?: string;
}

const templateId = ref<string | null>(null);
const portfolioUid = ref<string | null>(null); // To store the portfolio UID after first save
const portfolioTitle = ref('我的作品集'); // Default or fetch existing title
const projects = ref<Project[]>([]); // Use the Project type from api.d.ts
const selectedGroupIndex = ref(0); // Index of the selected project/group
const selectedWorkIndex = ref(0); // Index of the selected work within the project
const bleedEnabled = ref(false); // State for the bleed line toggle
const isUploading = ref(false); // Flag for upload status

// Function to generate a simple unique ID (replace with a proper UUID library if needed)
const generateUid = () => `uid_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;

onLoad((options) => {
  console.log('Edit page loaded with options:', options);
  if (options && options.templateId && options.groups) {
    templateId.value = options.templateId;
    try {
      const passedGroups: PassedGroup[] = JSON.parse(decodeURIComponent(options.groups));
      console.log('Parsed groups:', passedGroups);

      // Initialize projects based on passed groups
      projects.value = passedGroups.map((g, index) => ({
        uid: g.uid || generateUid(), // Generate UID if not present
        name: g.name,
        order: g.order,
        portfolio_uid: portfolioUid.value || '', // Will be set later if new
        works: [], // Start with empty works, user will add them
      }));

      console.log('Initialized projects:', projects.value);

      // Select the first group and first work by default
      selectGroup(0);

    } catch (e) {
      console.error('Failed to parse groups or initialize:', e);
      uni.showToast({ title: '加载分组信息失败', icon: 'none' });
      uni.navigateBack();
    }
  } else {
    console.error('Template ID or groups not provided to edit page.');
    uni.showToast({ title: '无法加载编辑信息', icon: 'none' });
    uni.navigateBack();
  }
});

// Computed property for the currently selected project's works
const currentWorks = computed(() => {
  if (projects.value.length > 0 && selectedGroupIndex.value < projects.value.length) {
    return projects.value[selectedGroupIndex.value].works;
  }
  return [];
});

// Computed property for the currently selected work
const selectedWork = computed(() => {
  const works = currentWorks.value;
  if (works.length > 0 && selectedWorkIndex.value < works.length) {
    return works[selectedWorkIndex.value];
  }
  return null; // Or a default placeholder work object
});

const selectGroup = (index: number) => {
  if (index >= 0 && index < projects.value.length) {
    selectedGroupIndex.value = index;
    selectedWorkIndex.value = 0; // Reset work selection when group changes
    console.log('Selected group:', projects.value[index].name);
  } else {
    console.warn(`Invalid group index: ${index}`);
  }
};

const selectWork = (index: number) => {
   const works = currentWorks.value;
   if (index >= 0 && index < works.length) {
       selectedWorkIndex.value = index;
       console.log('Selected work page:', works[index].page);
   } else {
       console.warn(`Invalid work index: ${index} for group ${selectedGroupIndex.value}`);
   }
};

const toggleBleed = () => {
  bleedEnabled.value = !bleedEnabled.value;
  console.log('Bleed enabled:', bleedEnabled.value);
};

// Modified addWork to trigger image selection and upload
const addWork = async () => {
  if (isUploading.value) {
    uni.showToast({ title: '正在上传...', icon: 'loading' });
    return;
  }
  if (projects.value.length === 0 || selectedGroupIndex.value >= projects.value.length) {
    console.error('Cannot add work, no project selected or projects array is empty.');
    uni.showToast({ title: '请先选择分组', icon: 'none' });
    return;
  }

  const currentProject = projects.value[selectedGroupIndex.value];

  try {
    // 1. Choose Image
    const chooseRes = await uni.chooseImage({
      count: 1, // Allow only one image selection
      sizeType: ['original', 'compressed'], // Specify image size types
      sourceType: ['album', 'camera'], // Specify image sources
    });

    if (chooseRes && chooseRes.tempFiles && chooseRes.tempFiles.length > 0) {
      const tempFile = chooseRes.tempFiles[0];
      const filePath = chooseRes.tempFilePaths[0];
      const fileSuffix = getFileSuffix(filePath);

      // Show loading state
      isUploading.value = true;
      uni.showLoading({ title: '上传中...' });

      // 2. Prepare Work Data (Generate UID now)
      const workUid = generateUid();
      const portfolioId = portfolioUid.value || 'new'; // Use 'new' or similar if portfolio not saved yet
      const projectUid = currentProject.uid;

      // 3. Upload File
      const ossKey = await uploadWorkFile(
        tempFile, // Pass the file object
        portfolioId, // Use placeholder if new, actual UID if exists
        projectUid,
        workUid,
        fileSuffix
      );

      // 4. Create New Work Object
      const nextPageNumber = currentProject.works.length + 1;
      const newWork: Work = {
        uid: workUid, // Use the generated UID
        oss_key: ossKey,
        project_uid: projectUid,
        size: "1920x1080", // TODO: Get actual size or use template default
        margin_top: "0",
        margin_left: "0",
        bleed: ["3", "3", "3", "3"], // TODO: Use template default
        page: nextPageNumber,
      };

      // 5. Add to Project and Select
      currentProject.works.push(newWork);
      selectedWorkIndex.value = currentProject.works.length - 1;
      console.log('Added and uploaded new work:', newWork);
      uni.showToast({ title: '添加成功', icon: 'success' });

    } else {
      console.log('User cancelled image selection.');
    }
  } catch (error) {
    console.error('Failed to add or upload work:', error);
    uni.showToast({ title: '添加或上传失败', icon: 'none' });
  } finally {
    isUploading.value = false;
    uni.hideLoading();
  }
};

// Modified savePortfolio to use the imported API function
const savePortfolio = async () => {
  console.log('Saving portfolio...');
  if (isUploading.value) {
      uni.showToast({ title: '请等待上传完成', icon: 'none' });
      return;
  }
  if (!templateId.value) {
    uni.showToast({ title: '模板信息丢失', icon: 'none' });
    return;
  }

  // Ensure all works have an oss_key before saving
  const hasMissingOssKey = projects.value.some(proj =>
    proj.works.some(work => !work.oss_key)
  );
  if (hasMissingOssKey) {
      uni.showToast({ title: '部分图片未上传或处理中', icon: 'none' });
      return;
  }

  // Construct the payload
  const payload = {
    uid: portfolioUid.value || "",
    title: portfolioTitle.value,
    template_uid: templateId.value,
    projects: projects.value.map(proj => ({
      uid: proj.uid,
      name: proj.name,
      order: proj.order,
      portfolio_uid: portfolioUid.value || '',
      works: proj.works.map(work => ({
        uid: work.uid,
        oss_key: work.oss_key,
        project_uid: work.project_uid,
        size: work.size,
        margin_top: work.margin_top,
        margin_left: work.margin_left,
        bleed: work.bleed,
        page: work.page,
      }))
    }))
  };

  console.log('Payload for save:', JSON.stringify(payload, null, 2));
  uni.showLoading({ title: '保存中...' });

  try {
    const response = await apiSavePortfolio(payload);
    console.log('Save response:', response);
    if (response && response.uid) {
        portfolioUid.value = response.uid;
        projects.value.forEach(proj => proj.portfolio_uid = response.uid);
        uni.showToast({ title: '保存成功', icon: 'success' });
    } else {
        uni.showToast({ title: '保存成功 (无UID返回)', icon: 'success' });
    }

  } catch (error: any) {
    console.error('Failed to save portfolio:', error);
    uni.showToast({ title: error?.message || '保存失败，请稍后重试', icon: 'none' });
  } finally {
      uni.hideLoading();
  }
};

// Handle updates from the editor component
const handleEditorUpdate = (updatedWorkData: Partial<Work>) => {
    if (selectedWork.value) {
        console.log('Updating work data from editor:', updatedWorkData);
        Object.assign(selectedWork.value, updatedWorkData);
    }
};

// Add logic to get preview URL for the editor
const editorImageUrl = computed(async () => {
    if (selectedWork.value && selectedWork.value.oss_key) {
        try {
            return await getPreviewUrl(selectedWork.value.oss_key);
        } catch (error) {
            console.error("Failed to get preview URL for editor:", error);
            return null;
        }
    }
    return null;
});

</script>

<style scoped>
.edit-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f0f0f0; /* Light background */
}

/* Group Selection */
.group-scroll {
  background-color: #ffffff;
  padding: 15rpx 20rpx;
  white-space: nowrap;
  box-shadow: 0 2rpx 5rpx rgba(0, 0, 0, 0.05);
}
.group-list {
  display: inline-flex; /* Use inline-flex for horizontal layout */
  gap: 20rpx;
}
.group-item {
  display: inline-block; /* Ensure items stay inline */
  text-align: center;
  padding: 10rpx;
  border-radius: 8rpx;
  border: 3rpx solid transparent;
}
.group-item.active {
  border-color: #57DD00; /* Green border for active */
}
.thumbnail-placeholder {
  background-color: #e0e0e0;
  border-radius: 8rpx;
  margin-bottom: 8rpx;
}
.group-thumb {
  width: 160rpx; /* Adjust size */
  height: 100rpx; /* Adjust size */
}
.group-name {
  font-size: 24rpx;
  color: #555;
  display: block;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 160rpx;
}
.group-item.active .group-name {
    color: #333;
    font-weight: bold;
}


/* Work Selection */
.work-scroll {
  background-color: #f8f8f8;
  padding: 15rpx 20rpx;
  white-space: nowrap;
}
.work-list {
  display: inline-flex;
  gap: 15rpx;
  align-items: center; /* Align items vertically */
}
.work-item {
  display: inline-block;
  text-align: center;
  padding: 8rpx;
  border-radius: 8rpx;
  border: 3rpx solid transparent;
}
.work-item.active {
  border-color: #57DD00;
}
.work-thumb {
  width: 120rpx; /* Smaller than group thumbs */
  height: 75rpx;
}
.work-page {
  font-size: 22rpx;
  color: #666;
}
.work-item.active .work-page {
    color: #333;
    font-weight: bold;
}
.add-work-thumb {
    cursor: pointer;
}
.add-placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #d0d0d0;
}
.add-icon-small {
    font-size: 40rpx;
    color: #ffffff;
    font-weight: bold;
}


/* Toolbar */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15rpx 30rpx;
  background-color: #ffffff;
  border-bottom: 1px solid #eee;
}
.tool-left {
  display: flex;
  align-items: center;
  gap: 40rpx;
}
.tool-icon {
  font-size: 40rpx;
  color: #333;
  font-weight: bold;
}
.image-icon {
    width: 40rpx;
    height: 40rpx;
    /* Assuming SVG needs fill */
    filter: invert(20%) sepia(0%) saturate(0%) hue-rotate(0deg) brightness(95%) contrast(85%);
}

.tool-right {
  display: flex;
  align-items: center;
  gap: 10rpx;
}
.bleed-label {
  font-size: 26rpx;
  color: #555;
}

/* Editor Area */
.editor-container {
  flex-grow: 1; /* Takes remaining space */
  padding: 20rpx;
  overflow-y: auto; /* Allow scrolling if content overflows */
  display: flex;
  justify-content: center;
  align-items: center; /* Center placeholder */
}
.editor-placeholder {
    color: #999;
    font-size: 28rpx;
}
.editor-component-wrapper {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

/* Action Buttons */
.action-buttons {
  position: fixed;
  bottom: 40rpx; /* Adjust as needed */
  right: 40rpx; /* Adjust as needed */
  display: flex;
  flex-direction: column;
  gap: 25rpx;
  z-index: 10;
}
.action-btn {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background-color: #333;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.2);
  border: none;
  padding: 0; /* Reset padding */
  line-height: 1; /* Reset line-height */
}
.action-btn:active {
    background-color: #555;
}
.btn-icon {
  font-size: 50rpx;
  font-weight: bold;
}
.save-icon {
    width: 45rpx;
    height: 45rpx;
    filter: brightness(0) invert(1); /* Make SVG white */
}

/* Add styles for loading/disabled states if needed */
.add-btn[disabled] {
    background-color: #999;
    cursor: not-allowed;
}
.save-btn[disabled] {
    background-color: #999;
    cursor: not-allowed;
}

</style>
