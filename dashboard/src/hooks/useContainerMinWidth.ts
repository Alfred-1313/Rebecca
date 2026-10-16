import { useCallback, useRef, useState } from "react";

export const useContainerMinWidth = (minWidth: number) => {
	const [matches, setMatches] = useState(false);
	const observerRef = useRef<ResizeObserver | null>(null);
	const ref = useCallback(
		(node: HTMLElement | null) => {
			observerRef.current?.disconnect();
			observerRef.current = null;
			if (!node) return;
			setMatches(node.getBoundingClientRect().width >= minWidth);
			const observer = new ResizeObserver(([entry]) =>
				setMatches(entry.contentRect.width >= minWidth),
			);
			observer.observe(node);
			observerRef.current = observer;
		},
		[minWidth],
	);
	return [ref, matches] as const;
};
