"use client";

import React from "react";
import styles from "./LineChart.module.css";

export interface LineDatum {
  label: string;
  value: number;
}

interface LineChartProps {
  data: LineDatum[];
  color?: string;
  height?: number;
}

export function LineChart({ data, color = "var(--action)", height = 120 }: LineChartProps) {
  if (data.length < 2) {
    return (
      <div className={styles.empty} style={{ height }}>
        <span>Not enough data</span>
      </div>
    );
  }

  const vals = data.map((d) => d.value);
  const min  = Math.min(...vals);
  const max  = Math.max(...vals);
  const range = max - min || 1;

  const W = 320;
  const H = height - 24; // reserve bottom for labels
  const stepX = W / (data.length - 1);

  const points = data.map((d, i) => ({
    x: i * stepX,
    y: H - ((d.value - min) / range) * H,
    label: d.label,
  }));

  const polyline = points.map((p) => `${p.x},${p.y}`).join(" ");

  // area fill below the line
  const areaPath = [
    `M ${points[0].x} ${H}`,
    ...points.map((p) => `L ${p.x} ${p.y}`),
    `L ${points[points.length - 1].x} ${H}`,
    "Z",
  ].join(" ");

  return (
    <div className={styles.wrap} style={{ height }}>
      <svg viewBox={`0 0 ${W} ${height}`} preserveAspectRatio="none" className={styles.svg}>
        <defs>
          <linearGradient id="lineGrad" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%"   stopColor={color} stopOpacity="0.25" />
            <stop offset="100%" stopColor={color} stopOpacity="0" />
          </linearGradient>
        </defs>
        <path d={areaPath} fill="url(#lineGrad)" />
        <polyline points={polyline} fill="none" stroke={color} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
        {points.map((p, i) => (
          <circle key={i} cx={p.x} cy={p.y} r="3" fill={color} />
        ))}
        {/* X-axis labels */}
        {points.map((p, i) => (
          <text
            key={`l${i}`}
            x={p.x}
            y={height - 4}
            textAnchor="middle"
            fontSize="10"
            fill="var(--text-tertiary)"
          >
            {p.label}
          </text>
        ))}
      </svg>
    </div>
  );
}
