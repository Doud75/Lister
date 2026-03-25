const LONG_PRESS_DELAY_MS = 1000;
const MOVEMENT_CANCEL_THRESHOLD_PX = 10;

export function longPressDragHandle(node: HTMLElement) {
	let timer: ReturnType<typeof setTimeout> | null = null;
	let startX = 0;
	let startY = 0;
	let startScreenX = 0;
	let startScreenY = 0;
	let activeTouchId: number | null = null;
	let isFiring = false;

	function applyFeedback() {
		node.style.transition = 'box-shadow 0.3s ease, background-color 0.3s ease';
		node.style.boxShadow = '0 0 0 3px rgba(99, 102, 241, 0.4)';
		node.style.backgroundColor = 'rgba(99, 102, 241, 0.1)';
		node.style.borderRadius = '6px';
	}

	function removeFeedback() {
		node.style.transition = '';
		node.style.boxShadow = '';
		node.style.backgroundColor = '';
		node.style.borderRadius = '';
	}

	function cancel() {
		if (timer !== null) {
			clearTimeout(timer);
			timer = null;
		}
		removeFeedback();
		activeTouchId = null;
	}

	function onTouchStart(e: TouchEvent) {
		if (isFiring) return;

		e.stopImmediatePropagation();

		const touch = e.changedTouches[0];
		activeTouchId = touch.identifier;
		startX = touch.clientX;
		startY = touch.clientY;
		startScreenX = touch.screenX;
		startScreenY = touch.screenY;

		applyFeedback();

		timer = setTimeout(() => {
			timer = null;
			if (activeTouchId === null) return;

			removeFeedback();

			isFiring = true;
			node.dispatchEvent(
				new MouseEvent('mousedown', {
					bubbles: true,
					cancelable: true,
					clientX: startX,
					clientY: startY,
					screenX: startScreenX,
					screenY: startScreenY,
					button: 0,
					buttons: 1
				})
			);
			isFiring = false;
			activeTouchId = null;
		}, LONG_PRESS_DELAY_MS);
	}

	function onTouchMove(e: TouchEvent) {
		if (activeTouchId === null) return;
		const touch = Array.from(e.changedTouches).find((t) => t.identifier === activeTouchId);
		if (!touch) return;

		const dx = Math.abs(touch.clientX - startX);
		const dy = Math.abs(touch.clientY - startY);
		if (dx > MOVEMENT_CANCEL_THRESHOLD_PX || dy > MOVEMENT_CANCEL_THRESHOLD_PX) {
			cancel();
		}
	}

	function onTouchEnd(e: TouchEvent) {
		if (activeTouchId === null) return;
		const touch = Array.from(e.changedTouches).find((t) => t.identifier === activeTouchId);
		if (touch) cancel();
	}

	node.addEventListener('touchstart', onTouchStart as EventListener, {
		capture: true,
		passive: false
	});
	window.addEventListener('touchmove', onTouchMove as EventListener, { passive: true });
	window.addEventListener('touchend', onTouchEnd as EventListener, { passive: true });

	return {
		destroy() {
			node.removeEventListener('touchstart', onTouchStart as EventListener, { capture: true });
			window.removeEventListener('touchmove', onTouchMove as EventListener);
			window.removeEventListener('touchend', onTouchEnd as EventListener);
			cancel();
		}
	};
}