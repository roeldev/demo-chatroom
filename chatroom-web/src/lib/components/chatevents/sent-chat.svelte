<script lang="ts">
	import type { SentChatEvent } from "$lib/chatevents/event-types.ts";
	import { cn } from "$lib/utils.js";

	let {
		class: className,
		date,
		chatId,
		text,
	}: SentChatEvent & {
		class?: string,
	} = $props()
</script>

<div
	data-chat-id={chatId}
	class={cn(
		"flex w-max max-w-[75%] mb-2 ml-auto",
		"text-sm text-primary-foreground",
		className,
	)}
>
	<div class="grid gap-2 px-3 py-2 bg-primary rounded-b-lg rounded-tl-lg">
		{text}
		<span class="text-xs text-right">{date.toLocaleTimeString()}</span>
	</div>
</div>
