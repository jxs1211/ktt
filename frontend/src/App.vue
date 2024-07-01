<script setup>
import { get } from 'lodash'
import { computed, nextTick, onMounted, provide, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { GetPreferences } from '../wailsjs/go/storage/PreferencesStorage.js'
import ContentPane from './components/ContentPane.vue'
import hljs from 'highlight.js/lib/core'
import json from 'highlight.js/lib/languages/json'
import plaintext from 'highlight.js/lib/languages/plaintext'
import { useThemeVars } from 'naive-ui'


const themeVars = useThemeVars()

hljs.registerLanguage('json', json)
hljs.registerLanguage('plaintext', plaintext)

const data = reactive({
	asideWith: 300,
	hoverResize: false,
	resizing: false,
})

const preferences = ref({})
provide('preferences', preferences)
const i18n = useI18n()

onMounted(async () => {
	preferences.value = await GetPreferences()
	await nextTick(() => {
		i18n.locale.value = get(preferences.value, 'general.language', 'en')
	})
})

// TODO: apply font size to all elements
const getFontSize = computed(() => {
	return get(preferences.value, 'general.font_size', 'en')
})

const themeOverrides = {
	common: {
		// primaryColor: '#409EFF',
		borderRadius: '4px',
		borderRadiusSmall: '3px',
		fontFamily: `"Nunito", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto",
  "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue"`,
		lineHeight: 1.5,
	},
	Tag: {
		// borderRadius: '3px'
	},
}


</script>

<template>
	<n-config-provider :hljs="hljs" :inline-theme-disabled="true" :theme-overrides="themeOverrides" class="fill-height">
  </n-config-provider>
</template>

<style>
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
</style>
