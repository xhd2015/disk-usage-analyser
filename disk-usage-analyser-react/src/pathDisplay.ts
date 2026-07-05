const ELLIPSIS = '…';

/**
 * Truncate a filesystem path by hiding the prefix with an ellipsis while keeping
 * the suffix fully visible. Prefers breaking at `/` boundaries.
 */
export function truncatePathKeepSuffix(path: string, maxVisibleChars: number): string {
    if (maxVisibleChars <= 0) {
        return '';
    }
    if (path.length <= maxVisibleChars) {
        return path;
    }

    const budget = maxVisibleChars - ELLIPSIS.length;
    if (budget <= 0) {
        return ELLIPSIS.slice(0, maxVisibleChars);
    }

    for (let i = 0; i < path.length; i++) {
        if (path[i] !== '/') {
            continue;
        }
        const suffix = path.slice(i);
        if (suffix.length <= budget) {
            return ELLIPSIS + suffix;
        }
    }

    const tail = path.slice(-budget);
    const slashIdx = tail.indexOf('/');
    const suffix = slashIdx >= 0 ? tail.slice(slashIdx) : tail;
    const display = ELLIPSIS + suffix;
    return display.length <= maxVisibleChars ? display : display.slice(0, maxVisibleChars);
}

export const PATH_VISIBLE_CHAR_LIMIT = 56;