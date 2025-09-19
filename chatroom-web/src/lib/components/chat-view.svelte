<script lang="ts">
	// noinspection ES6UnusedImports
	import * as ChatEvent from '$lib/components/chatevents';
	import type { ChatView } from "./chat-view.svelte.ts";

	let {
		view,
	}: {
		view: ChatView,
	} = $props();
</script>

<div data-view={view.owner ?? "global"}>
	{#each view.getEvents() as event}
		{#if event.case == "userJoin"}
			<ChatEvent.UserJoin {...event.value} />
		{:else if event.case == "userLeave"}
			<ChatEvent.UserLeave {...event.value} />
		{:else if event.case == "userUpdate"}
			<ChatEvent.UserUpdate {...event.value} />
		{:else if event.case == "userStatus"}
			<ChatEvent.UserStatus {...event.value} />
		{:else if event.case == "receivedChat"}
			<ChatEvent.ReceivedChat showAvatar {...event.value} />
		{:else if event.case == "sentChat"}
			<ChatEvent.SentChat {...event.value} />
		{/if}
	{/each}
</div>
