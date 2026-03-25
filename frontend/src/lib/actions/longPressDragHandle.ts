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

	function getListItem(): HTMLElement | null {
		return node.closest('li') as HTMLElement | null;
	}

	function applyFeedback() {
		const item = getListItem() ?? node;
		item.style.transition = 'transform 0.3s ease, box-shadow 0.3s ease';
		item.style.transform = 'scale(1.02)';
		item.style.boxShadow = '0 4px 16px rgba(0,0,0,0.15)';
		item.style.zIndex = '1';
		item.style.position = 'relative';
		item.style.borderRadius = '8px';
	}

	function removeFeedback() {
		const item = getListItem() ?? node;
		item.style.transition = 'transform 0.2s ease, box-shadow 0.2s ease';
		item.style.transform = '';
		item.style.boxShadow = '';
		item.style.zIndex = '';
		item.style.position = '';
		item.style.borderRadius = '';
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
			removeFeedback();

			if (activeTouchId === null) return;

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