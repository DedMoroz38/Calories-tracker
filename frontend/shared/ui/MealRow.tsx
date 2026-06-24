"use client";

import React from "react";
import { Trash2 } from "lucide-react";
import { IconButton } from "./IconButton";
import styles from "./MealRow.module.css";

interface MealRowProps {
  name:     string;
  calories: number;
  carbs?:   number;
  fat?:     number;
  protein?: number;
  onDelete?: () => void;
  onClick?:  () => void;
}

export function MealRow({ name, calories, carbs, fat, protein, onDelete, onClick }: MealRowProps) {
  return (
    <div className={styles.row}>
      <button className={styles.main} onClick={onClick} disabled={!onClick}>
        <span className={styles.name}>{name}</span>
        <span className={styles.cals}>{calories} kcal</span>
        {(carbs !== undefined || fat !== undefined || protein !== undefined) && (
          <span className={styles.macros}>
            {carbs  !== undefined && <span className={styles.carb}>C {carbs}g</span>}
            {fat    !== undefined && <span className={styles.fat}>F {fat}g</span>}
            {protein !== undefined && <span className={styles.protein}>P {protein}g</span>}
          </span>
        )}
      </button>
      {onDelete && (
        <IconButton label="Delete meal" size="sm" onClick={onDelete} className={styles.del}>
          <Trash2 size={16} />
        </IconButton>
      )}
    </div>
  );
}
