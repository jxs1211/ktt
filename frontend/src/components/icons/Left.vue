<script setup>
const props = defineProps({
	strokeWidth: {
		type: [Number, String],
		default: 3,
	},
})
</script>

<template>
	<svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
		<path
		:stroke-width="props.strokeWidth"
		d="M30 12L18 24L30 36"
		stroke="currentColor"
		stroke-linecap="round"
		stroke-linejoin="round" />
	</svg>
</template>

<style lang="scss" scoped></style>
