<script lang="ts">
	// noinspection ES6UnusedImports
	import BotIcon from "@lucide/svelte/icons/bot";
	// noinspection ES6UnusedImports
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";

	import { getUsersList, UserStatus, UserFlag } from "$lib/chatusers";
	// noinspection ES6UnusedImports
	import * as Avatar from "$lib/components/ui/avatar";
	// noinspection ES6UnusedImports
	import * as Sidebar from "$lib/components/ui/sidebar";

	const users = getUsersList();

	function statusString(typing: boolean, status: UserStatus): string {
		if (typing) {
			return "Typing..."
		}

		switch (status) {
			case UserStatus.DEFAULT:
				return ""
			case UserStatus.BUSY:
				return "Busy"
			case UserStatus.AWAY:
				return "Away"
			default:
				return "Invalid"
		}
	}

	function isBot(flags: UserFlag): boolean {
		return (flags & UserFlag.IS_BOT) === UserFlag.IS_BOT
	}
</script>

<Sidebar.Group>
	<Sidebar.GroupLabel>Active users</Sidebar.GroupLabel>
	<Sidebar.Menu>
		{#each users.entries() as user}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton tooltipContent={user.details.name} data-uuid={user.id}>
					<Avatar.Root class="size-7 group-data-[state=icon]:size-6!">
						<!--						<Avatar.Image src={user.details?.image} alt={user.details?.name}/>-->
						<Avatar.Fallback
							style="background-color:{user.details.color1}; color:{user.details.color2}"
							class="rounded-full group-data-[state=collapsed]:size-6"
						>
							{user.details.initials}
						</Avatar.Fallback>
					</Avatar.Root>
					<div class="grid flex-1 text-left text-sm leading-tight">
						<span class="truncate font-medium">{user.details.name}</span>
						<span class="truncate text-xs italic">{statusString(!!user.typing, user.status)}</span>
					</div>
					{#if isBot(user.flags)}
						<BotIcon/>
					{/if}
					<ChevronRightIcon/>
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		{/each}
	</Sidebar.Menu>
</Sidebar.Group>
