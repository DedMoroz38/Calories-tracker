"use client";

import React from "react";
import styles from "./MacroRing.module.css";

interface MacroRingProps {
  /** 0–1 fill fraction of the calorie goal */
  calories: number;
  /** Inner label — typically current calorie count */
  label?:   string;
  sublabel?: string;
  size?: number;
}

const RADIUS = 44;
const STROKE = 7;
const CX = 56;
const CY = 56;
const CIRCUMFERENCE = 2 * Math.PI * RADIUS;

/**
 * A single calorie-progress donut. The arc sweeps the full circle in
 * proportion to how much of the calorie goal has been consumed, and the
 * center shows the calorie count plus the percentage of goal reached.
 */
export function MacroRing({
  calories,
  label,
  sublabel,
  size = 112,
}: MacroRingProps) {
  const frac    = Math.max(0, Math.min(1, calories));
  const dashLen = frac * CIRCUMFERENCE;
  const percent = Math.round(calories * 100);

  return (
    <div className={styles.ring} style={{ width: size, height: size }}>
      <svg viewBox="0 0 112 112" fill="none">
        {/* Track */}
        <circle
          cx={CX} cy={CY} r={RADIUS}
          stroke="var(--surface-card-2)"
          strokeWidth={STROKE}
          fill="none"
        />
        {/* Calorie progress arc */}
        <circle
          cx={CX} cy={CY} r={RADIUS}
          stroke="var(--calorie-500)"
          strokeWidth={STROKE}
          strokeLinecap="round"
          fill="none"
          strokeDasharray={`${dashLen} ${CIRCUMFERENCE}`}
          style={{ transform: "rotate(-90deg)", transformOrigin: "50% 50%" }}
        />
      </svg>
      {(label || sublabel) && (
        <div className={styles.inner}>
          {label    && <span className={styles.label}>{label}</span>}
          {sublabel && <span className={styles.sub}>{sublabel}</span>}
          <span className={styles.percent}>{percent}%</span>
        </div>
      )}
    </div>
  );
}
