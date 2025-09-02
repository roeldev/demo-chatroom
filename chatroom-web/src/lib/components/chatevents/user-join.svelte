<script lang="ts">
	import type { UserJoinEvent } from "$lib/chatevents/event-types.ts";
	import { cn } from "$lib/utils.js";

	let {
		class: className,
		date,
		user,
	}: UserJoinEvent & {
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
		User {user.details?.name} joined
	</div>
</div>
