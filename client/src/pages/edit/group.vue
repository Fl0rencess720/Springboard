<template>
  <view class="group-management-page">
    <view class="title">编辑分组:</view>

    <view class="group-card">
      <view class="group-list">
        <!-- Fixed Cover Item -->
        <view class="group-item fixed-item">
          <text>封面</text>
        </view>

        <!-- Dynamic Group Items -->
        <view v-for="(group, index) in groups" :key="group.id" class="group-item editable-item">
          <input class="group-input" v-model="group.name" placeholder="分组名称" />
          <text class="delete-icon" @click="removeGroup(index)">×</text>
        </view>

        <!-- Add Group Button -->
        <view class="group-item add-item" @click="addGroup">
          <text class="add-icon">+</text>
          <text>添加分组</text>
        </view>

        <!-- Fixed Back Cover Item -->
        <view class="group-item fixed-item">
          <text>封底</text>
        </view>
      </view>
    </view>

    <!-- Confirm Button -->
    <button class="confirm-button" @click="confirmGroups">
      确定 <text class="arrow">&gt;</text>
    </button>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { onLoad } from '@dcloudio/uni-app';

interface Group {
  id: number; // Simple ID for v-for key
  name: string;
}

const templateId = ref<string | null>(null);
const groups = ref<Group[]>([
  { id: Date.now() + 1, name: '分组1' },
  { id: Date.now() + 2, name: '分组2' },
]);
let nextGroupId = Date.now() + 3;

onLoad((options) => {
  if (options && options.templateId) {
    templateId.value = options.templateId;
    console.log('Group management for template ID:', templateId.value);
    // TODO: Fetch existing groups for this template/portfolio if applicable
  } else {
    console.error('Template ID not provided to group management page.');
    uni.showToast({ title: '无法加载分组信息', icon: 'none' });
    uni.navigateBack();
  }
});

const addGroup = () => {
  groups.value.push({ id: nextGroupId++, name: `分组${groups.value.length + 1}` });
};

const removeGroup = (index: number) => {
  groups.value.splice(index, 1);
};

const confirmGroups = () => {
  // Validate group names (e.g., ensure not empty)
  if (groups.value.some(g => !g.name.trim())) {
    uni.showToast({ title: '分组名称不能为空', icon: 'none' });
    return;
  }

  console.log('Confirming groups:', groups.value);
  console.log('Navigating to edit page...');

  // TODO: Save groups (e.g., via API call or pass to next page)
  // Structure might need adjustment based on backend requirements (e.g., include order)
  const groupData = groups.value.map((group, index) => ({
      name: group.name,
      order: index + 1 // Assuming order is based on current array position
  }));
  console.log('Group data to save/pass:', groupData);


  // Navigate to the next step: Edit Page
  uni.navigateTo({
    // Pass templateId and potentially the group structure
    url: `/pages/edit/edit?templateId=${templateId.value}&groups=${encodeURIComponent(JSON.stringify(groupData))}`
  });
};
</script>

<style scoped>
.group-management-page {
  padding: 40rpx 30rpx;
  background-color: #f8f8f8;
  min-height: 100vh;
  box-sizing: border-box;
}

.title {
  font-size: 40rpx;
  font-weight: bold;
  margin-bottom: 30rpx;
  color: #333;
}

.group-card {
  background-color: #ffffff;
  padding: 20rpx 30rpx 40rpx; /* More padding bottom */
  border-radius: 16rpx;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
  margin-bottom: 40rpx;
}

.group-list {
  display: flex;
  flex-direction: column;
  gap: 25rpx; /* Spacing between items */
}

.group-item {
  background-color: #f0f0f0; /* Light gray background for items */
  padding: 25rpx 30rpx;
  border-radius: 10rpx;
  font-size: 30rpx;
  color: #333;
  display: flex;
  align-items: center;
  justify-content: space-between; /* Default for editable items */
}

.group-item.fixed-item {
  justify-content: center; /* Center text for fixed items */
  font-weight: 500;
  background-color: #e9e9e9; /* Slightly different bg for fixed */
}

.group-item.add-item {
  justify-content: center; /* Center content */
  color: #555;
  cursor: pointer;
  transition: background-color 0.2s;
}

.group-item.add-item:active {
  background-color: #e0e0e0;
}

.add-icon {
  font-size: 36rpx;
  font-weight: bold;
  margin-right: 15rpx;
  width: 40rpx;
  height: 40rpx;
  line-height: 36rpx;
  text-align: center;
  border: 1px solid #555;
  border-radius: 50%;
  display: inline-block;
}

.group-input {
  flex-grow: 1;
  margin-right: 20rpx;
  /* Basic input styling */
  background: none;
  border: none;
  outline: none;
  font-size: 30rpx;
  color: #333;
}

.delete-icon {
  color: #ff5555;
  font-size: 40rpx;
  font-weight: bold;
  cursor: pointer;
  padding: 0 10rpx; /* Easier to tap */
}

.confirm-button {
  background-color: #57DD00; /* Green color */
  color: #ffffff;
  border: none;
  border-radius: 10rpx;
  padding: 20rpx 0;
  font-size: 30rpx;
  text-align: center;
  width: auto; /* Fit content */
  min-width: 180rpx; /* Minimum width */
  margin-left: auto; /* Align to the right */
  display: block; /* Needed for margin-left: auto */
  /* Simulate the > arrow */
  position: relative;
  padding-right: 50rpx; /* Make space for the pseudo-element */
}

.confirm-button::after {
    content: '>';
    position: absolute;
    right: 30rpx; /* Adjust position */
    top: 50%;
    transform: translateY(-50%);
    font-weight: bold;
    font-size: 36rpx;
}

.confirm-button:active {
  background-color: #4cc200; /* Darker shade on press */
}

.arrow {
  /* Style for the arrow text if needed, currently handled by ::after */
  display: none; /* Hide the text arrow if using ::after */
}
</style>
