export interface FileInfo {
    name: string;
    size: number;
    isDir: boolean;
    status: 'pending' | 'done';
}

export interface UsageResponse {
    path: string;
    totalSize: number;
    items: FileInfo[];
    error?: string;
}

export class DiskUsageAPI {
    static streamUsage(dirPath: string, callbacks: {
        onPath: (path: string) => void;
        onItem: (item: FileInfo) => void;
        onDone: () => void;
        onError: (error: string) => void;
    }): EventSource {
        const url = dirPath ? `/api/usage?path=${encodeURIComponent(dirPath)}` : '/api/usage';
        const es = new EventSource(url);

        es.addEventListener('path', (e) => {
            const d = JSON.parse((e as MessageEvent).data);
            callbacks.onPath(d.path);
        });

        es.addEventListener('item', (e) => {
            const item: FileInfo = JSON.parse((e as MessageEvent).data);
            callbacks.onItem(item);
        });

        es.addEventListener('done', () => {
            callbacks.onDone();
            es.close();
        });

        es.addEventListener('server_error', (e) => {
            const d = JSON.parse((e as MessageEvent).data);
            callbacks.onError(d.error);
            es.close();
        });

        es.onerror = (e) => {
            console.error('SSE Error', e);
            if (es.readyState === EventSource.CLOSED) return;

            // Close on error
            es.close();
            callbacks.onError('Connection error');
        };

        return es;
    }

    static async moveToTrash(path: string): Promise<void> {
        const res = await fetch(`/api/moveToTrash?path=${encodeURIComponent(path)}`, {
            method: 'POST'
        });
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text);
        }
    }

    static async refresh(path: string): Promise<void> {
        const res = await fetch(`/api/refresh?path=${encodeURIComponent(path)}`, {
            method: 'POST'
        });
        if (!res.ok) {
            const text = await res.text();
            throw new Error(text);
        }
    }
}
