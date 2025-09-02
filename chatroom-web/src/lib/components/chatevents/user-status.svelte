<script lang="ts">
	import type { UserStatusEvent } from "$lib/chatevents/event-types.ts";
	import { cn } from "$lib/utils.js";

	let {
		class: className,
		date,
		user,
		status,
		before,
	}: UserStatusEvent & {
		class?: string,
	} = $props()
</script>

<div
	class={cn(
		"flex w-max max-w-[75%]",
		"gap-2 rounded-b-lg rounded-tl-lg mx-4 mb-2 px-3 py-2 ml-auto",
		"text-sm text-primary-foreground bg-primary",
		className,
	)}
>
	User {user.details!.name} status changed from {before} to {status}
	<span class="text-xs">{date}</span>
</div>
