"use client";

import { useState } from "react";
import { X } from "lucide-react";
import { Button, Card, IconButton, Input, SegmentedControl, MealRow } from "@/shared/ui";
import { createFood } from "@/entities/food/api/foodRepository";
import type { RecentDish } from "@/entities/food/api/foodRepository";
import styles from "./AddSheet.module.css";

type Mode = "quick" | "dish" | "recent";

const MODE_SEGMENTS = [
  { label: "Quick add", value: "quick"  },
  { label: "New dish",  value: "dish"   },
  { label: "Recent",    value: "recent" },
] as const;

interface AddSheetProps {
  isOpen:  boolean;
  onClose: () => void;
  onAdded: () => void;
  recent:  RecentDish[];
}

export function AddSheet({ isOpen, onClose, onAdded, recent }: AddSheetProps) {
  const [mode, setMode] = useState<Mode>("quick");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Quick add fields
  const [qCals, setQCals] = useState("");
  const [qCarbs, setQCarbs] = useState("");
  const [qFat, setQFat] = useState("");
  const [qProtein, setQProtein] = useState("");

  // New dish fields
  const [dName, setDName] = useState("");
  const [dCals, setDCals] = useState("");
  const [dCarbs, setDCarbs] = useState("");
  const [dFat, setDFat] = useState("");
  const [dProtein, setDProtein] = useState("");

  const reset = () => {
    setQCals(""); setQCarbs(""); setQFat(""); setQProtein("");
    setDName(""); setDCals(""); setDCarbs(""); setDFat(""); setDProtein("");
    setError(null);
  };

  const handleClose = () => { reset(); onClose(); };

  const addEntry = async (name: string, cals: number, carbs?: number, fat?: number, protein?: number) => {
    setLoading(true);
    setError(null);
    try {
      await createFood({ name, calories: cals, carbs, fat, protein });
      onAdded();
      handleClose();
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const submitQuick = () => {
    if (!qCals) return;
    addEntry(
      "Quick add",
      Number(qCals),
      qCarbs ? Number(qCarbs) : undefined,
      qFat   ? Number(qFat)   : undefined,
      qProtein ? Number(qProtein) : undefined,
    );
  };

  const submitDish = () => {
    if (!dName || !dCals) return;
    addEntry(
      dName,
      Number(dCals),
      dCarbs ? Number(dCarbs) : undefined,
      dFat   ? Number(dFat)   : undefined,
      dProtein ? Number(dProtein) : undefined,
    );
  };

  const reAddRecent = (dish: RecentDish) => {
    addEntry(dish.name, dish.calories, dish.carbs, dish.fat, dish.protein);
  };

  if (!isOpen) return null;

  return (
    <div className={styles.overlay} onClick={handleClose}>
      <div className={styles.sheet} onClick={(e) => e.stopPropagation()}>
        {/* Handle */}
        <div className={styles.handle} />

        {/* Header */}
        <div className={styles.header}>
          <h2 className={styles.title}>Add food</h2>
          <IconButton label="Close" size="sm" onClick={handleClose}>
            <X size={18} />
          </IconButton>
        </div>

        <SegmentedControl segments={MODE_SEGMENTS} value={mode} onChange={setMode} />

        {/* Quick add */}
        {mode === "quick" && (
          <div className={styles.fields}>
            <Input label="Calories *" type="number" placeholder="e.g. 350" value={qCals} onChange={(e) => setQCals(e.target.value)} />
            <div className={styles.macroRow}>
              <Input label="Carbs (g)" type="number" placeholder="0" value={qCarbs} onChange={(e) => setQCarbs(e.target.value)} />
              <Input label="Fat (g)"   type="number" placeholder="0" value={qFat}   onChange={(e) => setQFat(e.target.value)}   />
              <Input label="Protein(g)"type="number" placeholder="0" value={qProtein}onChange={(e) => setQProtein(e.target.value)} />
            </div>
          </div>
        )}

        {/* New dish */}
        {mode === "dish" && (
          <div className={styles.fields}>
            <Input label="Name *"     placeholder="e.g. Oat porridge" value={dName} onChange={(e) => setDName(e.target.value)} />
            <Input label="Calories *" type="number" placeholder="e.g. 350" value={dCals} onChange={(e) => setDCals(e.target.value)} />
            <div className={styles.macroRow}>
              <Input label="Carbs (g)" type="number" placeholder="0" value={dCarbs}   onChange={(e) => setDCarbs(e.target.value)}   />
              <Input label="Fat (g)"   type="number" placeholder="0" value={dFat}     onChange={(e) => setDFat(e.target.value)}     />
              <Input label="Protein(g)"type="number" placeholder="0" value={dProtein} onChange={(e) => setDProtein(e.target.value)} />
            </div>
          </div>
        )}

        {/* Recent */}
        {mode === "recent" && (
          <div className={styles.recentList}>
            {recent.length === 0 ? (
              <p className={styles.empty}>No recent dishes yet.</p>
            ) : (
              recent.map((d) => (
                <MealRow
                  key={d.name}
                  name={d.name}
                  calories={d.calories}
                  carbs={d.carbs}
                  fat={d.fat}
                  protein={d.protein}
                  onClick={() => reAddRecent(d)}
                />
              ))
            )}
          </div>
        )}

        {error && <p className={styles.error}>{error}</p>}

        {mode !== "recent" && (
          <Button
            fullWidth
            loading={loading}
            disabled={(mode === "quick" && !qCals) || (mode === "dish" && (!dName || !dCals))}
            onClick={mode === "quick" ? submitQuick : submitDish}
          >
            Add
          </Button>
        )}
      </div>
    </div>
  );
}
