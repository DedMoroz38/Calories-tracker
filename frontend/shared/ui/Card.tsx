import React from "react";
import styles from "./Card.module.css";

interface CardProps {
  children: React.ReactNode;
  padding?: "none" | "sm" | "md" | "lg";
  className?: string;
  onClick?: () => void;
}

export function Card({ children, padding = "md", className = "", onClick }: CardProps) {
  const Tag = onClick ? "button" : "div";
  return (
    <Tag
      className={[styles.card, styles[`pad-${padding}`], onClick ? styles.clickable : "", className]
        .filter(Boolean)
        .join(" ")}
      onClick={onClick}
    >
      {children}
    </Tag>
  );
}
