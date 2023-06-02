<script>
	import { quintOut } from 'svelte/easing';

	import { createEventDispatcher } from 'svelte';

	const modal = (/** @type {Element} */ node, { duration = 300 } = {}) => {
		const transform = getComputedStyle(node).transform;

		return {
			duration,
			easing: quintOut,
			css: (/** @type {any} */ t, /** @type {number} */ u) => {
				return `transform:
            ${transform}
            scale(${t})
            translateY(${u * -100}%)
          `;
			}
		};
	};

	const dispatch = createEventDispatcher();
	function closeModal() {
		dispatch('close', {});
	}
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class="modal-background">
	<div transition:modal={{ duration: 1000 }} class="modal" role="dialog" aria-modal="true">
		<!-- svelte-ignore a11y-missing-attribute -->
		<a title="Close" class="modal-close" on:click={closeModal}>
			<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 384 512">
				<path
					d="M342.6 150.6c12.5-12.5 12.5-32.8 0-45.3s-32.8-12.5-45.3 0L192 210.7 86.6 105.4c-12.5-12.5-32.8-12.5-45.3 0s-12.5 32.8 0 45.3L146.7 256 41.4 361.4c-12.5 12.5-12.5 32.8 0 45.3s32.8 12.5 45.3 0L192 301.3 297.4 406.6c12.5 12.5 32.8 12.5 45.3 0s12.5-32.8 0-45.3L237.3 256 342.6 150.6z"
				/>
			</svg>
		</a>
		<div class="container">
			<slot />
		</div>
	</div>
</div>

<style>
	.modal-background {
		width: 100%;
		height: 100%;
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.9);
		z-index: 9999;
	}

	.modal {
		position: absolute;
		left: 50%;
		top: 50%;
		width: 70%;
		box-shadow: 0 0 10px hsl(0 0% 0% / 10%);
		transform: translate(-50%, -50%);
	}
	@media (max-width: 990px) {
		.modal {
			width: 90%;
		}
	}
	.modal-close {
		border: none;
	}

	.modal-close svg {
		display: block;
		margin-left: auto;
		margin-right: auto;
		fill: rgb(14 165 233 /1);
		transition: all 0.5s;
	}
	.modal-close:hover svg {
		fill: rgb(225 29 72);
		transform: scale(1.5);
	}
	.modal .container {
		max-height: 90vh;
		overflow-y: auto;
	}
	@media (min-width: 680px) {
		.modal .container {
			flex-direction: column;
			left: 0;
			width: 100%;
		}
	}
</style>
