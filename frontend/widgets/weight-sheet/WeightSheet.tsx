"use client";

import { useState } from "react";
import { X } from "lucide-react";
import { Button, IconButton, Input } from "@/shared/ui";
import { logWeight } from "@/entities/weight/api/weightRepository";
import styles from "./WeightSheet.module.css";

interface WeightSheetProps {
  isOpen:      boolean;
  onClose:     () => void;
  onLogged:    () => void;
  currentWeight?: number;
  goalWeight?:    number;
}

export function WeightSheet({ isOpen, onClose, onLogged, currentWeight, goalWeight }: WeightSheetProps) {
  const [value, setValue] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleClose = () => { setValue(""); setError(null); onClose(); };

  const submit = async () => {
    if (!value) return;
    setLoading(true);
    setError(null);
    try {
      await logWeight(Number(value));
      onLogged();
      handleClose();
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  const distToGoal =
    currentWeight && goalWeight
      ? Math.abs(currentWeight - goalWeight).toFixed(1)
      : null;

  return (
    <div className={styles.overlay} onClick={handleClose}>
      <div className={styles.sheet} onClick={(e) => e.stopPropagation()}>
        <div className={styles.handle} />

        <div className={styles.header}>
          <h2 className={styles.title}>Log weight</h2>
          <IconButton label="Close" size="sm" onClick={handleClose}>
            <X size={18} />
          </IconButton>
        </div>

        {distToGoal && (
          <p className={styles.sub}>
            {distToGoal} kg to your {goalWeight} kg goal.
          </p>
        )}

        <Input
          label="Weight (kg)"
          type="number"
          placeholder={currentWeight ? String(currentWeight) : "e.g. 84.5"}
          value={value}
          onChange={(e) => setValue(e.target.value)}
        />

        {error && <p className={styles.error}>{error}</p>}

        <Button fullWidth loading={loading} disabled={!value} onClick={submit}>
          Save weight
        </Button>
      </div>
    </div>
  );
}
