<script lang="ts">
	import { onMount } from "svelte";

	import { EventStreamer } from "$lib/chatevents";
	import { getUsersList } from "$lib/chatusers";
	import { cn } from "$lib/utils";

	import ChatView from "./chat-view.svelte";
	import { initChatViewport } from "./chat-viewport.svelte.ts";
	import UserTyping from "./user-typing.svelte";

	const viewport = initChatViewport();

	onMount(() => {
		const streamer = new EventStreamer([viewport, getUsersList()]);
		streamer.stream(async () => {
			await viewport.updateAfterStreamEvent();
		})
	});

	let {
		class: className
	}: {
		class?: string,
	} = $props();
</script>

<div
	bind:this={viewport.scrollElem}
	onscrollend={() => viewport.onScrollEnd()}
	class="w-full h-full overflow-y-auto pb-4"
>
	<div
		class={cn(
			"w-full mx-auto",
			className,
		)}
	>
		<div id="previous-events">
			Loading previous events...
		</div>

		{#each viewport.getViewsEntries() as [user, view]}
			<ChatView list={view.getEvents()} user={user}/>
		{/each}

		{#if viewport.getUsersTyping().size > 3}
			<div>{viewport.getUsersTyping().size} users are typing...</div>
		{:else}
			{#each viewport.getUsersTyping().entries() as [uuid, user]}
				<UserTyping uuid={uuid} user={user}/>
			{/each}
		{/if}
	</div>
</div>
