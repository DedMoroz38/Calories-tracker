import React from "react";
import styles from "./ProgressBar.module.css";

export type ProgressColor = "action" | "carb" | "fat" | "protein" | "calorie";

interface ProgressBarProps {
  value: number;         // current value
  max: number;           // goal value
  color?: ProgressColor;
  showLabel?: boolean;
  label?: string;
  className?: string;
}

export function ProgressBar({
  value,
  max,
  color = "action",
  showLabel = false,
  label,
  className = "",
}: ProgressBarProps) {
  const pct = max > 0 ? Math.min(100, (value / max) * 100) : 0;

  return (
    <div className={[styles.wrap, className].filter(Boolean).join(" ")}>
      {(showLabel || label) && (
        <div className={styles.row}>
          {label && <span className={styles.labelText}>{label}</span>}
          {showLabel && (
            <span className={styles.values}>
              {value} / {max}
              <span className={styles.pct}>{Math.round(pct)}%</span>
            </span>
          )}
        </div>
      )}
      <div className={styles.track}>
        <div
          className={[styles.fill, styles[color]].join(" ")}
          style={{ width: `${pct}%` }}
          role="progressbar"
          aria-valuenow={value}
          aria-valuemax={max}
          aria-valuemin={0}
        />
      </div>
    </div>
  );
}
