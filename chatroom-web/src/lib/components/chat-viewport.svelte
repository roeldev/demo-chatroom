<script lang="ts">
	import { onMount } from "svelte";
	// noinspection ES6UnusedImports
	import LoaderIcon from "@lucide/svelte/icons/loader";

	import { EventStreamer } from "$lib/chatevents";
	import { getUsersList } from "$lib/chatusers";
	import { cn } from "$lib/utils";

	import ChatView from "./chat-view.svelte";
	import { initChatViewport } from "./chat-viewport.svelte.ts";
	import UserTyping from "./user-typing.svelte";

	let loading = $state(true);
	const viewport = initChatViewport();

	onMount(() => {
		const handler = viewport.getEventsHandler()
		const events = new EventStreamer();
		events.stream([handler, getUsersList()], async () => {
			await viewport.updateScrollPosition();
		});

		events.loadPrevious([handler]);
		viewport.updateScrollPosition();
		loading = false;
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
		{#if loading}
			<div class="flex mb-2 text-sm">
				<div class="flex mx-auto px-3 py-2">
					<LoaderIcon class="animate-spin"/>
					<span class="grid flex-1 pl-2">Loading previous events...</span>
				</div>
			</div>
		{/if}

		{#each viewport.getViewsEntries() as [_, view]}
			<ChatView view={view}/>
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
