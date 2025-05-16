<template>
  <view class="editor-component" :style="{ borderColor: bleedEnabled ? 'red' : 'transparent' }">
    <!-- Basic display of work data -->
    <view class="content-area">
        <!--
            This is where the actual page rendering and editing UI will go.
            It needs to handle different content types (images, text) based on workData.
            For now, just showing a placeholder.
        -->
        <image
            v-if="workData.oss_key.startsWith('dummy_')"
            class="placeholder-image"
            src="/static/images/template1.png"
            mode="aspectFit"
        />
         <image
            v-else
            class="actual-image"
            :src="workData.oss_key"
            mode="aspectFit"
        />
        <!-- Add text rendering, etc. -->
        <text class="page-info">Editing Page {{ workData.page }} ({{ workData.oss_key }})</text>
        <text class="size-info">Size: {{ displaySize }}</text>
        <text v-if="bleedEnabled" class="bleed-info">Bleed: {{ workData.bleed.join(', ') }}</text>

        <!-- Placeholder for editing controls -->
        <!-- <textarea v-if="isTextElement" @input="handleContentChange($event.detail.value)"></textarea> -->

    </view>
  </view>
</template>

<script setup lang="ts">
import { type PropType, computed } from 'vue'; // Use 'type' for type-only import
import type { Work } from '../../types/api';

const props = defineProps({
  workData: {
    type: Object as PropType<Work>,
    required: true,
  },
  bleedEnabled: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:workData']);

// Example: If the editor modifies the workData, emit an update
const handleContentChange = (newContent: string) => {
    // This is just a placeholder for actual editing logic
    // You would update the relevant part of workData based on user actions
    const updatedData: Partial<Work> = {
        // Example: Assuming 'content' is a field in Work for text
        // content: newContent
    };
    // emit('update:workData', updatedData);
};

// Example computed property based on props
const displaySize = computed(() => {
    // Add logic based on bleedEnabled if needed
    return props.workData.size;
});

</script>

<style scoped>
.editor-component {
  width: 100%;
  height: 100%; /* Adjust as needed, might need specific dimensions */
  background-color: #ffffff;
  border: 2px dashed; /* Use border to show bleed area */
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
  padding: 20rpx;
  position: relative; /* For absolute positioning of elements if needed */
}

.content-area {
    text-align: center;
    color: #888;
}

.placeholder-image, .actual-image {
    max-width: 80%;
    max-height: 300rpx; /* Limit initial display size */
    margin-bottom: 20rpx;
    border: 1px solid #eee;
}

.page-info, .size-info, .bleed-info {
    display: block;
    font-size: 24rpx;
    margin-top: 10rpx;
}

.bleed-info {
    color: red;
}
</style>
