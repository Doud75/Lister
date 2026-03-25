const LONG_PRESS_DELAY_MS = 500;
const MOVEMENT_CANCEL_THRESHOLD_PX = 10;

export function longPressDragHandle(node: HTMLElement) {
	let timer: ReturnType<typeof setTimeout> | null = null;
	let startX = 0;
	let startY = 0;
	let isLongPressActive = false;
	let isDispatching = false;

	function applyFeedback() {
		node.style.transition = 'transform 0.1s ease-out';
		node.style.transform = 'scale(1.2)';
		node.style.boxShadow = '0 0 0 4px rgba(99, 102, 241, 0.4)';
		node.style.backgroundColor = 'rgba(99, 102, 241, 0.1)';
		node.style.zIndex = '50';
	}

	function removeFeedback() {
		node.style.transition = 'transform 0.1s ease-in';
		node.style.transform = '';
		node.style.boxShadow = '';
		node.style.backgroundColor = '';
		node.style.borderRadius = '';
		node.style.zIndex = '';
		node.style.opacity = '1';
	}

	function blockNativePointerDown(e: PointerEvent) {
		if (e.pointerType === 'mouse') return;
		if (!isDispatching) {
			e.stopImmediatePropagation();
		}
	}

	function cancel() {
		if (timer) {
			clearTimeout(timer);
			timer = null;
		}
		isLongPressActive = false;
		removeFeedback();
	}

	function onTouchStart(e: TouchEvent) {
		if (e.touches.length > 1) return;
		e.stopImmediatePropagation();
		const touch = e.touches[0];
		startX = touch.clientX;
		startY = touch.clientY;
		isLongPressActive = true;
		node.style.opacity = '0.6';

		timer = setTimeout(() => {
			if (!isLongPressActive) return;

			if (window.navigator && window.navigator.vibrate) {
				window.navigator.vibrate(40);
			}
			applyFeedback();

			isDispatching = true;
			const opts = {
				bubbles: true,
				cancelable: true,
				clientX: startX,
				clientY: startY,
				view: window,
				button: 0,
				buttons: 1
			};
			node.dispatchEvent(new MouseEvent('mousedown', opts));
			isDispatching = false;
		}, LONG_PRESS_DELAY_MS);
	}

	function onTouchMove(e: TouchEvent) {
		if (!isLongPressActive) return;
		const touch = e.touches[0];
		const dx = Math.abs(touch.clientX - startX);
		const dy = Math.abs(touch.clientY - startY);
		if (dx > MOVEMENT_CANCEL_THRESHOLD_PX || dy > MOVEMENT_CANCEL_THRESHOLD_PX) {
			cancel();
		}
	}

	node.addEventListener('touchstart', onTouchStart, { capture: true, passive: true });
	node.addEventListener('touchmove', onTouchMove, { capture: true, passive: true });
	node.addEventListener('touchend', cancel, { capture: true, passive: true });
	node.addEventListener('touchcancel', cancel, { capture: true, passive: true });

	node.addEventListener('pointerdown', blockNativePointerDown, { capture: true });

	return {
		destroy() {
			node.removeEventListener('touchstart', onTouchStart, { capture: true });
			node.removeEventListener('touchmove', onTouchMove, { capture: true });
			node.removeEventListener('touchend', cancel, { capture: true });
			node.removeEventListener('touchcancel', cancel, { capture: true });
			node.removeEventListener('pointerdown', blockNativePointerDown, { capture: true });
			cancel();
		}
	};
}