"use client";

import { useState, useEffect, useCallback } from "react";
import { Card, SegmentedControl, Stat, BarChart, LineChart, TabBar } from "@/shared/ui";
import { RouteGuard } from "@/features/auth/ui/RouteGuard";
import { getStats, type StatsRange, type StatsData } from "@/entities/stats/api/statsRepository";
import { AddSheet } from "@/widgets/add-sheet/AddSheet";
import { getRecentFoods, type RecentDish } from "@/entities/food/api/foodRepository";
import styles from "./stats.module.css";

const RANGE_SEGMENTS = [
  { label: "Week",  value: "week"  },
  { label: "Month", value: "month" },
  { label: "Year",  value: "year"  },
] as const;

function StatsScreen() {
  const [range, setRange]     = useState<StatsRange>("week");
  const [data, setData]       = useState<StatsData | null>(null);
  const [recent, setRecent]   = useState<RecentDish[]>([]);
  const [addOpen, setAddOpen] = useState(false);
  const [loadErr, setLoadErr] = useState<string | null>(null);

  const loadStats = useCallback(async () => {
    setLoadErr(null);
    try {
      const [stats, r] = await Promise.all([getStats(range), getRecentFoods()]);
      setData(stats);
      setRecent(r);
    } catch (err) {
      setLoadErr((err as Error).message);
    }
  }, [range]);

  useEffect(() => { loadStats(); }, [loadStats]);

  return (
    <>
      <main className={styles.root}>
        <h1 className={styles.title}>Statistics</h1>

        <SegmentedControl
          segments={RANGE_SEGMENTS}
          value={range}
          onChange={(v) => setRange(v as StatsRange)}
        />

        {loadErr && <p className={styles.errMsg}>{loadErr}</p>}

        {/* Summary stats */}
        {data && (
          <Card>
            <div className={styles.statRow}>
              <Stat label="Avg / day"     value={Math.round(data.avg_per_day)} unit="kcal" />
              <Stat label="Goal"          value={data.goal}                    unit="kcal" />
              <Stat label="Days on track" value={data.days_under_goal} color="calorie" />
              <Stat label="Streak"        value={data.streak} unit="days" color="action" />
            </div>
          </Card>
        )}

        {/* Calories per day */}
        <Card>
          <p className={styles.chartLabel}>Calories per day</p>
          <BarChart
            data={data?.calories_per_day ?? []}
            goal={data?.goal}
            color="var(--action)"
            height={160}
          />
        </Card>

        {/* Weight trend */}
        {data && data.weight_trend.length >= 2 && (
          <Card>
            <p className={styles.chartLabel}>Weight trend</p>
            <LineChart
              data={data.weight_trend}
              color="var(--protein-500)"
              height={140}
            />
          </Card>
        )}

        <div style={{ height: 80 }} />
      </main>

      <TabBar onFabClick={() => setAddOpen(true)} />

      <AddSheet
        isOpen={addOpen}
        onClose={() => setAddOpen(false)}
        onAdded={loadStats}
        recent={recent}
      />
    </>
  );
}

export function StatsPage() {
  return (
    <RouteGuard>
      <StatsScreen />
    </RouteGuard>
  );
}
