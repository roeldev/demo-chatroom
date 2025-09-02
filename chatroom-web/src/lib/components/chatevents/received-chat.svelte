<script lang="ts">
	// noinspection ES6UnusedImports
	import * as Avatar from '$lib/components/ui/avatar';
	import type { ReceivedChatEvent } from "$lib/chatevents/event-types.ts";
	import { cn } from "$lib/utils";

	let {
		class: className,
		date,
		showAvatar,
		chatId,
		user,
		text,
	}: ReceivedChatEvent & {
		class?: string,
	} = $props();
</script>

<div
	data-chat-id={chatId.value}
	class={cn(
		"flex w-max max-w-[75%] gap-2 mb-2",
		className,
	)}
>
	{#if showAvatar}
		<Avatar.Root class="flex-none size-9" data-user-id={user.id?.value}>
			<!--						<Avatar.Image src={user.details?.image} alt={user.details?.name}/>-->
			<Avatar.Fallback
				style="background-color:{user.details?.color1?.value}; color:{user.details?.color2?.value}"
				class="rounded-full"
			>
				{user.details?.initials}
			</Avatar.Fallback>
		</Avatar.Root>
	{:else}
		<div class="flex-none size-9"></div>
	{/if}
	<div class="grid gap-2 px-3 py-2 text-sm bg-muted rounded-b-lg rounded-tr-lg">
		{text}
		<span class="text-xs">{date.toLocaleTimeString()}</span>
	</div>
</div>
