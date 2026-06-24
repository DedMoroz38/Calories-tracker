"use client";

import React from "react";
import styles from "./SegmentedControl.module.css";

interface Segment<T extends string> {
  label: string;
  value: T;
}

interface SegmentedControlProps<T extends string> {
  segments: readonly Segment<T>[];
  value: T;
  onChange: (value: T) => void;
  className?: string;
}

export function SegmentedControl<T extends string>({
  segments,
  value,
  onChange,
  className = "",
}: SegmentedControlProps<T>) {
  return (
    <div className={[styles.track, className].filter(Boolean).join(" ")} role="tablist">
      {segments.map((seg) => (
        <button
          key={seg.value}
          role="tab"
          aria-selected={seg.value === value}
          className={[styles.seg, seg.value === value ? styles.active : ""].filter(Boolean).join(" ")}
          onClick={() => onChange(seg.value)}
        >
          {seg.label}
        </button>
      ))}
    </div>
  );
}
