"use client";

import React from "react";
import styles from "./IconButton.module.css";

interface IconButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  label: string;
  size?: "sm" | "md" | "lg";
  variant?: "ghost" | "filled";
}

export function IconButton({
  label,
  size = "md",
  variant = "ghost",
  children,
  className = "",
  ...rest
}: IconButtonProps) {
  return (
    <button
      aria-label={label}
      className={[styles.iconBtn, styles[size], styles[variant], className].filter(Boolean).join(" ")}
      {...rest}
    >
      {children}
    </button>
  );
}
