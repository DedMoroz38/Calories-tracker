"use client";
/**
 * /showcase — temporary eyeball page for all shared UI components.
 * Remove or restrict before production.
 */

import { useState } from "react";
import {
  Button, Card, Badge, IconButton, Input, SegmentedControl,
  MacroRing, ProgressBar, Stat, DayDot, BarChart, LineChart, MealRow,
} from "@/shared/ui";
import { Settings, Star } from "lucide-react";

const SEGMENTS = [
  { label: "Week",  value: "week"  },
  { label: "Month", value: "month" },
  { label: "Year",  value: "year"  },
];

const BAR_DATA = [
  { label: "M", value: 2100, muted: false },
  { label: "T", value: 1850, muted: false },
  { label: "W", value: 2400, muted: false },
  { label: "T", value: 2200, muted: false },
  { label: "F", value: 1600, muted: false },
  { label: "S", value: 900,  muted: true  },
  { label: "S", value: 0,    muted: true  },
];

const LINE_DATA = [
  { label: "Wk1", value: 86.2 },
  { label: "Wk2", value: 85.5 },
  { label: "Wk3", value: 84.8 },
  { label: "Wk4", value: 84.1 },
];

export default function ShowcasePage() {
  const [seg, setSeg] = useState("week");
  const [inputVal, setInputVal] = useState("");

  return (
    <div style={{ padding: 24, maxWidth: 430, margin: "0 auto", display: "flex", flexDirection: "column", gap: 32 }}>
      <h1 style={{ fontSize: 24, fontWeight: 700 }}>UI Showcase</h1>

      {/* Buttons */}
      <section style={{ display: "flex", flexDirection: "column", gap: 8 }}>
        <h2 style={{ fontSize: 14, color: "var(--text-tertiary)" }}>BUTTONS</h2>
        <Button>Primary</Button>
        <Button variant="secondary">Secondary</Button>
        <Button variant="ghost">Ghost</Button>
        <Button variant="destructive">Destructive</Button>
        <Button size="sm">Small</Button>
        <Button loading>Loading…</Button>
        <Button leftIcon={<Star size={16} />}>With Icon</Button>
      </section>

      {/* Badges */}
      <section style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
        <h2 style={{ fontSize: 14, color: "var(--text-tertiary)", width: "100%" }}>BADGES</h2>
        <Badge>Default</Badge>
        <Badge color="carb">Carbs</Badge>
        <Badge color="fat">Fat</Badge>
        <Badge color="protein">Protein</Badge>
        <Badge color="calorie">Calorie</Badge>
        <Badge color="action">Action</Badge>
      </section>

      {/* IconButton */}
      <section style={{ display: "flex", gap: 8 }}>
        <IconButton label="settings" variant="ghost"><Settings size={20} /></IconButton>
        <IconButton label="settings" variant="filled"><Settings size={20} /></IconButton>
      </section>

      {/* Input */}
      <section style={{ display: "flex", flexDirection: "column", gap: 8 }}>
        <h2 style={{ fontSize: 14, color: "var(--text-tertiary)" }}>INPUTS</h2>
        <Input
          label="Calories"
          placeholder="Enter value"
          type="number"
          value={inputVal}
          onChange={(e) => setInputVal(e.target.value)}
          hint="kcal per serving"
        />
        <Input label="With error" placeholder="Enter value" error="Required field" />
      </section>

      {/* SegmentedControl */}
      <section>
        <h2 style={{ fontSize: 14, color: "var(--text-tertiary)", marginBottom: 8 }}>SEGMENTED CONTROL</h2>
        <SegmentedControl segments={SEGMENTS} value={seg} onChange={setSeg} />
      </section>

      {/* MacroRing */}
      <section style={{ display: "flex", gap: 24, alignItems: "center" }}>
        <h2 style={{ fontSize: 14, color: "var(--text-tertiary)", width: "100%" }}>MACRO RING</h2>
        <MacroRing calories={0.6} label="1440" sublabel="kcal" />
        <MacroRing calories={0} size={80} />
      </section>

      {/* ProgressBar */}
      <section style={{ display: "flex", flexDirection: "column", gap: 10 }}>
        <h2 style={{ fontSize: 14, color: "var(--text-tertiary)" }}>PROGRESS BARS</h2>
        <ProgressBar value={1440} max={2400} color="calorie" showLabel label="Calories" />
        <ProgressBar value={120}  max={270}  color="carb"    showLabel label="Carbs" />
        <ProgressBar value={45}   max={80}   color="fat"     showLabel label="Fat" />
        <ProgressBar value={95}   max={150}  color="protein" showLabel label="Protein" />
      </section>

      {/* Stats */}
      <section style={{ display: "flex", gap: 16 }}>
        <Stat label="Calories" value={1440} unit="kcal" />
        <Stat label="Streak"   value={12}   unit="days" color="action" />
        <Stat label="Weight"   value={84.2} unit="kg" />
      </section>

      {/* DayDots */}
      <section>
        <h2 style={{ fontSize: 14, color: "var(--text-tertiary)", marginBottom: 8 }}>DAY DOTS</h2>
        <div style={{ display: "flex", gap: 4 }}>
          {(["hit","hit","miss","today","future","future","future"] as const).map((s, i) => (
            <DayDot key={i} letter={["M","T","W","T","F","S","S"][i]} state={s} />
          ))}
        </div>
      </section>

      {/* BarChart */}
      <section>
        <h2 style={{ fontSize: 14, color: "var(--text-tertiary)", marginBottom: 8 }}>BAR CHART</h2>
        <Card>
          <BarChart data={BAR_DATA} goal={2000} />
        </Card>
      </section>

      {/* LineChart */}
      <section>
        <h2 style={{ fontSize: 14, color: "var(--text-tertiary)", marginBottom: 8 }}>LINE CHART</h2>
        <Card>
          <LineChart data={LINE_DATA} color="var(--protein-500)" />
        </Card>
      </section>

      {/* MealRow */}
      <section>
        <h2 style={{ fontSize: 14, color: "var(--text-tertiary)", marginBottom: 8 }}>MEAL ROWS</h2>
        <Card padding="none">
          <div style={{ padding: "0 16px" }}>
            <MealRow name="Oat porridge"   calories={350} carbs={55} fat={6}  protein={12} onDelete={() => {}} />
            <MealRow name="Chicken breast" calories={220} carbs={0}  fat={4}  protein={44} onDelete={() => {}} />
            <MealRow name="Quick add"      calories={500} />
          </div>
        </Card>
      </section>
    </div>
  );
}
