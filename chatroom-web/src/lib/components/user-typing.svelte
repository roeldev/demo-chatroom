<script lang="ts">
	import type { UserDetails, UUID } from "$lib/chatusers";
	// noinspection ES6UnusedImports
	import * as Avatar from '$lib/components/ui/avatar';
	import { cn } from "$lib/utils";

	let {
		class: className,
		uuid,
		user,
	}: {
		class?: string,
		uuid?: UUID,
		user: UserDetails,
	} = $props()
</script>

<div
	data-user-id={uuid}
	class={cn(
		"flex w-max max-w-[75%] gap-2 mb-2",
		className,
	)}
>
	<Avatar.Root class="flex-none size-9">
		<Avatar.Fallback
			style="background-color:{user.color1}; color:{user.color2}"
			class="rounded-full"
		>
			{user.initials}
		</Avatar.Fallback>
	</Avatar.Root>
	<div class="grid flex-auto gap-2 rounded-b-lg rounded-tr-lg px-3 py-2 text-sm bg-muted">
		...
	</div>
</div>
