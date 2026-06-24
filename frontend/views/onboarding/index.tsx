"use client";

import { Button, Input, SegmentedControl, Card } from "@/shared/ui";
import { RouteGuard } from "@/features/auth/ui/RouteGuard";
import { useOnboarding, type Direction } from "@/features/onboarding/model/useOnboarding";
import styles from "./onboarding.module.css";

const DIRECTION_SEGMENTS = [
  { label: "Lose",     value: "lose"     },
  { label: "Maintain", value: "maintain" },
  { label: "Gain",     value: "gain"     },
] as const;

function Step1({
  state,
  set,
  onNext,
}: ReturnType<typeof useOnboarding> & { onNext: () => void }) {
  const canNext = state.currentWeight !== "" && state.goalWeight !== "";
  return (
    <div className={styles.stepWrap}>
      <h2 className={styles.stepTitle}>Your weights</h2>
      <p className={styles.stepSub}>We use this to estimate your starting point.</p>
      <Input
        label="Current weight (kg)"
        type="number"
        placeholder="e.g. 85"
        value={state.currentWeight}
        onChange={(e) => set("currentWeight", e.target.value)}
      />
      <Input
        label="Goal weight (kg)"
        type="number"
        placeholder="e.g. 78"
        value={state.goalWeight}
        onChange={(e) => set("goalWeight", e.target.value)}
      />
      <div className={styles.fieldWrap}>
        <label className={styles.label}>Direction</label>
        <SegmentedControl
          segments={DIRECTION_SEGMENTS}
          value={state.direction}
          onChange={(v) => set("direction", v as Direction)}
        />
      </div>
      <Button fullWidth disabled={!canNext} onClick={onNext}>
        Continue
      </Button>
    </div>
  );
}

function Step2({
  state,
  set,
  onNext,
  onBack,
}: ReturnType<typeof useOnboarding> & { onNext: () => void; onBack: () => void }) {
  const canNext = state.calorieGoal !== "" && Number(state.calorieGoal) > 0;
  return (
    <div className={styles.stepWrap}>
      <h2 className={styles.stepTitle}>Daily calorie goal</h2>
      <p className={styles.stepSub}>How many calories do you want to eat each day?</p>
      <Input
        label="Calorie goal (kcal)"
        type="number"
        placeholder="e.g. 2200"
        value={state.calorieGoal}
        onChange={(e) => set("calorieGoal", e.target.value)}
      />
      <div className={styles.btnRow}>
        <Button variant="ghost" onClick={onBack}>Back</Button>
        <Button disabled={!canNext} onClick={onNext}>Continue</Button>
      </div>
    </div>
  );
}

function Step3({
  state,
  set,
  onBack,
  submit,
}: ReturnType<typeof useOnboarding> & { onBack: () => void }) {
  return (
    <div className={styles.stepWrap}>
      <h2 className={styles.stepTitle}>Macro goals</h2>
      <p className={styles.stepSub}>Optional — leave blank to skip.</p>
      <Input
        label="Carbs goal (g)"
        type="number"
        placeholder="e.g. 270"
        value={state.carbsGoal}
        onChange={(e) => set("carbsGoal", e.target.value)}
      />
      <Input
        label="Fat goal (g)"
        type="number"
        placeholder="e.g. 80"
        value={state.fatGoal}
        onChange={(e) => set("fatGoal", e.target.value)}
      />
      <Input
        label="Protein goal (g)"
        type="number"
        placeholder="e.g. 150"
        value={state.proteinGoal}
        onChange={(e) => set("proteinGoal", e.target.value)}
      />
      {state.error && <p className={styles.error}>{state.error}</p>}
      <div className={styles.btnRow}>
        <Button variant="ghost" onClick={onBack}>Back</Button>
        <Button loading={state.loading} onClick={submit}>Start tracking</Button>
      </div>
    </div>
  );
}

function OnboardingScreen() {
  const hook = useOnboarding();
  const { state, nextStep, prevStep } = hook;

  return (
    <main className={styles.root}>
      {/* Progress dots */}
      <div className={styles.dots}>
        {[1, 2, 3].map((n) => (
          <div key={n} className={[styles.dot, state.step >= n ? styles.activeDot : ""].join(" ")} />
        ))}
      </div>

      <Card className={styles.card}>
        {state.step === 1 && <Step1 {...hook} onNext={nextStep} />}
        {state.step === 2 && <Step2 {...hook} onNext={nextStep} onBack={prevStep} />}
        {state.step === 3 && <Step3 {...hook} onBack={prevStep} />}
      </Card>
    </main>
  );
}

export function OnboardingPage() {
  return (
    <RouteGuard>
      <OnboardingScreen />
    </RouteGuard>
  );
}
