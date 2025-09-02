<script lang="ts">
	// noinspection ES6UnusedImports
	import * as Card from "$lib/components/ui/card";
	import Logo from "$lib/components/logo.svelte";
	import JoinForm from "$lib/components/join-form.svelte";
</script>

<div class="flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
	<div class="flex w-full max-w-sm flex-col gap-6">
		<a href="##" class="flex items-center gap-2 self-center font-medium">
			<Logo />
		</a>
		<div class="flex flex-col gap-6">
			<Card.Root>
				<Card.Header class="text-center">
					<Card.Title class="text-xl">Welcome!</Card.Title>
					<Card.Description>Choose a unique name and join the chatroom!</Card.Description>
				</Card.Header>
				<Card.Content>
					<JoinForm />
				</Card.Content>
			</Card.Root>
		</div>
	</div>
</div>
