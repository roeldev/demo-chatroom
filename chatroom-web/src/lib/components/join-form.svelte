<script lang="ts">
	import AlertCircleIcon from "@lucide/svelte/icons/alert-circle";
	// noinspection ES6UnusedImports
	import * as Alert from "$lib/components/ui/alert";
	import { Input } from "$lib/components/ui/input";
	import { Button } from "$lib/components/ui/button";
	import { Label } from "$lib/components/ui/label";
	import { JoinForm, type JoinSubmitEvent } from "$lib/components/join-form.svelte.ts";

	let lock = $state(false)
	let error = $state("");

	const id = $props.id();
	const name = "name";
	const join = new JoinForm(name);

	function submit(e: JoinSubmitEvent) {
		lock = true;
		join.submitForm(e).catch((err) => {
			lock = false;
			error = err.message;
		});
	}

</script>

<form method="POST" onsubmit={submit}>
	<div class="grid gap-3">
		<div class="grid gap-3">
			<Label for="name-{id}">Name</Label>
			<Input id="name-{id}" name={name} type="text" required/>
		</div>
		{#if error !== ""}
			<Alert.Root variant="destructive">
				<AlertCircleIcon/>
				<Alert.Title>Unable to join</Alert.Title>
				<Alert.Description>
					<p>{error}</p>
				</Alert.Description>
			</Alert.Root>
		{/if}
		<Button type="submit" class="w-full" disabled={lock}>Join</Button>
	</div>
</form>
