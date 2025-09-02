<script lang="ts">
	import { ModeWatcher } from "mode-watcher";
	import { onMount } from "svelte";

	import { goto } from "$app/navigation";
	import { initConnectTransport } from "$lib/chatapi";
	import { isAuthenticated } from "$lib/chatauth";
	import "../app.css";

	initConnectTransport()
	onMount(() => {
		if (!isAuthenticated()) {
			goto("/join");
		}
	})

	let { children } = $props();
</script>

<ModeWatcher/>
{@render children()}
