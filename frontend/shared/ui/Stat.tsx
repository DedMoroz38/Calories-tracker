import React from "react";
import styles from "./Stat.module.css";

interface StatProps {
  label: string;
  value: string | number;
  unit?: string;
  sub?: string;
  color?: "default" | "carb" | "fat" | "protein" | "calorie" | "action";
}

export function Stat({ label, value, unit, sub, color = "default" }: StatProps) {
  return (
    <div className={styles.stat}>
      <span className={styles.label}>{label}</span>
      <span className={[styles.value, color !== "default" ? styles[color] : ""].filter(Boolean).join(" ")}>
        {value}
        {unit && <span className={styles.unit}>{unit}</span>}
      </span>
      {sub && <span className={styles.sub}>{sub}</span>}
    </div>
  );
}
