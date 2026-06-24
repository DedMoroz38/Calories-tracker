"use client";

import { useState, useEffect, useCallback } from "react";
import { Scale } from "lucide-react";
import {
  MacroRing, ProgressBar, Card, DayDot, Stat, IconButton, MealRow,
  TabBar,
} from "@/shared/ui";
import { RouteGuard } from "@/features/auth/ui/RouteGuard";
import { getSummary, type Summary } from "@/entities/stats/api/summaryRepository";
import { getFoods, deleteFood, getRecentFoods, type FoodEntry, type RecentDish } from "@/entities/food/api/foodRepository";
import { AddSheet } from "@/widgets/add-sheet/AddSheet";
import { WeightSheet } from "@/widgets/weight-sheet/WeightSheet";
import styles from "./home.module.css";

function HomeScreen() {
  const [summary, setSummary]     = useState<Summary | null>(null);
  const [foods, setFoods]         = useState<FoodEntry[]>([]);
  const [recent, setRecent]       = useState<RecentDish[]>([]);
  const [addOpen, setAddOpen]     = useState(false);
  const [weightOpen, setWeightOpen] = useState(false);
  const [loadErr, setLoadErr]     = useState<string | null>(null);

  const loadData = useCallback(async () => {
    setLoadErr(null);
    try {
      const [sum, f, r] = await Promise.all([
        getSummary(),
        getFoods(),
        getRecentFoods(),
      ]);
      setSummary(sum);
      setFoods(f);
      setRecent(r);
    } catch (err) {
      setLoadErr((err as Error).message);
    }
  }, []);

  useEffect(() => { loadData(); }, [loadData]);

  const handleDelete = async (id: number) => {
    try {
      await deleteFood(id);
      await loadData();
    } catch { /* silently retry */ }
  };

  if (loadErr) {
    return (
      <main className={styles.root}>
        <p className={styles.errMsg}>{loadErr}</p>
      </main>
    );
  }

  const totals = summary?.totals ?? { calories: 0, carbs: 0, fat: 0, protein: 0 };
  const goals  = summary?.goals  ?? { calories: 2000, carbs: 250, fat: 70, protein: 150 };
  const weight = summary?.weight ?? { current: 0, goal: 0 };
  const week   = summary?.week   ?? [];

  const caloriesFrac = goals.calories  > 0 ? totals.calories / goals.calories  : 0;

  return (
    <>
      <main className={styles.root}>
        {/* ── Header ── */}
        <div className={styles.header}>
          <div>
            <p className={styles.greeting}>Today</p>
            <h1 className={styles.date}>
              {new Date().toLocaleDateString("en-US", { weekday: "long", month: "long", day: "numeric" })}
            </h1>
          </div>
          <IconButton label="Log weight" variant="filled" onClick={() => setWeightOpen(true)}>
            <Scale size={20} />
          </IconButton>
        </div>

        {/* ── Hero ring + calorie stats ── */}
        <Card className={styles.heroCard}>
          <div className={styles.heroRow}>
            <MacroRing
              calories={caloriesFrac}
              label={String(totals.calories)}
              sublabel="kcal"
              size={128}
            />
            <div className={styles.statCol}>
              <Stat label="Goal"      value={goals.calories} unit="kcal" />
              <Stat label="Remaining" value={Math.max(0, goals.calories - totals.calories)} unit="kcal" />
              {summary?.streak ? <Stat label="Streak" value={summary.streak} unit="days" color="action" /> : null}
            </div>
          </div>

          {/* Macro progress bars */}
          <div className={styles.macros}>
            <ProgressBar value={totals.carbs}   max={goals.carbs}   color="carb"    showLabel label="Carbs" />
            <ProgressBar value={totals.fat}      max={goals.fat}     color="fat"     showLabel label="Fat" />
            <ProgressBar value={totals.protein}  max={goals.protein} color="protein" showLabel label="Protein" />
          </div>
        </Card>

        {/* ── Weight card ── */}
        {weight.current > 0 && (
          <Card padding="md">
            <div className={styles.weightRow}>
              <Stat label="Current weight" value={weight.current} unit="kg" />
              <Stat label="Goal weight"    value={weight.goal}    unit="kg" />
            </div>
          </Card>
        )}

        {/* ── Week streak dots ── */}
        {week.length > 0 && (
          <Card padding="md">
            <p className={styles.sectionLabel}>This week</p>
            <div className={styles.dots}>
              {week.map((d) => (
                <DayDot key={d.date} letter={d.letter} state={d.state} />
              ))}
            </div>
          </Card>
        )}

        {/* ── Today's food log ── */}
        <Card padding="none">
          <div className={styles.logHeader}>
            <p className={styles.sectionLabel}>Today&apos;s log</p>
          </div>
          {foods.length === 0 ? (
            <p className={styles.empty}>No entries yet. Tap + to add food.</p>
          ) : (
            <div style={{ paddingInline: 16 }}>
              {foods.map((f) => (
                <MealRow
                  key={f.id}
                  name={f.name}
                  calories={f.calories}
                  carbs={f.carbs}
                  fat={f.fat}
                  protein={f.protein}
                  onDelete={() => handleDelete(f.id)}
                />
              ))}
            </div>
          )}
        </Card>

        {/* Bottom padding so content clears the TabBar */}
        <div style={{ height: 80 }} />
      </main>

      <TabBar onFabClick={() => setAddOpen(true)} />

      <AddSheet
        isOpen={addOpen}
        onClose={() => setAddOpen(false)}
        onAdded={loadData}
        recent={recent}
      />

      <WeightSheet
        isOpen={weightOpen}
        onClose={() => setWeightOpen(false)}
        onLogged={loadData}
        currentWeight={weight.current || undefined}
        goalWeight={weight.goal || undefined}
      />
    </>
  );
}

export function HomePage() {
  return (
    <RouteGuard>
      <HomeScreen />
    </RouteGuard>
  );
}
