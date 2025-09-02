<script lang="ts">
	import { type Event } from '$lib/chatevents';
	import { type UUID } from '$lib/chatusers';
	// noinspection ES6UnusedImports
	import * as ChatEvent from '$lib/components/chatevents';

	let {
		list,
		user
	}: {
		list: Event[],
		user?: UUID,
	} = $props();
</script>

<div data-view={user ?? "global"}>
	{#each list as event}
		{#if event.case == "userJoin"}
			<ChatEvent.UserJoin {...event.value} />
		{:else if event.case == "userLeave"}
			<ChatEvent.UserLeave {...event.value} />
		{:else if event.case == "userUpdate"}
			<ChatEvent.UserUpdate {...event.value} />
		{:else if event.case == "userStatus"}
			<ChatEvent.UserStatus {...event.value} />
		{:else if event.case == "receivedChat"}
			<ChatEvent.ReceivedChat {...event.value} />
		{:else if event.case == "sentChat"}
			<ChatEvent.SentChat {...event.value} />
		{/if}
	{/each}
</div>
