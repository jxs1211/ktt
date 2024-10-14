<script setup>
import { computed, h, nextTick, reactive, ref, onUnmounted, watchEffect } from 'vue'
import { types, typesColor } from '@/consts/support_redis_type.js'
import useDialog from 'stores/dialog'
import { endsWith, get, isEmpty, keys, map, trim } from 'lodash'
import { useI18n } from 'vue-i18n'
import useTabStore from 'stores/tab.js'
import useConnectionStore from "stores/connections.js";
import useBrowserStore from 'stores/browser.js'
import { useSessionStore } from '@/stores/session.js';
import {
  CreateSession,
  GetSessionsByClusterName,
} from 'wailsjs/go/db/DBService.js';

const i18n = useI18n()
const connectionStore = useConnectionStore();
const sessionStore = useSessionStore();
const browserStore = useBrowserStore()
const tabStore = useTabStore()
// const newForm = reactive({
// 	cluster_name: "",
// 	address: "",
// 	port: "",
// 	cmds: "",
// })
const formRules = computed(() => {
	const requiredMsg = i18n.t('dialogue.field_required')
	return {
		cluster_name: { required: true, message: requiredMsg, trigger: 'input' },
		address: { required: true, message: requiredMsg, trigger: 'input' },
		port: { required: true, message: requiredMsg, trigger: 'input' },
		cmds: { required: true, message: requiredMsg, trigger: 'input' },
	}
})
// const dbOptions = computed(() =>
// map(keys(browserStore.getDBList(newForm.server)), (key) => ({
// 	label: key,
// 	value: parseInt(key),
// })),
// )
const newFormRef = ref(null)
// const subFormRef = ref(null)

const options = computed(() => {
	return Object.keys(types).map((t) => ({
		value: t,
		label: t,
	}))
})

const dialogStore = useDialog()
const scrollRef = ref(null)
watchEffect(() => {
	if (dialogStore.newCliDialogVisible) {
		console.log("changed cli dialog visible: ", dialogStore.newCliDialogVisible)
		// const { prefix, server, db } = dialogStore.newKeyParam
		// const separator = browserStore.getSeparator(server)
		// newForm.server = server
		// if (isEmpty(prefix)) {
		// 	newForm.key = ''
		// } else {
		// 	if (!endsWith(prefix, separator)) {
		// 		newForm.key = prefix + separator
		// 	} else {
		// 		newForm.key = prefix
		// 	}
		// }
		// newForm.db = db
		// newForm.type = options.value[0].value
		// newForm.ttl = -1
		// newForm.value = null
	}
});
onUnmounted(()=>{
	console.log("new cli dialog onUnmounted")
})
const startPort = 12110;
const endPort = 22110;
const randomPort = (start, end, excludes = []) => {
	if (start < 1 || end > 65535 || start > end) {
		throw new Error("Invalid port range");
	}
	let port;
	do {
		port = Math.floor(Math.random() * (end - start + 1)) + start;
	} while (excludes.includes(port)); // Keep generating until found a port not in excludes
	return port;
};

const onAdd = async () => {
	console.log("onAdd")
	try {
		const resp = await CreateSession(
			sessionStore.formValue.cluster_name,
			sessionStore.formValue.address,
			String(sessionStore.formValue.port),
			sessionStore.formValue.cmds);
		if (!resp.success) {
			console.error("add cli failed: ", resp.msg)
			$message.error("add cli failed: ", resp.msg)
			return
		}
		console.log("create session: ", resp.data)
		const resp2 = await GetSessionsByClusterName(sessionStore.formValue.cluster_name);
		if (!resp2.success) {
			console.error("refresh sessions by cluster name failed: ", resp2.msg)
			$message.error("get sessions failed: ", resp2.msg)
			return
		}
		console.log("get sessions by cluster name: ", resp2.data)
		sessionStore.setResults(resp2.data);
		dialogStore.closeNewCliDialog()
		$message.success("add cli ok")
	} catch (e) {
		$message.error(e)
		return false;
	}
	return true
}

const onClose = () => {
	dialogStore.closeNewCliDialog()
}
</script>

<template>
	<n-modal
	v-model:show="dialogStore.newCliDialogVisible"
	:closable="false"
	:mask-closable="false"
	:show-icon="false"
	:title="$t('dialogue.cli.new')"
	close-on-esc
	preset="dialog"
	style="width: 600px"
	transform-origin="center"
	@esc="onClose">
		<n-scrollbar ref="scrollRef" style="max-height: 500px">
			<n-form
				ref="newFormRef"
				:model="sessionStore.formValue"
				:rules="formRules"
				:show-require-mark="false"
				label-placement="top"
				style="padding-right: 15px">
				<n-form-item :label="$t('dialogue.cli.cluster_name')" path="key" required>
					<n-input disabled :value="sessionStore.formValue.cluster_name" placeholder="" />
				</n-form-item>
				<n-form-item :label="$t('dialogue.cli.address')" path="db" required>
					<n-input disabled :value="sessionStore.formValue.address" placeholder="" />
				</n-form-item>
				<n-form-item :label="$t('dialogue.cli.port')" path="db" required>
					<n-input disabled :value="sessionStore.formValue.port" placeholder="" />
				</n-form-item>
				<n-form-item :label="$t('dialogue.cli.cmds')" path="type" required>
					<n-input v-model:value="sessionStore.formValue.cmds" placeholder="" />
				</n-form-item>
				<div class="flex-item n-dialog__action">
					<n-button :focusable="false" size="medium" @click="onClose">
						{{ $t('common.cancel') }}
					</n-button>
					<n-button :focusable="false" size="medium" type="primary" @click="onAdd">
						{{ $t('common.confirm') }}
					</n-button>
				</div>
			</n-form>
		</n-scrollbar>
	</n-modal>
</template>

<style lang="scss" scoped>
</style>
