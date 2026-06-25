"use client";

import React from "react";
import styles from "./BarChart.module.css";

export interface BarDatum {
  label: string;
  value: number;
  muted?: boolean;
}

interface BarChartProps {
  data: BarDatum[];
  goal?: number;
  color?: string;
  /** Height of the chart area in px */
  height?: number;
}

export function BarChart({ data, goal, color = "var(--action)", height = 140 }: BarChartProps) {
  const max = Math.max(...data.map((d) => d.value), goal ?? 0, 1);

  return (
    <div className={styles.wrap} style={{ height }}>
      {/* Goal line */}
      {goal && (
        <div
          className={styles.goalLine}
          style={{ bottom: `${(goal / max) * 100}%` }}
          aria-label={`Goal: ${goal}`}
        />
      )}
      <div className={styles.bars}>
        {data.map((d, i) => (
          <div key={i} className={styles.barCol}>
            <div
              className={[styles.bar, d.muted ? styles.muted : ""].filter(Boolean).join(" ")}
              style={{
                height: `${(d.value / max) * 100}%`,
                background: d.muted ? "var(--surface-card-2)" : color,
              }}
            />
            <span className={styles.barLabel}>{d.label}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
