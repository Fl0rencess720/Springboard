<template>
  <view class="tabbar-container">
    <view
      class="tab-item"
      v-for="item in tabList"
      :key="item.pagePath"
      @click="switchTab(item)"
    >
      <image
        class="icon"
        :src="getIconPath(item)"
        mode="aspectFit"
      />
      <!-- <text :class="['text', { active: currentPage === item.pagePath }]">
        {{ item.text }}
      </text> -->
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

// Define props to receive the current page path
const props = defineProps({
  currentPage: {
    type: String,
    required: true,
  },
});

interface TabItem {
  pagePath: string;
  text: string;
  icon: string; // Base icon symbol ID
  iconCurrent: string; // Current icon symbol ID
  isSpecial?: boolean; // Flag for special buttons like 'add'
}

const tabList = ref<TabItem[]>([
  {
    pagePath: '/pages/index/index',
    text: '首页',
    icon: 'index',
    iconCurrent: 'index-current',
  },
  {
    pagePath: '/pages/select/size', // Navigate to size selection page first
    text: '添加',
    icon: 'add',
    iconCurrent: 'add', // 'add' might not have a 'current' state visually
    isSpecial: true,
  },
  {
    pagePath: '/pages/profile/profile', // Assuming a 'profile' page
    text: '我的',
    icon: 'profile',
    iconCurrent: 'profile-current',
  },
]);

// Function to determine the correct icon path with fragment identifier
const getIconPath = (item: TabItem): string => {
  const iconId = props.currentPage === item.pagePath && !item.isSpecial
    ? item.iconCurrent
    : item.icon;
  // Note: SVG fragment identifiers might have varying support across platforms in uni-app.
  // Consider using separate SVG files per icon if issues arise.
  return `/static/icons/${iconId}.svg`;
};

// Function to handle tab switching
const switchTab = (item: TabItem) => {
  if (item.isSpecial) {
    // Handle special button action, navigate to the start of the creation flow
    uni.navigateTo({
       url: item.pagePath // Navigate to '/pages/select/size'
    });
    console.log('Navigating to creation flow:', item.pagePath);
  } else {
    uni.switchTab({
      url: item.pagePath,
    });
  }
};
</script>

<style scoped>
.tabbar-container {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 50px; /* Adjust height as needed */
  background-color: #ffffff;
  display: flex;
  justify-content: space-around;
  align-items: center;
  border-top: 1px solid #e0e0e0;
  box-shadow: 0 -1px 5px rgba(0, 0, 0, 0.1);
  padding-bottom: constant(safe-area-inset-bottom); /* iOS safe area */
  padding-bottom: env(safe-area-inset-bottom); /* iOS safe area */
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  height: 100%;
}

.icon {
  width: 65px; /* Adjust icon size */
  height: 24px; /* Adjust icon size */
  margin-bottom: 2px;
}

/* Special styling for the 'add' button if needed */
.tab-item:nth-child(2) .icon {
   width: 36px; /* Make 'add' icon larger */
   height: 36px;
}
.tab-item:nth-child(2) .text {
    display: none; /* Optionally hide text for the special button */
}


.text {
  font-size: 10px; /* Adjust text size */
  color: #8a8a8a; /* Default text color */
}

.text.active {
  color: #57DD00; /* Active text color - match the green in icons */
}
</style>