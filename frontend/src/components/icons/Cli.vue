<script setup>
const props = defineProps({
  inverse: {
    type: Boolean,
    default: false,
  },
  strokeWidth: {
    type: [Number, String],
    default: 3,
  },
});
</script>

<template>
  <svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
    <rect
      :fill="props.inverse ? 'currentColor' : 'none'"
      :stroke-width="props.strokeWidth"
      x="4"
      y="8"
      width="40"
      height="32"
      stroke="currentColor"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
    <path
      :stroke="props.inverse ? '#FFF' : 'currentColor'"
      :stroke-width="props.strokeWidth"
      d="M12 18L20 24L12 30"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
    <path
      :stroke="props.inverse ? '#FFF' : 'currentColor'"
      :stroke-width="props.strokeWidth"
      d="M24 32H36"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </svg>
</template>

<style lang="scss" scoped></style>
