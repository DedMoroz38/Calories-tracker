import React from "react";
import { Flame } from "lucide-react";
import styles from "./DayDot.module.css";

export type DayState = "hit" | "miss" | "today" | "future";

interface DayDotProps {
  letter: string;
  state:  DayState;
  /** When true the dot shows the streak flame instead of a plain fill. */
  fire?:  boolean;
}

export function DayDot({ letter, state, fire = false }: DayDotProps) {
  return (
    <div className={styles.wrap}>
      <div className={[styles.dot, styles[state], fire ? styles.fire : ""].join(" ")}>
        {fire && <Flame size={18} className={styles.flame} fill="currentColor" strokeWidth={2.5} />}
      </div>
      <span className={styles.letter}>{letter}</span>
    </div>
  );
}
