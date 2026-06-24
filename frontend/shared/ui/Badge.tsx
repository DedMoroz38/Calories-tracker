import React from "react";
import styles from "./Badge.module.css";

export type BadgeColor = "default" | "carb" | "fat" | "protein" | "calorie" | "action";

interface BadgeProps {
  children: React.ReactNode;
  color?: BadgeColor;
  className?: string;
}

export function Badge({ children, color = "default", className = "" }: BadgeProps) {
  return (
    <span className={[styles.badge, styles[color], className].filter(Boolean).join(" ")}>
      {children}
    </span>
  );
}
