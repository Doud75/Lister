const LONG_PRESS_DELAY_MS = 500;
const MOVEMENT_CANCEL_THRESHOLD_PX = 10;

export function longPressDragHandle(node: HTMLElement) {
	let timer: ReturnType<typeof setTimeout> | null = null;
	let startX = 0;
	let startY = 0;
	let startScreenX = 0;
	let startScreenY = 0;
	let activeTouchId: number | null = null;
	let dragTouchId: number | null = null;
	let isFiring = false;

	function applyFeedback() {
		node.style.transition = 'transform 0.1s ease-out, box-shadow 0.1s ease-out';
		node.style.transform = 'scale(1.1)';
		node.style.boxShadow = '0 0 0 4px rgba(99, 102, 241, 0.5)';
		node.style.backgroundColor = 'rgba(99, 102, 241, 0.2)';
		node.style.borderRadius = '6px';
		node.style.zIndex = '50';
	}

	function removeFeedback() {
		node.style.transition = 'transform 0.2s ease-in, box-shadow 0.2s ease-in';
		node.style.transform = '';
		node.style.boxShadow = '';
		node.style.backgroundColor = '';
		node.style.borderRadius = '';
		node.style.zIndex = '';
		node.style.opacity = '1';
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

		const touch = e.changedTouches[0];
		activeTouchId = touch.identifier;
		startX = touch.clientX;
		startY = touch.clientY;
		startScreenX = touch.screenX;
		startScreenY = touch.screenY;

		node.style.opacity = '0.7';

		timer = setTimeout(() => {
			timer = null;
			if (activeTouchId === null) return;

			dragTouchId = activeTouchId;
			activeTouchId = null;

			if (window.navigator && window.navigator.vibrate) {
				window.navigator.vibrate(30);
			}

			applyFeedback();

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
		}, LONG_PRESS_DELAY_MS);
	}

	function onTouchMove(e: TouchEvent) {
		if (activeTouchId !== null) {
			const touch = Array.from(e.changedTouches).find((t) => t.identifier === activeTouchId);
			if (!touch) return;
			const dx = Math.abs(touch.clientX - startX);
			const dy = Math.abs(touch.clientY - startY);
			if (dx > MOVEMENT_CANCEL_THRESHOLD_PX || dy > MOVEMENT_CANCEL_THRESHOLD_PX) {
				cancel();
			}
		} else if (dragTouchId !== null) {
			const touch = Array.from(e.changedTouches).find((t) => t.identifier === dragTouchId);
			if (touch) e.preventDefault();
		}
	}

	function onTouchEnd(e: TouchEvent) {
		if (activeTouchId !== null) {
			const touch = Array.from(e.changedTouches).find((t) => t.identifier === activeTouchId);
			if (touch) cancel();
		} else if (dragTouchId !== null) {
			const touch = Array.from(e.changedTouches).find((t) => t.identifier === dragTouchId);
			if (touch) {
				dragTouchId = null;
				removeFeedback();
			}
		}
	}

	node.addEventListener('touchstart', onTouchStart as EventListener, {
		capture: true,
		passive: false
	});
	window.addEventListener('touchmove', onTouchMove as EventListener, {
		capture: true,
		passive: false
	});
	window.addEventListener('touchend', onTouchEnd as EventListener, {
		capture: true,
		passive: false
	});

	return {
		destroy() {
			node.removeEventListener('touchstart', onTouchStart as EventListener, { capture: true });
			window.removeEventListener('touchmove', onTouchMove as EventListener, { capture: true });
			window.removeEventListener('touchend', onTouchEnd as EventListener, { capture: true });
			cancel();
		}
	};
}
