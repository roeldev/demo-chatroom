<script lang="ts">
	import { Input } from "$lib/components/ui/input";
	import { cn } from "$lib/utils.ts";

	import { ChatInput } from "./chat-input.svelte.ts";

	const chat = new ChatInput();
	let {
		class: className
	}: {
		class?: string,
	} = $props();
</script>

<div
	class={cn(
		"w-full mx-auto shrink-0 items-center pb-4",
		className,
	)}
>
	<form method="POST" onsubmit={e => chat.onSubmitForm(e)} class="w-full">
		<Input bind:value={chat.value}
					 oninput={e => chat.onInput(e)}
					 onblur={() => chat.indicateTyping(false)}
					 name="input"
					 type="text"
					 placeholder="Type your chat message here..."
					 class="h-12"
					 required
		/>
	</form>
</div>
