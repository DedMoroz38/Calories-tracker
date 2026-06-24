"use client";

import React from "react";
import styles from "./Input.module.css";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  hint?:  string;
  error?: string;
  leftSlot?:  React.ReactNode;
  rightSlot?: React.ReactNode;
}

export function Input({
  label,
  hint,
  error,
  leftSlot,
  rightSlot,
  className = "",
  id,
  ...rest
}: InputProps) {
  const fieldId = id ?? label?.toLowerCase().replace(/\s+/g, "-");

  return (
    <div className={styles.wrapper}>
      {label && (
        <label htmlFor={fieldId} className={styles.label}>
          {label}
        </label>
      )}
      <div className={[styles.inputWrap, error ? styles.hasError : ""].filter(Boolean).join(" ")}>
        {leftSlot && <span className={styles.slot}>{leftSlot}</span>}
        <input
          id={fieldId}
          className={[styles.input, className].filter(Boolean).join(" ")}
          {...rest}
        />
        {rightSlot && <span className={[styles.slot, styles.slotRight].join(" ")}>{rightSlot}</span>}
      </div>
      {error && <p className={styles.error}>{error}</p>}
      {hint && !error && <p className={styles.hint}>{hint}</p>}
    </div>
  );
}
