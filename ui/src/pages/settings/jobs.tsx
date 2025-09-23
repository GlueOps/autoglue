import {useCallback, useEffect, useMemo, useState} from "react";
import {api} from "@/lib/api.ts";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table.tsx";

export interface KPI {
    RunningNow: number;
    DueNow: number;
    ScheduledFuture: number;
    Succeeded24h: number;
    Failed24h: number;
    Retryable: number;
}

export interface QueueRollup {
    QueueName: string;
    Running: number;
    QueuedDue: number;
    QueuedFuture: number;
    Success24h: number;
    Failed24h: number;
    AvgDurationSecs: number;
}


export interface JobListItem {
    id: string;
    queue_name: string;
    status: string;
    retry_count: number;
    max_retry: number;
    scheduled_at: string; // ISO
    started_at?: string; // ISO
    updated_at: string; // ISO
    last_error?: string;
}


const fmtNumber = (n: number | undefined) => (n ?? 0).toLocaleString();


const fmtSeconds = (secs: number) => {
    if (!isFinite(secs) || secs <= 0) return "–";
    if (secs < 60) return `${secs.toFixed(0)}s`;
    if (secs < 3600) return `${Math.floor(secs / 60)}m ${Math.floor(secs % 60)}s`;
    const h = Math.floor(secs / 3600);
    const m = Math.floor((secs % 3600) / 60);
    return `${h}h ${m}m`;
};

export default function JobsDashboard() {
    const [kpi, setKpi] = useState<KPI | null>(null);
    const [queues, setQueues] = useState<QueueRollup[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [autoRefresh, setAutoRefresh] = useState(true);
    const [refreshMs, setRefreshMs] = useState(5000);

    const loadAll = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const [k, q] = await Promise.all([
                api.get<KPI>("/api/v1/jobs/kpi"),
                api.get<QueueRollup[]>("/api/v1/jobs/queues"),
            ])
            setKpi(k)
            setQueues(q)
        } catch (e: any) {
            setError(e.message || String(e));
        } finally {
            setLoading(false);
        }

    }, [])

    useEffect(() => { void loadAll(); }, [loadAll]);


    useEffect(() => {
        if (!autoRefresh) return;
        const id = setInterval(loadAll, refreshMs);
        return () => clearInterval(id);
    }, [autoRefresh, refreshMs, loadAll]);

    const totals = useMemo(() => ({
        queues: queues.length,
        running: queues.reduce((s, q) => s + q.Running, 0),
        due: queues.reduce((s, q) => s + q.QueuedDue, 0),
        future: queues.reduce((s, q) => s + q.QueuedFuture, 0),
    }), [queues]);


    return (
        <div className="p-6 space-y-6">
            <header className="flex items-center justify-between gap-3">
                <h1 className="text-2xl font-semibold tracking-tight">Jobs Dashboard</h1>
                <div className="flex items-center gap-3">
                    <label className="flex items-center gap-2 text-sm">
                        <input type="checkbox" className="h-4 w-4" checked={autoRefresh} onChange={(e)=>setAutoRefresh(e.target.checked)} />
                        Auto refresh
                    </label>
                    <select
                        className="border rounded px-2 py-1 text-sm"
                        value={refreshMs}
                        onChange={(e)=>setRefreshMs(parseInt(e.target.value))}
                    >
                        <option value={3000}>3s</option>
                        <option value={5000}>5s</option>
                        <option value={10000}>10s</option>
                        <option value={30000}>30s</option>
                    </select>

                    <button
                        className="px-3 py-1.5 rounded bg-slate-900 text-white text-sm hover:opacity-90"
                        onClick={loadAll}
                        disabled={loading}
                    >{loading ? "Refreshing…" : "Refresh"}</button>
                </div>
            </header>

            {error && (
                <div className="rounded border border-red-300 bg-red-50 text-red-800 p-3 text-sm">
                    {error}
                </div>
            )}

            {/* KPI cards */}
            <section className="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-6">
                <KpiCard label="Running" value={fmtNumber(kpi?.RunningNow)} />
                <KpiCard label="Due now" value={fmtNumber(kpi?.DueNow)} />
                <KpiCard label="Scheduled" value={fmtNumber(kpi?.ScheduledFuture)} />
                <KpiCard label="Succeeded (24h)" value={fmtNumber(kpi?.Succeeded24h)} />
                <KpiCard label="Failed (24h)" value={fmtNumber(kpi?.Failed24h)} />
                <KpiCard label="Retryable" value={fmtNumber(kpi?.Retryable)} />
            </section>

            {/* Per-queue table */}
            <section className="space-y-2">
                <div className="flex items-center justify-between">
                    <h2 className="text-lg font-medium">Queues <span className="text-slate-500 text-sm">({totals.queues})</span></h2>
                </div>
                <div className="overflow-x-auto rounded border">
                    <Table className="min-w-full text-sm">
                        <TableHeader>
                            <TableRow>
                                <TableHead>Queue</TableHead>
                                <TableHead className='text-right'>Running</TableHead>
                                <TableHead className='text-right'>Due</TableHead>
                                <TableHead className='text-right'>Future</TableHead>
                                <TableHead className='text-right'>Success 24h</TableHead>
                                <TableHead className='text-right'>Failed 24h</TableHead>
                                <TableHead className='text-right'>Avg Duration</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {queues.map((q) => (
                                <TableRow key={q.QueueName} className="border-t">
                                    <TableCell>{q.QueueName}</TableCell>
                                    <TableCell className='text-right'>{q.Running}</TableCell>
                                    <TableCell className='text-right'>{q.QueuedDue}</TableCell>
                                    <TableCell className='text-right'>{q.QueuedFuture}</TableCell>
                                    <TableCell className='text-right'>{q.Success24h}</TableCell>
                                    <TableCell className='text-right'>{q.Failed24h}</TableCell>
                                    <TableCell className='text-right'>{fmtSeconds(q.AvgDurationSecs)}</TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </div>
            </section>

            {/* <ManualEnqueue onSubmitted={loadAll} /> */}
        </div>
    )
}

function KpiCard({ label, value }: { label: string; value: string }) {
    return (
        <Card>
            <CardHeader>
                <CardTitle>{label}</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="mt-1 text-2xl font-semibold">{value}</div>
            </CardContent>
        </Card>
    );
}
