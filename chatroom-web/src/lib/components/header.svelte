<script lang="ts">
	import { toggleMode } from "mode-watcher";
	// noinspection ES6UnusedImports
	import MoonIcon from "@lucide/svelte/icons/moon";
	// noinspection ES6UnusedImports
	import SunIcon from "@lucide/svelte/icons/sun";

	import { Button } from "$lib/components/ui/button";
	import { Separator } from '$lib/components/ui/separator';
	import { SidebarTrigger } from "$lib/components/ui/sidebar";
	import { cn } from "$lib/utils.ts";

	let {
		class: className
	}: {
		class?: string,
	} = $props();
</script>

<header
	class={cn(
		"bg-background flex h-16 shrink-0 items-center gap-2 border-b",
		"group-has-data-[collapsible=icon]/sidebar-wrapper:h-12",
		"transition-[width,height] ease-linear",
		className,
	)}
>
	<div class="flex grow items-center gap-2 px-4">
		<SidebarTrigger class="-ml-1"/>
	</div>
	<div class="flex items-center gap-2 px-4">
		<Button onclick={toggleMode} variant="ghost" size="icon" class="size-7">
			<MoonIcon class="rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0"/>
			<SunIcon class="absolute rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100"/>
			<span class="sr-only">Toggle light/dark mode</span>
		</Button>
		<Separator orientation="vertical" class="data-[orientation=vertical]:h-4"/>
	</div>
</header>
