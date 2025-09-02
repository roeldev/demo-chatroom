<script lang="ts">
	import { onMount } from "svelte";

	import { isAuthenticated } from "$lib/chatauth";
	import { getUsersList } from "$lib/chatusers";
	// noinspection ES6UnusedImports
	import * as Sidebar from "$lib/components/ui/sidebar";
	import ChatInput from "$lib/components/chat-input.svelte";
	import ChatViewport from "$lib/components/chat-viewport.svelte";
	import Header from "$lib/components/header.svelte";
	import Logo from "$lib/components/logo.svelte";
	import NavUsers from "$lib/components/users-list.svelte";

	const users = getUsersList();

	onMount(() => {
		if (!isAuthenticated()) {
			return;
		}
		users.fetch().catch((err) => {
			console.log("failed to get active users", err.message);
		});
	});
</script>

<Sidebar.Provider class="h-svh">
	<Sidebar.Root collapsible="icon">
		<Sidebar.Header>
			<Sidebar.MenuButton
				size="lg"
				class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
			>
				<Logo/>
			</Sidebar.MenuButton>
		</Sidebar.Header>
		<Sidebar.Content>
			<NavUsers/>
		</Sidebar.Content>
	</Sidebar.Root>
	<Sidebar.Inset>
		<Header/>
		<ChatViewport class="px-4 lg:max-w-4xl"/>
		<ChatInput class="px-4 lg:max-w-4xl"/>
	</Sidebar.Inset>
</Sidebar.Provider>
