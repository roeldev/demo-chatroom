<script lang="ts">
	import { API } from "$lib/chatapi"
	import type { UserLeaveEvent } from "$lib/chatevents/event-types.ts";
	import { cn } from "$lib/utils.js";

	let {
		class: className,
		date,
		user,
		reason,
	}: UserLeaveEvent & {
		class?: string,
	} = $props()
</script>

<div
	class={cn(
		"flex mb-2",
		"text-xs text-primary-foreground",
		className,
	)}
>
	<div class="grid gap-2 mx-auto px-3 py-2 bg-primary rounded-lg">
		{#if reason == API.LeaveReason.DISCONNECTED}
			User {user.details?.name} disconnected
		{:else}
			User {user.details?.name} left
		{/if}
	</div>
</div>
