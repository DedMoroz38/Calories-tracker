"use client";

import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Home, BarChart2, PlusCircle } from "lucide-react";
import styles from "./TabBar.module.css";

const TABS = [
  { href: "/home",  icon: Home,      label: "Home" },
  { href: "/stats", icon: BarChart2, label: "Stats" },
] as const;

interface TabBarProps {
  onFabClick: () => void;
}

export function TabBar({ onFabClick }: TabBarProps) {
  const pathname = usePathname();

  return (
    <nav className={styles.bar}>
      {/* Left tab */}
      <TabItem href={TABS[0].href} icon={TABS[0].icon} label={TABS[0].label} active={pathname === TABS[0].href} />

      {/* Centre FAB */}
      <button className={styles.fab} onClick={onFabClick} aria-label="Add food entry">
        <PlusCircle size={28} strokeWidth={1.8} />
      </button>

      {/* Right tab */}
      <TabItem href={TABS[1].href} icon={TABS[1].icon} label={TABS[1].label} active={pathname === TABS[1].href} />
    </nav>
  );
}

function TabItem({
  href,
  icon: Icon,
  label,
  active,
}: {
  href: string;
  icon: React.ElementType;
  label: string;
  active: boolean;
}) {
  return (
    <Link
      href={href}
      className={[styles.tab, active ? styles.activeTab : ""].filter(Boolean).join(" ")}
      aria-current={active ? "page" : undefined}
    >
      <Icon size={22} strokeWidth={active ? 2.2 : 1.8} />
      <span className={styles.tabLabel}>{label}</span>
    </Link>
  );
}
