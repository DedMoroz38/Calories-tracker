import React from "react";
import styles from "./DayDot.module.css";

export type DayState = "hit" | "miss" | "today" | "future";

interface DayDotProps {
  letter: string;
  state:  DayState;
}

export function DayDot({ letter, state }: DayDotProps) {
  return (
    <div className={styles.wrap}>
      <div className={[styles.dot, styles[state]].join(" ")} />
      <span className={styles.letter}>{letter}</span>
    </div>
  );
}
